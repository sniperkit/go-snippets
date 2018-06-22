package placer

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warpfork/go-errcat"

	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	. "go.polydawn.net/rio/testutil"
	. "go.polydawn.net/rio/transmat/mixins/tests"
)

func TestPlacers(t *testing.T) {
	Convey("Copy placer spec tests:", t, Requires(RequiresCanManageOwnership, func() {
		WithTmpdir(func(tmpDir fs.AbsolutePath) {
			specPlacerGood(CopyPlacer, tmpDir)
		})
	}))
	Convey("Bind placer spec tests:", t, Requires(RequiresCanMountBind, func() {
		WithTmpdir(func(tmpDir fs.AbsolutePath) {
			specPlacerGood(BindPlacer, tmpDir)
		})
	}))
	Convey("AUFS placer spec tests:", t, Requires(RequiresCanMountAny, RequiresEnvBlank("RIO_TEST_SKIP_AUFS"), func() {
		WithTmpdir(func(tmpDir fs.AbsolutePath) {
			aufsPlacer, err := NewAufsPlacer(tmpDir.Join(fs.MustRelPath("aufs")))
			So(err, ShouldBeNil)
			specPlacerGood(aufsPlacer, tmpDir)
		})
	}))
	Convey("Overlay placer spec tests:", t, Requires(RequiresCanMountAny, func() {
		WithTmpdir(func(tmpDir fs.AbsolutePath) {
			overlayPlacer, err := NewOverlayPlacer(tmpDir.Join(fs.MustRelPath("overlay")))
			So(err, ShouldBeNil)
			specPlacerGood(overlayPlacer, tmpDir)
		})
	}))
}

func specPlacerGood(placeFunc Placer, tmpDir fs.AbsolutePath) {
	afs := osfs.New(tmpDir)
	Convey("Placement of a dir should work, and maintain parent props", func() {
		PlaceFixture(afs, []FixtureFile{
			{fs.Metadata{Name: fs.MustRelPath("srcParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2004, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("srcParent/content"), Type: fs.Type_Dir, Uid: 4000, Perms: 0755, Mtime: time.Date(2005, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("srcParent/content/file"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2006, 01, 15, 0, 0, 0, 0, time.UTC)}, []byte("asdf")},
			{fs.Metadata{Name: fs.MustRelPath("dstParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
		})

		janitor, err := placeFunc(tmpDir.Join(fs.MustRelPath("srcParent/content")), tmpDir.Join(fs.MustRelPath("dstParent/content")), true)
		So(err, ShouldBeNil)

		// First check the content files and dirs.
		So(ShouldStat(afs, fs.MustRelPath("dstParent/content")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent/content"), Type: fs.Type_Dir, Uid: 4000, Perms: 0755, Mtime: time.Date(2005, 01, 15, 0, 0, 0, 0, time.UTC)})
		So(ShouldStat(afs, fs.MustRelPath("dstParent/content/file")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent/content/file"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2006, 01, 15, 0, 0, 0, 0, time.UTC), Size: 4})
		// Last (because you're most likely to screw this up) check the parent dir didn't get boinked.
		So(ShouldStat(afs, fs.MustRelPath("dstParent")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC)})

		So(janitor.Teardown(), ShouldBeNil)
	})
	Convey("Placement of a file should work, and maintain parent props", func() {
		PlaceFixture(afs, []FixtureFile{
			{fs.Metadata{Name: fs.MustRelPath("srcParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2004, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("srcParent/file"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2006, 01, 15, 0, 0, 0, 0, time.UTC)}, []byte("asdf")}, {fs.Metadata{Name: fs.MustRelPath("srcParent/content"), Type: fs.Type_Dir, Uid: 4000, Perms: 0755, Mtime: time.Date(2005, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("dstParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
		})

		janitor, err := placeFunc(tmpDir.Join(fs.MustRelPath("srcParent/file")), tmpDir.Join(fs.MustRelPath("dstParent/file")), true)
		So(err, ShouldBeNil)

		// First check the content files and dirs.
		So(ShouldStat(afs, fs.MustRelPath("dstParent/file")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent/file"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2006, 01, 15, 0, 0, 0, 0, time.UTC), Size: 4})
		// Last (because you're most likely to screw this up) check the parent dir didn't get boinked.
		So(ShouldStat(afs, fs.MustRelPath("dstParent")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC)})

		So(janitor.Teardown(), ShouldBeNil)
	})
	Convey("Placements overlapping existing content should work, and obscure it", func() {
		PlaceFixture(afs, []FixtureFile{
			{fs.Metadata{Name: fs.MustRelPath("srcParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2004, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("srcParent/content"), Type: fs.Type_Dir, Uid: 4000, Perms: 0755, Mtime: time.Date(2005, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("srcParent/content/file"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2006, 01, 15, 0, 0, 0, 0, time.UTC)}, []byte("asdf")},
			{fs.Metadata{Name: fs.MustRelPath("dstParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("dstParent/content"), Type: fs.Type_Dir, Perms: 0600, Mtime: time.Date(2029, 01, 15, 0, 0, 0, 0, time.UTC)}, nil},
			{fs.Metadata{Name: fs.MustRelPath("dstParent/content/chump"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2106, 01, 15, 0, 0, 0, 0, time.UTC)}, []byte("qwer")},
		})

		janitor, err := placeFunc(tmpDir.Join(fs.MustRelPath("srcParent/content")), tmpDir.Join(fs.MustRelPath("dstParent/content")), true)
		So(err, ShouldBeNil)

		// First check the content files and dirs.
		So(ShouldStat(afs, fs.MustRelPath("dstParent/content")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent/content"), Type: fs.Type_Dir, Uid: 4000, Perms: 0755, Mtime: time.Date(2005, 01, 15, 0, 0, 0, 0, time.UTC)})
		So(ShouldStat(afs, fs.MustRelPath("dstParent/content/file")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent/content/file"), Type: fs.Type_File, Perms: 0640, Mtime: time.Date(2006, 01, 15, 0, 0, 0, 0, time.UTC), Size: 4})
		_, err = afs.LStat(fs.MustRelPath("dstParent/content/chump"))
		So(err, errcat.ErrorShouldHaveCategory, fs.ErrNotExists)
		// Last (because you're most likely to screw this up) check the parent dir didn't get boinked.
		So(ShouldStat(afs, fs.MustRelPath("dstParent")), ShouldResemble, fs.Metadata{Name: fs.MustRelPath("dstParent"), Type: fs.Type_Dir, Perms: 0755, Mtime: time.Date(2019, 01, 15, 0, 0, 0, 0, time.UTC)})

		So(janitor.Teardown(), ShouldBeNil)
	})
}
