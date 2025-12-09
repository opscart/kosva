# kosva Policy Schema

Policies are YAML files that define security validation rules.

## Structure
```yaml
name: "Policy Name"
description: "What this policy checks"
cis_mapping: "CIS 5.7"  # Optional: CIS Kubernetes Benchmark reference
severity: "HIGH"         # LOW, MEDIUM, HIGH, CRITICAL
enabled: true

rules:
  - id: "rule-001"
    name: "Rule name"
    check_type: "spot-instance"  # or "right-size", "multi-tenancy"
    conditions:
      - field: "workload"
        operator: "contains"
        values: ["payment", "auth"]
      - field: "namespace"
        operator: "equals"
        values: ["production"]
    action: "block"
    risk_score: 8.5
    message: "Why this is blocked"
    remediation: "What to do instead"
```

## Field Reference

### check_type
- `spot-instance`: Validates Spot instance recommendations
- `right-size`: Validates resource limit changes
- `multi-tenancy`: Validates workload isolation
- `storage`: Validates storage class security

### operator
- `contains`: Field contains any of the values
- `equals`: Field equals any of the values
- `starts_with`: Field starts with any of the values
- `regex`: Field matches regex pattern

### action
- `block`: Fail validation
- `warn`: Show warning but pass
- `info`: Informational only