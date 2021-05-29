package app

import "github.com/jasonblanchard/di-notebook/pkg/store"

// App dependency injection container for app
type App struct {
	StoreWriter store.Writer
	StoreReader store.Reader
}
