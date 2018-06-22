// Implements HTTP 1.1 Conditional Requests
package conditional

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTP headers for conditional requests
const (
	IfModifiedSince   = "If-Modified-Since"
	IfUnmodifiedSince = "If-Unmodified-Since"
	IfMatch           = "If-Match"
	IfNoneMatch       = "If-None-Match"
	IfRange           = "If-Range"
)

// HTTP Methods that are checked for when calculating conditional requests
const (
	Get  = "GET"
	Head = "HEAD"
)

// HTTP header used when calculating Range based conditional requests
const Range = "Range"

type Etagger interface {
	// Etag
	Etag() (string, error)
}

type LastModifier interface {
	LastModified() time.Time
}

var (
	// An error that the Etag function can return to signify that
	// no resource exists at the given Location.
	ErrNoResource = errors.New("No Resource Available")

	// Server MUST respond with either:
	// a) the 412 (Precondition Failed) status code
	// b) one of the 2xx (Successful) status codes if the
	//    origin server has verified that a state change is being
	//    requested and the final state is already reflected
	//    in the current state of the target resource
	ErrWasModified = errors.New("Resource was modified since header, check if final state would match")

	ErrRangeMismatch = errors.New("Calculating If-Range failed, respond with entire resource")
)

func Conditional(c *gin.Context, resource interface{}) (bool, error) {
	etagger, canCheckEtag := resource.(Etagger)
	modifier, canCheckModifier := resource.(LastModifier)

	if header := c.Request.Header.Get(IfMatch); canCheckEtag && header != "" {

		// Does the request have an If-Match header?
		if handleIfMatch(etagger, header) == false {
			return false, ErrWasModified
		}

	} else if header := c.Request.Header.Get(IfUnmodifiedSince); canCheckModifier && header != "" {

		// Does the request have an If-Unmodified-Since header?
		if handleIfUnmodifiedSince(modifier, header) == false {
			return false, ErrWasModified
		}

	}

	if header := c.Request.Header.Get(IfNoneMatch); canCheckEtag && header != "" {

		// Does the request have an If-None-Match header?
		if handleIfNoneMatch(etagger, header) == false {
			if c.Request.Method == Get || c.Request.Method == Head {
				c.AbortWithStatus(http.StatusNotModified)
				return true, nil
			} else {
				c.AbortWithStatus(http.StatusPreconditionFailed)
				return true, nil
			}
		}

	} else if c.Request.Method != Get || c.Request.Method != Head {
		return false, nil
	} else if header := c.Request.Header.Get(IfModifiedSince); canCheckModifier && header != "" {
		if handleIfModifiedSince(modifier, header) == false {
			c.AbortWithStatus(http.StatusNotModified)
			return true, nil
		}
	}

	if header := c.Request.Header.Get(IfRange); c.Request.Method == Get &&
		c.Request.Header.Get(Range) != "" && header != "" {
		if handleIfRange(etagger, header) == false {
			return false, ErrRangeMismatch
		}
	}

	return false, nil
}

// Implements the Section 3.1 from RFC7232
// https://tools.ietf.org/html/rfc7232#section-3.1
func handleIfMatch(resource Etagger, clientEtag string) bool {
	serverEtag, err := resource.Etag()
	if err != nil && err != ErrNoResource {
		return false
	}

	if clientEtag == "*" && err == ErrNoResource {
		return false
	}

	if serverEtag == clientEtag || clientEtag == "*" {
		return true
	} else {
		return false
	}
}

// Implements the Section 3.4 from RFC7232
// https://tools.ietf.org/html/rfc7232#section-3.4
func handleIfUnmodifiedSince(resource LastModifier, date string) bool {
	clientDate, err := time.Parse(http.TimeFormat, date)
	if err != nil {
		// A recipient MUST ignore the If-Unmodified-Since header field if the
		// received field-value is not a valid HTTP-date.
		return false
	}

	serverDate := resource.LastModified()
	if clientDate.Before(serverDate) {
		return true
	}

	return false
}

// Implements the Section 3.2 from RFC7232
// https://tools.ietf.org/html/rfc7232#section-3.2
func handleIfNoneMatch(resource Etagger, clientEtag string) bool {
	serverEtag, err := resource.Etag()
	if err != nil {
		if clientEtag == "*" && err == ErrNoResource {
			return true
		}
		return false
	}

	if clientEtag == serverEtag {
		return false
	}

	return true
}

// Implements the Section 3.3 from RFC7232
// https://tools.ietf.org/html/rfc7232#section-3.3
func handleIfModifiedSince(resource LastModifier, date string) bool {
	clientDate, err := time.Parse(http.TimeFormat, date)
	if err != nil {
		return false
	}

	serverDate := resource.LastModified()

	// A date which is later than the server's current time is invalid.
	// If the header is earlier than the current date, request should continue
	if clientDate.After(time.Now()) || clientDate.Before(serverDate) {
		return false
	}

	return true
}

func handleIfRange(resource Etagger, clientEtag string) bool {

	return false
}
