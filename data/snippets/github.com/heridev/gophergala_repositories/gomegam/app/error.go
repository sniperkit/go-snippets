package app


import (
	"fmt"
)

type AppLifecycleError struct {
	app string
	Err error
}

func (e *AppLifecycleError) Error() string {
	return fmt.Sprintf("gulpd failed to apply the lifecle to the app %q: %s", e.app, e.Err)
}

