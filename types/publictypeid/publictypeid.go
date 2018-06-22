package publictypeid

import (
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/version/versiontype"
)

var ResolverDomain = domain.MustDecodeId(encodertype.BASE62, "T", "0", uint32(0), uint(0), versiontype.SEMANTIC)
