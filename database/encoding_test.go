package database

import (
	"testing"
)

func TestUint64(t *testing.T) {
	testcases := []uint64{18446744073709551615, 0}

	for _, num := range testcases {
		t.Run("Uint64", func(t *testing.T) {
			numByte := uint64ToByte(num)
			if num != byteToUint64(numByte) {
				t.Log(num, byteToUint64(numByte))
				t.Fail()
			}
		})
	}

}

func TestInt(t *testing.T) {
	testcases := []int{2147483647, -2147483648}

	for _, num := range testcases {
		t.Run("Int", func(t *testing.T) {
			numByte := intToByte(num)
			if num != byteToInt(numByte) {
				t.Log(num, byteToInt(numByte))
				t.Fail()
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	testcases := []float64{-5.66, +6.55}

	for _, num := range testcases {
		t.Run("Float64", func(t *testing.T) {
			numByte := float64ToByte(num)
			if num != byteToFloat64(numByte) {
				t.Log(num, byteToFloat64(numByte))
				t.Fail()
			}
		})
	}

}
