# Contributing to terraform-provider-filess

Thank you for your interest in contributing to the filess Terraform Provider! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions. We're all here to build something great together.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the behavior
- **Expected behavior**
- **Actual behavior**
- **Terraform/OpenTofu version**
- **Provider version**
- **Relevant configuration files** (sanitized)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear title**
- **Provide detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful**
- **List alternative solutions** you've considered

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, descriptive commits
3. **Add or update tests** if applicable
4. **Update documentation** if you're changing functionality
5. **Run tests** to ensure nothing breaks
6. **Submit a pull request**

## Development Setup

### Prerequisites

- Go >= 1.21
- Terraform >= 1.0 or OpenTofu >= 1.0
- Git

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/filess/terraform-provider-filess
cd terraform-provider-filess
```

2. Build the provider:
```bash
go build -o terraform-provider-filess .
```

3. Set up local provider override:

Create `~/.terraformrc` (or `~/.tofurc` for OpenTofu):
```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/filess/filess" = "/path/to/terraform-provider-filess"
  }
  direct {}
}
```

4. Test your changes:
```bash
cd examples/create-mysql-database
terraform init
terraform plan
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v ./internal/resources -run TestResourceDatabase
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and small

### Commit Messages

Write clear commit messages:

```
type(scope): brief description

Longer description if needed explaining:
- What changed
- Why it changed
- Any breaking changes
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

Examples:
```
feat(database): add support for MongoDB
fix(regions): handle null region descriptions
docs: update quick start guide
```

## Documentation

### Updating Documentation

Documentation is generated from:
- Comments in Go code (schema descriptions)
- Templates in `templates/` directory
- Examples in `examples/` directory

To regenerate documentation:

```bash
tfplugindocs generate
```

### Writing Good Documentation

- **Be clear and concise**
- **Include examples** for every feature
- **Explain the "why"**, not just the "what"
- **Keep examples up to date**
- **Use proper formatting** (code blocks, lists, etc.)

## Adding New Resources

When adding a new resource:

1. Create the resource file in `internal/resources/`
2. Implement CRUD operations:
   - `Create`: Create the resource
   - `Read`: Read and update state
   - `Update`: Update the resource
   - `Delete`: Delete the resource
3. Add schema with descriptions
4. Handle errors gracefully
5. Add examples in `examples/resources/`
6. Create documentation template in `templates/resources/`
7. Update CHANGELOG.md

## Adding New Data Sources

When adding a new data source:

1. Create the data source file in `internal/datasources/`
2. Implement `Read` operation
3. Add schema with descriptions
4. Add examples in `examples/data-sources/`
5. Create documentation template in `templates/data-sources/`
6. Update CHANGELOG.md

## Testing Guidelines

### Manual Testing Checklist

Before submitting a PR, test:

- [ ] Resource creation
- [ ] Resource updates
- [ ] Resource deletion
- [ ] Resource import
- [ ] Error handling
- [ ] Documentation accuracy
- [ ] Examples work correctly

### Integration Testing

If you have access to a filess.io account:

```bash
export TF_VAR_filess_api_token="your-token"
export TF_VAR_organization_slug="your-org"
export TF_VAR_namespace_slug="test"

cd examples/create-mysql-database
terraform apply
terraform destroy
```

## Release Process

Maintainers follow this process for releases:

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create and push a git tag:
   ```bash
   git tag -a v1.1.0 -m "Release v1.1.0"
   git push origin v1.1.0
   ```
4. GitHub Actions automatically builds and releases

## Security

### Reporting Security Issues

**DO NOT** open public issues for security vulnerabilities.

Instead, email: security@filess.io

### Security Best Practices

- Never commit sensitive data (tokens, passwords, etc.)
- Always mark sensitive fields as `Sensitive: true` in schemas
- Use environment variables for credentials
- Review code for potential security issues

## Getting Help

- 📧 Email: support@filess.io
- 💬 Discord: [filess.io Discord](https://discord.gg/filess)
- 📚 Documentation: https://docs.filess.io

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes (for significant contributions)

Thank you for contributing to make filess better! 🚀

