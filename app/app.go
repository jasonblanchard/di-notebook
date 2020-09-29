package app

import "github.com/jasonblanchard/di-notebook/store"

// App dependency injection container for app
type App struct {
	StoreWriter store.Writer
	StoreReader store.Reader
}
