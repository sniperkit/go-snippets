package osfs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/tests"
	"go.polydawn.net/rio/testutil"
)

func TestAll(t *testing.T) {
	Convey("osfs spec compliance tests", t, func() {
		testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
			tfs := New(tmpDir)
			boxPath := fs.MustRelPath("sandbox")
			tfs.Mkdir(boxPath, 0755)
			afs := New(tmpDir.Join(boxPath))

			tests.CheckBaseLstat(afs)
			tests.CheckMkdirLstatRoundtrip(afs)
			tests.CheckDeepMkdirError(afs)
			tests.CheckMklinkLstatRoundtrip(afs)
			tests.CheckSymlinks(afs)
			tests.CheckPerniciousSymlinks(afs)
			tests.CheckOpsTraversingSymlinks(afs)
		})
	})
}
