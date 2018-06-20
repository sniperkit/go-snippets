package main

// these imports auto-register the components that will be used
import (
	_ "github.com/gliderlabs/stdcom/web/auth/init"
	_ "github.com/gliderlabs/stdcom/web/console/init"
	_ "github.com/gliderlabs/stdcom/web/init"
	_ "github.com/gliderlabs/stdcom/web/sessions/init"

	_ "github.com/gl-prototypes/wiki/app/wiki"
	_ "github.com/gl-prototypes/wiki/app/wiki/store/memory"
)
