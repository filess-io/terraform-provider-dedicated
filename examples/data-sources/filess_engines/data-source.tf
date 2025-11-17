terraform {
  required_providers {
    filess = {
      source = "app.terraform.io/filess/provider/filessdedicated"
      version = ">=1.0.0"
    }
  }
}

data "filess_engines" "all" {}

# Find MySQL 8.0
locals {
  mysql_engine = [
    for engine in data.filess_engines.all.engines :
    engine if engine.name == "MySQL" && engine.version == "8.0"
  ][0]
}

output "mysql_engine_id" {
  value = local.mysql_engine.id
}

output "all_engines" {
  value = data.filess_engines.all.engines
}

