output "s3_bucket_names" {
  description = "Map of S3 bucket names"
  value       = { for k, v in aws_s3_bucket.buckets : k => v.id }
}

output "s3_public_url" {
  value = aws_s3_bucket_website_configuration.vite_site.website_endpoint
}