package main

import (
	"testing"
)

func TestformatToInt64(t *testing.T) {
	formatter := newFormatter()
	str := "23437645125"
	integer := int(23437645125)
	float := float64(23437645125)
	form1 := int64(23437645125)
	if formatter.formatToInt64(str) != form1 {
		t.Fatal("Formmated string", str, "is not equal to the int64", form1)
	} else if formatter.formatToInt64(integer) != form1 {
		t.Fatal("Integer", integer, " is not equal to the int64", form1)
	} else if formatter.formatToInt64(float) != form1 {
		t.Fatal("Float", float, " is not equal to the int64", form1)
	}
}
