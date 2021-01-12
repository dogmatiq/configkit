package app

// NonApp does not implement dogma.Application interface.
type NonApp struct{}

// Foo is dummy method of App.
func (a *NonApp) Foo() {}
