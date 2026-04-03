output "s3_bucket_names" {
  description = "Map of S3 bucket names"
  value       = { for k, v in aws_s3_bucket.buckets : k => v.id }
}