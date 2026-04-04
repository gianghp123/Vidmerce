output "lambda_functions" {
  value = { for k, v in aws_lambda_function.this: k => {
    arn: v.arn
    invoke_arn: v.invoke_arn
    path: ""
  }}
}