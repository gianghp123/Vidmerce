resource "aws_apigatewayv2_api" "this" {
  name = "${var.project}-${var.environment}-rest-api"
  protocol_type = "HTTP"
}
resource "aws_apigatewayv2_integration" "asset_lambda" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  description               = "Asset service"
  integration_method        = "POST"
  integration_uri           = try(var.lambda_functions["assets"].invoke_arn)
}


resource "aws_apigatewayv2_route" "asset_route" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/assets/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.asset_lambda.id}"
}


resource "aws_apigatewayv2_integration" "video_lambda" {
  api_id           = aws_apigatewayv2_api.this.id
  integration_type = "AWS_PROXY"

  connection_type           = "INTERNET"
  content_handling_strategy = "CONVERT_TO_TEXT"
  description               = "Asset service"
  integration_method        = "POST"
  integration_uri           = try(var.lambda_functions["videos"].invoke_arn)
  passthrough_behavior      = "WHEN_NO_MATCH"
}


resource "aws_apigatewayv2_route" "video_route" {
  api_id    = aws_apigatewayv2_api.this.id
  route_key = "ANY /api/videos/{proxy+}"

  target = "integrations/${aws_apigatewayv2_integration.video_lambda.id}"
}