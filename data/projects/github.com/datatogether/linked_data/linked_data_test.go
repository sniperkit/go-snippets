package ld_test

import (
	"github.com/datatogether/linked_data"
	"github.com/datatogether/linked_data/dcat"
	"github.com/datatogether/linked_data/pod"
)

// assert that each interface is properly implemented
var _ ld.Catalog = &pod.Catalog{}
var _ ld.Catalog = &dcat.Catalog{}
