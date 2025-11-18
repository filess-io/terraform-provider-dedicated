# MySQL Database Example

This example demonstrates how to create a MySQL 8.0 database on filess.io using the Terraform provider.

## Overview

**Provider Source**: `filess-io/dedicated`

This example creates a complete MySQL database instance with:
- MySQL 8.0 engine
- Configurable resources (CPU, memory, storage, bandwidth)
- Automatic region selection
- Connection credentials output
- Full database lifecycle management

## Prerequisites

1. **filess.io Account**: Create an account at [filess.io](https://filess.io)
2. **API Token**: Generate an API token from your filess.io dashboard
3. **Organization & Namespace**: Have an organization and namespace created
4. **Terraform or OpenTofu**: Install [Terraform](https://www.terraform.io/downloads) >= 1.0 or [OpenTofu](https://opentofu.org/) >= 1.0

## Quick Start

### 1. Configure Variables

Copy the example variables file:

```bash
cp terraform.tfvars.example terraform.tfvars
```

Edit `terraform.tfvars` with your values:

```hcl
filess_api_token   = "your-api-token-here"
filess_api_url     = "https://backend.filess.io"
organization_slug  = "your-organization"
namespace_slug     = "your-namespace"
```

⚠️ **Important**: Never commit `terraform.tfvars` to version control!

### 2. Initialize Terraform

```bash
terraform init
```

or with OpenTofu:

```bash
tofu init
```

### 3. Review the Plan

```bash
terraform plan
```

This will show you what resources will be created.

### 4. Apply the Configuration

```bash
terraform apply
```

Confirm by typing `yes` when prompted.

### 5. Get Connection Details

After successful creation, view the outputs:

```bash
# View all outputs
terraform output

# View specific output
terraform output database_hostname
terraform output database_service_port
terraform output database_username

# View sensitive password
terraform output database_password
```

## What Gets Created

This example creates:

1. **MySQL 8.0 Database** with the following resources:
   - **Storage**: 0.5 GiB
   - **CPU**: 0.25 vCore
   - **Memory**: 0.5 GiB
   - **Bandwidth**: 100 MB/s
   - **Region**: Automatically selected from available regions

2. **Outputs** providing:
   - Database ID
   - Database name and status
   - Connection hostname and port
   - Database username and password (sensitive)
   - Engine and region information

## Configuration

### Resource Allocation

You can modify the resource allocation in `main.tf` by adjusting the `quantity` values:

```hcl
database_plan {
  billable_items {
    billable_item_id = "6"   # Storage 0.5GiB
    quantity         = 2     # Double storage to 1 GiB
  }
  
  billable_items {
    billable_item_id = "10"  # CPU Core 0.25
    quantity         = 4     # 1 full vCore
  }
  
  billable_items {
    billable_item_id = "13"  # Memory 0.5GiB
    quantity         = 4     # 2 GiB RAM
  }
  
  # ... other items
}
```

### Billable Items Reference

| Item ID | Description | Unit | Min | Max |
|---------|-------------|------|-----|-----|
| 6 | Storage | 0.5 GiB | 1 | 50 |
| 7 | Network Bandwidth | 1 MB/s | 100 | 400 |
| 8 | Region Selection | - | 1 | 1 |
| 10 | CPU Core | 0.25 vCore | 1 | 16 |
| 12 | Database Setup | - | 1 | 1 |
| 13 | Memory | 0.5 GiB | 1 | 8 |

### Engine Selection

By default, this example selects MySQL 8.0. To use a different engine or version:

```hcl
locals {
  # Use PostgreSQL instead
  selected_engine = [
    for engine in data.filess_engines.all.engines :
    engine if engine.name == "PostgreSQL" && engine.version == "15.2"
  ][0]
}
```

### Region Selection

The example uses the first available region. To select a specific region:

```hcl
locals {
  # Select Spain region
  selected_region = [
    for region in data.filess_regions.all.regions :
    region if region.name == "Spain (pre)"
  ][0]
}
```

## Connecting to Your Database

After creation, use the outputs to connect:

### MySQL Command Line

```bash
mysql -h $(terraform output -raw database_hostname) \
      -P $(terraform output -raw database_service_port) \
      -u $(terraform output -raw database_username) \
      -p$(terraform output -raw database_password)
```

### Connection String

Get the full connection string:

```bash
terraform output database_connection_string
```

Example output:
```
mysql://root:password123@abc123.h.filess.io:12345
```

### In Your Application

```python
# Python example
import mysql.connector

config = {
    'host': 'your-hostname.h.filess.io',
    'port': 12345,
    'user': 'root',
    'password': 'your-password',
    'database': 'your_database'
}

connection = mysql.connector.connect(**config)
```

```javascript
// Node.js example
const mysql = require('mysql2');

const connection = mysql.createConnection({
  host: 'your-hostname.h.filess.io',
  port: 12345,
  user: 'root',
  password: 'your-password',
  database: 'your_database'
});
```

## Payment Required

If your organization doesn't have a payment method configured, you'll see a message like:

```
╔════════════════════════════════════════════════════════════╗
║  ⚠️  PAYMENT REQUIRED                                      ║
╚════════════════════════════════════════════════════════════╝

The database requires payment to continue provisioning.
Please open this URL to complete the Stripe checkout:

  https://checkout.stripe.com/...

Waiting for payment completion...
```

Simply open the URL in your browser, complete the payment, and the provisioning will continue automatically.

## Outputs

This example provides the following outputs:

| Output | Description | Sensitive |
|--------|-------------|-----------|
| `database_id` | Unique database identifier | No |
| `database_name` | Name of the database | No |
| `database_status` | Current status (e.g., "deployed") | No |
| `database_engine` | Database engine name | No |
| `database_version` | Engine version | No |
| `database_region` | Deployment region | No |
| `database_hostname` | Connection hostname | No |
| `database_service_port` | Connection port | No |
| `database_username` | Database username | No |
| `database_password` | Database password | Yes |

## Cleanup

To destroy the database and all resources:

```bash
terraform destroy
```

Confirm by typing `yes` when prompted.

⚠️ **Warning**: This will permanently delete your database and all data!

## Troubleshooting

### Database Not Found Error

If you see "Database with id X not found", the database may have been deleted outside of Terraform. Run:

```bash
terraform apply
```

Terraform will detect the missing resource and recreate it.

### Authentication Errors

Verify your API token:

```bash
# Test the token manually
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://backend.filess.io/api/v1/engines
```

### Region or Engine Not Available

List available options:

```bash
# After terraform plan, check the data sources
terraform console
> data.filess_engines.all.engines
> data.filess_regions.all.regions
```

## Advanced Usage

### With IP Whitelist

```hcl
resource "filess_database" "mysql_test" {
  # ... basic configuration ...
  
  ip_whitelist_ids = ["whitelist-id-1", "whitelist-id-2"]
}
```

### With SSH Keys

```hcl
resource "filess_database" "mysql_test" {
  # ... basic configuration ...
  
  ssh_key_ids = ["ssh-key-id-1"]
}
```

### With Tailscale

```hcl
resource "filess_database" "mysql_test" {
  # ... basic configuration ...
  
  tailscale_config_id = "tailscale-config-id"
}
```

## Files in This Example

- `main.tf` - Main Terraform configuration
- `terraform.tfvars.example` - Example variables file
- `README.md` - This file

## Security Best Practices

1. ✅ Never commit `terraform.tfvars`
2. ✅ Use environment variables for sensitive data
3. ✅ Store state files securely (e.g., S3 with encryption)
4. ✅ Rotate credentials regularly
5. ✅ Use IP whitelisting for production databases
6. ✅ Enable backups (managed by filess.io)

## Next Steps

- Explore other database engines (PostgreSQL, MariaDB, MongoDB)
- Configure IP whitelisting for enhanced security
- Set up monitoring and alerts via filess.io dashboard
- Integrate with your CI/CD pipeline
- Scale resources as your needs grow

## Support

- 📧 Email: support@filess.io
- 📚 Documentation: https://docs.filess.io
- 🐛 Issues: https://github.com/filess/terraform-provider-dedicated/issues
- 💬 Community: [filess.io Discord](https://discord.gg/filess)

## License

This example is part of the terraform-provider-dedicated project, licensed under the MIT License.

