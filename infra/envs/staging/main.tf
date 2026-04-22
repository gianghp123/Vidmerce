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
  frontend_folder = "${path.module}/../../../frontend/dist"
  region          = var.region
}

module "lambda" {
  source = "../../modules/lambdas"

  project          = var.project
  environment      = var.environment
  dynamodb_tables  = module.dynamodb.dynamodb_tables
  lambda_functions = var.lambda_functions
  s3_buckets       = module.s3.s3_bucket_names
  region           = var.region
  cdn_url          = module.s3.s3_asset_storage_url
}

module "api_gateway" {
  source           = "../../modules/api-gateway"
  project          = var.project
  environment      = var.environment
  lambda_functions = module.lambda[0].lambda_functions
  api_routes       = var.api_routes
}