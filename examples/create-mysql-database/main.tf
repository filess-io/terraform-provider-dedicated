terraform {
  required_providers {
    filess = {
      source = "registry.terraform.io/filess-io/dedicated"
      version = ">=1.0.0"
    }
  }
}

provider "filess" {
  api_token = var.filess_api_token
  api_url   = var.filess_api_url
}

variable "filess_api_token" {
  type        = string
  description = "API token for filess.io"
  sensitive   = true
}

variable "filess_api_url" {
  type        = string
  description = "Base URL for filess.io API"
  default     = "https://backend.filess.io"
}

variable "organization_slug" {
  type        = string
  description = "Organization slug"
}

variable "namespace_slug" {
  type        = string
  description = "Namespace slug"
}

# Data sources para obtener engines y regions
data "filess_engines" "all" {}
data "filess_regions" "all" {}

# Encontrar MySQL 8.0 engine
locals {
  mysql_engine = [
    for engine in data.filess_engines.all.engines : engine
    if engine.name == "MySQL" && engine.version == "8.0"
  ][0]
  
  # Usar la primera región disponible
  selected_region = data.filess_regions.all.regions[0]
}

# Crear base de datos MySQL con todos los billable items requeridos
resource "filess_database" "mysql_test" {
  organization_slug = var.organization_slug
  namespace_slug    = var.namespace_slug
  
  name        = "terraform-mysql-db"
  description = "MySQL database created by Terraform provider"
  
  engine_id = local.mysql_engine.id
  region_id = local.selected_region.id
  
  database_plan {
    # Required billable items (según /api/v1/databases/create/metadata)
    billable_items {
      billable_item_id = "6"   # db_storage_500MiB - Storage 0.5GiB (mínimo 1, máximo 50)
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "7"   # network_bandwidth_1M - Network Bandwidth 1MB/s (mínimo 100, máximo 400)
      quantity         = 100
    }
    
    billable_items {
      billable_item_id = "8"   # choose_region - Choose Region (mínimo 1, máximo 1)
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "10"  # cpu_core_500m - CPU Core 0.25 (mínimo 1, máximo 16)
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "12"  # db_instance - Database Setup (mínimo 1, máximo 1)
      quantity         = 1
    }
    
    billable_items {
      billable_item_id = "13"  # memory_500MiB - Memory 0.5GiB (mínimo 1, máximo 8)
      quantity         = 1
    }
  }
}

output "database_id" {
  value = filess_database.mysql_test.id
}

output "database_name" {
  value = filess_database.mysql_test.name
}

output "database_status" {
  value = filess_database.mysql_test.status
}

output "database_engine" {
  value = local.mysql_engine.name
}

output "database_version" {
  value = local.mysql_engine.version
}

output "database_region" {
  value = local.selected_region.name
}

output "database_hostname" {
  value = filess_database.mysql_test.database_hostname
}

output "database_service_port" {
  value = filess_database.mysql_test.database_service_port
}

output "database_username" {
  value = filess_database.mysql_test.database_username
}

output "database_password" {
  value     = filess_database.mysql_test.database_password
  sensitive = true
}

