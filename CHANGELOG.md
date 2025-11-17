# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-17

### Added
- **Default API URL**: Changed default `api_url` to `https://backend.filess.io`
- Initial release of the filess Terraform Provider

### Changed
- **Provider Source**: Updated to `app.terraform.io/filess/provider/filessdedicated` for Terraform Cloud/Enterprise compatibility

### Fixed
- **GoReleaser**: Fixed duplicate archive names for ARM architectures (armv6/armv7 now generate unique files)
- **Resource: `filess_database`**
  - Create, read, update, and delete database instances
  - Automatic waiting for database credentials before completion
  - Integrated Stripe checkout handling with user notification
  - Support for IP whitelist and SSH key configuration
  - Tailscale integration support
  - Computed attributes: hostname, port, username, password
  - State reconciliation for externally deleted resources (404 handling)

- **Data Source: `filess_engines`**
  - List all available database engines
  - Filter by name, version, and active status
  - Support for MySQL, PostgreSQL, MariaDB, MongoDB, and more

- **Data Source: `filess_regions`**
  - List all available deployment regions
  - Filter by region name and status

### Features
- 🗄️ Full database lifecycle management
- 🌍 Multi-region deployment support
- 🔧 Multiple database engines
- 💳 Automatic Stripe checkout integration
- 🔐 Automatic credential provisioning
- ✅ State reconciliation
- 🔄 Compatible with Terraform >= 1.0 and OpenTofu >= 1.0

### Documentation
- Complete provider documentation with examples
- Resource and data source documentation
- Quick start guide
- Security best practices
- Contributing guidelines

### Developer Experience
- GoReleaser configuration for multi-platform builds
- Comprehensive `.gitignore` for security
- MIT License
- GitHub-ready repository structure

## [Unreleased]

### Planned
- Support for database backup configuration
- Support for database monitoring and alerts
- Support for database scaling operations
- Additional data sources for metadata
- Enhanced import capabilities

---

## Version History

### Version 1.0.0 (Initial Release)
First stable release of the filess Terraform Provider with full CRUD support for databases, automatic credential management, and Stripe checkout integration.

