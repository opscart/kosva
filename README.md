# kosva - Kubernetes Optimization Security Validator

**Stop implementing unsafe cost optimizations.**

`kosva` validates Kubernetes cost optimization recommendations against security policies and compliance frameworks before you implement them.

## Quick Start
```bash
# Install
go install github.com/opscart/kosva/cmd/kosva@latest

# Export Kubecost recommendations
kubecost recommendations export > recs.json

# Validate
kosva check --input recs.json
```

## Problem

Cost tools like Kubecost suggest optimizations to reduce spend. But many optimizations create security vulnerabilities:

- Spot instances → Service interruptions during attacks
- Under-provisioned resources → DoS vulnerabilities  
- Multi-tenancy → Compliance violations

Security teams block these recommendations, wasting engineering time.

## Solution

`kosva` validates recommendations BEFORE implementation:
```bash
$ kosva check --input kubecost-recs.json --policy cis-benchmark

✅ SAFE (3 recommendations):
1. Right-size logging-service → $400/month (Risk: 0.2/10)
2. Enable cluster autoscaler → $600/month (Risk: 0.3/10)

❌ BLOCKED (2 recommendations):
1. Spot instances for payment-api → $800/month
   ⚠️  CIS 5.7: Sensitive data workload
   ⚠️  High availability requirement conflict
   
   Suggestion: Use Reserved Instances → $240/month savings (compliant)
```

## Features

- Validates against CIS Kubernetes Benchmark
- Policy-based compliance checking (PCI-DSS, HIPAA, SOC2)
- Risk scoring for optimizations
- Safe alternative suggestions
- Integration with Kubecost

## Status

**Alpha** - Under active development

## Roadmap

- [x] Project setup
- [ ] Kubecost JSON parser
- [ ] Basic CIS checks
- [ ] Policy engine
- [ ] CLI interface
- [ ] Documentation

## Contributing

Built by DevOps engineers who've been burned by unsafe optimizations.

## License

MIT
