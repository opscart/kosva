package checks

// CheckResult represents the outcome of a security validation
type CheckResult struct {
	Passed      bool
	Severity    string  // "LOW", "MEDIUM", "HIGH", "CRITICAL"
	RiskScore   float64 // 0.0 to 10.0
	CheckName   string
	Message     string
	Remediation string // Suggested fix
}

// ValidationResult contains all checks for a single recommendation
type ValidationResult struct {
	WorkloadName       string
	RecommendationType string
	Savings            float64
	Approved           bool
	Checks             []CheckResult
	Alternative        string // If blocked, what's a safe alternative?
}
