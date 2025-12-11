cat > CONTRIBUTING.md << 'EOF'
# Contributing to kosva

Thanks for your interest in contributing! kosva exists because cost optimization tools don't validate security, and security tools don't consider cost.

## Ways to Contribute

### 1. Report Issues

Found a bug? False positive? False negative?

**Please include:**
- Input data (anonymized Kubecost JSON)
- Expected behavior
- Actual behavior
- kosva version: `kosva version`

### 2. Add Security Checks

We need more CIS Kubernetes Benchmark checks! Current coverage:

- [x] CIS 5.7 - Spot Instance Suitability
- [x] CIS 5.10 - Resource Limit Safety
- [x] CIS 5.7.3 - Multi-Tenancy Isolation
- [x] CIS 5.4.1 - Storage Security
- [ ] CIS 5.3 - Network Policy Requirements
- [ ] CIS 5.4.2 - Immutable Root Filesystem
- [ ] CIS 5.2 - Pod Security Standards
- [ ] CIS 5.1 - RBAC Configuration

**To add a check:**

1. Create policy file: `policies/cis-benchmark/your-check.yaml`
2. Follow the [Policy Creation Guide](docs/policy-guide.md)
3. Add test cases: `examples/your-check-test.json`
4. Test: `kosva check --input examples/your-check-test.json`
5. Submit PR

### 3. Add Compliance Templates

Need PCI-DSS, HIPAA, SOC2 policies? We're building a library!

**Structure:**
```
policies/compliance/
├── pci-dss.yaml
├── hipaa.yaml
├── soc2.yaml
└── gdpr.yaml
```

### 4. Improve Documentation

- Tutorial videos
- Blog posts
- Integration guides
- Real-world examples

### 5. Integrate with Other Tools

Current: Kubecost  
Wanted: OpenCost, Cloud provider cost tools, FinOps platforms

## Development Setup
```bash
# Clone
git clone https://github.com/opscart/kosva.git
cd kosva

# Build
go build -o kosva cmd/kosva/main.go

# Run tests
go test ./...

# Test with examples
./kosva check --input examples/kubecost-api-format.json
```

## Code Style

- Follow standard Go conventions
- Add comments for complex logic
- Keep functions small and focused
- Use descriptive variable names

## Pull Request Process

1. Fork the repo
2. Create feature branch: `git checkout -b feature/your-feature`
3. Make changes
4. Test thoroughly
5. Update docs if needed
6. Commit: `git commit -m "Add feature: description"`
7. Push: `git push origin feature/your-feature`
8. Open Pull Request

**PR Title Format:**
- `feat: Add multi-tenancy validation`
- `fix: Correct Spot instance risk scoring`
- `docs: Update Kubecost integration guide`
- `refactor: Simplify policy engine`

## Testing Checklist

Before submitting PR:

- [ ] Code compiles: `go build`
- [ ] Tests pass: `go test ./...`
- [ ] Example files work: `./kosva check --input examples/*.json`
- [ ] Policy listing works: `./kosva policies`
- [ ] Documentation updated
- [ ] No hardcoded values

## Questions?

- Open an issue
- Email: [Your contact]
- Website: [opscart.com](https://opscart.com)

## License

By contributing, you agree your contributions will be licensed under the MIT License.
EOF