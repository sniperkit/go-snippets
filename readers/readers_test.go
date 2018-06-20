package readers_test

import (
	"testing"

	"github.com/sahilm/go-tour/readers"
)

func TestMyReaderReturnsAStreamOfAs(t *testing.T) {
	b := make([]byte, 1024, 2048)
	i, o := 0, 0
	for ; i < 1<<20 && o < 1<<20; i++ { // test 1mb
		n, err := readers.MyReader{}.Read(b)
		if n != 1024 {
			t.Errorf("invalid read length: %v, want %v", n, 1024)
		}
		for i, v := range b[:n] {
			if v != 'A' {
				t.Errorf("got byte %x at offset %v, want 'A'\n", v, o+i)
			}
		}
		o += n
		if err != nil {
			t.Errorf("read error: %v\n", err)
			return
		}
	}
	if o == 0 {
		t.Errorf("read zero bytes after %d Read calls\n", i)
	}
}
