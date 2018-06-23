package basic

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math"
	"testing"
	"time"
)

// ioutil.ReadAll is slower than io.Copy

func TestReadAllvsCopy(t *testing.T) {
	bs := make([]byte, 1024*1024*512)
	log.Println(len(bs))

	startTimeForUtil := time.Now()
	log.Println("READALL Start ====> ", startTimeForUtil)
	ioutil.ReadAll(bytes.NewReader(bs))
	log.Println("READALL Process ====> ", time.Now().Sub(startTimeForUtil))

	startTimeForIO := time.Now()
	log.Println("COPY Start ====> ", startTimeForIO)
	io.Copy(bytes.NewBuffer(nil), bytes.NewReader(bs))
	log.Println("COPY Process ====> ", time.Now().Sub(startTimeForIO))

}

func readAll(n int) {
	bs := make([]byte, 1024*1024*n)
	ioutil.ReadAll(bytes.NewReader(bs))
}

func copy(n int) {
	bs := make([]byte, 1024*1024*n)
	io.Copy(bytes.NewBuffer(nil), bytes.NewReader(bs))
}

func BenchmarkCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		size := math.Min(1, float64(1024*i/b.N))
		// copy(rand.Intn(1024))
		copy(int(size))
	}
}

func BenchmarkReadAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		size := math.Min(1, float64(1024*i/b.N))
		// copy(rand.Intn(1024))
		readAll(int(size))
	}
}
