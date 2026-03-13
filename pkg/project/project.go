package project

var (
	description = "Sonar deploys a debugging container to a Kubernetes cluster."
	name        = "sonar"
	source      = "https://github.com/a7d-corp/sonar"
	version     = "1.0.0"
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
