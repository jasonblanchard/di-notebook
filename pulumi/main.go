package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/codedeploy"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v4/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		stack := ctx.Stack()
		caller, err := aws.GetCallerIdentity(ctx, nil, nil)
		if err != nil {
			return err
		}

		role, err := iam.NewRole(ctx, "di-apilambda-task-exec-role", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Sid": "",
						"Effect": "Allow",
						"Principal": {
							"Service": "lambda.amazonaws.com"
						},
						"Action": "sts:AssumeRole"
					}]
				}`),
		})
		if err != nil {
			return err
		}

		logPolicy, err := iam.NewRolePolicy(ctx, "di-apilambda-log-policy", &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: pulumi.String(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": [
							"logs:CreateLogGroup",
							"logs:CreateLogStream",
							"logs:PutLogEvents"
						],
						"Resource": "arn:aws:logs:*:*:*"
					}]
				}`),
		})

		eniPolicy, err := iam.NewRolePolicy(ctx, "di-apilambda-eni-policy", &iam.RolePolicyArgs{
			Role: role.Name,
			Policy: pulumi.String(`{
					"Version": "2012-10-17",
					"Statement": [{
						"Effect": "Allow",
						"Action": [
							"ec2:CreateNetworkInterface",
							"ec2:DescribeNetworkInterfaces",
							"ec2:DeleteNetworkInterface"
						],
						"Resource": "*"
					}]
				}`),
		})

		lambdaSourceBucket, err := s3.NewBucket(ctx, fmt.Sprintf("di-notebook-%s", stack), &s3.BucketArgs{
			Acl: pulumi.String("private"),
		})
		if err != nil {
			return err
		}

		ctx.Export("lambdaSourceBucket", lambdaSourceBucket.Bucket)

		// This will eventually be overwritten by CD pipeline
		// Presumes the project has been built
		initialLambdaBuild, err := s3.NewBucketObject(ctx, "examplebucketObject", &s3.BucketObjectArgs{
			Key:    pulumi.String("initial"),
			Bucket: lambdaSourceBucket.ID(),
			Source: pulumi.NewFileAsset("./bin/apilambda.zip"),
		}, pulumi.IgnoreChanges([]string{"source"}))
		if err != nil {
			return err
		}

		sgid := conf.Require("sgid")
		snuseast1a := conf.Require("snuseast1a")
		snuseast1b := conf.Require("snuseast1b")

		apilambdafn, err := lambda.NewFunction(ctx, fmt.Sprintf("di-api-%s", stack), &lambda.FunctionArgs{
			Name:     pulumi.String(fmt.Sprintf("di-notebook-%s", stack)),
			Handler:  pulumi.String("apilambda"),
			Role:     role.Arn,
			Runtime:  pulumi.String("go1.x"),
			S3Bucket: lambdaSourceBucket.ID(),
			S3Key:    initialLambdaBuild.Key,
			Publish:  pulumi.BoolPtr(true),
			VpcConfig: &lambda.FunctionVpcConfigArgs{
				SecurityGroupIds: pulumi.ToStringArray([]string{sgid}),
				SubnetIds:        pulumi.ToStringArray([]string{snuseast1a, snuseast1b}),
			},
		},
			pulumi.DependsOn([]pulumi.Resource{logPolicy, eniPolicy}),
		)
		if err != nil {
			return err
		}

		ctx.Export("apiLambdaName", apilambdafn.Name)

		releaseAlias, err := lambda.NewAlias(ctx, "release", &lambda.AliasArgs{
			Name:            pulumi.String("release"),
			FunctionName:    apilambdafn.Name,
			FunctionVersion: pulumi.String("$LATEST"),
		})
		if err != nil {
			return err
		}

		apigwid := conf.Require("apigwid")

		apigw, err := apigatewayv2.GetApi(ctx, fmt.Sprintf("di-%s", stack), pulumi.ID(apigwid), &apigatewayv2.ApiState{})
		if err != nil {
			return err
		}

		_, err = lambda.NewPermission(ctx, "apilambda-permission", &lambda.PermissionArgs{
			Action:    pulumi.String("lambda:InvokeFunction"),
			Function:  apilambdafn.Name,
			Qualifier: releaseAlias.Name,
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("arn:aws:execute-api:%s:%s:%s/*/*/*", "us-east-1", caller.AccountId, apigw.ID()),
		})
		if err != nil {
			return err
		}

		apilambdaIntegration, err := apigatewayv2.NewIntegration(ctx, "apilambda-integration", &apigatewayv2.IntegrationArgs{
			ApiId:                apigw.ID(),
			IntegrationType:      pulumi.String("AWS_PROXY"),
			IntegrationUri:       releaseAlias.InvokeArn,
			IntegrationMethod:    pulumi.String("POST"),
			PayloadFormatVersion: pulumi.String("1.0"),
		})
		if err != nil {
			return err
		}

		authorizerid := conf.Require("authorizerid")

		target := apilambdaIntegration.ID().OutputState.ApplyT(func(id pulumi.ID) string {
			return fmt.Sprintf("integrations/%s", id)
		}).(pulumi.StringOutput)

		_, err = apigatewayv2.NewRoute(ctx, "routev2", &apigatewayv2.RouteArgs{
			ApiId:             apigw.ID(),
			RouteKey:          pulumi.String("ANY /api/{proxy+}"),
			Target:            target,
			AuthorizerId:      pulumi.ID(authorizerid),
			AuthorizationType: pulumi.String("JWT"),
		})
		if err != nil {
			return err
		}

		// CodeDeploy
		codeDeployApplication, err := codedeploy.NewApplication(ctx, "app", &codedeploy.ApplicationArgs{
			Name:            pulumi.String(fmt.Sprintf("di-notebook-%s", stack)),
			ComputePlatform: pulumi.String("Lambda"),
		})
		if err != nil {
			return err
		}

		codeDeployRole, err := iam.NewRole(ctx, fmt.Sprintf("di-notebook-codedeploy-%s", stack), &iam.RoleArgs{
			AssumeRolePolicy: pulumi.Any(fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v", "{\n", "  \"Version\": \"2012-10-17\",\n", "  \"Statement\": [\n", "    {\n", "      \"Sid\": \"\",\n", "      \"Effect\": \"Allow\",\n", "      \"Principal\": {\n", "        \"Service\": \"codedeploy.amazonaws.com\"\n", "      },\n", "      \"Action\": \"sts:AssumeRole\"\n", "    }\n", "  ]\n", "}\n")),
		})
		if err != nil {
			return err
		}

		_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("di-notebook-codedeploy-%s", stack), &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSCodeDeployRoleForLambda"),
			Role:      codeDeployRole.Name,
		})
		if err != nil {
			return err
		}

		// TODO: Scope this down
		_, err = iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("di-notebook-codedeploy-s3-full%s", stack), &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
			Role:      codeDeployRole.Name,
		})
		if err != nil {
			return err
		}

		_, err = codedeploy.NewDeploymentGroup(ctx, fmt.Sprintf("di-notebook-codedeploy-%s", stack), &codedeploy.DeploymentGroupArgs{
			AppName:              codeDeployApplication.Name,
			DeploymentGroupName:  pulumi.String("release"),
			DeploymentConfigName: pulumi.String("CodeDeployDefault.LambdaAllAtOnce"),
			ServiceRoleArn:       codeDeployRole.Arn,
			DeploymentStyle: &codedeploy.DeploymentGroupDeploymentStyleArgs{
				DeploymentOption: pulumi.String("WITH_TRAFFIC_CONTROL"),
				DeploymentType:   pulumi.String("BLUE_GREEN"),
			},
		})

		if err != nil {
			return err
		}

		deployspecBucket, err := s3.NewBucket(ctx, fmt.Sprintf("di-notebook-codedeploy-deployspec-%s", stack), &s3.BucketArgs{
			Acl: pulumi.String("private"),
			Versioning: &s3.BucketVersioningArgs{
				Enabled: pulumi.Bool(true),
			},
		})

		ctx.Export("deployspecbucket", deployspecBucket.Bucket)

		return nil
	})
}
