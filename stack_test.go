// Copyright 2016 Chao Wang <hit9@icloud.com>

package diskstack

import (
	"bytes"
	"os"
	"runtime"
	"testing"
)

// Must asserts the given value is True for testing.
func Must(t *testing.T, v bool) {
	if !v {
		_, fileName, line, _ := runtime.Caller(1)
		t.Errorf("\n unexcepted: %s:%d", fileName, line)
	}
}

func TestOpenEmpty(t *testing.T) {
	fileName := "stack.db"
	s, err := Open(fileName, nil)
	// Must open without errors
	Must(t, err == nil)
	Must(t, s != nil)
	defer os.Remove(fileName)
	info, err := os.Stat(fileName)
	// Must size be 4+8
	Must(t, err == nil && info.Size() == headSize)
}

func TestReOpen(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	// Put one item.
	data := []byte{1, 2, 3}
	s.Put(data)
	// Close stack.
	s.Close()
	// Reopen.
	s, _ = Open(fileName, nil)
	// Must offset be correct.
	Must(t, s.offset == int64(len(data))+4+8+4)
}

func TestTopEmpty(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	data, err := s.Top()
	// Must be nil,nil
	Must(t, data == nil && err == nil)
	data, err = s.Pop()
	// Must be nil,nil
	Must(t, data == nil && err == nil)
}

func TestOperations(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	data1 := []byte{1, 2, 3, 4}
	data2 := []byte{5, 6, 7, 8}
	data3 := []byte{9, 10, 11, 12}
	// Must put ok.
	Must(t, s.Put(data1) == nil)
	// Top should be data1.
	data, err := s.Top()
	Must(t, err == nil && bytes.Compare(data, data1) == 0)
	// Must put ok.
	Must(t, s.Put(data2) == nil)
	Must(t, s.Put(data3) == nil)
	// Pop should be data3
	data, err = s.Pop()
	Must(t, err == nil && bytes.Compare(data, data3) == 0)
	// Pop again should be data2
	data, err = s.Pop()
	Must(t, err == nil && bytes.Compare(data, data2) == 0)
	// Pop again should be data1
	data, err = s.Pop()
	Must(t, err == nil && bytes.Compare(data, data1) == 0)
	// Pop again should be nil
	data, err = s.Pop()
	Must(t, err == nil && bytes.Compare(data, nil) == 0)
}

func TestOperationsBetweenOpens(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	data1 := []byte{1, 2, 3, 4}
	data2 := []byte{5, 6, 7, 8}
	data3 := []byte{9, 10, 11, 12}
	// Must put ok.
	Must(t, s.Put(data1) == nil)
	Must(t, s.Put(data2) == nil)
	Must(t, s.Put(data3) == nil)
	Must(t, s.Len() == 3)
	// Pop one. offset should persist between opens.
	s.Pop()
	// Close.
	s.Close()
	// Reopen.
	s, _ = Open(fileName, nil)
	// Must length be persist
	Must(t, s.Len() == 2)
	// Must offset be correct.
	Must(t, s.offset == int64(headSize+s.Len()*(4+4)))
	// Pops should be correct.
	data, err := s.Pop()
	Must(t, err == nil && bytes.Compare(data, data2) == 0 && s.Len() == 1)
	data, err = s.Pop()
	Must(t, err == nil && bytes.Compare(data, data1) == 0 && s.Len() == 0)
	data, err = s.Pop()
	Must(t, err == nil && bytes.Compare(data, nil) == 0 && s.Len() == 0)
}

func TestCompact(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, &Options{FragmentsThreshold: 48})
	defer os.Remove(fileName)
	// Put some items.
	data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	for i := 0; i < 6; i++ {
		s.Put(data)
	} // (4+8)+(8+4)*6=84
	// Must file size is excepted
	info, _ := os.Stat(fileName)
	except := int64((4 + 8) + (4+8)*6)
	Must(t, info.Size() == except)
	// Pop some items.
	for i := 0; i < 4; i++ {
		s.Pop()
		if i < 3 {
			// Must no compaction for first 3 times
			info, _ = os.Stat(fileName)
			Must(t, info.Size() == except)
		}
	} // Poped (8+4) x 4 = 48 bytes
	// Must be compacted now.
	info, _ = os.Stat(fileName)
	Must(t, info.Size() == (4+8)+2*(8+4))
	Must(t, s.frags == 0)
	// Put/Pop should work fine after the compaction.
	data1 := []byte{5, 6, 7, 8}
	Must(t, s.Put(data1) == nil)
	v, err := s.Pop()
	Must(t, err == nil && bytes.Compare(v, data1) == 0)
	v, err = s.Pop()
	Must(t, err == nil && bytes.Compare(v, data) == 0)
}

func TestClear(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, &Options{FragmentsThreshold: 64})
	defer os.Remove(fileName)
	// Put some items.
	data := []byte{1, 2, 3, 4}
	for i := 0; i < 1024; i++ {
		s.Put(data)
	}
	// Clear
	Must(t, s.Clear() == nil)
	// Must file be cleared.
	info, _ := os.Stat(fileName)
	Must(t, info.Size() == headSize)
	// Must offset be 0
	Must(t, s.offset == headSize)
	// Pop should be nil.
	v, _ := s.Pop()
	Must(t, v == nil)
	// Length must be 0
	Must(t, s.Len() == 0)
}

func TestLen(t *testing.T) {
	fileName := "stack.db"
	s, _ := Open(fileName, &Options{FragmentsThreshold: 64})
	defer os.Remove(fileName)
	// Put some items.
	n := 1024
	data := []byte{1, 2, 3, 4}
	for i := 0; i < n; i++ {
		s.Put(data)
	}
	Must(t, s.Len() == 1024)
}

func BenchmarkPut(b *testing.B) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	data := []byte("12345678910")
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Put(data)
	}
}

func BenchmarkPutLargeItem(b *testing.B) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 255)
	} // data length is now 1kb
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Put(data)
	}
}

func BenchmarkPop(b *testing.B) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	data := []byte("12345678910")
	for i := 0; i < b.N; i++ {
		s.Put(data)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Pop()
	}
}

func BenchmarkPopLargeItem(b *testing.B) {
	fileName := "stack.db"
	s, _ := Open(fileName, nil)
	defer os.Remove(fileName)
	var data []byte
	for i := 0; i < 1024; i++ {
		data = append(data, 255)
	} // data length is now 1kb
	for i := 0; i < b.N; i++ {
		s.Put(data)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s.Pop()
	}
}
