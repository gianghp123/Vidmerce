variable "environment" {
  type = string
  default = "staging"
  validation {
    condition = contains(["dev", "staging", "prod"], var.environment)
    error_message = "environment must be one of 'dev', 'staging', or 'prod'."
  }
}

variable "project" {
  type = string
  default = "vidmerce"
}

variable "region" {
  type = string
  default = "ap-southeast-1"
}

variable "dynamodb_table_configs" {
  description = "List of dynamo db table's configurations"
  type = map(object({
    hash_key       = string
    range_key      = string
    read_capacity  = optional(number, 20)
    write_capacity = optional(number, 20)
    attributes = list(object({
      name = string
      type = string #S, N
    }))
    global_secondary_index = optional(object({
      name = string
      key_schemas = list(object({
        attribute_name = string
        key_type       = string #HASH, RANGE
      }))
      write_capacity  = optional(number, 10)
      read_capacity   = optional(number, 10)
      projection_type = optional(string, "ALL")
    }))

    tags = optional(object({
      Name        = string
    }))
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

variable "api_routes" {
  type = map(object({
    path       = string
    lambda_key = string
  }))
}
