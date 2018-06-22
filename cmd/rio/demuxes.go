package main

import (
	. "github.com/warpfork/go-errcat"

	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/transmat/git"
	"go.polydawn.net/rio/transmat/tar"
)

func demuxPackTool(packType string) (rio.PackFunc, error) {
	switch packType {
	case "tar":
		return tartrans.Pack, nil
	default:
		return nil, Errorf(rio.ErrUsage, "unsupported packtype %q", packType)
	}
}

func demuxUnpackTool(packType string) (rio.UnpackFunc, error) {
	switch packType {
	case "tar":
		return tartrans.Unpack, nil
	case "git":
		return git.Unpack, nil
	default:
		return nil, Errorf(rio.ErrUsage, "unsupported packtype %q", packType)
	}
}

func demuxScanTool(packType string) (rio.ScanFunc, error) {
	switch packType {
	case "tar":
		return tartrans.Scan, nil
	default:
		return nil, Errorf(rio.ErrUsage, "unsupported packtype %q", packType)
	}
}

func demuxMirrorTool(packType string) (rio.MirrorFunc, error) {
	switch packType {
	case "tar":
		return tartrans.Mirror, nil
	default:
		return nil, Errorf(rio.ErrUsage, "unsupported packtype %q", packType)
	}
}
