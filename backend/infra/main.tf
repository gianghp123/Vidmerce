provider "aws" {
  region = "ap-southeast-1"
}

locals {
  db_vars = yamldecode(file("${path.module}/modules/dynamodb/config.yaml"))
}

module "dynamodb" {
  source = "./modules/dynamodb"

  dynamodb_table_configs = local.db_vars.dynamodb_table_configs
  project = var.project
  environment = var.environment
}

module "s3" {
  source = "./modules/s3"
  project = var.project
  environment = var.environment
  frontend_folder = "${path.module}/../../../frontend"
}