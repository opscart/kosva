package kubecost

// Recommendation represents a Kubecost cost optimization suggestion
type Recommendation struct {
	Type        string   `json:"type"`     // "spot-instance", "right-size", "multi-tenancy", "storage-security"
	Workload    string   `json:"workload"` // "deployment/payment-api"
	Namespace   string   `json:"namespace"`
	Savings     float64  `json:"savings"` // Monthly savings in dollars
	Current     Resource `json:"current"`
	Recommended Resource `json:"recommended"`

	// Additional fields for new check types
	NodePool     string `json:"node_pool,omitempty"`     // For multi-tenancy checks
	StorageClass string `json:"storage_class,omitempty"` // For storage checks
	AccessMode   string `json:"access_mode,omitempty"`   // For storage checks
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
