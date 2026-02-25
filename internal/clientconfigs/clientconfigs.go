package clientconfigs

type SonarConfig struct {
	DryRun              bool
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
}
