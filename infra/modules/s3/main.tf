locals {
  files  = fileset(var.frontend_folder, "**")
}

resource "aws_s3_bucket" "buckets" {
  for_each = {
    static-hosting = "Static Hosting Bucket"
    asset-storage  = "Asset Bucket"
  }

  bucket = "${var.project}-${var.environment}-${each.key}"

  tags = {
    Name        = each.value
    Environment = var.environment
  }
}

resource "aws_s3_bucket_public_access_block" "bucket-policies" {
  for_each = aws_s3_bucket.buckets

  bucket = each.value.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_website_configuration" "cdn" {
  bucket = aws_s3_bucket.buckets["static-hosting"].id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html" # Using index.html for errors to support SPA routing
  }
}


resource "aws_s3_object" "upload_files" {
  for_each = { for file in local.files: file => file}

  bucket = aws_s3_bucket.buckets["static-hosting"].id

  key = each.key

  source = "${var.frontend_folder}/${each.value}"
  etag = filemd5("${var.frontend_folder}/${each.value}")

  content_type = lookup(
    {
      html = "text/html"
      css  = "text/css"
      js   = "application/javascript"
      json = "application/json"
      png  = "image/png"
      jpg  = "image/jpeg"
      jpeg = "image/jpeg"
      svg  = "image/svg+xml"
    },
    lower(element(split(".", each.value), length(split(".", each.value)) - 1)),
    "binary/octet-stream"
  )
}