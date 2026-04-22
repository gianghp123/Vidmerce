variable "project" {
  
}

variable "environment" {
  
}

variable "frontend_folder" {
  default  = null
  type     = string
  nullable =true
  description = "Path to frontend dist folder. Null = backend only."
}

variable "region" {
  
}