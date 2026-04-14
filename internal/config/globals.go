package config

// Globals contains the global user-provided configuration
type Globals struct {
	KubeConfig  string
	KubeContext string
	FullName    string
	Name        string
	Namespace   string
	Labels      map[string]string
}
