terraform {
  required_providers {
    filess = {
      source = "app.terraform.io/filess/provider/filessdedicated"
      version = ">=1.0.0"
    }
  }
}

data "filess_regions" "all" {}

# Use first available region
output "first_region" {
  value = data.filess_regions.all.regions[0]
}

output "all_regions" {
  value = data.filess_regions.all.regions
}

# Filter EU regions
locals {
  eu_regions = [
    for region in data.filess_regions.all.regions :
    region if can(regex("Europe|Spain", region.name))
  ]
}

output "eu_regions" {
  value = local.eu_regions
}

