package impl

import (
	"github.com/bblfsh/go-driver/driver/golang"
	"gopkg.in/bblfsh/sdk.v2/sdk/driver"
)

func init() {
	driver.DefaultDriver = golang.NewDriver()
}
