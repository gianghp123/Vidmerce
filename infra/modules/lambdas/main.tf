resource "aws_lambda_function" "this" {
  for_each = var.lambda_functions

  function_name = "${var.project}-${var.environment}-${each.value.function_name}"
  handler       = each.value.handler
  runtime       = each.value.runtime
  filename      = each.value.filename
  role          = aws_iam_role.iam_for_lambda.arn

  source_code_hash = filebase64sha256(each.value.filename)

  environment {
    variables = {
      ENVIRONMENT         = var.environment

      S3_BUCKET_NAME = try(var.s3_buckets[each.value.s3_bucket_key])

      DYNAMODB_TABLE = try(
        var.dynamodb_tables[each.value.db_table_key].name,
        ""
      )

      CDN_URL = var.cdn_url
    }
  }
  tags = {
    Application = each.key
  }
}
