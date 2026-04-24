package app

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// GetViper retrieves the Viper instance from the command's context.
func GetViper(cmd *cobra.Command) (*viper.Viper, error) {
	val := cmd.Context().Value("viperKey")
	if val == nil {
		return nil, fmt.Errorf("Viper not initialised")
	}

	v, ok := val.(*viper.Viper)
	if !ok {
		return nil, fmt.Errorf("invalid app type")
	}

	return v, nil
}
