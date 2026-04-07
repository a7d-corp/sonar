package config

// DeleteConfig contains the delete-specific user-provided configuration
type DeleteConfig struct {
	SearchLabels []string
	Name         string
	Namespace    string
}
