package sonarconfig

type SonarConfig struct {
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
	PodSecurityPolicy   bool
	PodGroup            int64
	PodUser             int64
	Privileged          bool
	PrivilegeEscalation bool
}
