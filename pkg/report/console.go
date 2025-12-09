package report

import (
	"fmt"
	"strings"

	"github.com/opscart/kosva/pkg/checks"
)

// PrintConsoleReport outputs validation results to console
func PrintConsoleReport(results []checks.ValidationResult) {
	approved := []checks.ValidationResult{}
	blocked := []checks.ValidationResult{}

	// Separate approved and blocked
	for _, r := range results {
		if r.Approved {
			approved = append(approved, r)
		} else {
			blocked = append(blocked, r)
		}
	}

	// Calculate totals
	totalSavings := 0.0
	safeSavings := 0.0
	blockedSavings := 0.0

	for _, r := range results {
		totalSavings += r.Savings
		if r.Approved {
			safeSavings += r.Savings
		} else {
			blockedSavings += r.Savings
		}
	}

	// Print summary
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("  KOSVA VALIDATION REPORT")
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("\nTotal Recommendations: %d\n", len(results))
	fmt.Printf("Potential Savings: $%.0f/month\n", totalSavings)
	fmt.Println()

	// Print approved recommendations
	if len(approved) > 0 {
		fmt.Println("SAFE RECOMMENDATIONS (" + fmt.Sprintf("%d", len(approved)) + "):")
		fmt.Println(strings.Repeat("-", 70))
		for i, r := range approved {
			fmt.Printf("%d. %s -> $%.0f/month\n", i+1, r.WorkloadName, r.Savings)
			fmt.Printf("   Type: %s\n", r.RecommendationType)
			for _, check := range r.Checks {
				fmt.Printf("   [PASS] %s (Risk: %.1f/10)\n", check.CheckName, check.RiskScore)
				fmt.Printf("          %s\n", check.Message)
			}
			fmt.Println()
		}
		fmt.Printf("Safe Savings: $%.0f/month\n\n", safeSavings)
	}

	// Print blocked recommendations
	if len(blocked) > 0 {
		fmt.Println("BLOCKED RECOMMENDATIONS (" + fmt.Sprintf("%d", len(blocked)) + "):")
		fmt.Println(strings.Repeat("-", 70))
		for i, r := range blocked {
			fmt.Printf("%d. %s -> Would save $%.0f/month\n", i+1, r.WorkloadName, r.Savings)
			fmt.Printf("   Type: %s\n", r.RecommendationType)
			for _, check := range r.Checks {
				fmt.Printf("   [FAIL] %s [%s] (Risk: %.1f/10)\n", check.CheckName, check.Severity, check.RiskScore)
				fmt.Printf("          %s\n", check.Message)
				fmt.Printf("          Remediation: %s\n", check.Remediation)
			}
			if r.Alternative != "" {
				fmt.Printf("\n   SAFE ALTERNATIVE:\n")
				fmt.Printf("   %s\n", r.Alternative)
			}
			fmt.Println()
		}
		fmt.Printf("Blocked Savings: $%.0f/month\n\n", blockedSavings)
	}

	// Final summary
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("SUMMARY: $%.0f/$%.0f safe (%.0f%% of potential savings)\n",
		safeSavings, totalSavings, (safeSavings/totalSavings)*100)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println()
}
