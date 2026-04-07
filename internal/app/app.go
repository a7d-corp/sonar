package app

import "github.com/glitchcrab/sonar/internal/config"

// appKey is a custom type to avoid potential collisions.
const appKey contextKey = "app"

type contextKey string

// App is the initialised and validated runtime state.
type App struct {
	Globals config.Globals
}

// RetrieveContextKey returns the context key for the App.
func (a *App) RetrieveAppKey() contextKey {
	return appKey
}
