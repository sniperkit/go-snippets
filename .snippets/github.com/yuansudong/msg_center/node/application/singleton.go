package application

var (
	gAppInstance *Application
)

func init() {
	gAppInstance = createApplication()
}

// GetAppInstance is used as obtain a application instance
func GetAppInstance() *Application {
	return gAppInstance
}
