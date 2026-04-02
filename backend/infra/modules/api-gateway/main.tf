resource "aws_apigatewayv2_api" "example" {
  name          = "${var.project}-${var.environment}-http-api"
  protocol_type = "HTTP"
}

