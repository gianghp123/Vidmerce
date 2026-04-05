resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

# ---------------------------------------------------------
# 2. S3 MANAGED POLICY
# ---------------------------------------------------------
resource "aws_iam_policy" "s3_policy" {
  name        = "lambda_s3_managed_policy"
  description = "Allow lambda to access S3 buckets"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:GetObject",
          "s3:DeleteObject",
          "s3:ListBucket",
          "s3:PutObject",
          "s3:HeadObject"
        ]
        Effect   = "Allow"
        Resource = "*" 
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "s3_attach" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.s3_policy.arn
}

resource "aws_iam_policy" "dynamodb_policy" {
  name        = "lambda_dynamodb_managed_policy"
  description = "Allow lambda to access DynamoDB tables and indexes"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
          "dynamodb:Query",
          "dynamodb:Scan"
        ]
        Resource = flatten([
          for table in var.dynamodb_tables : [
            table.arn,
            "${table.arn}/index/*"
          ]
        ])
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "dynamodb_attach" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.dynamodb_policy.arn
}

resource "aws_iam_role_policy_attachment" "lambda_basic_logging" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}