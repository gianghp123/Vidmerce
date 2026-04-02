output "dynamodb_tables" {
  description = "Map of table names to their ARNs"
  value = {
    for k, table in aws_dynamodb_table.dynamodb-tables : k => {
      arn  = table.arn
      id   = table.id
      name = table.name
    }
  }
}