package policyengine

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadPolicies loads all policies from a directory
func LoadPolicies(policyDir string) ([]Policy, error) {
	var policies []Policy

	// Walk through policy directory
	err := filepath.Walk(policyDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-YAML files
		if info.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}

		// Load policy file
		policy, err := LoadPolicy(path)
		if err != nil {
			return fmt.Errorf("failed to load %s: %w", path, err)
		}

		// Only add enabled policies
		if policy.Enabled {
			policies = append(policies, *policy)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return policies, nil
}

// LoadPolicy loads a single policy from a file
func LoadPolicy(filepath string) (*Policy, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	var policy Policy
	if err := yaml.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse policy YAML: %w", err)
	}

	return &policy, nil
}
