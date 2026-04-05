resource "aws_apigatewayv2_api" "this" {
  name = "${var.project}-${var.environment}-rest-api"
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
resource "aws_apigatewayv2_integration" "asset_lambda" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  payload_format_version = "2.0" 
  connection_type           = "INTERNET"
  description               = "Asset service"
  integration_method        = "POST"
  integration_uri           = try(var.lambda_functions["assets"].invoke_arn)
}


resource "aws_apigatewayv2_route" "asset_route_proxy" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/assets/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.asset_lambda.id}"
}

resource "aws_apigatewayv2_route" "asset_route" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/assets"

  target = "integrations/${aws_apigatewayv2_integration.asset_lambda.id}"
}

resource "aws_apigatewayv2_integration" "video_lambda" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  payload_format_version = "2.0" 
  connection_type           = "INTERNET"
  description               = "Asset service"
  integration_method        = "POST"
  integration_uri           = try(var.lambda_functions["videos"].invoke_arn)
  passthrough_behavior      = "WHEN_NO_MATCH"
}


resource "aws_apigatewayv2_route" "video_route_proxy" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/videos/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.video_lambda.id}"
}

resource "aws_apigatewayv2_route" "video_route" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/videos"

  target = "integrations/${aws_apigatewayv2_integration.video_lambda.id}"
}

resource "aws_apigatewayv2_stage" "this" {
  api_id = aws_apigatewayv2_api.this.id
  name   = "$default"
  auto_deploy = true
}