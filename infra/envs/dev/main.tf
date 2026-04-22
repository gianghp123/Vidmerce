provider "aws" {
  region = var.region
}

locals {
  config_path = "${path.module}/configs"

  db_vars     = yamldecode(file("${local.config_path}/dynamodb.yaml"))
  # lambda_vars = yamldecode(file("${local.config_path}/lambdas_api.yaml"))
  # api_routes  = yamldecode(file("${local.config_path}/api_routes.yaml"))
  # lambda_functions = {
  #   for k, v in local.lambda_vars.lambda_functions : k => merge(v, {
  #     filename = "${path.module}/../../${v.filename}"
  #   })
  # }
}

module "dynamodb" {
  source = "../../modules/dynamodb"

  dynamodb_table_configs = local.db_vars.dynamodb_table_configs
  project                = var.project
  environment            = var.environment
}

module "s3" {
  source = "../../modules/s3"

  project         = var.project
  environment     = var.environment
  # frontend_folder = "${path.module}/../../../frontend/dist"
  region          = var.region
}

# module "lambda" {
#   source = "../../modules/lambdas"

#   project          = var.project
#   environment      = var.environment
#   dynamodb_tables  = module.dynamodb.dynamodb_tables
#   lambda_functions = local.lambda_functions
#   s3_buckets       = module.s3.s3_bucket_names
#   region           = var.region
#   count            = var.is_local ? 0 : 1
#   cdn_url          = module.s3.s3_asset_storage_url
# }

# module "api_gateway" {
#   source           = "../../modules/api-gateway"
#   project          = var.project
#   environment      = var.environment
#   lambda_functions = module.lambda[0].lambda_functions
#   api_routes       = local.api_routes.routes
#   count            = var.is_local ? 0 : 1
# }