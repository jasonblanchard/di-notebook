package app

import "github.com/jasonblanchard/di-notebook/store"

// App dependency injection container for app
type App struct {
	Writer store.Writer
	Reader store.Reader
}
