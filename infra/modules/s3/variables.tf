variable "project" {
  
}

variable "environment" {
  type = string
  # Valid values: "dev", "staging", "prod"
  validation {
    condition = contains(["dev", "staging", "prod"], var.environment)
    error_message = "environment must be one of 'dev', 'staging', or 'prod'."
  }
}


variable "frontend_folder" {
  default  = null
  type     = string
  nullable =true
  description = "Path to frontend dist folder. Null = backend only."
}

variable "region" {
  
}