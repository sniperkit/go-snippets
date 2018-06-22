// Copyright 2016 Chao Wang <hit9@icloud.com>

/*

Package diskstack implements the disk-based stack.

Abstract

Stack is stored on a file, contains head and body:

	+------------------+------+
	| [offset 8 bytes] | head |
	| [length 4 bytes] |  12  |
	+------------------+------+
	| [data   X bytes] |      |
	| [size   4 bytes] | body | -> offset TOP
	| [data   X bytes] |  X   |
	| [size   4 bytes] |      |
	| ...              |      |

Put steps: (IO: 2)

	1. Write data and data size on offset.
	2. Adjust offset (and frags, length etc.).
	3. Write head (including offset, length).

Pop steps: (IO: worst 4, best 3)

	1. Read data size at offset-4.
	2. Read data at offset-4-X.
	3. Adjust offset (and frags, length etc.).
	4. Write Head (including offset, length).
	5. If the fragments size is larger enough, truncate
	   the file.

*/
package diskstack

import (
	"encoding/binary"
	"os"
	"sync"
)

// Size units
const (
	KB int64 = 1024
	MB int64 = 1024 * KB
	GB int64 = 1024 * MB
)

// Head size
const (
	offsetSize = 8
	lengthSize = 4
	headSize   = offsetSize + lengthSize
)

// Default options.
const (
	DefaultFragmentsThreshold int64 = 256 * MB
)

// Options is the options to open Stack.
type Options struct {
	// FragmentsThreshold is the fragments size threshold to trigger the
	// file compaction.
	FragmentsThreshold int64
}

// Stack is the disk-based stack abstraction.
type Stack struct {
	file   *os.File     // os file handle
	offset int64        // top offset (real offset is 8+4+offset)
	size   int64        // file size
	frags  int64        // fragments size
	length int          // length of stack
	lock   sync.RWMutex // protects offset,frags,filesize,size
	opts   *Options
}

// Open opens or creates a Stack for given path, will create if not exist.
func Open(path string, opts *Options) (s *Stack, err error) {
	// Open or create the file.
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		return
	}
	// Create Stack.
	if opts == nil {
		opts = &Options{DefaultFragmentsThreshold}
	}
	s = &Stack{opts: opts, file: file}
	// Get file size.
	info, err := file.Stat()
	if err != nil {
		return
	}
	fileSize := info.Size()
	if fileSize < headSize {
		if err = s.file.Truncate(0); err != nil {
			// Force truncate the file to be empty.
			return
		}
		s.offset = headSize
		s.size = headSize
		s.frags = 0
		s.length = 0
		if err = s.writeHead(); err != nil {
			return
		}
		return
	}
	// Read offset.
	b := make([]byte, offsetSize)
	if _, err = file.ReadAt(b, 0); err != nil {
		return
	}
	s.offset = int64(binary.BigEndian.Uint64(b))
	// Read length.
	b = make([]byte, 4)
	if _, err = file.ReadAt(b, offsetSize); err != nil {
		return
	}
	s.length = int(binary.BigEndian.Uint32(b))
	// File Size
	s.size = fileSize
	// Frags
	if err = s.truncate(); err != nil { // Remove the fragements
		return
	}
	s.frags = 0
	return s, nil
}

// Put an item onto the Stack.
func (s *Stack) Put(data []byte) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// Make sure offset<=filesize
	// offset is read from db file, may be not safe.
	offset := s.offset
	if offset > s.size {
		offset = s.size
	}
	// Write data and data size at offset.
	buf := make([]byte, len(data)+4)
	copy(buf, data)                                                // data
	binary.BigEndian.PutUint32(buf[len(data):], uint32(len(data))) // size
	if _, err = s.file.WriteAt(buf, offset); err != nil {
		return
	}
	// Adjust offset, frags, size and length.
	s.offset += int64(len(buf))
	if s.frags > int64(len(buf)) {
		// Written buffer is smaller than fragements, the file didn't grow its
		// size.
		s.frags -= int64(len(buf))
	} else {
		// Written buffer is larger than fragements, the file grew its size.
		s.frags = 0
		s.size += int64(len(buf)) - s.frags
	}
	s.length++
	return s.writeHead()
}

// top returns the top item.
func (s *Stack) top() (data []byte, err error) {
	if s.offset < headSize+4 {
		return nil, nil // Empty stack.
	}
	b := make([]byte, 4)
	if _, err = s.file.ReadAt(b, s.offset-4); err != nil { // size
		return
	}
	size := binary.BigEndian.Uint32(b)
	data = make([]byte, size)
	if _, err = s.file.ReadAt(data, s.offset-4-int64(size)); err != nil { // data
		return
	}
	return
}

// Top returns the top item on the Stack, returns nil on empty.
func (s *Stack) Top() (data []byte, err error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.top()
}

// Len returns the Stack length.
func (s *Stack) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.length
}

// Size returns the Stack file size.
func (s *Stack) Size() int64 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.size
}

// Pop an item from the Stack, returns nil on empty.
func (s *Stack) Pop() (data []byte, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if data, err = s.top(); err != nil {
		return
	}
	if data == nil {
		return // Do nothing on stack empty.
	}
	s.offset -= int64(len(data)) + 4
	s.length--
	s.frags += int64(len(data)) + 4 // fragments should be larger
	if err = s.writeHead(); err != nil {
		return
	}
	if err = s.compact(); err != nil {
		return
	}
	return
}

// Clear the Stack.
func (s *Stack) Clear() (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.offset = headSize
	s.length = 0
	s.frags = s.size - headSize
	s.size = headSize
	if err = s.writeHead(); err != nil {
		return
	}
	return s.truncate()
}

// Close the Stack.
func (s *Stack) Close() (err error) {
	if err = s.file.Close(); err != nil {
		return
	}
	return
}

// compact truncates the file if the fragments is greater than the threshold.
func (s *Stack) compact() (err error) {
	if s.frags >= s.opts.FragmentsThreshold {
		err = s.truncate()
		s.frags = 0
		return
	}
	return nil
}

// truncate the file to size the offset.
func (s *Stack) truncate() (err error) {
	// Make sure offset<=filesize
	// offset is read from db file, may be not safe.
	size := s.offset
	if s.size < s.offset {
		size = s.size
	}
	err = s.file.Truncate(size)
	// Adjust the file size, operation file.Truncate must make the file size be
	// size.
	s.size = size
	return
}

// writeHead writes the head.
func (s *Stack) writeHead() (err error) {
	// [offset 8bytes | length 4bytes]
	b := make([]byte, headSize)
	binary.BigEndian.PutUint64(b, uint64(s.offset))
	binary.BigEndian.PutUint32(b[offsetSize:], uint32(s.length))
	if _, err = s.file.WriteAt(b, 0); err != nil {
		return err
	}
	return nil
}
