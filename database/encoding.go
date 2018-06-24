package database

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/rsmohamad/comp4321/models"
	"math"
)

func uint64ToByte(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func byteToUint64(arr []byte) uint64 {
	return binary.BigEndian.Uint64(arr)
}

func intToByte(v int) []byte {
	b := make([]byte, 8)
	binary.PutVarint(b, int64(v))
	return b
}

func byteToInt(arr []byte) int {
	val, _ := binary.Varint(arr)
	return int(val)
}

func float64ToByte(f float64) []byte {
	return uint64ToByte(math.Float64bits(f))
}

func byteToFloat64(arr []byte) float64 {
	return math.Float64frombits(byteToUint64(arr))
}

func historyToByte(arr []models.SearchHistory) []byte {
	var byteBuffer bytes.Buffer
	encoder := gob.NewEncoder(&byteBuffer)
	encoder.Encode(arr)
	return byteBuffer.Bytes()
}

func byteToHistory(arr []byte) []models.SearchHistory {
	var byteBuffer bytes.Buffer
	byteBuffer.Write(arr)

	decoder := gob.NewDecoder(&byteBuffer)
	var rv []models.SearchHistory
	decoder.Decode(&rv)
	return rv
}

func docToByte(doc *models.Document) []byte {
	var byteBuffer bytes.Buffer
	encoder := gob.NewEncoder(&byteBuffer)
	encoder.Encode(doc)
	return byteBuffer.Bytes()
}

func byteToDoc(arr []byte) *models.Document {
	var byteBuffer bytes.Buffer
	byteBuffer.Write(arr)

	decoder := gob.NewDecoder(&byteBuffer)
	var rv models.Document
	decoder.Decode(&rv)
	return &rv
}
