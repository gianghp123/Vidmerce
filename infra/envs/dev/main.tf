provider "aws" {
  access_key                  = "test"
  secret_key                  = "test"
  region                      = var.region


  # only required for non virtual hosted-style endpoint use case.
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs#s3_use_path_style
  s3_use_path_style           = true
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true


  endpoints {
    s3             = "http://localhost:4566"
    dynamodb       = "http://localhost:4566"
  }
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