package idsinit

import (
	"github.com/distributed-vision/go-resources/ids/signaturedomain"
	"github.com/distributed-vision/go-resources/init/schemeinit"
)

func Init() {
	schemeinit.Init()
	signaturedomain.Init()
}
