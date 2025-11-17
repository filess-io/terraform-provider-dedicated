# Terraform Provider filess.io - Examples

This directory contains practical examples demonstrating how to use the filess Terraform provider.

## Available Examples

### 📁 [create-mysql-database](./create-mysql-database/)

Complete example showing how to create a MySQL 8.0 database with:
- Automatic engine and region selection
- Configurable resource allocation
- Connection credentials output
- Full lifecycle management

**Use this example to**: Learn the basics and get started quickly with a production-ready MySQL database.

## Quick Start

Each example includes:
- ✅ Complete Terraform configuration
- ✅ Variables file template
- ✅ Detailed README
- ✅ Output examples

### Running an Example

1. Navigate to the example directory:
   ```bash
   cd create-mysql-database
   ```

2. Copy and configure variables:
   ```bash
   cp terraform.tfvars.example terraform.tfvars
   # Edit terraform.tfvars with your values
   ```

3. Initialize and apply:
   ```bash
   terraform init
   terraform apply
   ```

4. View outputs:
   ```bash
   terraform output
   ```

## Common Patterns

### Finding Database Engines

```hcl
terraform {
  required_providers {
    filess = {
      source = "app.terraform.io/filess/provider/filessdedicated"
      version = ">=1.0.0"
    }
  }
}

data "filess_engines" "all" {}

# Filter by name and version
locals {
  mysql_8 = [
    for engine in data.filess_engines.all.engines :
    engine if engine.name == "MySQL" && engine.version == "8.0"
  ][0]
  
  postgres_latest = [
    for engine in data.filess_engines.all.engines :
    engine if engine.name == "PostgreSQL"
  ][0]
}
```

### Finding Regions

```hcl
terraform {
  required_providers {
    filess = {
      source = "app.terraform.io/filess/provider/filessdedicated"
      version = ">=1.0.0"
    }
  }
}

data "filess_regions" "all" {}

# Use first available
resource "filess_database" "example" {
  region_id = data.filess_regions.all.regions[0].id
}

# Filter by name
locals {
  eu_region = [
    for region in data.filess_regions.all.regions :
    region if can(regex("Spain|Europe", region.name))
  ][0]
}
```

### Resource Configuration

```hcl
terraform {
  required_providers {
    filess = {
      source = "app.terraform.io/filess/provider/filessdedicated"
      version = ">=1.0.0"
    }
  }
}

resource "filess_database" "example" {
  organization_slug = "my-org"
  namespace_slug    = "production"
  
  name        = "my-database"
  description = "Production database"
  
  engine_id = local.mysql_engine.id
  region_id = local.selected_region.id
  
  database_plan {
    # Storage: 0.5 GiB units (min: 1, max: 50)
    billable_items {
      billable_item_id = "6"
      quantity         = 2  # 1 GiB
    }
    
    # Bandwidth: 1 MB/s units (min: 100, max: 400)
    billable_items {
      billable_item_id = "7"
      quantity         = 100  # 100 MB/s
    }
    
    # Region (required, always 1)
    billable_items {
      billable_item_id = "8"
      quantity         = 1
    }
    
    # CPU: 0.25 vCore units (min: 1, max: 16)
    billable_items {
      billable_item_id = "10"
      quantity         = 4  # 1 vCore
    }
    
    # Database setup (required, always 1)
    billable_items {
      billable_item_id = "12"
      quantity         = 1
    }
    
    # Memory: 0.5 GiB units (min: 1, max: 8)
    billable_items {
      billable_item_id = "13"
      quantity         = 4  # 2 GiB
    }
  }
}
```

### Outputs Pattern

```hcl
output "connection_details" {
  value = {
    host     = filess_database.example.database_hostname
    port     = filess_database.example.database_service_port
    username = filess_database.example.database_username
    password = filess_database.example.database_password
  }
  sensitive = true
}

output "connection_string" {
  value = "${filess_database.example.database_username}:${filess_database.example.database_password}@${filess_database.example.database_hostname}:${filess_database.example.database_service_port}"
  sensitive = true
}
```

## Billable Items Reference

| Item ID | Resource | Unit | Minimum | Maximum |
|---------|----------|------|---------|---------|
| 6 | Storage | 0.5 GiB | 1 (0.5 GiB) | 50 (25 GiB) |
| 7 | Network Bandwidth | 1 MB/s | 100 | 400 |
| 8 | Region Selection | - | 1 | 1 |
| 10 | CPU Core | 0.25 vCore | 1 (0.25 vCore) | 16 (4 vCores) |
| 12 | Database Setup | - | 1 | 1 |
| 13 | Memory | 0.5 GiB | 1 (0.5 GiB) | 8 (4 GiB) |

## Best Practices

### 1. Use Data Sources

Always fetch engines and regions dynamically:

```hcl
data "filess_engines" "all" {}
data "filess_regions" "all" {}
```

Don't hardcode IDs as they may change.

### 2. Mark Sensitive Outputs

Always mark password outputs as sensitive:

```hcl
output "database_password" {
  value     = filess_database.example.database_password
  sensitive = true
}
```

### 3. Use Variables

Make your configuration reusable:

```hcl
variable "environment" {
  type = string
}

resource "filess_database" "app" {
  name = "${var.environment}-database"
  # ...
}
```

### 4. Remote State

For production, use remote state:

```hcl
terraform {
  backend "s3" {
    bucket = "my-terraform-state"
    key    = "filess/databases.tfstate"
    region = "us-east-1"
  }
}
```

### 5. Organize by Environment

```
environments/
├── production/
│   ├── main.tf
│   └── terraform.tfvars
├── staging/
│   ├── main.tf
│   └── terraform.tfvars
└── development/
    ├── main.tf
    └── terraform.tfvars
```

## Security Considerations

1. **Never Commit Credentials**
   - Add `*.tfvars` to `.gitignore`
   - Use environment variables or secret managers

2. **Use IP Whitelisting**
   ```hcl
   resource "filess_database" "secure" {
     ip_whitelist_ids = ["whitelist-id"]
     # ...
   }
   ```

3. **Rotate Credentials**
   - Regularly rotate API tokens
   - Change database passwords periodically

4. **Secure State Files**
   - State files contain sensitive data
   - Use encryption at rest
   - Restrict access with IAM policies

## Testing Examples

Before using in production:

1. Test in a development namespace
2. Verify all outputs are correct
3. Test connection to the database
4. Test destroy operation
5. Review costs in filess.io dashboard

## Troubleshooting

### Common Issues

**"Database not found" after recreation**
- Solution: Run `terraform refresh` or `terraform apply`

**Authentication errors**
- Solution: Verify API token is correct and not expired

**Resource limits exceeded**
- Solution: Check your organization's limits in filess.io dashboard

**Payment required during apply**
- Solution: Complete Stripe checkout at the provided URL

## Contributing Examples

Have a useful example? Contributions are welcome!

1. Create a new directory in `examples/`
2. Include:
   - `main.tf` - Complete configuration
   - `terraform.tfvars.example` - Variables template
   - `README.md` - Detailed documentation
3. Test thoroughly
4. Submit a pull request

## Additional Resources

- 📚 [Provider Documentation](https://registry.terraform.io/providers/filess/filess/latest/docs)
- 🌐 [filess.io Dashboard](https://filess.io/dashboard)
- 📖 [filess.io Documentation](https://docs.filess.io)
- 💬 [Community Discord](https://discord.gg/filess)
- 🐛 [Report Issues](https://github.com/filess/terraform-provider-dedicated/issues)

## Support

Need help? Reach out:
- Email: support@filess.io
- Documentation: https://docs.filess.io
- GitHub Issues: https://github.com/filess/terraform-provider-dedicated/issues

---

**Happy Infrastructure as Code! 🚀**

