# filess.io Dedicated Database Provider

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/filess/terraform-provider-dedicated)](https://goreportcard.com/report/github.com/filess/terraform-provider-dedicated)

This repository contains the **only** working example you need to provision a dedicated MySQL database on [filess.io](https://filess.io) using Terraform or OpenTofu.  
The README below mirrors the exact configuration from `examples/create-mysql-database`, so you can copy/paste and get a successful deployment.

---

## 1. Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| Terraform **or** OpenTofu | ≥ 1.0 | All commands below use `tofu`, but Terraform works the same way |
| Go (optional) | ≥ 1.21 | Only required if you want to build the provider locally |
| filess.io account | – | You need an organization, namespace and API token |

### Required inputs

- `filess_api_token` – generate it from the filess.io dashboard
- `organization_slug` – your organization identifier (e.g. `acme`)
- `namespace_slug` – the namespace where databases will live (e.g. `production`)

We recommend storing those values in a `terraform.tfvars` file (see step 3).

---

## 2. Clone the repo & open the working example

```bash
git clone https://github.com/filess-io/terraform-provider-dedicated.git
cd terraform-provider-dedicated/examples/create-mysql-database
```

The directory already contains a **fully working** `main.tf` plus outputs.

---

## 3. Provide your credentials

You can either export environment variables or create a `terraform.tfvars` file.

### Option A – Environment variables (recommended for CI)

```bash
export TF_VAR_filess_api_token="your-api-token"
export TF_VAR_organization_slug="your-org"
export TF_VAR_namespace_slug="your-namespace"
```

### Option B – `terraform.tfvars`

Copy the template and fill your values:

```bash
cp terraform.tfvars.example terraform.tfvars
```

```hcl
filess_api_token = "your-api-token"
organization_slug = "your-org"
namespace_slug = "your-namespace"
```

> `filess_api_url` defaults to `https://backend.filess.io`, override it only if filess support asks you to point to another environment.

---

## 4. Understand the configuration

This is the **exact** configuration that works today (abridged for clarity):

```hcl
terraform {
  required_providers {
    filess = {
      source  = "filess-io/dedicated"
      version = ">= 1.0.0"
    }
  }
}

provider "filess" {
  api_token = var.filess_api_token
  api_url   = var.filess_api_url
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

resource "filess_database" "mysql_test" {
  organization_slug = var.organization_slug
  namespace_slug    = var.namespace_slug

  name        = "terraform-mysql-db"
  description = "MySQL database created by Terraform provider"

  engine_id = local.mysql_engine.id
  region_id = local.selected_region.id

  database_plan {
    billable_items { billable_item_id = "6"  quantity = 1   } # Storage 0.5GiB
    billable_items { billable_item_id = "7"  quantity = 100 } # Network 100 MB/s
    billable_items { billable_item_id = "8"  quantity = 1   } # Choose Region
    billable_items { billable_item_id = "10" quantity = 1   } # CPU 0.25 core
    billable_items { billable_item_id = "12" quantity = 1   } # DB setup
    billable_items { billable_item_id = "13" quantity = 1   } # Memory 0.5 GiB
  }
}

output "database_hostname"        { value = filess_database.mysql_test.database_hostname }
output "database_service_port"    { value = filess_database.mysql_test.database_service_port }
output "database_username"        { value = filess_database.mysql_test.database_username }
output "database_password"        { value = filess_database.mysql_test.database_password  sensitive = true }
```

Key aspects:

- **Provider source**: `filess-io/dedicated`
- **Billable items**: Use the exact IDs shown above or consult `/api/v1/databases/create/metadata`
- **Credential outputs**: hostname, port, username and password are computed once the DB is deployed
- **Stripe checkout**: if payment is required, the provider halts and prints the checkout URL directly in the terminal (no `TF_LOG` required)

---

## 5. Run it

From `examples/create-mysql-database/`:

```bash
tofu init    # or terraform init
tofu plan    # verify configuration
tofu apply   # provision the database
```

During `apply` you may see a message similar to:

```
⚠️  PAYMENT REQUIRED
Please open this URL to complete the Stripe checkout:
  https://checkout.stripe.com/...
```

Open the URL, pay, and the provider will keep waiting until filess marks the database as `deployed`.

---

## 6. Outputs & next steps

After `tofu apply` finishes you will get:

```
Outputs:

database_hostname        = "mysql-1234.backend.filess.io"
database_service_port    = 3306
database_username        = "root"
database_password        = (sensitive value)
database_id              = "db_abc123"
database_status          = "deployed"
database_region          = "us-east-1"
```

Use those values directly in your applications or CI pipelines to connect to the database.

---

## 7. Troubleshooting checklist

| Symptom | Fix |
|---------|-----|
| `Invalid provider source string` | Ensure `source = "filess-io/dedicated"` |
| `Could not retrieve list of available versions` | Run `tofu init -upgrade` or check internet access |
| `401 Unauthorized` | Confirm `filess_api_token` is valid and belongs to the target organization |
| Stripe URL only shows with `TF_LOG` | Fixed – the provider prints the message to `/dev/tty` automatically |
| Database deleted outside Terraform | Provider detects 404s and will recreate on next `apply` |

If you run into issues, run with `TF_LOG=DEBUG tofu apply` and/or open an [issue](https://github.com/filess-io/terraform-provider-dedicated/issues) including the log excerpt.

---

## 8. Developing the provider locally (optional)

```bash
git clone https://github.com/filess-io/terraform-provider-dedicated.git
cd terraform-provider-dedicated
go build -o terraform-provider-dedicated
```

Create `~/.terraformrc` or `~/.tofurc`:

```hcl
provider_installation {
  dev_overrides {
    "filess-io/dedicated" = "/absolute/path/to/terraform-provider-dedicated"
  }
  direct {}
}
```

Then re-run `tofu init` inside your example directory and Terraform/OpenTofu will use the local binary.

---

## 9. Contributing & Support

- 📚 Examples: see the `examples/` folder (the MySQL one is production-ready)
- 🐛 Issues & feature requests: [github.com/filess-io/terraform-provider-dedicated/issues](https://github.com/filess-io/terraform-provider-dedicated/issues)
- 📧 Support: support@filess.io
- 💬 Community: [Discord](https://discord.gg/filess)

Pull requests are welcome! Please open an issue first if you plan to work on new resources or data sources so we can coordinate efforts.

---

## 10. License

MIT © filess.io – see [LICENSE](LICENSE).

---

Made with ❤️ by the filess.io team.

