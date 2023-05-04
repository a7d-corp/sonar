package project

var (
	description = "Sonar deploys a debugging container to a Kubernetes cluster."
	name        = "sonar"
	source      = "https://github.com/glitchcrab/sonar"
	version     = "0.8.1-dev"
)

func Description() string {
	return description
}

func Name() string {
	return name
}

func Source() string {
	return source
}

func Version() string {
	return version
}
