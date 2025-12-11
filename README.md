# kosva - Kubernetes Optimization Security Validator

**Validate Kubecost recommendations against CIS benchmarks before implementation.**

kosva prevents unsafe cost optimizations by checking them against security policies and compliance frameworks. Stop wasting time on recommendations that security will reject.

## The Problem

Cost optimization tools like Kubecost suggest ways to reduce cloud spend, but many optimizations create security vulnerabilities:

- **Spot instances for critical services** → 2-minute interruption notices violate SLA requirements
- **Aggressive resource reductions** → OOM kills create denial-of-service vulnerabilities  
- **Shared infrastructure for regulated workloads** → PCI-DSS/HIPAA compliance violations
- **Unencrypted storage for sensitive data** → Data-at-rest encryption requirements not met

**Result:** Security teams block recommendations, engineering time is wasted, and cost savings opportunities are lost.

## The Solution

kosva validates cost recommendations BEFORE proposing them to security teams:
```bash
$ kosva check --input kubecost-recommendations.json

Loading: kubecost-recommendations.json
Loaded 7 recommendations
Loading policies from: ./policies
Loaded 4 policies
Running validation...

======================================================================
  KOSVA VALIDATION REPORT
======================================================================

Total Recommendations: 7
Potential Savings: $6600/month

SAFE RECOMMENDATIONS (1):
----------------------------------------------------------------------
1. deployment/batch-processor -> $1200/month
   Type: spot-instance
   [PASS] CIS 5.7 - Spot Instance Suitability (Risk: 2.0/10)

BLOCKED RECOMMENDATIONS (6):
----------------------------------------------------------------------
1. deployment/payment-processor -> Would save $2500/month
   Type: multi-tenancy
   [FAIL] CIS 5.7.3 - Multi-Tenancy Isolation [HIGH] (Risk: 9.0/10)
          Regulated workload cannot share infrastructure
          Remediation: Use dedicated node pools with taints/tolerations

[... more blocked recommendations ...]

======================================================================
SUMMARY: $1200/$6600 safe (18% of potential savings)
======================================================================
```

**Time saved:** 3+ weeks of security review  
**Clear answer:** Which optimizations are safe to implement

## Quick Start

### Installation
```bash
# Clone repository
git clone https://github.com/opscart/kosva.git
cd kosva

# Build
go build -o kosva cmd/kosva/main.go

# Verify
./kosva version
```

### Basic Usage
```bash
# Export recommendations from Kubecost
# (Via UI: Savings -> Export as JSON)
# (Via API: curl http://kubecost-url/model/savings > recs.json)

# Validate
./kosva check --input kubecost-recommendations.json

# List available policies
./kosva policies

# Use custom policy directory
./kosva check --input recs.json --policy-dir /path/to/policies
```

### Direct API Integration
```bash
# Port-forward Kubecost
kubectl port-forward -n kubecost svc/kubecost-cost-analyzer 9090:9090

# Validate directly from API
./kosva check --kubecost-url http://localhost:9090
```

## Features

### Current (v0.2.0-alpha)

- **YAML-based policy engine** - Customize without editing code
- **CIS Kubernetes Benchmark policies** - Pre-built security checks
- **Risk scoring** - 0-10 scale for each recommendation
- **Safe alternatives** - Suggestions for blocked optimizations
- **Kubecost integration** - File-based and API-based validation

### Security Checks (4 Policies, 11 Rules)

| Policy | CIS Mapping | Checks |
|--------|-------------|--------|
| **Spot Instance Suitability** | CIS 5.7 | Critical workload detection, production namespace validation, API service warnings |
| **Resource Limit Safety** | CIS 5.10 | Aggressive reduction blocking (>50%), moderate reduction warnings (30-50%) |
| **Multi-Tenancy Isolation** | CIS 5.7.3 | Regulated workload isolation, prod/non-prod separation, database isolation |
| **Storage Security** | CIS 5.4.1 | Encryption requirements, access mode validation, storage class compliance |

## Real-World Use Case

**Scenario:** FinOps exports Kubecost recommendations showing $15K/month in potential savings.

**Without kosva:**
1. FinOps proposes all 10 recommendations (3 weeks)
2. Security reviews manually (2 weeks)  
3. 60% rejected for security reasons
4. Back-and-forth on alternatives (1 week)
5. **Total:** 6 weeks, $6K/month approved

**With kosva:**
1. Export Kubecost recommendations (5 min)
2. Run: `kosva check --input recs.json` (30 sec)
3. Get instant feedback: $8K safe, $7K blocked with alternatives
4. Submit only safe recommendations (1 day)
5. Security approves (2 days)
6. **Total:** 3 days, $8K/month approved

**Result:** 93% time saved, 33% more savings achieved

### Comparison with Existing Tools

| Tool | What It Does | What It Doesn't | Why kosva? |
|------|--------------|-----------------|------------|
| **OPA/Gatekeeper** | General policy enforcement | Not cost-aware, requires Rego expertise | Pre-built cost+security policies, no Rego needed |
| **Kyverno** | Kubernetes policy engine | Not integrated with cost tools | Direct Kubecost integration with CIS mapping |
| **Polaris/kube-score** | K8s best practices | No cost optimization focus | Validates cost recommendations specifically |
| **kube-bench** | CIS cluster auditing | Audits configuration, not cost changes | Validates proposed cost optimizations |
| **Kubecost alone** | Cost recommendations | No security validation | Adds security validation layer |
| **Manual review** | Human judgment | Slow, inconsistent, doesn't scale | Automated, consistent, instant feedback |

### What About General Security Scanning?

**kosva is NOT a replacement for cluster security scanners.** It has a specific mission:

**kosva's Mission:** Validate cost optimization recommendations against security policies

**NOT kosva's Mission:** General Kubernetes cluster security auditing

### Tool Positioning

| Tool | Focus | When to Use |
|------|-------|-------------|
| **kube-bench** | Cluster CIS compliance audit | Run weekly to audit cluster security posture |
| **Polaris** | Workload best practices | Check running workloads for misconfigurations |
| **Kubecost** | Cost analysis | Find cost optimization opportunities |
| **kosva** | Cost recommendation validation | Before implementing Kubecost recommendations |
| **OPA/Kyverno** | Runtime policy enforcement | Enforce policies on deployment |

**Workflow:**
```
1. Kubecost finds $15K/month savings opportunities
2. kosva validates which are security-compliant ($8K safe)
3. You propose the safe ones to security
4. OPA/Kyverno enforces policies at runtime
5. kube-bench audits cluster monthly
```

### What kosva DOES Validate

✅ **Cost Optimization Recommendations:**
- Spot instance suitability (CIS 5.7)
- Resource limit safety (CIS 5.10)
- Multi-tenancy isolation (CIS 5.7.3)
- Storage security (CIS 5.4.1)

### What kosva Does NOT Validate

❌ **General Cluster Configuration:**
- Network policies (use kube-bench)
- Pod security standards (use Polaris)
- RBAC configuration (use kube-bench)
- Immutable root filesystems (use kube-score)

**Why?** Those checks audit cluster configuration, not cost optimization recommendations. Use specialized tools for comprehensive cluster scanning.

### Key Differentiator

kosva fills the gap between cost tools and security tools - it validates that cost optimizations won't violate security policies BEFORE you waste time proposing them."


## CI/CD Integration

### Exit Codes
```bash
0 = All recommendations safe
1 = Some blocked, alternatives available  
2 = Critical violations, no safe alternatives
3 = Error (invalid input, API failure)
```

### GitHub Actions Example
```yaml
name: Validate Cost Recommendations

on:
  pull_request:
    paths:
      - 'cost-optimization/**'

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install kosva
        run: |
          wget https://github.com/opscart/kosva/releases/latest/kosva
          chmod +x kosva
      
      - name: Validate recommendations
        run: |
          ./kosva check \
            --input cost-optimization/kubecost-recs.json \
            --policy-dir ./policies \
            --fail-on high
      
      - name: Comment PR
        if: failure()
        run: echo "Security validation failed - blocked recommendations found"
```

### Azure DevOps Example
```yaml
- task: Bash@3
  displayName: 'Validate Cost Recommendations'
  inputs:
    targetType: 'inline'
    script: |
      ./kosva check --input $(Build.SourcesDirectory)/recommendations.json
      if [ $? -eq 2 ]; then
        echo "##vso[task.logissue type=error]Critical security violations found"
        exit 1
      fi
```

## Customization

### Creating Custom Policies

Policies are YAML files in the `policies/` directory:
```yaml
name: "Custom Policy Name"
description: "What this policy checks"
cis_mapping: "CIS X.Y"
severity: "HIGH"
enabled: true

rules:
  - id: "custom-001"
    name: "Rule name"
    check_type: "spot-instance"  # or right-size, multi-tenancy, storage-security
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["my-app", "critical-service"]
    action: "block"
    risk_score: 8.5
    message: "Why this is blocked"
    remediation: "What to do instead"
```

**Supported operators:** `contains`, `equals`, `starts_with`, `regex`

**Example: Block Spot for your specific apps**
```yaml
name: "Company-Specific Spot Policy"
description: "Blocks Spot instances for our critical applications"
cis_mapping: "CIS 5.7"
severity: "HIGH"
enabled: true

rules:
  - id: "company-001"
    name: "Block Spot for revenue-generating services"
    check_type: "spot-instance"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["checkout", "order-processing", "payment-gateway"]
    action: "block"
    risk_score: 9.0
    message: "Revenue-generating services require guaranteed availability"
    remediation: "Use Reserved Instances or Savings Plans"
```

Save to `policies/custom/company-policy.yaml` and kosva will load it automatically.

## Production Status

**Current Status:** Alpha (v0.2.0)

**Safe for:**
- Pre-production validation
- CI/CD advisory checks  
- Manual cost review workflows

**Not yet for:**
- Automated enforcement (advisory only)
- Real-time blocking (no operator mode yet)

**Production users:** Testing internally at Fortune 500 pharmaceutical company (500+ AKS cores)

**Looking for:** Early adopters to validate real-world scenarios

## Roadmap

### v0.3 (Next 2 weeks)
- [ ] Compliance templates: PCI-DSS, HIPAA, SOC2
- [ ] JSON output format for automation
- [ ] Additional security checks (network policies, RBAC)

### v0.4 (Month 2)
- [ ] Kubernetes Operator mode
- [ ] Continuous validation
- [ ] Slack/Teams notifications

### v0.5 (Month 3)
- [ ] GitHub Action
- [ ] Azure DevOps extension
- [ ] GitLab CI integration

### Future
- [ ] OPA/Rego integration (use kosva with existing policies)
- [ ] Helm chart support
- [ ] Policy marketplace

## Contributing

Contributions welcome! This project exists because cost tools don't validate security, and security tools don't consider cost.

**What we need:**
- Additional CIS benchmark checks
- Compliance framework policies (PCI-DSS, HIPAA, SOC2)
- Integration with other cost tools
- Real-world validation scenarios

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Architecture
```
kosva/
├── cmd/kosva/              # CLI entry point
├── pkg/
│   ├── policyengine/       # YAML policy loader and evaluator
│   ├── kubecost/           # Kubecost API integration
│   ├── checks/             # Validation result types
│   ├── validator/          # Validation orchestration
│   └── report/             # Output formatting
├── policies/               # Security policy definitions
│   └── cis-benchmark/
│       ├── spot-instances.yaml
│       ├── resource-limits.yaml
│       ├── multi-tenancy.yaml
│       └── storage-security.yaml
└── examples/               # Sample data
```

## License

MIT License - See [LICENSE](LICENSE) for details.

## Author

Built by [@opscart](https://github.com/opscart) - Senior DevOps Engineer with 10+ years managing cloud-native infrastructure at scale.

**Website:** [opscart.com](https://opscart.com)  
**Problem:** Cost optimizations kept getting blocked by security  
**Solution:** Validate before proposing

---

**If kosva saves you from implementing an unsafe optimization, please star the repo!**