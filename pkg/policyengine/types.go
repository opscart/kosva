package policyengine

// Policy represents a security policy
type Policy struct {
    Name        string `yaml:"name"`
    Description string `yaml:"description"`
    CISMapping  string `yaml:"cis_mapping"`
    Severity    string `yaml:"severity"`
    Enabled     bool   `yaml:"enabled"`
    Rules       []Rule `yaml:"rules"`
}

// Rule represents a validation rule
type Rule struct {
    ID          string      `yaml:"id"`
    Name        string      `yaml:"name"`
    CheckType   string      `yaml:"check_type"`
    Conditions  []Condition `yaml:"conditions"`
    Action      string      `yaml:"action"`
    RiskScore   float64     `yaml:"risk_score"`
    Message     string      `yaml:"message"`
    Remediation string      `yaml:"remediation"`
}

// Condition represents a rule condition
type Condition struct {
    Field    string   `yaml:"field"`
    Operator string   `yaml:"operator"`
    Values   []string `yaml:"values"`
}

// PolicyResult represents the outcome of a policy check
type PolicyResult struct {
    PolicyName  string
    RuleName    string
    Passed      bool
    Action      string
    Severity    string
    RiskScore   float64
    Message     string
    Remediation string
}
