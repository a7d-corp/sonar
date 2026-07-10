package app

import "github.com/glitchcrab/sonar/internal/config"

const (
	appKey   contextKey = "app"
	viperKey contextKey = "viper"
)

type contextKey string

// App is the initialised and validated runtime state.
type App struct {
	Globals config.Globals
}

// RetrieveAppKey returns the context key for the App.
func (a *App) RetrieveAppKey() contextKey {
	return appKey
}

// RetrieveViperKey returns the context key for Viper.
func (a *App) RetrieveViperKey() contextKey {
	return viperKey
}
