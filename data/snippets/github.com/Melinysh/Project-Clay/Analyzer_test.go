package main

import (
	"reflect"
	"testing"
)

func TestdetectTweetDifferences(t *testing.T) {
	//	newTweet =
	// Check new tweet string JSON, marshal into a map then copy it, modify it and see if the changes are detected.
}

func TestDiff(t *testing.T) {
	set1 := []int64{int64(1234567890), int64(9876577654), int64(765734299), int64(786837575), int64(287364758), int64(1043746), int64(1234663110), int64(34987539)}
	set2 := []int64{int64(1234567890), int64(9876577654), int64(765734299), int64(786837575), int64(287364758), int64(46565656)}

	d1 := []int64{int64(1043746), int64(1234663110), int64(34987539)} //new in set1 -set2
	d2 := []int64{int64(46565656)}
	diff1 := Diff(set1, set2)
	diff2 := Diff(set2, set1)

	if len(diff1) != len(d1) {
		t.Fatal("First diff length is", len(diff1), " which doesn't equal predicted diff of length", len(d1))
	} else if len(diff2) != len(d2) {
		t.Fatal("Second diff length is", len(diff2), " which doesn't equal predicted diff of length", len(d2))
	}

	if reflect.DeepEqual(diff1, d1) == false {
		t.Fatal("The first diff doesn't match the expected diff.", diff1, d1)
	} else if reflect.DeepEqual(diff2, d2) == false {
		t.Fatal("The second diff doesn't match the expected diff.", diff2, d2)
	}
}

func TestConvertInts(t *testing.T) {
	ints := []int64{int64(6755), int64(87687), int64(8565543), int64(8554433266324)}
	interfs := ConvertInts(ints)
	if reflect.TypeOf(interfs) != reflect.TypeOf([]interface{}{}) {
		t.Fatal("Interfs not of type []interface{}, and is actually", reflect.TypeOf(interfs).Name())
	}
}

func TestConvertInterfaces(t *testing.T) {
	interfs := []interface{}{int64(6755), int64(87687), int64(8565543), int64(8554433266324)}
	ints := ConvertInterfaces(interfs)
	if reflect.TypeOf(ints) != reflect.TypeOf([]int64{}) {
		t.Fatal("Ints is not of type []int64 and is instead", reflect.TypeOf(ints))
	}
}
