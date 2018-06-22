package domain_test

import (
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/init/schemeinit"
)

func TestMain(m *testing.M) {

	schemeinit.Init()

	os.Exit(m.Run())
}
