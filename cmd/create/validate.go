package create

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/glitchcrab/sonar/internal/config"
	"github.com/spf13/cobra"
)

const (
	imageRegex = "^[a-z0-9/.-]*[:][a-z0-9.-]*$"
)

func validateCreateConfig(command *cobra.Command, c *config.CreateConfig) error {
	var errs []error

	if command.Flags().Lookup("image").Changed {
		// Validate image to see if a tag has been provided; if not then
		// use :latest. Does not validate full image name, just whether a
		// tag was provided.
		ok, _ := regexp.MatchString(imageRegex, image)
		if !ok {
			image = fmt.Sprintf("%s:latest", image)
		}
	}

	// Set sane options if we're exec-ing into a node.
	if c.NodeExec {
		// Error out if node name was not provided.
		if c.NodeName == "" {
			errs = append(errs, fmt.Errorf("--node-exec also requires --node-name to be provided"))
		}

		// Set options that don't make sense for node exec to their defaults.
		c.NetworkPolicy = false
		c.Privileged = true
	}

	// If there were any validation errors, return them as a single error.
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
