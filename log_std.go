// +build !appengine

package rest

import (
	"github.com/bradberger/context"

	log "github.com/sirupsen/logrus"
)

// Criticalf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Criticalf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Debugf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Debugf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Errorf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Errorf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Infof is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Infof LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warningf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Warningf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Warningf(format, args...)
}

// Fatalf is shorthand for the appeninge/log func with the same name, with the added advantage that it's a variable and can be overridden if needed
var Fatalf LogFunc = func(ctx context.Context, format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
