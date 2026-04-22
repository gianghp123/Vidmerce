output "s3_bucket_names" {
  description = "Map of S3 bucket names"
  value       = { for k, v in aws_s3_bucket.buckets : k => v.id }
}

output "s3_public_url" {
  description = "The website endpoint for static hosting (null if no frontend)"
  # Using try() or one() here is best practice for resources using 'count'
  value       = try(aws_s3_bucket_website_configuration.vite_site[0].website_endpoint, null)
}

output "s3_asset_storage_url" {
  description = "The regional URL for the asset storage bucket"
  # Safe to access directly because asset-storage always exists in your bucket_config
  value       = "https://${aws_s3_bucket.buckets["asset-storage"].bucket}.s3.${var.region}.amazonaws.com"
}