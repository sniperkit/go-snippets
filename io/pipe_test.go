package io

import "testing"

func checkWrite(t *testing.T, w Writer, data []byte, c chan int) {
	n, err := w.Write(data)
	if err != nil {
		t.Errorf("write error: %s", err)
	}
	if len(data) != n {
		t.Errorf("short write: %d != %d", n, len(data))
	}
	c <- 0
}

func TestPipe1(t *testing.T) {
	c := make(chan int)
	r, w := Pipe()
	var buf = make([]byte, 64)
	go checkWrite(t, w, []byte("Hello, world"), c)
	n, err := r.Read(buf)
	if err != nil {
		t.Errorf("read error: %s", err)
	} else if n != 12 || string(buf[:12]) != "Hello, world" {
		t.Errorf("bad read: got %q", buf[:n])
	}
	<-c
	r.Close()
	w.Close()
}

func reader(t *testing.T, r Reader, c chan int) {
	var buf = make([]byte, 64)
	for {
		n, err := r.Read(buf)
		if err == EOF {
			c <- 0
			break
		}
		if err != nil {
			t.Errorf("read error: %s", err)
		}
		c <- n
	}
}

func TestPipe2(t *testing.T) {
	var c = make(chan int)
	r, w := Pipe()
	go reader(t, r, c)
	var buf = make([]byte, 64)
	for i := 0; i < 5; i++ {
		data := buf[0 : 5+i*10]
		n, err := w.Write(data)
		if len(data) != n {
			t.Errorf("write: %d got %d", len(data), n)
		}
		if err != nil {
			t.Errorf("write error: %s", err)
		}
		nn := <-c
		if nn != n {
			t.Errorf("write: %d read: %d", n, nn)
		}
	}
	w.Close()
	nn := <-c
	if nn != 0 {
		t.Errorf("final read got %d", nn)
	}
}
