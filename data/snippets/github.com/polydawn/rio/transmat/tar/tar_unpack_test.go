package tartrans

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/go-timeless-api/util"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/testutil"
	"go.polydawn.net/rio/transmat/mixins/tests"
)

func TestTarUnpack(t *testing.T) {
	Convey("Spec compliance: Tar unpack", t,
		testutil.Requires(testutil.RequiresCanManageOwnership, func() {
			Convey("Using kvfs warehouse, in content-addressable mode:", func() {
				testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
					osfs.New(tmpDir).Mkdir(fs.MustRelPath("bounce"), 0755)
					tests.CheckRoundTrip(PackType, Pack, Unpack, api.WarehouseAddr(fmt.Sprintf("ca+file://%s/bounce", tmpDir)))
					// Following tests could be done in all modes, but isn't about warehouses, so would be redundant to do so.
					tests.CheckCachePopulation(PackType, Pack, Unpack, api.WarehouseAddr(fmt.Sprintf("ca+file://%s/bounce", tmpDir)))
				})
			})
			Convey("Using kvfs warehouse, in *non*-content-addressable mode:", func() {
				testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
					tests.CheckRoundTrip(PackType, Pack, Unpack, api.WarehouseAddr(fmt.Sprintf("file://%s/bounce", tmpDir)))
				})
			})
		}),
	)
}

/*
	Tests against pre-generated, known fixtures of tar binary blobs.

	These tests allow us to cover compat with other tar impls, compression, etc.
*/
func TestTarFixtureUnpack(t *testing.T) {
	Convey("Tar transmat: unpacking of fixtures", t,
		testutil.Requires(testutil.RequiresCanManageOwnership, func() {
			testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
				Convey("Unpack a fixture from gnu tar which includes a base dir", func() {
					wareID := api.WareID{"tar", "5y6NvK6GBPQ6CcuNyJyWtSrMAJQ4LVrAcZSoCRAzMSk5o53pkTYiieWyRivfvhZwhZ"}
					gotWareID, err := Unpack(
						context.Background(),
						wareID,
						tmpDir.String(),
						api.Filter_NoMutation,
						rio.Placement_Direct,
						[]api.WarehouseAddr{"file://./fixtures/tar_withBase.tgz"},
						rio.Monitor{},
					)
					So(err, ShouldBeNil)
					So(gotWareID, ShouldResemble, wareID)

					fmeta, reader, err := fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("ab"))
					So(err, ShouldBeNil)
					So(fmeta.Name, ShouldResemble, fs.MustRelPath("ab"))
					So(fmeta.Type, ShouldResemble, fs.Type_File)
					So(fmeta.Uid, ShouldEqual, 7000)
					So(fmeta.Gid, ShouldEqual, 7000)
					So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 53, 35, 0, time.UTC))
					body, err := ioutil.ReadAll(reader)
					So(string(body), ShouldResemble, "")

					fmeta, reader, err = fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("bc"))
					So(err, ShouldBeNil)
					So(fmeta.Name, ShouldResemble, fs.MustRelPath("bc"))
					So(fmeta.Type, ShouldResemble, fs.Type_Dir)
					So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 53, 35, 0, time.UTC))
					So(reader, ShouldBeNil)

					fmeta, reader, err = fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("."))
					So(err, ShouldBeNil)
					So(fmeta.Name, ShouldResemble, fs.MustRelPath("."))
					So(fmeta.Type, ShouldResemble, fs.Type_Dir)
					So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 53, 35, 0, time.UTC))
					So(reader, ShouldBeNil)
				})
				Convey("Unpack a fixture from gnu tar which lacks a base dir", func() {
					wareID := api.WareID{"tar", "2RLHdc3am6tMCFy56vfcHm5kWLoAtYBfiaQcq17vDm1tEzQn9CC6tcF2yzpAJvehPC"}
					gotWareID, err := Unpack(
						context.Background(),
						wareID,
						tmpDir.String(),
						api.Filter_NoMutation,
						rio.Placement_Direct,
						[]api.WarehouseAddr{"file://./fixtures/tar_sansBase.tgz"},
						rio.Monitor{},
					)
					So(err, ShouldBeNil)
					So(gotWareID, ShouldResemble, wareID)

					fmeta, reader, err := fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("ab"))
					So(err, ShouldBeNil)
					So(fmeta.Name, ShouldResemble, fs.MustRelPath("ab"))
					So(fmeta.Type, ShouldResemble, fs.Type_File)
					So(fmeta.Uid, ShouldEqual, 7000)
					So(fmeta.Gid, ShouldEqual, 7000)
					So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 11, 23, 0, time.UTC))
					body, err := ioutil.ReadAll(reader)
					So(string(body), ShouldResemble, "")

					fmeta, reader, err = fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("."))
					So(err, ShouldBeNil)
					So(fmeta.Name, ShouldResemble, fs.MustRelPath("."))
					So(fmeta.Type, ShouldResemble, fs.Type_Dir)
					So(fmeta.Mtime.UTC(), ShouldResemble, apiutil.DefaultMtime)
					So(reader, ShouldBeNil)
				})
			})
		}),
	)
}
