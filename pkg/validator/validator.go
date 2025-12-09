package validator

import (
	"fmt"

	"github.com/opscart/kosva/pkg/checks"
	"github.com/opscart/kosva/pkg/kubecost"
	"github.com/opscart/kosva/pkg/policyengine"
)

// Validator performs security validation
type Validator struct {
	PolicyEngine *policyengine.Engine
}

// NewValidator creates a new validator
func NewValidator(policyDir string) (*Validator, error) {
	// Load policies
	policies, err := policyengine.LoadPolicies(policyDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	if len(policies) == 0 {
		return nil, fmt.Errorf("no policies loaded from %s", policyDir)
	}

	fmt.Printf("Loaded %d policies\n", len(policies))

	return &Validator{
		PolicyEngine: policyengine.NewEngine(policies),
	}, nil
}

// ValidateRecommendation runs all policies on a recommendation
func (v *Validator) ValidateRecommendation(rec kubecost.Recommendation) checks.ValidationResult {
	result := checks.ValidationResult{
		WorkloadName:       rec.Workload,
		RecommendationType: rec.Type,
		Savings:            rec.Savings,
		Checks:             []checks.CheckResult{},
	}

	// Run policy engine
	policyResults := v.PolicyEngine.Evaluate(rec)

	// Convert policy results to check results
	for _, pr := range policyResults {
		checkResult := checks.CheckResult{
			Passed:      pr.Passed,
			Severity:    pr.Severity,
			RiskScore:   pr.RiskScore,
			CheckName:   fmt.Sprintf("%s - %s", pr.PolicyName, pr.RuleName),
			Message:     pr.Message,
			Remediation: pr.Remediation,
		}
		result.Checks = append(result.Checks, checkResult)
	}

	// Determine overall approval (fail if any check blocks)
	result.Approved = true
	for _, check := range result.Checks {
		if !check.Passed {
			result.Approved = false
			break
		}
	}

	// Generate alternative if blocked
	if !result.Approved {
		result.Alternative = generateAlternative(rec)
	}

	return result
}

// ValidateAll runs validation on all recommendations
func (v *Validator) ValidateAll(recList *kubecost.RecommendationList) []checks.ValidationResult {
	results := []checks.ValidationResult{}

	for _, rec := range recList.Recommendations {
		result := v.ValidateRecommendation(rec)
		results = append(results, result)
	}

	return results
}

func generateAlternative(rec kubecost.Recommendation) string {
	switch rec.Type {
	case "spot-instance":
		return fmt.Sprintf("Use Reserved Instances -> ~$%.0f/month savings (30%% of Spot savings, compliant)", rec.Savings*0.3)
	case "right-size":
		return "Implement gradual reduction (25% max) with monitoring"
	default:
		return "Review with security team for compliant alternative"
	}
}
