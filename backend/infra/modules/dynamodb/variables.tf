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

variable "environment" {
  
}

variable "project" {
  
}