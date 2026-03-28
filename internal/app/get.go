package app

import (
	"fmt"

	"github.com/spf13/cobra"
)

// GetApp retrieves the App instance from the command's context.
func GetApp(cmd *cobra.Command) (*App, error) {
	val := cmd.Context().Value(appKey)
	if val == nil {
		return nil, fmt.Errorf("app not initialised")
	}

	app, ok := val.(*App)
	if !ok {
		return nil, fmt.Errorf("invalid app type")
	}

	return app, nil
}
