provider "aws" {
  region = var.region
}

module "dynamodb" {
  source = "../../modules/dynamodb"

  dynamodb_table_configs = var.dynamodb_table_configs
  project                = var.project
  environment            = var.environment
}

module "s3" {
  source = "../../modules/s3"

  project         = var.project
  environment     = var.environment
  region          = var.region
}