package application

import (
	"node/sentinel"
)

type (
	// Application is used as  description all application
	Application struct {
		sentinel *sentinel.Sentinel
	}
)

// createApplication is used as create a new's object which is of application
func createApplication() *Application {
	pApp := &Application{}
	pApp.sentinel = sentinel.CreateSentinel(gPullAddress)
	return pApp
}

// Run is used as run all progress
func (app *Application) Run() {
	app.sentinel.Run()
}
