package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"bytes"
	"encoding/binary"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:12345")
	//CheckError(err)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error", err.Error())
			continue
		}
		go Handle(conn)
	}
}
func Handle(conn net.Conn) {
	buffer := make([]byte, 1024)
	tmpBuffer := make([]byte, 0)
	readerChannal := make(chan []byte, 1024)
	go func(read chan []byte) {
		for {
			select {
			case data := <-read:
				fmt.Println("recv : %d = %s", len(data), string(data))
			}
		}
	}(readerChannal)
	for {
		time.Sleep(1 * time.Second)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err.Error())
		}
		tmpBuffer = Depack(append(tmpBuffer, buffer[:n]...), readerChannal)
		fmt.Println("server recv : %n = %s", n, string(buffer))
		n, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err.Error())
		}
	}
}

const (
	ConstHeader = "Headers"
	ConstHeaderLength = 7
	ConstLength = 4
)

func Enpack(message []byte)[]byte{
	return append(append([]byte(ConstHeader),IntToBytes(len(message))...),message...)
}

func Depack(buffer []byte,readerChannel chan []byte) []byte{
	length := len(buffer)

	var i int
	for i = 0 ; i < length ; i++{
		if length < i + ConstHeaderLength + ConstLength{
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == ConstHeader{
			messageLength := BytesToInt(buffer[i+ConstHeaderLength:i+ConstHeaderLength+ConstLength])
			if length < i+ConstHeaderLength+messageLength{
				break
			}
			data := buffer[i+ConstHeaderLength+ConstLength : i+ConstHeaderLength+ConstLength+messageLength]
			readerChannel <- data
		}
	}
	if i == length{
		return make([]byte,0)
	}

	return buffer[i:]
}

func IntToBytes(n int)[]byte{
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer,binary.BigEndian,x)
	return bytesBuffer.Bytes()
}

func BytesToInt(b []byte) int{
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer,binary.BigEndian,&x)
	return int(x)
}
