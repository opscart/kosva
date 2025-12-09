package policyengine

import (
    "regexp"
    "strings"
    
    "github.com/opscart/kosva/pkg/kubecost"
)

// Engine executes policies against recommendations
type Engine struct {
    Policies []Policy
}

// NewEngine creates a new policy engine
func NewEngine(policies []Policy) *Engine {
    return &Engine{
        Policies: policies,
    }
}

// Evaluate runs all policies against a recommendation
func (e *Engine) Evaluate(rec kubecost.Recommendation) []PolicyResult {
    var results []PolicyResult
    
    for _, policy := range e.Policies {
        for _, rule := range policy.Rules {
            // Only evaluate rules matching this recommendation type
            if rule.CheckType != rec.Type {
                continue
            }
            
            // Evaluate all conditions
            allConditionsMet := true
            for _, condition := range rule.Conditions {
                if !e.evaluateCondition(condition, rec) {
                    allConditionsMet = false
                    break
                }
            }
            
            // If all conditions met, create result
            if allConditionsMet {
                result := PolicyResult{
                    PolicyName:  policy.Name,
                    RuleName:    rule.Name,
                    Passed:      rule.Action != "block",
                    Action:      rule.Action,
                    Severity:    policy.Severity,
                    RiskScore:   rule.RiskScore,
                    Message:     rule.Message,
                    Remediation: rule.Remediation,
                }
                results = append(results, result)
            }
        }
    }
    
    return results
}

// evaluateCondition checks if a condition matches
func (e *Engine) evaluateCondition(cond Condition, rec kubecost.Recommendation) bool {
    var fieldValue string
    
    // Extract field value from recommendation
    switch cond.Field {
    case "workload":
        fieldValue = strings.ToLower(rec.Workload)
    case "namespace":
        fieldValue = strings.ToLower(rec.Namespace)
    case "type":
        fieldValue = rec.Type
    default:
        return false
    }
    
    // Check condition based on operator
    switch cond.Operator {
    case "contains":
        for _, val := range cond.Values {
            if strings.Contains(fieldValue, strings.ToLower(val)) {
                return true
            }
        }
        return false
        
    case "equals":
        for _, val := range cond.Values {
            if fieldValue == strings.ToLower(val) {
                return true
            }
        }
        return false
        
    case "starts_with":
        for _, val := range cond.Values {
            if strings.HasPrefix(fieldValue, strings.ToLower(val)) {
                return true
            }
        }
        return false
        
    case "regex":
        for _, pattern := range cond.Values {
            matched, _ := regexp.MatchString(pattern, fieldValue)
            if matched {
                return true
            }
        }
        return false
        
    default:
        return false
    }
}
