provider "aws" {
  region = "ap-southeast-1"
}

locals {
  config_path = "${path.module}/configs"

  db_vars     = yamldecode(file("${local.config_path}/dynamodb.yaml"))
  lambda_vars = yamldecode(file("${local.config_path}/lambdas.yaml"))
  lambda_functions = {
    for k, v in local.lambda_vars.lambda_functions : k => merge(v, {
      filename = "${path.module}/../../${v.filename}"
    })
  }
}

module "dynamodb" {
  source = "../../modules/dynamodb"

  dynamodb_table_configs = local.db_vars.dynamodb_table_configs
  project     = var.project
  environment = var.environment
}

module "s3" {
  source = "../../modules/s3"

  project          = var.project
  environment      = var.environment
  frontend_folder  = "${path.module}/../../../frontend"
}

module "lambda" {
  source = "../../modules/lambdas"

  project           = var.project
  environment       = var.environment
  dynamodb_tables   = module.dynamodb.dynamodb_tables
  lambda_functions  = local.lambda_functions
  s3_buckets = module.s3.s3_bucket_names
  region = var.region
  count = var.is_local ? 0 : 1
}

module "api_gateway" {
  source = "../../modules/api-gateway"
  project = var.project
  environment = var.environment
  lambda_functions = module.lambda.lambda_functions
  count = var.is_local ? 0 : 1
}