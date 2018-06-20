package file_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/go-mixins/loader/file"
)

func TestLoader_Changes(t *testing.T) {
	td, err := ioutil.TempFile("", "loader")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.Remove(td.Name())
	var loaded bool
	mockUnmarshaler := func(data []byte, dest interface{}) error {
		loaded = true
		return nil
	}
	loader := file.New(td.Name(), mockUnmarshaler)
	defer loader.Close()
	if err := loader.Load(nil); err != nil {
		t.Errorf("%+v", err)
	}
	if !loaded {
		t.Error("should have called mockUnmarshaler")
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		ioutil.WriteFile(td.Name(), []byte("change"), 0644)
	}()
	now := time.Now()
	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for change")
		break
	case <-loader.Changes():
		break
	}
	if dt := time.Since(now); dt < time.Duration(file.DebounceTimeout)*time.Millisecond {
		t.Errorf("debounce should not be %v", dt)
	}
}
