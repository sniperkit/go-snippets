// +build appengine

package rest

import (
	"github.com/bradberger/context"

	"google.golang.org/appengine/log"
)

// Criticalf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Criticalf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Criticalf(ctx, format, args...)
}

// Debugf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Debugf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(ctx, format, args...)
}

// Errorf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Errorf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}

// Infof is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Infof LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

// Warningf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Warningf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(ctx, format, args...)
}
