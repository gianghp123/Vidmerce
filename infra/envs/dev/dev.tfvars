dynamodb_tables = {
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