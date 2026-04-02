resource "aws_dynamodb_table" "dynamodb-tables" {
  for_each = var.dynamodb_table_configs

  name         = each.key
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = each.value.hash_key
  range_key = try(each.value.range_key, null)

  # Attributes
  dynamic "attribute" {
    for_each = each.value.attributes
    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }

  # Global Secondary Index (optional)
  dynamic "global_secondary_index" {
    for_each = each.value.global_secondary_index != null ? [each.value.global_secondary_index] : []
    content {
      name            = global_secondary_index.value.name
      hash_key        = lookup(
        { for ks in global_secondary_index.value.key_schemas : ks.key_type => ks.attribute_name },
        "HASH"
      )
      range_key       = lookup(
        { for ks in global_secondary_index.value.key_schemas : ks.key_type => ks.attribute_name },
        "RANGE",
        null
      )
      projection_type = global_secondary_index.value.projection_type

      # Only used if PROVISIONED (not PAY_PER_REQUEST, but safe to include conditionally)
      read_capacity  = try(global_secondary_index.value.read_capacity, null)
      write_capacity = try(global_secondary_index.value.write_capacity, null)
    }
  }

  # Tags (optional)
  tags = try({
    Name        = each.value.tags.Name
    Environment = var.environment
  }, null)
}