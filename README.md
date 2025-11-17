# Terraform Provider for filess.io

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/filess/terraform-provider-filess)](https://goreportcard.com/report/github.com/filess/terraform-provider-filess)

The filess provider enables you to create, manage, and configure dedicated databases on the [filess.io](https://filess.io) platform using Infrastructure as Code.

## Features

- 🗄️ **Database Management**: Create and manage dedicated databases with custom configurations
- 🌍 **Multi-Region**: Deploy databases across multiple regions
- 🔧 **Multiple Engines**: Support for MySQL, PostgreSQL, MariaDB, MongoDB, and more
- 💳 **Automatic Billing**: Integrated Stripe checkout for seamless payment processing
- 🔐 **Credential Management**: Automatic retrieval of connection credentials
- ✅ **State Reconciliation**: Handles external changes and deleted resources gracefully
- 🔄 **Terraform & OpenTofu**: Compatible with both tools

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0 or [OpenTofu](https://opentofu.org/) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)
- filess.io account with API token

## Using the Provider

### Installation

The provider is available on the [Terraform Registry](https://registry.terraform.io/providers/filess/filess/latest).

```hcl
terraform {
  required_providers {
    filess = {
      source  = "app.terraform.io/filess/provider/filessdedicated"
      version = "~> 1.0"
    }
  }
}

provider "filess" {
  api_token = var.filess_api_token
  api_url   = "https://backend.filess.io" # Optional, defaults to production API
}
```

### Quick Start

```hcl
# Get available engines and regions
data "filess_engines" "all" {}
data "filess_regions" "all" {}

# Find MySQL 8.0
locals {
  mysql_engine = [
    for engine in data.filess_engines.all.engines :
    engine if engine.name == "MySQL" && engine.version == "8.0"
  ][0]
}

# Create a MySQL database
resource "filess_database" "production" {
  organization_slug = "my-org"
  namespace_slug    = "production"
  
  name        = "prod-mysql-db"
  description = "Production MySQL database"
  
  engine_id = local.mysql_engine.id
  region_id = data.filess_regions.all.regions[0].id
  
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

# Output connection details
output "database_connection" {
  value = {
    host     = filess_database.production.database_hostname
    port     = filess_database.production.database_service_port
    username = filess_database.production.database_username
    password = filess_database.production.database_password
  }
  sensitive = true
}
```

## Documentation

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/filess/filess/latest/docs).

### Resources

- [filess_database](https://registry.terraform.io/providers/filess/filess/latest/docs/resources/database) - Manage database instances

### Data Sources

- [filess_engines](https://registry.terraform.io/providers/filess/filess/latest/docs/data-sources/engines) - Get available database engines
- [filess_regions](https://registry.terraform.io/providers/filess/filess/latest/docs/data-sources/regions) - Get available deployment regions

## Development

### Building from Source

```bash
git clone https://github.com/filess/terraform-provider-filess
cd terraform-provider-filess
go build -o terraform-provider-filess
```

### Local Development

Create a `.terraformrc` file in your home directory:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/filess/filess" = "/path/to/terraform-provider-filess"
  }
  direct {}
}
```

Then build the provider:

```bash
go build -o terraform-provider-filess .
```

### Running Tests

```bash
go test ./...
```

### Generating Documentation

Documentation is generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs):

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
tfplugindocs generate
```

## Examples

See the [examples](./examples) directory for complete working examples:

- [create-mysql-database](./examples/create-mysql-database) - Complete MySQL database setup

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

### Development Guidelines

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-feature`)
3. Commit your changes (`git commit -am 'Add new feature'`)
4. Push to the branch (`git push origin feature/my-feature`)
5. Open a Pull Request

## Security

### Reporting Security Issues

Please report security vulnerabilities to: security@filess.io

### Best Practices

- **Never commit `terraform.tfvars`** - Use environment variables or secret managers
- **Use remote state** - Store state securely in S3, Terraform Cloud, etc.
- **Rotate API tokens** - Regularly rotate your filess.io API tokens
- **Review state files** - State files contain sensitive data; handle with care

## Support

- 📧 Email: support@filess.io
- 📚 Documentation: https://docs.filess.io
- 🐛 Issues: https://github.com/filess/terraform-provider-filess/issues
- 💬 Community: [filess.io Discord](https://discord.gg/filess)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with:
- [Terraform Plugin SDK](https://github.com/hashicorp/terraform-plugin-sdk)
- [Go](https://golang.org/)

---

Made with ❤️ by the [filess.io](https://filess.io) team

