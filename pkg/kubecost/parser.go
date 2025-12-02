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

	var recList RecommendationList
	if err := json.Unmarshal(data, &recList); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &recList, nil
}
