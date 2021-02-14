resource "aws_s3_bucket" "destination" {
  bucket = "di-entry-revisions-production"
  acl    = "private"

  tags = {
    Name        = "di-entry-revisions"
    Environment = "production"
  }
}

resource "aws_iam_policy" "firehose_policy" {
  name = "di_entry_revisions_firehose_stream_policy"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid = "",
        Action = [
          "s3:AbortMultipartUpload",
          "s3:GetBucketLocation",
          "s3:GetObject",
          "s3:ListBucket",
          "s3:ListBucketMultipartUploads",
          "s3:PutObject"
        ],
        Resource = [
          "arn:aws:s3:::di-entry-revisions-production",
          "arn:aws:s3:::di-entry-revisions-production/*"
        ]
        Effect = "Allow"
      },
    ]
  })
}

resource "aws_iam_role" "firehose_role" {
  name = "di_entry_revisions_firehose_stream_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "test-attach" {
  role       = aws_iam_role.firehose_role.name
  policy_arn = aws_iam_policy.firehose_policy.arn
}

resource "aws_kinesis_firehose_delivery_stream" "extended_s3_stream" {
  name        = "di-entry-revisions-stream-production"
  destination = "extended_s3"

  extended_s3_configuration {
    bucket_arn          = aws_s3_bucket.destination.arn
    role_arn            = aws_iam_role.firehose_role.arn
    prefix              = "year=!{timestamp:yyyy}/month=!{timestamp:MM}/"
    error_output_prefix = "error/!{firehose:random-string}/!{firehose:error-output-type}/!{timestamp:yyyy/MM/dd}/"
    buffer_interval     = "60"
    buffer_size         = 64
    # cloudwatch_logging_options {
    #   enabled = true
    #   log_group_name = "/aws/kinesisfirehose/di-entry-revisions-stream-production"
    #   log_stream_name = "stream"
    # }

    data_format_conversion_configuration {
      input_format_configuration {
        deserializer {
          open_x_json_ser_de {}
        }
      }

      output_format_configuration {
        serializer {
          parquet_ser_de {}
        }
      }

      schema_configuration {
        database_name = "di"
        role_arn      = "arn:aws:iam::076797644834:role/service-role/AWSGlueServiceRole-test"
        table_name    = "di_entry_revisions_production"
        region        = "us-east-1"
      }
    }
  }
}