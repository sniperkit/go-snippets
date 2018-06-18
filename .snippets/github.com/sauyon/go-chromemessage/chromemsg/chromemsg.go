package chromemsg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"os"
	"unsafe"
)

var nativeEndian binary.ByteOrder = endianness()

var defaultMsgr = Messenger{bufio.NewReadWriter(
	bufio.NewReader(os.Stdin),
	bufio.NewWriter(os.Stdout))}

type Messenger struct {
	port *bufio.ReadWriter
}

func New(port *bufio.ReadWriter) *Messenger {
	return &Messenger{port}
}

func Read(data interface{}) error {
	return defaultMsgr.Read(data)
}

func Write(msg interface{}) error {
	return defaultMsgr.Write(msg)
}

func (msgr *Messenger) Read(data interface{}) error {
	lengthBits := make([]byte, 4)
	_, err := msgr.port.Read(lengthBits)
	if err != nil {
		return err
	}
	length := nativeToInt(lengthBits)
	content := make([]byte, length)
	_, err = msgr.port.Read(content)
	if err != nil {
		return err
	}
	json.Unmarshal(content, data)
	return nil
}

func (msgr *Messenger) Write(msg interface{}) error {
	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	length := len(json)
	bits := make([]byte, 4)
	buf := bytes.NewBuffer(bits)
	err = binary.Write(buf, nativeEndian, length)
	if err != nil {
		return err
	}
	_, err = msgr.port.Write(bits)
	if err != nil {
		return err
	}

	_, err = msgr.port.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func nativeToInt(bits []byte) int {
	var length uint32
	buf := bytes.NewBuffer(bits)
	binary.Read(buf, nativeEndian, &length)
	return int(length)
}

func endianness() binary.ByteOrder {
	var i int = 1
	bs := (*[unsafe.Sizeof(0)]byte)(unsafe.Pointer(&i))
	if bs[0] == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}
