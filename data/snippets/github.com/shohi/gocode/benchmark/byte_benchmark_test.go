package benchmark

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"testing"
	"unsafe"
)

func convertUint64WithString(val uint64) []byte {
	// return []byte(strconv.Itoa(val))
	return []byte(fmt.Sprintf("%v", val))
}

func TestConvertString(t *testing.T) {
	var bb uint64
	bb = 100000000000000
	log.Println(convertUint64WithString(bb))
	log.Println(convertUint64WithBinary(bb))
	log.Println(convertUint64WithGob(bb))
	log.Println(convertUint64WithWriter(bb))
}

func TestParseString(t *testing.T) {
	bb := 100000000000000
	log.Println(strconv.ParseInt(fmt.Sprintf("%v", bb), 10, 64))
	var bs []byte
	log.Println(strconv.ParseInt(string(bs), 10, 64))
}

func convertUint64WithBinary(val uint64) []byte {
	bs := make([]byte, unsafe.Sizeof(val))
	binary.LittleEndian.PutUint64(bs, val)

	return bs
}

func convertUint64WithWriter(val uint64) []byte {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, val)
	return buf.Bytes()
}

func convertUint64WithGob(val uint64) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func BenchmarkByteConvert(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func(uint64) []byte
	}{
		{"ConvertWithString", convertUint64WithString},
		{"ConvertWithGob", convertUint64WithGob},
		{"ConvertWithBinary", convertUint64WithBinary},
		{"ConvertWithWriter", convertUint64WithWriter},
	}

	var val uint64
	var r *rand.Rand
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			r = rand.New(rand.NewSource(99))
			for k := 0; k < b.N; k++ {
				val = rand.Uint64()
				bm.fn(val)
			}
		})
	}

}
