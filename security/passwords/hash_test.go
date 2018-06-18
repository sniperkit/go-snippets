package passwords

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTestPasswordHasing(t *testing.T) {
	Convey("TestPasswordHasing", t, func() {
		hash1, err := Hash("mypassword1")
		So(err, ShouldBeNil)
		So(hash1, ShouldNotEqual, "")
		So(hash1, ShouldNotEqual, "mypassword1")

		hash2, err := Hash("mypassword1") //same password
		So(err, ShouldBeNil)
		So(hash2, ShouldNotEqual, "")
		So(hash2, ShouldNotEqual, "mypassword1")

		hash3, err := Hash("mypassword2") //different password
		So(err, ShouldBeNil)
		So(hash3, ShouldNotEqual, "")
		So(hash3, ShouldNotEqual, "mypassword2")

		So(hash2, ShouldNotEqual, hash1)
		So(Verify("mypassword1", hash1), ShouldBeNil)
		So(Verify("mypassword1", hash2), ShouldBeNil)
		So(Verify("mypassword1", "mypassword1"), ShouldNotBeNil)
		So(Verify("mypassword1", hash3), ShouldNotBeNil)
		So(Verify("mypassword2", hash3), ShouldBeNil)
	})
}
