package bin

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	X, Y *int32
	Name string
}

func txt_write() {
	data := []byte("Hello World!\n")

	file1, _ := os.Create("data2")
	defer file1.Close()

	bytes, _ := file1.Write(data)
	fmt.Printf("Wrote %d bytes to file \n", bytes)

	file2, _ := os.Open("data2")
	defer file2.Close()

	read2 := make([]byte, len(data))

	bytes, _ = file2.Read(read2)

	fmt.Printf("Read %d bytes from file\n", bytes)
	fmt.Println(read2, string(read2))
}

func main() {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	dec := gob.NewDecoder(&network)

	// Encode (send) the value.
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}

	// Decode (receive) the value.
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Println(q)
	fmt.Printf("%q: {%d,%d}\n", q.Name, *q.X, *q.Y)

	outputFile, outputError := os.OpenFile("output.dat", os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		fmt.Printf("An error occurred with file opening or creation\n")
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	outputString := "hello world!\n"

	for i := 0; i < 10; i++ {
		outputWriter.WriteString(outputString)
	}
	outputWriter.Flush()

}
