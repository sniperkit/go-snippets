package tx

import (
	"encoding/binary"
	"io"
	"math"
)

// Varint .
type Varint uint64

func (varint *Varint) Write(writer io.Writer) error {

	if uint64(*varint) < 0xFD {
		_, err := writer.Write([]byte{byte(*varint)})
		return err
	} else if uint64(*varint) < math.MaxUint16 {
		buff := make([]byte, 3)
		buff[0] = 0xFD
		binary.LittleEndian.PutUint16(buff[1:], uint16(*varint))
		_, err := writer.Write(buff)

		return err
	} else if uint64(*varint) < math.MaxUint32 {
		buff := make([]byte, 5)
		buff[0] = 0xFE
		binary.LittleEndian.PutUint32(buff[1:], uint32(*varint))
		_, err := writer.Write(buff)

		return err
	}

	buff := make([]byte, 9)
	buff[0] = 0xFF
	binary.LittleEndian.PutUint64(buff[1:], uint64(*varint))
	_, err := writer.Write(buff)

	return err
}

func (varint *Varint) Read(reader io.Reader) error {
	flag := make([]byte, 1)

	_, err := reader.Read(flag)

	if err != nil {
		return err
	}

	if flag[0] == 0xFD {
		buff := make([]byte, 2)
		_, err := reader.Read(buff)

		(*varint) = Varint(binary.LittleEndian.Uint16(buff))

		if err != nil {
			return err
		}
	} else if flag[0] == 0xFE {
		buff := make([]byte, 4)
		_, err := reader.Read(buff)

		(*varint) = Varint(binary.LittleEndian.Uint32(buff))

		if err != nil {
			return err
		}
	} else if flag[0] == 0xFF {
		buff := make([]byte, 8)
		_, err := reader.Read(buff)

		(*varint) = Varint(binary.LittleEndian.Uint64(buff))

		if err != nil {
			return err
		}

	} else {
		*varint = Varint(flag[0])
	}

	return nil
}
