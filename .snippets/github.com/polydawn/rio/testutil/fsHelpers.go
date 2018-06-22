package testutil

import (
	"github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/rio/fs"
)

func ShouldStat(afs fs.FS, path fs.RelPath) fs.Metadata {
	stat, err := afs.LStat(path)
	convey.So(err, convey.ShouldBeNil)
	stat.Mtime = stat.Mtime.UTC()
	return *stat
}
