output "asset_storage_bucket_name" {
  value = module.s3.s3_bucket_names["asset-storage"]
}