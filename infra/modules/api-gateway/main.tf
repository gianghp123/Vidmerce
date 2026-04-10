resource "aws_apigatewayv2_api" "this" {
  name          = "${var.project}-${var.environment}-rest-api"
  protocol_type = "HTTP"

  cors_configuration {
    allow_origins = [
      "http://localhost:3000",
      "http://localhost:5173"
    ]

    allow_methods = [
      "GET",
      "POST",
      "PUT",
      "DELETE",
      "OPTIONS"
    ]

    allow_headers = ["*"]

    expose_headers = [
      "content-length"
    ]

    max_age = 86400
  }
}

resource "aws_apigatewayv2_integration" "this" {
  for_each = var.api_routes

  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  payload_format_version = "2.0"
  connection_type        = "INTERNET"
  description            = each.key
  integration_method     = "POST"
  integration_uri        = var.lambda_functions[each.value.lambda_key].invoke_arn
  passthrough_behavior   = "WHEN_NO_MATCH"
}

resource "aws_apigatewayv2_route" "this" {
  for_each = var.api_routes

  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY ${each.value.path}"

  target = "integrations/${aws_apigatewayv2_integration.this[each.key].id}"
}

resource "aws_apigatewayv2_stage" "this" {
  api_id      = aws_apigatewayv2_api.this.id
  name        = "$default"
  auto_deploy = true
}