# staging.tfvars
api_routes = {
  assets_root = {
    path       = "/api/assets",
    lambda_key = "assets"
  }
  assets = {
    path       = "/api/assets/{proxy+}",
    lambda_key = "assets"
  }
  webhooks = {
    path       = "/api/webhooks/{proxy+}",
    lambda_key = "webhooks"
  }
}

dynamodb_table_configs = {
  MediaProjectTable = {
    hash_key = "Pk"
    range_key = "Sk"
    attributes = [
      { name = "Pk",    type = "S" },
      { name = "Sk",    type = "S" },
      { name = "Gsi1Pk", type = "S" },
      { name = "Gsi1Sk", type = "S" },
    ]
    global_secondary_index = {
      name          = "GSI1"
      key_schemas = [
        { attribute_name = "Gsi1Pk", key_type = "HASH" },
        { attribute_name = "Gsi1Sk", key_type = "RANGE" },
      ]
      projection_type = "ALL"
    }
    tags = {
      Name = "MediaProjectTable"
    }
  }
}

lambda_functions = {
  assets = {
    function_name = "assets_lambda"
    handler       = "bootstrap"
    filename      = "artifacts/lambdas/assets/function.zip"
    runtime       = "provided.al2023"
    s3_bucket_key = "asset-storage"
    db_table_key  = "MediaProjectTable"
  }
  webhooks = {
    function_name = "webhooks_lambda"
    handler       = "bootstrap"
    filename      = "artifacts/lambdas/webhooks/function.zip"
    runtime       = "provided.al2023"
    s3_bucket_key = null
    db_table_key  = null
  }
}