package goxa

import (
	"testing"
)

func Test_001_Connect_Disconnect(t *testing.T) {

	xa, err := Connect("tcp", "www.google.com:80")
	if err != nil {
		t.Errorf(err.Error())
	}

	err = Disconnect(xa)
	if err != nil {
		t.Errorf(err.Error())
	}

	return
}
