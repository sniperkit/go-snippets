package fsOp

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/testutil"
)

func TestPlaceFile(t *testing.T) {
	Convey("PlaceFile suite:", t, func() {
		testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
			Convey("Simple file placements should work...", func() {
				afs := osfs.New(tmpDir)
				Convey("Placing a file with read bits should work", func() {
					fsErr := PlaceFile(afs, fs.Metadata{
						Name:  fs.MustRelPath("thing"),
						Type:  fs.Type_File,
						Perms: 0644,
					}, bytes.NewBuffer([]byte("abc\n")), true)
					So(fsErr, ShouldBeNil)
					bs, err := ioutil.ReadFile(tmpDir.Join(fs.MustRelPath("thing")).String())
					So(err, ShouldBeNil)
					So(string(bs), ShouldResemble, "abc\n")
				})
				Convey("Placing a file with *no* read bits should work", func() {
					fsErr := PlaceFile(afs, fs.Metadata{
						Name:  fs.MustRelPath("thing"),
						Type:  fs.Type_File,
						Perms: 0, // this is a meaningful zero!
					}, bytes.NewBuffer([]byte("abc\n")), true)
					So(fsErr, ShouldBeNil)
					// Skip attempt to read.  If low privilege, will fail.
				})
				Convey("File placements missing parent dirs should fail", func() {
					fsErr := PlaceFile(afs, fs.Metadata{
						Name: fs.MustRelPath("deeper/thing"),
						Type: fs.Type_File,
					}, bytes.NewBuffer([]byte("abc\n")), true)
					So(fsErr.Error(), ShouldContainSubstring, "no such")
				})
			})
			Convey("Simple dir placements should work", func() {
				// TODO
			})
			Convey("Placements that would traverse a symlink out of the base path should fail", func() {
				// TODO
			})
		})
	})
}
