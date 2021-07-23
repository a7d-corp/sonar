package sonarconfig

type SonarConfig struct {
	Image             string
	Name              string
	Namespace         string
	NetworkPolicy     bool
	PodSecurityPolicy bool
	Privileged        bool
}
