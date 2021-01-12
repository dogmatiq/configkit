package app

// App does not implement dogma.Application interface.
type NonApp struct{}

// Foo is dummy method of App.
func (a *App) Foo() {}
