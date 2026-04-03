variable "environment" {
  
}

variable "project" {
  
}

variable "dynamodb_tables" {
  type = map(object({
    name = string
    arn = string
    id = string
  }))
}

variable "lambda_functions" {
  type = map(object({
    function_name = string
    handler = string
    filename = string
    runtime = string
  }))
}