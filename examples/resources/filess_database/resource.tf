terraform {
  required_providers {
    filess = {
      source = "app.terraform.io/filess/provider/filessdedicated"
      version = ">=1.0.6"
    }
  }
}

data "filess_engines" "all" {}
data "filess_regions" "all" {}

locals {
  mysql_engine = [
    for engine in data.filess_engines.all.engines :
    engine if engine.name == "MySQL" && engine.version == "8.0"
  ][0]
  
  selected_region = data.filess_regions.all.regions[0]
}

resource "filess_database" "example" {
  organization_slug = "my-org"
  namespace_slug    = "production"
  
  name        = "example-database"
  description = "Example MySQL database"
  
  engine_id = local.mysql_engine.id
  region_id = local.selected_region.id
  
  database_plan {
    billable_items {
      billable_item_id = "6"   # Storage 0.5GiB
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "7"   # Network Bandwidth 1MB/s
      quantity         = 100
    }
    
    billable_items {
      billable_item_id = "8"   # Choose Region
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "10"  # CPU Core 0.25
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "12"  # Database Setup
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "13"  # Memory 0.5GiB
      quantity         = 1
    }
  }
}

output "database_id" {
  value = filess_database.example.id
}

output "connection_string" {
  value = "mysql://${filess_database.example.database_username}:${filess_database.example.database_password}@${filess_database.example.database_hostname}:${filess_database.example.database_service_port}"
  sensitive = true
}

