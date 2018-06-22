package tx

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

// Fixed8 fixed point number
type Fixed8 int64

// MakeFixed8 .
func MakeFixed8(val float64) Fixed8 {
	return trunc(val)
}

// Int convert to big.Int object
func (fixed8 *Fixed8) Int() *big.Int {
	return big.NewInt(int64(*fixed8))
}

func (fixed8 *Fixed8) Write(writer io.Writer) error {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(*fixed8))

	_, err := writer.Write(data)

	return err
}

func (fixed8 *Fixed8) Read(reader io.Reader) error {

	data := make([]byte, 8)

	_, err := reader.Read(data)

	if err != nil {
		return err
	}

	*fixed8 = Fixed8(binary.LittleEndian.Uint64(data))

	return nil
}

// Float64 convert fixe8 to float64
func (fixed8 *Fixed8) Float64() float64 {
	valstr := fmt.Sprintf("%.8f", float64(*fixed8)/100000000)

	r, _ := strconv.ParseFloat(valstr, 8)

	return r
}

func (fixed8 *Fixed8) String() string {
	return fmt.Sprintf("%.8f", float64(*fixed8)/100000000)
}

func trunc(val float64) Fixed8 {

	val = val * 100000000

	valstr := fmt.Sprintf("%.0f", val)

	r, _ := strconv.ParseFloat(valstr, 8)

	return Fixed8(r)
}
