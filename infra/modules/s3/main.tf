locals {
  env_needs_frontend = var.environment == "staging" || var.environment == "prod"

  files = local.env_needs_frontend ? fileset(var.frontend_folder, "**") : []

  bucket_config = merge(
    { asset-storage = "Asset Bucket" },
    local.env_needs_frontend ? { static-hosting = "Static Hosting Bucket" } : {}
  )

  mime_types = {
    html  = "text/html"
    css   = "text/css"
    js    = "application/javascript"
    mjs   = "application/javascript"
    json  = "application/json"
    png   = "image/png"
    jpg   = "image/jpeg"
    jpeg  = "image/jpeg"
    gif   = "image/gif"
    webp  = "image/webp"
    svg   = "image/svg+xml"
    ico   = "image/x-icon"
    woff  = "font/woff"
    woff2 = "font/woff2"
    ttf   = "font/ttf"
    otf   = "font/otf"
    txt   = "text/plain"
    map   = "application/json"
  }
}
resource "aws_s3_bucket" "buckets" {
  for_each = local.bucket_config

  bucket        = "${var.project}-${var.environment}-${each.key}"
  force_destroy = true

  tags = {
    Name        = each.value
    Environment = var.environment
  }
}

resource "aws_s3_bucket_cors_configuration" "asset_storage_cors" {
  bucket = aws_s3_bucket.buckets["asset-storage"].id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET"]
    allowed_origins = ["*"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

resource "aws_s3_bucket_public_access_block" "bucket_policies" {
  for_each = aws_s3_bucket.buckets

  bucket = each.value.id

  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}

resource "aws_s3_bucket_policy" "asset_storage_policy" {
  bucket = aws_s3_bucket.buckets["asset-storage"].id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = "*"
      Action    = ["s3:GetObject"]
      Resource  = "${aws_s3_bucket.buckets["asset-storage"].arn}/*"
    }]
  })
}

resource "aws_s3_bucket_website_configuration" "vite_site" {
  count  = local.env_needs_frontend ? 1 : 0
  bucket = aws_s3_bucket.buckets["static-hosting"].id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"
  }
}

resource "aws_s3_object" "upload_files" {
  for_each = { for file in local.files : file => file }

  bucket       = aws_s3_bucket.buckets["static-hosting"].id
  key          = each.key
  source       = "${var.frontend_folder}/${each.value}"
  etag         = filemd5("${var.frontend_folder}/${each.value}")

  content_type = lookup(
    local.mime_types,
    lower(element(reverse(split(".", each.value)), 0)),
    "application/octet-stream"
  )

  cache_control = can(regex("\\.[a-f0-9]{8,}\\.", each.value)) ? "public, max-age=31536000, immutable" : "no-cache"
}