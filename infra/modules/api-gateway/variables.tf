variable "environment" {
  
}

variable "project" {
  
}

variable "lambda_functions" {
  type = map(object({
    arn = string
    invoke_arn = string
    path = string
  }))
}
