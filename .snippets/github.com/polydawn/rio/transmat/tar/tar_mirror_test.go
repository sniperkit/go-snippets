package tartrans

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/testutil"
	"go.polydawn.net/rio/transmat/mixins/tests"
)

func TestTarMirror(t *testing.T) {
	Convey("Spec compliance: Tar mirror", t,
		testutil.Requires(testutil.RequiresCanManageOwnership, func() {
			Convey("Populating kvfs warehouse, in content-addressable mode, from kvfs warehouse, in content-addressable mode:", func() {
				testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
					osfs.New(tmpDir).Mkdir(fs.MustRelPath("src"), 0755)
					osfs.New(tmpDir).Mkdir(fs.MustRelPath("dst"), 0755)
					srcAddr := api.WarehouseAddr(fmt.Sprintf("ca+file://%s/src", tmpDir))
					dstAddr := api.WarehouseAddr(fmt.Sprintf("ca+file://%s/dst", tmpDir))

					tests.CheckMirror(PackType, Mirror, Pack, Unpack, dstAddr, srcAddr)
				})
			})
		}),
	)
}
