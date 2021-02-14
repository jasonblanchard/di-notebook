resource "aws_iam_user" "notebook_production" {
  name = "di-notebook-production"
}

output "iam_user_name" {
    value = aws_iam_user.notebook_production.name
}

output "iam_user_name_arn" {
    value = aws_iam_user.notebook_production.arn
}