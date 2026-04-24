package config

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/glitchcrab/sonar/internal/k8sclient"
)

const (
	nameMaxLength  = 50
	nameRegex      = "^[a-zA-Z0-9-.]*$"
	nameStub       = "sonar"
	namespaceRegex = "^[a-zA-Z0-9-]*$"
)

var (
	labels = map[string]string{
		"created-by": "sonar",
		"owner":      "sonar",
	}
)

func ValidateGlobalConfig(g *Globals) error {
	var errs []error

	// If the user has provided a name then validate that it looks sane.
	if g.Name != "" {
		// Restrict deployment name to 50 characters. 50 is a relatively
		// arbitrary choice of length, but it should be sufficient.
		if len(g.Name) > nameMaxLength {
			errs = append(errs, fmt.Errorf("deployment name must be 50 characters or less"))
		}

		// Validate the provided name is suitable for a Kubernetes resource name.
		ok, _ := regexp.MatchString(nameRegex, g.Name)
		if !ok {
			errs = append(errs, fmt.Errorf("deployment name can only contain alphanumeric characters, hyphens and periods"))
		}
	}

	// If a name was provided, prepend with 'sonar-' for ease of identifying Sonar deployments.
	if g.Name != "" {
		g.FullName = fmt.Sprintf("%s-%s", nameStub, g.Name)
	} else {
		g.FullName = nameStub
		g.Name = nameStub
	}

	// Add the provided name to the labels map for tagging generated resources
	g.Labels["name"] = g.Name

	// Add default labels to global Labels map
	for k, v := range labels {
		g.Labels[k] = v
	}

	// If the user has provided a namespace then validate that it looks sane.
	if g.Namespace != "" {
		// Validate the provided namespace is suitable for a Kubernetes namespace.
		ok, _ := regexp.MatchString(namespaceRegex, g.Namespace)
		if !ok {
			errs = append(errs, fmt.Errorf("namespaces can only contain alphanumeric characters and hyphens"))
		}
	}

	// If the namespace was not provided, get the namespace from the context
	var err error
	if g.Namespace == "" {
		g.Namespace, err = k8sclient.GetNamespace(g.KubeConfig, g.KubeContext)
		if err != nil {
			errs = append(errs, err)
		}
	}

	// If there were any validation errors, return them as a single error.
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
