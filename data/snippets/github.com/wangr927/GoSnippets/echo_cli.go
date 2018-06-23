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
	conn,err := net.Dial("tcp","127.0.0.1:12345")
	CheckError(err)
	defer conn.Close()
	readerChannal := make(chan []byte,1024)
	go func(read chan []byte){
		for{
			select{
			case data := <- read:
				fmt.Println("client recv : %d = %s", len(data), string(data))
			}
		}
	}(readerChannal)
	go func(read chan []byte){
		tmpBuf := make([]byte,0)
		buf := make([]byte, 1024)
		for{
			time.Sleep(time.Second*100)
			n, err := conn.Read(buf)
			fmt.Println("client recv:",string(buf))
			CheckError(err)
			tmpBuf = Depack(append(tmpBuf,buf[:n]...),read)
		}
	}(readerChannal)
	for {
		_, err := conn.Write(Enpack([]byte("Hello world")))
		CheckError(err)
		time.Sleep(1*time.Second)
	}
}
func CheckError(err error){
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
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
