package tests

import (
	"os"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/warpfork/go-errcat"

	"go.polydawn.net/rio/fs"
)

func CheckBaseLstat(afs fs.FS) {
	// yes, this is making certain assumptions about your setup.
	Convey("SPEC: lstat on the base should see dir", func() {
		d1 := fs.MustRelPath(".")
		stat, err := afs.LStat(d1)
		So(err, ShouldBeNil)
		So(stat.Type, ShouldEqual, fs.Type_Dir)
	})
}

func CheckMkdirLstatRoundtrip(afs fs.FS) {
	Convey("SPEC: mkdir and lstat should roundtrip", func() {
		d1 := fs.MustRelPath("d1")
		So(afs.Mkdir(d1, 0755), ShouldBeNil)
		stat, err := afs.LStat(d1)
		So(err, ShouldBeNil)
		So(stat.Type, ShouldEqual, fs.Type_Dir)
	})
}

func CheckDeepMkdirError(afs fs.FS) {
	Convey("SPEC: deep mkdir should error", func() {
		d1d2 := fs.MustRelPath("d1/d2")
		So(afs.Mkdir(d1d2, 0755), errcat.ErrorShouldHaveCategory, fs.ErrNotExists)
		_, err := afs.LStat(d1d2)
		So(err, errcat.ErrorShouldHaveCategory, fs.ErrNotExists)
	})
}

func CheckMklinkLstatRoundtrip(afs fs.FS) {
	Convey("SPEC: mklink and lstat should roundtrip", func() {
		l1 := fs.MustRelPath("l1")
		So(afs.Mklink(l1, "./target"), ShouldBeNil)
		stat, err := afs.LStat(l1)
		So(err, ShouldBeNil)
		So(stat.Type, ShouldEqual, fs.Type_Symlink)
	})
}

func CheckSymlinks(afs fs.FS) {
	Convey("SPEC: symlink resolve", func() {
		Convey("symlinks to files resolve correctly", func() {
			Convey("short relative case", func() {
				l1 := fs.MustRelPath("l1")
				linkStr := "./target"
				target := fs.MustRelPath(linkStr)

				So(afs.Mklink(l1, linkStr), ShouldBeNil)
				So(makeFile(afs, target, "body"), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
			Convey("long relative case", func() {
				d1l1 := fs.MustRelPath("d1/l1")
				linkStr := ".././/d2/target"
				target := fs.MustRelPath("d2/target")

				So(afs.Mkdir(d1l1.Dir(), 0755), ShouldBeNil)
				So(afs.Mklink(d1l1, linkStr), ShouldBeNil)
				So(afs.Mkdir(target.Dir(), 0755), ShouldBeNil)
				So(makeFile(afs, target, "body"), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, d1l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
			Convey("loooong relative case", func() {
				d1d2d3l1 := fs.MustRelPath("d1/d2/d3/l1")
				linkStr := "../.././..//d2/target"
				target := fs.MustRelPath("d2/target")

				So(afs.Mkdir(d1d2d3l1.Dir().Dir().Dir(), 0755), ShouldBeNil)
				So(afs.Mkdir(d1d2d3l1.Dir().Dir(), 0755), ShouldBeNil)
				So(afs.Mkdir(d1d2d3l1.Dir(), 0755), ShouldBeNil)
				So(afs.Mklink(d1d2d3l1, linkStr), ShouldBeNil)
				So(afs.Mkdir(target.Dir(), 0755), ShouldBeNil)
				So(makeFile(afs, target, "body"), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, d1d2d3l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
			Convey("rooted case", func() {
				l1 := fs.MustRelPath("l1")
				linkStr := "/target"
				target := fs.MustRelPath("./target")

				So(afs.Mklink(l1, linkStr), ShouldBeNil)
				So(makeFile(afs, target, "body"), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
			Convey("dotdot-overload case", func() {
				d1l1 := fs.MustRelPath("d1/l1")
				linkStr := ".././/../../../d2/target"
				target := fs.MustRelPath("d2/target")

				So(afs.Mkdir(d1l1.Dir(), 0755), ShouldBeNil)
				So(afs.Mklink(d1l1, linkStr), ShouldBeNil)
				So(afs.Mkdir(target.Dir(), 0755), ShouldBeNil)
				So(makeFile(afs, target, "body"), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, d1l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
		})
		Convey("dangling symlinks resolve correctly", func() {
			Convey("short relative case", func() {
				l1 := fs.MustRelPath("l1")
				linkStr := "./target"
				target := fs.MustRelPath(linkStr)

				So(afs.Mklink(l1, linkStr), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
			Convey("long relative case", func() {
				d1l1 := fs.MustRelPath("d1/l1")
				linkStr := ".././/d2/target"
				target := fs.MustRelPath("d2/target")

				So(afs.Mkdir(d1l1.Dir(), 0755), ShouldBeNil)
				So(afs.Mklink(d1l1, linkStr), ShouldBeNil)
				So(afs.Mkdir(target.Dir(), 0755), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, d1l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
			Convey("rooted case", func() {
				l1 := fs.MustRelPath("l1")
				linkStr := "/target"
				target := fs.MustRelPath("./target")

				So(afs.Mklink(l1, linkStr), ShouldBeNil)

				resolved, err := afs.ResolveLink(linkStr, l1)
				So(err, ShouldBeNil)
				So(resolved, ShouldResemble, target)
			})
		})
	})
}

func CheckPerniciousSymlinks(afs fs.FS) {
	Convey("SPEC: perverse symlinks error (and do not hang!)", func() {
		Convey("symlink to self should error", func() {
			l1 := fs.MustRelPath("l1")
			linkStr := "./l1"

			So(afs.Mklink(l1, linkStr), ShouldBeNil)

			_, err := afs.ResolveLink(linkStr, l1)
			So(err, errcat.ErrorShouldHaveCategory, fs.ErrRecursion)
		})
		Convey("symlink to self via dotdot should error", func() {
			d1l1 := fs.MustRelPath("d1/l1")
			linkStr := "../d1/l1"

			So(afs.Mkdir(d1l1.Dir(), 0755), ShouldBeNil)
			So(afs.Mklink(d1l1, linkStr), ShouldBeNil)

			_, err := afs.ResolveLink(linkStr, d1l1)
			So(err, errcat.ErrorShouldHaveCategory, fs.ErrRecursion)
		})
		Convey("symlink to cycle pair should error", func() {
			l1 := fs.MustRelPath("l1")
			link1Str := "./l2"
			l2 := fs.MustRelPath("l2")
			link2Str := "./l1"

			So(afs.Mklink(l1, link1Str), ShouldBeNil)
			So(afs.Mklink(l2, link2Str), ShouldBeNil)

			_, err := afs.ResolveLink(link1Str, l1)
			So(err, errcat.ErrorShouldHaveCategory, fs.ErrRecursion)
			_, err = afs.ResolveLink(link2Str, l2)
			So(err, errcat.ErrorShouldHaveCategory, fs.ErrRecursion)
		})
	})
}

func CheckOpsTraversingSymlinks(afs fs.FS) {
	Convey("SPEC: regular ops traversing symlinks must behave", func() {

		// This is a basic, but low-stress test that traversal happens at all.
		//  It's weak because if you *failed* to resolve, the OS will do roughly
		//  the same; we're not testing that our semantic sandbox was exercised.
		Convey("ops crossing symlink to dir should fly", func() {
			So(afs.Mkdir(fs.MustRelPath("zone"), 0755), ShouldBeNil)
			So(afs.Mklink(fs.MustRelPath("lnk"), "zone"), ShouldBeNil)

			So(afs.Mkdir(fs.MustRelPath("lnk/test"), 0755), ShouldBeNil)

			stat, err := afs.LStat(fs.MustRelPath("zone/test"))
			So(err, ShouldBeNil)
			So(stat.Type, ShouldEqual, fs.Type_Dir)
			So(stat.Name, ShouldResemble, fs.MustRelPath("zone/test"))
			stat, err = afs.LStat(fs.MustRelPath("lnk/test"))
			So(err, ShouldBeNil)
			So(stat.Type, ShouldEqual, fs.Type_Dir)
			So(stat.Name, ShouldResemble, fs.MustRelPath("lnk/test"))
		})

		// This test actually tests that *we're* the ones doing resolve,
		//  because if it gets to the os calls, they'll see a dangle.
		Convey("ops crossing bounded over-dotted symlink to dir should fly", func() {
			So(afs.Mkdir(fs.MustRelPath("zone"), 0755), ShouldBeNil)
			So(afs.Mklink(fs.MustRelPath("lnk"), "../../../zone"), ShouldBeNil)

			So(afs.Mkdir(fs.MustRelPath("lnk/test"), 0755), ShouldBeNil)

			stat, err := afs.LStat(fs.MustRelPath("zone/test"))
			So(err, ShouldBeNil)
			So(stat.Type, ShouldEqual, fs.Type_Dir)
			So(stat.Name, ShouldResemble, fs.MustRelPath("zone/test"))
			stat, err = afs.LStat(fs.MustRelPath("lnk/test"))
			So(err, ShouldBeNil)
			So(stat.Type, ShouldEqual, fs.Type_Dir)
			So(stat.Name, ShouldResemble, fs.MustRelPath("lnk/test"))
		})
	})
}

func makeFile(afs fs.FS, path fs.RelPath, body string) error {
	f, err := afs.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(body))
	return err
}
