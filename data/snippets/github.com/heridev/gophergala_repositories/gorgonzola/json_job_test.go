package gorgonzola

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJob(t *testing.T) {
	Convey("Given invalid JsonJob", t, func() {
		jsonJob := "{}"
		err := validateDoc(jsonJob)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, `(root) : "company" property is missing and required, given {}; (root) : "url" property is missing and required, given {}; (root) : "remoteFriendly" property is missing and required, given {}; (root) : "market" property is missing and required, given {}; (root) : "size" property is missing and required, given {}; (root) : "jobs" property is missing and required, given {}; `)
	})
}
