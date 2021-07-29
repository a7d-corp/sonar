package sonarconfig

type SonarConfig struct {
	Image             string
	Labels            map[string]string
	Name              string
	Namespace         string
	NetworkPolicy     bool
	PodSecurityPolicy bool
	Privileged        bool
}
