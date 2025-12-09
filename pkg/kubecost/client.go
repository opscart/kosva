package kubecost

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a Kubecost API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Kubecost client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// KubecostSavingsResponse represents the API response structure
type KubecostSavingsResponse struct {
	RecommendedRightSizing []struct {
		ClusterID       string `json:"clusterId"`
		Namespace       string `json:"namespace"`
		ControllerKind  string `json:"controllerKind"`
		ControllerName  string `json:"controllerName"`
		Container       string `json:"container"`
		Recommendations struct {
			CPU    RecommendationDetail `json:"cpu"`
			Memory RecommendationDetail `json:"memory"`
		} `json:"recommendations"`
		TotalMonthlySavings float64 `json:"totalMonthlySavings"`
	} `json:"recommendedRightSizing"`

	UnderutilizedNodes []struct {
		NodeName        string  `json:"node"`
		MonthlySavings  float64 `json:"monthlySavings"`
		ProviderID      string  `json:"providerId"`
		ReplacementType string  `json:"recommendedInstanceType"`
	} `json:"underutilizedNodes"`
}

type RecommendationDetail struct {
	Current        float64 `json:"current"`
	Recommended    float64 `json:"recommended"`
	MonthlySavings float64 `json:"monthlySavings"`
}

// GetRecommendations fetches recommendations from Kubecost API
func (c *Client) GetRecommendations() (*RecommendationList, error) {
	url := fmt.Sprintf("%s/model/savings", c.BaseURL)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("kubecost API returned %d: %s", resp.StatusCode, string(body))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var apiResp KubecostSavingsResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our internal format
	return c.convertToInternalFormat(apiResp), nil
}

// convertToInternalFormat converts Kubecost API format to our internal format
func (c *Client) convertToInternalFormat(apiResp KubecostSavingsResponse) *RecommendationList {
	recList := &RecommendationList{
		Cluster:         "from-kubecost-api",
		Timestamp:       time.Now().Format(time.RFC3339),
		Recommendations: []Recommendation{},
	}

	// Convert right-sizing recommendations
	for _, rs := range apiResp.RecommendedRightSizing {
		workload := fmt.Sprintf("%s/%s", rs.ControllerKind, rs.ControllerName)

		rec := Recommendation{
			Type:      "right-size",
			Workload:  workload,
			Namespace: rs.Namespace,
			Savings:   rs.TotalMonthlySavings,
			Current: Resource{
				CPU:    formatCPU(rs.Recommendations.CPU.Current),
				Memory: formatMemory(rs.Recommendations.Memory.Current),
			},
			Recommended: Resource{
				CPU:    formatCPU(rs.Recommendations.CPU.Recommended),
				Memory: formatMemory(rs.Recommendations.Memory.Recommended),
			},
		}

		recList.Recommendations = append(recList.Recommendations, rec)
	}

	// Convert node recommendations to spot-instance suggestions
	for _, node := range apiResp.UnderutilizedNodes {
		rec := Recommendation{
			Type:      "spot-instance",
			Workload:  fmt.Sprintf("node/%s", node.NodeName),
			Namespace: "kube-system",
			Savings:   node.MonthlySavings,
			Current: Resource{
				CPU:    "varies",
				Memory: "varies",
			},
			Recommended: Resource{
				CPU:    "varies",
				Memory: "varies",
			},
		}

		recList.Recommendations = append(recList.Recommendations, rec)
	}

	return recList
}

// Helper functions to format resources
func formatCPU(millicores float64) string {
	return fmt.Sprintf("%.0fm", millicores)
}

func formatMemory(megabytes float64) string {
	if megabytes >= 1024 {
		return fmt.Sprintf("%.1fGi", megabytes/1024)
	}
	return fmt.Sprintf("%.0fMi", megabytes)
}
