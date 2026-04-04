variable "environment" {
  type = string
  default = "development"
}

variable "project" {
  type = string
  default = "vidmerce"
}

variable "region" {
  type = string
  default = "ap-southeast-1"
}

variable "is_local" {
  type = bool
  default = false
}