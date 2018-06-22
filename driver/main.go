package main

import (
	"github.com/bblfsh/typescript-driver/driver/normalizer"

	"gopkg.in/bblfsh/sdk.v1/sdk/driver"
)

func main() {
	driver.Run(normalizer.ToNode, normalizer.Transformers)
}
