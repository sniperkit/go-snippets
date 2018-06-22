package tests

import (
	"bytes"
	"time"

	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fsOp"
)

type FixtureFile struct {
	Metadata fs.Metadata
	Body     []byte
}

// Because golang's time.Time zero value causes Nonsense to occur.
var defaultTime = time.Date(1990, 1, 14, 12, 30, 0, 0, time.UTC)

var FixtureAlpha = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
}

var FixtureAlphaDiffContent = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("qwe")},
}

var FixtureAlphaDiffTime = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: time.Date(2004, 10, 14, 4, 3, 2, 0, time.UTC), Size: 3}, []byte("zyx")},
}

var FixtureAlphaDiffPerm = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0600, Mtime: defaultTime, Size: 3}, []byte("zyx")},
}

var FixtureAlphaDiffPerm2 = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	// Perms+=0020  makes a good test that umasks aren't screwing with us.
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0664, Mtime: defaultTime, Size: 3}, []byte("zyx")},
}

var FixtureAlphaDiffPerm3 = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	// Add in all sticky (setuid,setgid,sticky) bits.  These sometimes get stripped in weird places;
	//  keep 0644 underneath, so if they do, we fail the collision test on pack (in addition to obviously failing the roundtrip test).
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 07644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
}

var FixtureAlphaDiffUidGid = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3, Uid: 444, Gid: 444}, []byte("zyx")},
}

var FixtureEmpty = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
}

var FixtureMultifile = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
	{fs.Metadata{Name: fs.MustRelPath("./b"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("qwe")},
}

var FixtureDepth1 = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
	{fs.Metadata{Name: fs.MustRelPath("./d"), Type: fs.Type_Dir, Perms: 0750, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./d/c"), Type: fs.Type_File, Perms: 0664, Mtime: defaultTime, Size: 4}, []byte("asdf")},
}

var FixtureDepth3 = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
	{fs.Metadata{Name: fs.MustRelPath("./d"), Type: fs.Type_Dir, Perms: 0750, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./d/d2"), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./d/d2/c"), Type: fs.Type_File, Perms: 0664, Mtime: defaultTime, Size: 4}, []byte("asdf")},
}

var FixtureSymlinks = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./a"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
	// Perms on the link are set to 777, not because that works, but because *that's what you get* on a linux system.
	{fs.Metadata{Name: fs.MustRelPath("./ln"), Type: fs.Type_Symlink, Perms: 0777, Mtime: defaultTime, Linkname: "./a"}, nil},
}

// deep and varied structures.  files and dirs.
// subtle: a dir with a sibling that's a suffix of its name (can trip up dir/child adjacency sorting).
// subtle: a file with a sibling that's a suffix of its name (other half of the test, to make sure the prefix doesn't create an incorrect tree node).
var FixtureGamma = []FixtureFile{
	{fs.Metadata{Name: fs.MustRelPath("."), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./etc"), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./etc/init.d/"), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./etc/init.d/service-p"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 2}, []byte("p!")},
	{fs.Metadata{Name: fs.MustRelPath("./etc/init.d/service-q"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 2}, []byte("q!")},
	{fs.Metadata{Name: fs.MustRelPath("./etc/init/"), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./etc/init/zed"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 4}, []byte("grue")},
	{fs.Metadata{Name: fs.MustRelPath("./etc/trick"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("sib")},
	{fs.Metadata{Name: fs.MustRelPath("./etc/tricky"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("sob")},
	{fs.Metadata{Name: fs.MustRelPath("./var"), Type: fs.Type_Dir, Perms: 0755, Mtime: defaultTime}, nil},
	{fs.Metadata{Name: fs.MustRelPath("./var/fun"), Type: fs.Type_File, Perms: 0644, Mtime: defaultTime, Size: 3}, []byte("zyx")},
}

var AllFixtures = []struct {
	Name  string
	Files []FixtureFile
}{
	{"Alpha", FixtureAlpha},
	{"AlphaDiffContent", FixtureAlphaDiffContent},
	{"AlphaDiffTime", FixtureAlphaDiffTime},
	{"AlphaDiffPerm", FixtureAlphaDiffPerm},
	{"AlphaDiffPerm2", FixtureAlphaDiffPerm2},
	{"AlphaDiffPerm3", FixtureAlphaDiffPerm3},
	{"AlphaDiffUidGid", FixtureAlphaDiffUidGid},
	{"Empty", FixtureEmpty},
	{"Multifile", FixtureMultifile},
	{"Depth1", FixtureDepth1},
	{"Depth3", FixtureDepth3},
	{"Symlinks", FixtureSymlinks},
	{"Gamma", FixtureGamma},
}

/*
	Create files described by the fixtures on the filesystem given.
	Any errors will be panicked, since this is meant to be used in test setup.
*/
func PlaceFixture(afs fs.FS, fixture []FixtureFile) {
	// Range over fixture slice, making files.
	for _, ff := range fixture {
		if err := fsOp.PlaceFile(afs, ff.Metadata, bytes.NewBuffer(ff.Body), false); err != nil {
			panic(err)
		}
	}
	// Range again: ... in reverse order, to re-do time enforcement, covering our own tracks.
	for i := len(fixture) - 1; i >= 0; i-- {
		ff := fixture[i]
		if ff.Metadata.Type == fs.Type_Dir {
			if err := afs.SetTimesNano(ff.Metadata.Name, ff.Metadata.Mtime, fs.DefaultAtime); err != nil {
				panic(err)
			}
		}
	}
}
