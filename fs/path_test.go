package fs

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//--------------
// RelPath
//--------------

func TestRelPath(t *testing.T) {
	Convey("RelPath stringer suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    RelPath
			str   string
		}{
			{"zero values",
				RelPath{},
				"."},
			{"dot value",
				MustRelPath("."),
				"."},
			{"short value",
				MustRelPath("aa"),
				"./aa"},
			{"long value",
				MustRelPath("a/bb/ccc"),
				"./a/bb/ccc"},
			{"denormalized value",
				MustRelPath("../a/bb/../ccc"),
				"../a/ccc"},
			{"lone doubledot value",
				MustRelPath("../"),
				".."},
			{"dotted value",
				MustRelPath(".aa"),
				"./.aa"},
			{"dotted2 value",
				MustRelPath("..aa"),
				"./..aa"},
			{"dotted3 value",
				MustRelPath("...aa"),
				"./...aa"},
		} {
			Convey(tr.title, func() {
				v := fmt.Sprintf("%s", tr.p1)
				So(v, ShouldResemble, tr.str)
			})
		}
	})
}

func TestRelPathDir(t *testing.T) {
	Convey("RelPath.Dir suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    RelPath
			pdir  RelPath
		}{
			{"zero values",
				RelPath{},
				RelPath{}},
			{"dot value",
				MustRelPath("."),
				RelPath{}},
			{"short value",
				MustRelPath("aa"),
				RelPath{}},
			{"long value",
				MustRelPath("a/bb/ccc"),
				MustRelPath("a/bb")},
			{"denormalized value",
				MustRelPath("../a/bb/../ccc"),
				MustRelPath("../a")}, // cleans, then drops
			{"lone doubledot value",
				MustRelPath("../"),
				MustRelPath(".")}, // yep.  matches what stdlib 'path.Dir' does.
			{"double doubledot value",
				MustRelPath("../.."),
				MustRelPath("..")}, // yep.  matches what stdlib 'path.Dir' does.
		} {
			Convey(tr.title, func() {
				v := tr.p1.Dir()
				So(v, ShouldResemble, tr.pdir)
			})
		}
	})
}

func TestRelPathLast(t *testing.T) {
	Convey("RelPath.Last suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    RelPath
			last  string
		}{
			{"zero values",
				RelPath{},
				"."},
			{"dot value",
				MustRelPath("."),
				"."},
			{"short value",
				MustRelPath("aa"),
				"aa"},
			{"long value",
				MustRelPath("a/bb/ccc"),
				"ccc"},
			{"denormalized value",
				MustRelPath("../a/bb/../ccc"),
				"ccc"},
			{"lone doubledot value",
				MustRelPath("../"),
				".."},
		} {
			Convey(tr.title, func() {
				v := tr.p1.Last()
				So(v, ShouldResemble, tr.last)
			})
		}
	})
}

func TestRelPathJoins(t *testing.T) {
	Convey("RelPath.Join suite:", t, func() {
		for _, tr := range []struct {
			title  string
			p1, p2 RelPath
			pj     RelPath
		}{
			{"zero values",
				RelPath{}, RelPath{},
				RelPath{}},
			{"regular values",
				MustRelPath("rel"), MustRelPath("pth"),
				MustRelPath("rel/pth")},
			{"zero,short",
				MustRelPath("."), MustRelPath("pth"),
				MustRelPath("pth")},
			{"long,short",
				MustRelPath("r/el"), MustRelPath("pth"),
				MustRelPath("r/el/pth")},
			{"long,zero",
				MustRelPath("a/bb/ccc"), MustRelPath("."),
				MustRelPath("a/bb/ccc")},
			{"long,long",
				MustRelPath("a/bb/ccc"), MustRelPath("dd/e"),
				MustRelPath("a/bb/ccc/dd/e")},
			{"zero,up",
				MustRelPath("."), MustRelPath(".."),
				MustRelPath("..")},
			{"short,up",
				MustRelPath("rel"), MustRelPath(".."),
				MustRelPath(".")},
			{"long,up",
				MustRelPath("r/el"), MustRelPath(".."),
				MustRelPath("r")},
			{"dotted,short",
				MustRelPath(".dot"), MustRelPath("wonk"),
				MustRelPath(".dot/wonk")},
			{"dotted,dotted",
				MustRelPath(".dot"), MustRelPath(".wonk"),
				MustRelPath(".dot/.wonk")},
		} {
			Convey(tr.title, func() {
				v := tr.p1.Join(tr.p2)
				So(v, ShouldResemble, tr.pj)
			})
		}
	})
}

func TestRelPathSplits(t *testing.T) {
	Convey("RelPath.Split suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    RelPath
			ps    []RelPath
		}{
			{"zero values",
				RelPath{},
				[]RelPath{RelPath{}}},
			{"len=1 values",
				MustRelPath("./a"),
				[]RelPath{RelPath{}, MustRelPath("a")}},
			{"len=3 values",
				MustRelPath("./a/bb/c"),
				[]RelPath{RelPath{}, MustRelPath("a"), MustRelPath("a/bb"), MustRelPath("a/bb/c")}},
			{"dotted values",
				MustRelPath("./.a/bb/.c"),
				[]RelPath{RelPath{}, MustRelPath(".a"), MustRelPath(".a/bb"), MustRelPath(".a/bb/.c")}},
		} {
			Convey(tr.title, func() {
				v := tr.p1.Split()
				So(v, ShouldResemble, tr.ps)
			})
		}
	})
}

func TestRelPathSplitParent(t *testing.T) {
	Convey("RelPath.SplitParent suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    RelPath
			ps    []RelPath
		}{
			{"zero values",
				RelPath{},
				[]RelPath{}},
			{"len=1 values",
				MustRelPath("./a"),
				[]RelPath{RelPath{}}},
			{"len=3 values",
				MustRelPath("./a/bb/c"),
				[]RelPath{RelPath{}, MustRelPath("a"), MustRelPath("a/bb")}},
			{"dotted values",
				MustRelPath("./.a/bb/.c"),
				[]RelPath{RelPath{}, MustRelPath(".a"), MustRelPath(".a/bb")}},
		} {
			Convey(tr.title, func() {
				v := tr.p1.SplitParent()
				So(v, ShouldResemble, tr.ps)
			})
		}
	})
}

//--------------
// AbsolutePath
//--------------

func TestAbsolutePath(t *testing.T) {
	Convey("AbsolutePath stringer suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    AbsolutePath
			str   string
		}{
			{"zero values",
				AbsolutePath{},
				"/"},
			{"root value",
				MustAbsolutePath("/"),
				"/"},
			{"short value",
				MustAbsolutePath("/aa"),
				"/aa"},
			{"long value",
				MustAbsolutePath("/a/bb/ccc"),
				"/a/bb/ccc"},
			{"dotted value",
				MustAbsolutePath("/.aa"),
				"/.aa"},
			{"dotted2 value",
				MustAbsolutePath("/..aa"),
				"/..aa"},
			{"dotted3 value",
				MustAbsolutePath("/...aa"),
				"/...aa"},
		} {
			Convey(tr.title, func() {
				v := fmt.Sprintf("%s", tr.p1)
				So(v, ShouldResemble, tr.str)
			})
		}
	})
}

func TestAbsolutePathDir(t *testing.T) {
	Convey("AbsolutePath.Dir suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    AbsolutePath
			pdir  AbsolutePath
		}{
			{"zero values",
				AbsolutePath{},
				AbsolutePath{}},
			{"root value",
				MustAbsolutePath("/"),
				AbsolutePath{}},
			{"short value",
				MustAbsolutePath("/aa"),
				AbsolutePath{}},
			{"long value",
				MustAbsolutePath("/a/bb/ccc"),
				MustAbsolutePath("/a/bb")},
		} {
			Convey(tr.title, func() {
				v := tr.p1.Dir()
				So(v, ShouldResemble, tr.pdir)
			})
		}
	})
}

func TestAbsolutePathLast(t *testing.T) {
	Convey("AbsolutePath.Last suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    AbsolutePath
			last  string
		}{
			{"zero values",
				AbsolutePath{},
				"/"},
			{"root value",
				MustAbsolutePath("/"),
				"/"},
			{"short value",
				MustAbsolutePath("/aa"),
				"aa"},
			{"long value",
				MustAbsolutePath("/a/bb/ccc"),
				"ccc"},
		} {
			Convey(tr.title, func() {
				v := tr.p1.Last()
				So(v, ShouldResemble, tr.last)
			})
		}
	})
}

func TestAbsolutePathJoins(t *testing.T) {
	Convey("AbsolutePath.Join suite:", t, func() {
		for _, tr := range []struct {
			title string
			p1    AbsolutePath
			p2    RelPath
			pj    AbsolutePath
		}{
			{"zero values",
				AbsolutePath{}, RelPath{},
				AbsolutePath{}},
			{"regular values",
				MustAbsolutePath("/root/"), MustRelPath("pth"),
				MustAbsolutePath("/root/pth")},
			{"root,short",
				MustAbsolutePath("/"), MustRelPath("pth"),
				MustAbsolutePath("/pth")},
			{"long,short",
				MustAbsolutePath("/r/oot"), MustRelPath("pth"),
				MustAbsolutePath("/r/oot/pth")},
			{"long,zero",
				MustAbsolutePath("/a/bb/ccc"), MustRelPath("."),
				MustAbsolutePath("/a/bb/ccc")},
			{"long,long",
				MustAbsolutePath("/a/bb/ccc"), MustRelPath("dd/e"),
				MustAbsolutePath("/a/bb/ccc/dd/e")},
			{"root,up",
				MustAbsolutePath("/"), MustRelPath(".."),
				MustAbsolutePath("/")},
			{"short,up",
				MustAbsolutePath("/root"), MustRelPath(".."),
				MustAbsolutePath("/")},
			{"long,up",
				MustAbsolutePath("/r/oot/pth"), MustRelPath(".."),
				MustAbsolutePath("/r/oot")},
		} {
			Convey(tr.title, func() {
				v := tr.p1.Join(tr.p2)
				So(v, ShouldResemble, tr.pj)
			})
		}
	})
}
