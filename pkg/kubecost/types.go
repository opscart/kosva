package kubecost

// Recommendation represents a Kubecost cost optimization suggestion
type Recommendation struct {
	Type        string   `json:"type"`     // "spot-instance", "right-size", etc.
	Workload    string   `json:"workload"` // "deployment/payment-api"
	Namespace   string   `json:"namespace"`
	Savings     float64  `json:"savings"` // Monthly savings in dollars
	Current     Resource `json:"current"`
	Recommended Resource `json:"recommended"`
}

type Resource struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// RecommendationList is the top-level structure from Kubecost export
type RecommendationList struct {
	Recommendations []Recommendation `json:"recommendations"`
	Cluster         string           `json:"cluster"`
	Timestamp       string           `json:"timestamp"`
}
