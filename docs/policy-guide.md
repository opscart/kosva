# Policy Creation Guide

Learn how to create custom security policies for kosva.

## Policy Structure

Policies are YAML files that define validation rules. Each policy contains:

- **Metadata:** Name, description, CIS mapping, severity
- **Rules:** One or more validation rules
- **Conditions:** When the rule applies
- **Actions:** What to do when conditions match

## Basic Policy Anatomy
```yaml
name: "Policy Name"
description: "What this policy validates"
cis_mapping: "CIS X.Y"
severity: "HIGH"
enabled: true

rules:
  - id: "unique-id-001"
    name: "Descriptive rule name"
    check_type: "spot-instance"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["keyword1", "keyword2"]
    action: "block"
    risk_score: 8.5
    message: "Why this is blocked"
    remediation: "What to do instead"
```

---

## Field Reference

### Metadata Fields

| Field | Required | Values | Description |
|-------|----------|--------|-------------|
| `name` | Yes | String | Policy display name |
| `description` | No | String | What the policy checks |
| `cis_mapping` | No | String | CIS Kubernetes Benchmark reference |
| `severity` | Yes | `LOW`, `MEDIUM`, `HIGH`, `CRITICAL` | Policy severity |
| `enabled` | Yes | `true`, `false` | Whether policy is active |

### Rule Fields

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier (e.g., `company-001`) |
| `name` | Yes | Human-readable rule name |
| `check_type` | Yes | Type of recommendation to validate |
| `conditions` | Yes | List of conditions that must match |
| `action` | Yes | What happens when rule matches |
| `risk_score` | Yes | Risk level (0.0 to 10.0) |
| `message` | Yes | Why the rule triggered |
| `remediation` | Yes | How to fix or work around |

---

## Check Types

kosva supports these recommendation types:

| Check Type | Description | Use For |
|------------|-------------|---------|
| `spot-instance` | Spot instance recommendations | Validating workload suitability for Spot |
| `right-size` | Resource limit changes | Memory/CPU reduction safety |
| `multi-tenancy` | Shared infrastructure | Workload isolation requirements |
| `storage-security` | Storage configurations | Encryption, access modes |

---

## Condition Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `contains` | Field contains any value | Namespace contains "prod" |
| `equals` | Field equals any value exactly | Type equals "deployment" |
| `starts_with` | Field starts with any value | Workload starts with "critical-" |
| `regex` | Field matches regex pattern | Workload matches `^api-.*-prod$` |

---

## Condition Fields

| Field | Available For | Description |
|-------|---------------|-------------|
| `workload` | All types | Workload name (e.g., `deployment/payment-api`) |
| `namespace` | All types | Kubernetes namespace |
| `type` | All types | Recommendation type |
| `node_pool` | `multi-tenancy` | Target node pool name |
| `storage_class` | `storage-security` | Storage class name |
| `access_mode` | `storage-security` | PVC access mode |

---

## Actions

| Action | Effect | Use When |
|--------|--------|----------|
| `block` | Fail validation | Recommendation violates security policy |
| `warn` | Pass with warning | Recommendation carries risk but isn't prohibited |
| `info` | Pass with note | Informational only, no risk |

---

## Examples

### Example 1: Block Spot for Your Critical Apps
```yaml
name: "Company Critical Services Policy"
description: "Blocks Spot instances for revenue-generating services"
cis_mapping: "CIS 5.7"
severity: "CRITICAL"
enabled: true

rules:
  - id: "company-spot-001"
    name: "Block Spot for checkout and payment services"
    check_type: "spot-instance"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["checkout", "payment", "billing", "invoice"]
    action: "block"
    risk_score: 9.5
    message: "Revenue-generating services require guaranteed availability - Spot instances can be interrupted with 2-minute notice"
    remediation: "Use Reserved Instances (30% savings) or Savings Plans (up to 72% savings) with guaranteed capacity"
```

**Save as:** `policies/custom/critical-services.yaml`

---

### Example 2: Enforce Memory Reduction Limits
```yaml
name: "Conservative Resource Reduction Policy"
description: "Stricter limits on memory reductions than CIS default"
cis_mapping: "CIS 5.10"
severity: "HIGH"
enabled: true

rules:
  - id: "company-mem-001"
    name: "Block any reduction over 25%"
    check_type: "right-size"
    conditions:
      - field: "namespace"
        operator: "contains"
        values: ["prod", "production", "staging"]
    action: "block"
    risk_score: 7.0
    message: "Memory reduction exceeds company policy (25% max) - risk of OOM kills"
    remediation: "Reduce by maximum 25%, monitor for 2 weeks, then re-evaluate"

  - id: "company-mem-002"
    name: "Warn on any production reduction"
    check_type: "right-size"
    conditions:
      - field: "namespace"
        operator: "equals"
        values: ["production"]
    action: "warn"
    risk_score: 5.0
    message: "Any memory reduction in production requires approval"
    remediation: "Submit change request with performance test results"
```

---

### Example 3: PCI-DSS Compliance Policy
```yaml
name: "PCI-DSS Workload Isolation"
description: "Enforces PCI-DSS requirement for payment workload isolation"
cis_mapping: "CIS 5.7.3"
severity: "CRITICAL"
enabled: true

rules:
  - id: "pci-001"
    name: "Block PCI workloads on shared infrastructure"
    check_type: "multi-tenancy"
    conditions:
      - field: "namespace"
        operator: "contains"
        values: ["pci", "payment", "cardholder"]
    action: "block"
    risk_score: 10.0
    message: "PCI-DSS requires dedicated infrastructure for cardholder data environment - sharing nodes violates segmentation requirements"
    remediation: "Deploy to dedicated node pool with taints: 'pci-workload=true:NoSchedule'"

  - id: "pci-002"
    name: "Require encrypted storage for PCI data"
    check_type: "storage-security"
    conditions:
      - field: "namespace"
        operator: "contains"
        values: ["pci", "payment"]
    action: "block"
    risk_score: 10.0
    message: "PCI-DSS 3.4 requires encryption for cardholder data at rest"
    remediation: "Use storage class with encryption: 'encrypted-premium' or 'azure-disk-encrypted'"
```

---

### Example 4: Regex Pattern Matching
```yaml
name: "API Gateway Policy"
description: "Special rules for API gateway workloads"
cis_mapping: "CIS 5.7"
severity: "HIGH"
enabled: true

rules:
  - id: "api-001"
    name: "Block Spot for public-facing APIs"
    check_type: "spot-instance"
    conditions:
      - field: "workload"
        operator: "regex"
        values: ["^api-gateway-.*", ".*-public-api$"]
    action: "block"
    risk_score: 8.0
    message: "Public-facing API gateways require high availability - Spot interruptions impact customer experience"
    remediation: "Use on-demand instances or Reserved capacity for API gateways"
```

---

### Example 5: Multi-Condition Rules
```yaml
name: "Database Protection Policy"
description: "Comprehensive database workload protection"
cis_mapping: "CIS 5.7.3, CIS 5.4.1"
severity: "CRITICAL"
enabled: true

rules:
  - id: "db-001"
    name: "Block database multi-tenancy in production"
    check_type: "multi-tenancy"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["postgres", "mysql", "mongodb", "redis", "elasticsearch"]
      - field: "namespace"
        operator: "contains"
        values: ["prod", "production"]
    action: "block"
    risk_score: 9.0
    message: "Production databases require dedicated nodes for performance isolation and security"
    remediation: "Deploy databases to dedicated node pool with label: 'workload=database'"
```

**Note:** ALL conditions must match for the rule to trigger (AND logic).

---

## Best Practices

### 1. Start with Warnings

When creating new policies, start with `action: "warn"` to see what would be blocked:
```yaml
action: "warn"  # Start here
risk_score: 5.0
```

After validating the policy catches the right cases, change to `block`:
```yaml
action: "block"  # After validation
risk_score: 8.0
```

### 2. Use Descriptive IDs
```yaml
# Good
id: "company-spot-payment-001"
id: "pci-isolation-001"
id: "prod-memory-limit-001"

# Bad
id: "rule1"
id: "check-001"
id: "policy-a"
```

### 3. Provide Actionable Remediation
```yaml
# Good
remediation: "Use Reserved Instances (30% savings) or deploy to dedicated node pool with taint 'critical=true:NoSchedule'"

# Bad
remediation: "Don't use Spot instances"
remediation: "Fix this"
```

### 4. Match Risk Scores to Severity

| Severity | Risk Score Range |
|----------|------------------|
| CRITICAL | 9.0 - 10.0 |
| HIGH | 7.0 - 8.9 |
| MEDIUM | 4.0 - 6.9 |
| LOW | 0.1 - 3.9 |

### 5. Organize Policies by Purpose
```
policies/
├── cis-benchmark/          # Official CIS checks
│   ├── spot-instances.yaml
│   └── resource-limits.yaml
├── compliance/             # Regulatory requirements
│   ├── pci-dss.yaml
│   └── hipaa.yaml
└── custom/                 # Company-specific
    ├── critical-services.yaml
    └── database-policy.yaml
```

---

## Testing Your Policy

### 1. Create Test Data

Create `test-recommendations.json`:
```json
{
  "cluster": "test",
  "timestamp": "2024-12-10T00:00:00Z",
  "recommendations": [
    {
      "type": "spot-instance",
      "workload": "deployment/payment-api",
      "namespace": "production",
      "savings": 800.00,
      "current": {"cpu": "1000m", "memory": "2Gi"},
      "recommended": {"cpu": "1000m", "memory": "2Gi"}
    }
  ]
}
```

### 2. Test the Policy
```bash
# Place your policy in policies/custom/
mv my-policy.yaml policies/custom/

# Run validation
kosva check --input test-recommendations.json

# Should show your rule triggering
```

### 3. Verify Rule Logic
```bash
# List loaded policies
kosva policies

# Your policy should appear in the list
```

---

## Common Patterns

### Pattern 1: Environment-Based Rules
```yaml
# Block aggressive changes in prod, allow in dev
rules:
  - id: "env-001"
    name: "Strict production limits"
    check_type: "right-size"
    conditions:
      - field: "namespace"
        operator: "equals"
        values: ["production"]
    action: "block"
    risk_score: 8.0
    message: "Production requires conservative reductions"
    remediation: "Reduce by 25% max and monitor"

  - id: "env-002"
    name: "Lenient dev limits"
    check_type: "right-size"
    conditions:
      - field: "namespace"
        operator: "contains"
        values: ["dev", "test"]
    action: "warn"
    risk_score: 3.0
    message: "Development environment - acceptable risk"
    remediation: "Monitor for issues, adjust if needed"
```

### Pattern 2: Workload Classification
```yaml
# Different rules for different workload types
rules:
  - id: "workload-001"
    name: "Stateful workloads need dedicated nodes"
    check_type: "multi-tenancy"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["statefulset", "database", "redis", "kafka"]
    action: "block"
    risk_score: 9.0

  - id: "workload-002"
    name: "Stateless workloads can share"
    check_type: "multi-tenancy"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["deployment", "job", "cronjob"]
    action: "info"
    risk_score: 2.0
```

### Pattern 3: Graduated Risk Levels
```yaml
# Progressive warnings based on risk level
rules:
  - id: "risk-001"
    name: "Critical - requires security approval"
    check_type: "spot-instance"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["payment", "auth"]
    action: "block"
    risk_score: 9.0

  - id: "risk-002"
    name: "High - requires manager approval"
    check_type: "spot-instance"
    conditions:
      - field: "namespace"
        operator: "equals"
        values: ["production"]
    action: "block"
    risk_score: 7.0

  - id: "risk-003"
    name: "Medium - team lead approval"
    check_type: "spot-instance"
    conditions:
      - field: "namespace"
        operator: "equals"
        values: ["staging"]
    action: "warn"
    risk_score: 5.0
```

---

## Troubleshooting

### Policy Not Loading
```bash
# Check syntax
cat policies/custom/my-policy.yaml | yaml-lint

# Verify file is valid YAML
python3 -c "import yaml; yaml.safe_load(open('policies/custom/my-policy.yaml'))"

# Check kosva can see it
kosva policies --policy-dir ./policies
```

### Rule Not Triggering
```bash
# Check field values
echo '{"workload": "deployment/test"}' | jq .

# Verify condition logic
# Remember: ALL conditions must match (AND logic)
```

### Multiple Rules Triggering

This is expected! Multiple rules can match the same recommendation. kosva reports all matching rules.

---

## Next Steps

- [Kubecost Integration](kubecost-setup.md)
- [Contributing Custom Policies](../CONTRIBUTING.md)
- [Example Policies](../policies/)