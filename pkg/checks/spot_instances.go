package checks

import (
	"strings"

	"github.com/opscart/kosva/pkg/kubecost"
)

// CheckSpotInstanceSafety validates if a workload is suitable for Spot instances
func CheckSpotInstanceSafety(rec kubecost.Recommendation) CheckResult {
	result := CheckResult{
		CheckName: "CIS 5.7 - Spot Instance Suitability",
	}

	// Check 1: Is this a critical workload based on naming?
	criticalKeywords := []string{
		"payment", "auth", "database", "api-gateway",
		"billing", "order", "checkout", "production",
	}

	workloadLower := strings.ToLower(rec.Workload)

	for _, keyword := range criticalKeywords {
		if strings.Contains(workloadLower, keyword) {
			result.Passed = false
			result.Severity = "HIGH"
			result.RiskScore = 8.5
			result.Message = "Workload appears to be critical based on name. Spot instances can be interrupted with 2-minute notice."
			result.Remediation = "Use Reserved Instances for 30% savings with guaranteed availability, or use Spot only for stateless, fault-tolerant workloads."
			return result
		}
	}

	// Check 2: Production namespace?
	if strings.Contains(strings.ToLower(rec.Namespace), "prod") {
		result.Passed = false
		result.Severity = "MEDIUM"
		result.RiskScore = 6.5
		result.Message = "Production namespace detected. Spot instances may impact SLA compliance."
		result.Remediation = "Consider using Spot instances only for dev/staging or implement proper handling of interruptions."
		return result
	}

	// If passed all checks
	result.Passed = true
	result.Severity = "LOW"
	result.RiskScore = 2.0
	result.Message = "Workload appears suitable for Spot instances (non-critical, non-production)."
	result.Remediation = "Ensure proper handling of interruption notices and test failover mechanisms."

	return result
}
