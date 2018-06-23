package nqm

import (
	"github.com/crosserclaws/intest/common/conform"
	"github.com/crosserclaws/intest/common/utils"
	. "gopkg.in/check.v1"
)

type TestPingtaskSuite struct{}

var _ = Suite(&TestPingtaskSuite{})

// Tests validation of pingtasks
func (suite *TestPingtaskSuite) TestPingtaskModify(c *C) {
	testCase := &PingtaskModify{
		Name:    utils.PointerOfCloneString(" 台灣 "),
		Comment: utils.PointerOfCloneString(" 測試用 "),
	}

	conform.MustConform(testCase)

	c.Assert(testCase.Name, DeepEquals, utils.PointerOfCloneString("台灣"))
	c.Assert(testCase.Comment, DeepEquals, utils.PointerOfCloneString("測試用"))
}
