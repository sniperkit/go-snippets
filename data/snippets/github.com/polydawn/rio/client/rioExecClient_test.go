package rioexecclient_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/client"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/testutil"
	"go.polydawn.net/rio/transmat/mixins/tests"
)

// This test is moderately terrifying because it does indeed and really exec.
// This means the rio binary must have already been built.
// We set the path to the project's build dir; any commands on your regular host PATH should not interfere.

func Test(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Setenv("PATH", filepath.Join(cwd, "../bin/"))
	if err != nil {
		panic(err)
	}

	Convey("Spec compliance: exec-RPC Tar pack", t,
		testutil.Requires(testutil.RequiresCanManageOwnership, func() {
			tests.CheckPackProducesConsistentHash("tar", rioexecclient.PackFunc)
			tests.CheckPackHashVariesOnVariations("tar", rioexecclient.PackFunc)
		}),
	)

	Convey("exec client tests", t, func() {
		Convey("unpacking tar fixtures (happy path)",
			testutil.Requires(testutil.RequiresCanManageOwnership, func() {
				testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
					_, err := rioexecclient.UnpackFunc(
						context.Background(),
						api.WareID{"tar", "5y6NvK6GBPQ6CcuNyJyWtSrMAJQ4LVrAcZSoCRAzMSk5o53pkTYiieWyRivfvhZwhZ"},
						tmpDir.String(),
						api.Filter_NoMutation,
						rio.Placement_Direct,
						[]api.WarehouseAddr{"file://../transmat/tar/fixtures/tar_withBase.tgz"},
						rio.Monitor{},
					)
					So(err, ShouldBeNil)
				})
			}),
		)
		Convey("unpacking with a lack of warehouses (error path)",
			func() {
				testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
					_, err := rioexecclient.UnpackFunc(
						context.Background(),
						api.WareID{"tar", "5y6NvK6GBPQ6CcuNyJyWtSrMAJQ4LVrAcZSoCRAzMSk5o53pkTYiieWyRivfvhZwhZ"},
						tmpDir.String(),
						api.Filter_NoMutation,
						rio.Placement_Direct,
						nil,
						rio.Monitor{},
					)
					So(err, ShouldResemble, &rio.Error{rio.ErrWarehouseUnavailable, "no warehouses were available!", nil})
				})
			},
		)
	})
}
