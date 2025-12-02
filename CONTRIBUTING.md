# Contributing to kosva

Thanks for your interest! `kosva` is built by DevOps and Security engineers who've experienced the pain of blocked cost optimizations.

## How to Contribute

### Reporting Issues

Found a false positive? False negative? Please open an issue with:

1. **Input data** (anonymized Kubecost JSON)
2. **Expected behavior**
3. **Actual behavior**
4. **Your reasoning** (why should this pass/fail?)

### Adding Security Checks

We need more checks! Common scenarios:

- Multi-tenancy safety validation
- Network policy requirements
- Storage class security
- Image registry security
- RBAC misconfigurations

**To add a check:**

1. Create file in `pkg/checks/your_check.go`
2. Follow existing patterns (`spot_instances.go` is a good example)
3. Add test cases
4. Update validator to use your check

### Adding Compliance Policies

Need PCI-DSS? HIPAA? SOC2?

Create policy templates in `policies/` directory (coming in Week 3-4).

### Code Style

- Follow standard Go conventions
- Add comments for complex logic
- Keep functions small and focused
- Risk scores: 0-10 scale (0=no risk, 10=critical)

## Development Setup
```bash
# Clone
git clone https://github.com/opscart/kosva.git
cd kosva

# Install dependencies
go mod download

# Build
go build -o kosva cmd/kosva/main.go

# Run tests (coming soon)
go test ./...

# Test with examples
./kosva check --input examples/kubecost-sample.json
```

## Pull Request Process

1. Fork the repo
2. Create a feature branch
3. Make your changes
4. Test with example data
5. Update README if needed
6. Submit PR with clear description
