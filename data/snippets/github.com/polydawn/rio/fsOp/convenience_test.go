package fsOp

import (
	"bytes"
	"io"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warpfork/go-errcat"

	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	. "go.polydawn.net/rio/testutil"
)

func TestMkdirAll(t *testing.T) {
	// Note that all of these are assuming PlaceFile already works just fine.
	Convey("MkdirAll:", t, func() {
		WithTmpdir(func(tmpDir fs.AbsolutePath) {
			afs := osfs.New(tmpDir)
			Convey("MkdirAll on an existing path should work...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Perms: 0755}, nil)

				So(MkdirAll(afs, fs.MustRelPath("dir"), 0755), ShouldBeNil)
			})
			Convey("MkdirAll creating one node should work...", func() {
				So(MkdirAll(afs, fs.MustRelPath("dir"), 0755), ShouldBeNil)
				stat, err := afs.LStat(fs.MustRelPath("dir"))
				So(err, ShouldBeNil)
				So(stat.Type, ShouldEqual, fs.Type_Dir)
			})
			Convey("MkdirAll creating several nodes should work...", func() {
				So(MkdirAll(afs, fs.MustRelPath("dir/2/3"), 0755), ShouldBeNil)
				stat, err := afs.LStat(fs.MustRelPath("dir/2/3"))
				So(err, ShouldBeNil)
				So(stat.Type, ShouldEqual, fs.Type_Dir)
			})
			Convey("MkdirAll on an existing file should error...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("womp"), Type: fs.Type_File, Perms: 0755}, nil)

				So(MkdirAll(afs, fs.MustRelPath("womp"), 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotDir)
			})
			Convey("MkdirAll traversing existing file should error...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("womp"), Type: fs.Type_File, Perms: 0755}, nil)

				So(MkdirAll(afs, fs.MustRelPath("womp/2/3"), 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotDir)
			})
			Convey("MkdirAll traversing symlinks should work...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Perms: 0755}, nil)
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("lnk"), Type: fs.Type_Symlink, Linkname: "./dir"}, nil)

				So(MkdirAll(afs, fs.MustRelPath("lnk/woo"), 0755), ShouldBeNil)
				stat, err := afs.LStat(fs.MustRelPath("dir/woo"))
				So(err, ShouldBeNil)
				So(stat.Type, ShouldEqual, fs.Type_Dir)
			})
			Convey("MkdirAll with a dangling symlink should error...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("lnk"), Type: fs.Type_Symlink, Linkname: "./dir"}, nil)

				// In case you're wondering, yes, this is the same as
				//  `os.MkdirAll(afs.BasePath().String()+"/lnk/woo", 0755)`
				//  which would also error with "file exists".
				So(MkdirAll(afs, fs.MustRelPath("lnk/woo"), 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotDir)
			})
			Convey("MkdirAll traversing symlink to a file should error...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("file"), Type: fs.Type_File, Perms: 0644}, nil)
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("lnk"), Type: fs.Type_Symlink, Linkname: "./file"}, nil)

				So(MkdirAll(afs, fs.MustRelPath("lnk/woo"), 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotDir)
			})
			Convey("MkdirAll onto a symlink to a file should error...", func() {
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("file"), Type: fs.Type_File, Perms: 0644}, nil)
				mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("lnk"), Type: fs.Type_Symlink, Linkname: "./file"}, nil)

				So(MkdirAll(afs, fs.MustRelPath("lnk"), 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotDir)
			})
			Convey("MkdirAll when the entire filesystem DNE should error...", func() {
				afs := osfs.New(tmpDir.Join(fs.MustRelPath("nope")))

				So(MkdirAll(afs, fs.MustRelPath("dir"), 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotExists)
			})
		})
	})
}

func TestMkdirUsable(t *testing.T) {
	// Note that all of these are assuming PlaceFile already works just fine.
	Convey("MkdirUsable:", t,
		Requires(RequiresCanManageOwnership, func() {
			var (
				preferredTime  = time.Unix(500004440, 0).UTC()
				otherTime      = time.Unix(500000220, 0).UTC()
				preferredProps = fs.Metadata{
					Perms: 0755,
					Uid:   1,
					Gid:   2,
					Mtime: preferredTime,
				}
			)
			WithTmpdir(func(tmpDir fs.AbsolutePath) {
				afs := osfs.New(tmpDir)
				Convey("MkdirUsable on an empty path should work...", func() {
					So(MkdirUsable(afs, fs.MustRelPath("dir"), preferredProps), ShouldBeNil)

					stat, err := afs.LStat(fs.MustRelPath("dir"))
					So(err, ShouldBeNil)
					So(stat.Type, ShouldEqual, fs.Type_Dir)
					So(stat.Uid, ShouldEqual, 1)
					So(stat.Gid, ShouldEqual, 2)
				})
				Convey("MkdirUsable on a deep path should work...", func() {
					So(MkdirUsable(afs, fs.MustRelPath("dir/a/b"), preferredProps), ShouldBeNil)

					stat, err := afs.LStat(fs.MustRelPath("dir/a/b"))
					So(err, ShouldBeNil)
					So(stat.Type, ShouldEqual, fs.Type_Dir)
					So(stat.Uid, ShouldEqual, 1)
					So(stat.Gid, ShouldEqual, 2)
				})
				Convey("MkdirUsable on a mostly extant path (add last dir only) should work...", func() {
					mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Perms: 0755, Mtime: otherTime}, nil)

					So(MkdirUsable(afs, fs.MustRelPath("dir/tgt"), preferredProps), ShouldBeNil)

					So(ShouldStat(afs, fs.MustRelPath("dir/tgt")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir/tgt"), Type: fs.Type_Dir, Uid: 1, Gid: 2, Perms: 0755, Mtime: preferredTime})
					So(ShouldStat(afs, fs.MustRelPath("dir")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Uid: 0, Gid: 0, Perms: 0755, Mtime: otherTime})
				})
				Convey("MkdirUsable on a partially extant path (add several dirs) should work...", func() {
					mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Perms: 0755, Mtime: otherTime}, nil)

					So(MkdirUsable(afs, fs.MustRelPath("dir/a/b/tgt"), preferredProps), ShouldBeNil)

					So(ShouldStat(afs, fs.MustRelPath("dir/a/b/tgt")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir/a/b/tgt"), Type: fs.Type_Dir, Uid: 1, Gid: 2, Perms: 0755, Mtime: preferredTime})
					So(ShouldStat(afs, fs.MustRelPath("dir/a/b")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir/a/b"), Type: fs.Type_Dir, Uid: 1, Gid: 2, Perms: 0755, Mtime: preferredTime})
					So(ShouldStat(afs, fs.MustRelPath("dir/a")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir/a"), Type: fs.Type_Dir, Uid: 1, Gid: 2, Perms: 0755, Mtime: preferredTime})
					So(ShouldStat(afs, fs.MustRelPath("dir")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Uid: 0, Gid: 0, Perms: 0755, Mtime: otherTime})
				})
				Convey("MkdirUsable on a needing to repair parent perms should work...", func() {
					mustPlaceFile(afs, fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Perms: 0700, Mtime: otherTime}, nil)

					So(MkdirUsable(afs, fs.MustRelPath("dir/tgt"), preferredProps), ShouldBeNil)

					So(ShouldStat(afs, fs.MustRelPath("dir/tgt")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir/tgt"), Type: fs.Type_Dir, Uid: 1, Gid: 2, Perms: 0755, Mtime: preferredTime})
					So(ShouldStat(afs, fs.MustRelPath("dir")), ShouldResemble,
						fs.Metadata{Name: fs.MustRelPath("dir"), Type: fs.Type_Dir, Uid: 0, Gid: 0, Perms: 0701, Mtime: otherTime})
				})
			})
		}),
	)
}

func mustPlaceFile(afs fs.FS, fmeta fs.Metadata, body io.Reader) {
	if fmeta.Type == fs.Type_File && body == nil {
		body = &bytes.Buffer{}
	}
	if err := PlaceFile(afs, fmeta, body, true); err != nil {
		panic(err)
	}
}
