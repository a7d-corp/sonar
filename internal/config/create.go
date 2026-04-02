package config

// CreateConfig contains the create-specific user-provided configuration
type CreateConfig struct {
	DryRun              bool
	FullName            string
	Image               string
	Labels              map[string]string
	Name                string
	Namespace           string
	NetworkPolicy       bool
	NodeExec            bool
	NodeName            string
	NonRoot             bool
	PodArgs             string
	PodCommand          string
	PodGroup            int64
	PodUser             int64
	Privileged          bool
	PrivilegeEscalation bool
	UnprivilegedPing    bool
}
