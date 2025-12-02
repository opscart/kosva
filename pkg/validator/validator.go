package validator

import (
	"fmt"

	"github.com/opscart/kosva/pkg/checks"
	"github.com/opscart/kosva/pkg/kubecost"
)

// ValidateRecommendation runs all security checks on a single recommendation
func ValidateRecommendation(rec kubecost.Recommendation) checks.ValidationResult {
	result := checks.ValidationResult{
		WorkloadName:       rec.Workload,
		RecommendationType: rec.Type,
		Savings:            rec.Savings,
		Checks:             []checks.CheckResult{},
	}

	// Run appropriate checks based on recommendation type
	switch rec.Type {
	case "spot-instance":
		spotCheck := checks.CheckSpotInstanceSafety(rec)
		result.Checks = append(result.Checks, spotCheck)

		// Determine approval
		result.Approved = spotCheck.Passed

		// Suggest alternative if blocked
		if !spotCheck.Passed {
			result.Alternative = fmt.Sprintf("Use Reserved Instances → ~$%.0f/month savings (30%% of Spot savings, but compliant)", rec.Savings*0.3)
		}

	case "right-size":
		limitsCheck := checks.CheckResourceLimits(rec)
		result.Checks = append(result.Checks, limitsCheck)

		result.Approved = limitsCheck.Passed

		if !limitsCheck.Passed {
			result.Alternative = "Reduce resources by 25% instead of recommended amount"
		}

	default:
		// Unknown type - allow but warn
		result.Approved = true
		result.Checks = append(result.Checks, checks.CheckResult{
			Passed:    true,
			Severity:  "INFO",
			RiskScore: 0,
			CheckName: "Unknown recommendation type",
			Message:   fmt.Sprintf("No specific security checks for type: %s", rec.Type),
		})
	}

	return result
}

// ValidateAll runs validation on all recommendations
func ValidateAll(recList *kubecost.RecommendationList) []checks.ValidationResult {
	results := []checks.ValidationResult{}

	for _, rec := range recList.Recommendations {
		result := ValidateRecommendation(rec)
		results = append(results, result)
	}

	return results
}
