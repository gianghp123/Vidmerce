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
    s3_bucket_key    = string
    db_table_key  = string
  }))
}

variable "s3_buckets" {
  type        = map(string)
  description = "Map of S3 bucket names passed from S3 module"
  default     = {}
}

variable "cdn_url" {
  type    = string
  default = "http://localhost:4566" # Mặc định cho LocalStack
}

variable "region" {
  type = string
}

variable "localstack_host" {
  type = string
}