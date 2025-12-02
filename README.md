# kosva - Kubernetes Optimization Security Validator

**Stop implementing unsafe cost optimizations.**

`kosva` validates Kubernetes cost optimization recommendations against security policies and compliance frameworks before you implement them.

## The Problem

Cost optimization tools like Kubecost suggest ways to reduce cloud spend. But many optimizations create security vulnerabilities:

- **Spot instances** → Critical services can be interrupted with 2-minute notice
- **Under-provisioned resources** → OOM kills create denial-of-service vulnerabilities
- **Multi-tenancy** → Compliance violations in regulated industries

Security teams block these recommendations, wasting engineering time.

## The Solution

`kosva` validates recommendations **BEFORE** implementation:
```bash
$ ./kosva check --input kubecost-recommendations.json

══════════════════════════════════════════════════════════════════════
  KOSVA VALIDATION REPORT
══════════════════════════════════════════════════════════════════════

Total Recommendations: 2
Potential Savings: $1200/month

❌ BLOCKED RECOMMENDATIONS (1):
──────────────────────────────────────────────────────────────────────
1. deployment/payment-api → Would save $800/month
   Type: spot-instance
   ✗ CIS 5.7 - Spot Instance Suitability [HIGH] (Risk: 8.5/10)
     ⚠️  Workload appears to be critical. Spot instances can be interrupted.
     💡 Use Reserved Instances for guaranteed availability.

   🔄 SAFE ALTERNATIVE:
      Use Reserved Instances → ~$240/month savings (compliant)

✅ SAFE RECOMMENDATIONS (1):
──────────────────────────────────────────────────────────────────────
1. deployment/logging-service → $400/month
   Type: right-size
   ✓ CIS 5.10 - Resource Limit Safety (Risk: 1.5/10)

══════════════════════════════════════════════════════════════════════
SUMMARY: $400/$1200 safe (33% of potential savings)
══════════════════════════════════════════════════════════════════════
```

## Quick Start
```bash
# Install Go 1.21+
# Clone and build
git clone https://github.com/opscart/kosva.git
cd kosva
go build -o kosva cmd/kosva/main.go

# Export Kubecost recommendations
kubecost recommendations export > recommendations.json

# Validate
./kosva check --input recommendations.json
```

## Features

### Current (v0.1.0-alpha)

- ✅ Validates against CIS Kubernetes Benchmark checks
- ✅ Spot instance safety validation (CIS 5.7)
- ✅ Resource limit safety checks (CIS 5.10)
- ✅ Risk scoring (0-10 scale)
- ✅ Safe alternative suggestions
- ✅ Integration with Kubecost JSON exports

### Coming Soon

- [ ] YAML-based policy engine
- [ ] PCI-DSS compliance policies
- [ ] HIPAA compliance policies
- [ ] Multi-tenancy security checks
- [ ] Kubernetes Operator mode
- [ ] CI/CD pipeline integration
- [ ] JSON/HTML report formats

## Use Cases

### DevOps Engineers
**Problem:** "Security keeps blocking our cost optimizations"  
**Solution:** Validate recommendations before proposing to security team

### Security Teams
**Problem:** "DevOps implements changes that create vulnerabilities"  
**Solution:** Automated pre-implementation security gates

### FinOps Teams
**Problem:** "Don't know which cost savings are safe to implement"  
**Solution:** Risk-scored recommendations with compliance validation

## Project Status

🚧 **Alpha Release** - Core functionality working, under active development.

Built by DevOps engineers managing 500+ cores in production who've been burned by unsafe optimizations.

## Contributing

Contributions welcome! This project started because:
1. Cost tools don't validate security
2. Security tools don't consider cost
3. Nobody bridges the gap

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Roadmap

- [x] **Week 1-2:** Core validation engine (✅ DONE)
- [ ] **Week 3-4:** Policy engine with compliance templates
- [ ] **Week 5-6:** Kubernetes Operator mode
- [ ] **Week 7-8:** CI/CD integration and documentation

## License

MIT License - See [LICENSE](LICENSE) for details.

## Author

Built by [@opscart](https://github.com/opscart) - Senior DevOps Engineer with 10+ years managing cloud-native infrastructure at scale.

**Blog:** [opscart.com](https://opscart.com)

---

**⭐ If this tool saves you from implementing an unsafe optimization, please star the repo!**