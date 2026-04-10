variable "environment" {

}

variable "project" {

}

variable "api_routes" {
  type = map(object({
    path       = string
    lambda_key = string
  }))
}

variable "lambda_functions" {
  type = map(object({
    arn        = string
    invoke_arn = string
    path       = string
  }))
}
