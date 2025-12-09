package kubecost

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadRecommendations reads and parses Kubecost recommendations from a JSON file
func LoadRecommendations(filepath string) (*RecommendationList, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Try to parse as our internal format first
	var recList RecommendationList
	if err := json.Unmarshal(data, &recList); err == nil {
		// Successfully parsed as internal format
		if len(recList.Recommendations) > 0 {
			return &recList, nil
		}
	}

	// Try to parse as Kubecost API format
	var apiResp KubecostSavingsResponse
	if err := json.Unmarshal(data, &apiResp); err == nil {
		// Successfully parsed as Kubecost format
		client := &Client{}
		return client.convertToInternalFormat(apiResp), nil
	}

	return nil, fmt.Errorf("failed to parse as either internal or Kubecost format")
}
