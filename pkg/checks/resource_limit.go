package checks

import (
	"strconv"
	"strings"

	"github.com/opscart/kosva/pkg/kubecost"
)

// CheckResourceLimits validates if resource reduction is safe
func CheckResourceLimits(rec kubecost.Recommendation) CheckResult {
	result := CheckResult{
		CheckName: "CIS 5.10 - Resource Limit Safety",
	}

	// Parse current and recommended memory
	currentMem := parseMemory(rec.Current.Memory)
	recommendedMem := parseMemory(rec.Recommended.Memory)

	// Calculate reduction percentage
	reduction := (float64(currentMem-recommendedMem) / float64(currentMem)) * 100

	// If reducing by more than 50%, flag as risky
	if reduction > 50 {
		result.Passed = false
		result.Severity = "HIGH"
		result.RiskScore = 7.5
		result.Message = "Memory reduction exceeds 50%. High risk of OOM kills and service disruption."
		result.Remediation = "Consider a more gradual reduction (25-30% max) and monitor for OOM events."
		return result
	}

	// If reducing by 30-50%, warn
	if reduction > 30 {
		result.Passed = false
		result.Severity = "MEDIUM"
		result.RiskScore = 5.0
		result.Message = "Memory reduction of 30-50% carries moderate risk."
		result.Remediation = "Implement gradual rollout and monitor memory usage closely."
		return result
	}

	// Reasonable reduction
	result.Passed = true
	result.Severity = "LOW"
	result.RiskScore = 1.5
	result.Message = "Memory reduction appears safe (< 30%)."
	result.Remediation = "Monitor memory usage for 1-2 weeks after implementation."

	return result
}

// parseMemory converts Kubernetes memory strings to MB
func parseMemory(memStr string) int64 {
	memStr = strings.TrimSpace(memStr)

	if strings.HasSuffix(memStr, "Gi") {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Gi"), 64)
		return int64(val * 1024) // Convert to Mi
	}
	if strings.HasSuffix(memStr, "Mi") {
		val, _ := strconv.ParseInt(strings.TrimSuffix(memStr, "Mi"), 10, 64)
		return val
	}
	if strings.HasSuffix(memStr, "M") {
		val, _ := strconv.ParseInt(strings.TrimSuffix(memStr, "M"), 10, 64)
		return val
	}

	return 0
}
