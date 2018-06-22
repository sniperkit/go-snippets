package tests

import (
	"context"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/go-timeless-api"
	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/testutil"
)

func CheckMirror(packType api.PackType, mirror rio.MirrorFunc, pack rio.PackFunc, unpack rio.UnpackFunc, target api.WarehouseAddr, source api.WarehouseAddr) {
	testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
		// Initialization: make a pack to test against, put it in source warehouse.
		fixturePath := tmpDir.Join(fs.MustRelPath("fixture"))
		// Set up fixture.
		PlaceFixture(osfs.New(fixturePath), FixtureGamma)
		// Pack up into our warehouseaddr.
		wareID, err := pack(
			context.Background(),
			packType,
			fixturePath.String(),
			api.Filter_NoMutation,
			source,
			rio.Monitor{},
		)
		So(err, ShouldBeNil)

		// Okay, now mirror:
		Convey("mirror should succeed", func() {
			mirroredWareID, err := mirror(
				context.Background(),
				wareID,
				target,
				[]api.WarehouseAddr{source},
				rio.Monitor{},
			)
			So(err, ShouldBeNil)
			So(mirroredWareID, ShouldResemble, wareID)

			Convey("unpack from mirror should succeed", func() {
				unpackedWareID, err := unpack(
					context.Background(),
					wareID,
					"-",
					api.Filter_NoMutation,
					rio.Placement_None,
					[]api.WarehouseAddr{target},
					rio.Monitor{},
				)
				So(err, ShouldBeNil)
				So(unpackedWareID, ShouldResemble, wareID)
			})
		})
	})
}
