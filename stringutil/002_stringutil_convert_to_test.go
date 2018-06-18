package stringutil

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToFloat(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Float, nil)
	assert.NoError(err)
	assert.Equal(float64(0), v)

	v, err = ConvertTo(Float, "1.5")
	assert.NoError(err)
	assert.Equal(float64(1.5), v)

	v, err = ConvertTo(Float, "1")
	assert.NoError(err)
	assert.Equal(float64(1.0), v)

	v, err = ConvertToFloat("1.5")
	assert.NoError(err)
	assert.Equal(float64(1.5), v)

	v, err = ConvertToFloat("1.0")
	assert.NoError(err)
	assert.Equal(float64(1.0), v)

	for _, fail := range []string{`potato`, `true`, `2015-05-01 00:15:16`} {
		_, err := ConvertTo(Float, fail)
		assert.Error(err)

		_, err = ConvertToFloat(fail)
		assert.Error(err)
	}
}

func TestConvertToInteger(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Integer, nil)
	assert.NoError(err)
	assert.Equal(int64(0), v)

	v, err = ConvertTo(Integer, "7")
	assert.NoError(err)
	assert.Equal(int64(7), v)

	v, err = ConvertToInteger("7")
	assert.NoError(err)
	assert.Equal(int64(7), v)

	tm := time.Date(2010, 2, 21, 15, 14, 13, 0, time.UTC)

	v, err = ConvertTo(Integer, tm)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	v, err = ConvertToInteger(tm)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	v, err = ConvertTo(Integer, `2010-02-21 15:14:13`)
	assert.NoError(err)
	assert.Equal(tm.UnixNano(), v)

	for _, fail := range []string{`0.0`, `1.5`, `potato`, `true`} {
		_, err := ConvertTo(Integer, fail)
		assert.Error(err)

		_, err = ConvertToInteger(fail)
		assert.Error(err)
	}
}

func TestConvertToBytes(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Bytes, nil)
	assert.NoError(err)
	assert.Equal([]byte{}, v)

	v, err = ConvertTo(Bytes, []byte{})
	assert.NoError(err)
	assert.Equal([]byte{}, v)

	v, err = ConvertTo(Bytes, []byte{1, 2, 3})
	assert.NoError(err)
	assert.Equal([]byte{1, 2, 3}, v)

	v, err = ConvertTo(Bytes, `test`)
	assert.NoError(err)
	assert.Equal([]byte{0x74, 0x65, 0x73, 0x74}, v)

	v, err = ConvertTo(Bytes, []int{0x74, 0x65, 0x73, 0x74})
	assert.NoError(err)
	assert.Equal(`test`, string(v.([]byte)))
}

func TestConvertToBoolean(t *testing.T) {
	assert := require.New(t)

	v, err := ConvertTo(Boolean, nil)
	assert.Equal(false, v)

	v, err = ConvertTo(Boolean, `true`)
	assert.NoError(err)
	assert.Equal(true, v)

	v, err = ConvertTo(Boolean, `false`)
	assert.NoError(err)
	assert.Equal(false, v)

	v, err = ConvertToBool(`true`)
	assert.NoError(err)
	assert.Equal(true, v)

	v, err = ConvertToBool(`false`)
	assert.NoError(err)
	assert.Equal(false, v)

	for _, fail := range []string{`1.5`, `potato`, `01`, `2015-05-01 00:15:16`} {
		_, err := ConvertTo(Boolean, fail)
		assert.Error(err)

		_, err = ConvertToBool(fail)
		assert.Error(err)
	}
}

func TestConvertToTime(t *testing.T) {
	assert := require.New(t)

	atLeastNow := time.Now()

	values := map[string]time.Time{
		`2015-05-01 00:15:16`:         time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		`Fri May 1 00:15:16 UTC 2015`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `Fri May 01 00:15:16 +0000 2015`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `01 May 15 00:15 UTC`:            time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `01 May 15 00:15 +0000`:          time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		// `Friday, 01-May-15 00:15:16 UTC`: time.Date(2015, 5, 1, 0, 15, 16, 0, time.UTC),
		`1136239445`: time.Date(2006, 1, 2, 17, 4, 5, 0, time.Now().Location()),
	}

	v, err := ConvertToTime(`now`)
	assert.Nil(err)
	assert.True(v.After(atLeastNow))

	v, err = ConvertToTime(time.Now())
	assert.Nil(err)
	assert.True(v.After(atLeastNow))

	v, err = ConvertToTime(`0000-00-00 00:00:00`)
	assert.Nil(err)
	assert.Zero(v)

	for in, out := range values {
		if v, err := ConvertTo(Time, in); err == nil {
			switch v.(type) {
			case time.Time:
				assert.Equal(out, v.(time.Time))
			default:
				t.Errorf("Conversion yielded an incorrect result type: expected time.Time, got: %T", v)
			}
		} else {
			t.Errorf("Error during conversion: %v", err)
		}

		if v, err := ConvertToTime(in); err == nil {
			assert.Equal(out, v)
		} else {
			t.Errorf("Error during conversion: %v", err)
		}
	}

	for _, fail := range []string{`1.5`, `potato`, `false`} {
		if _, err := ConvertTo(Time, fail); err == nil {
			t.Errorf("Conversion should have failed for value '%s', but didn't", fail)
		}

		if _, err := ConvertToTime(fail); err == nil {
			t.Errorf("Conversion should have failed for value '%s', but didn't", fail)
		}
	}
}

func TestAutotypeNil(t *testing.T) {
	assert := assert.New(t)

	for _, testValue := range []string{
		``,
		`nil`,
		`null`,
		`Nil`,
		`NULL`,
		`None`,
		`undefined`,
	} {
		assert.Nil(Autotype(testValue), fmt.Sprintf("%q was not autotyped to nil", testValue))
	}
}

func TestAutotypeFloat(t *testing.T) {
	for _, testValue := range []string{
		`1.5`,
		`-1.5`,
	} {
		v := Autotype(testValue)

		switch v.(type) {
		case float64:
			return
		default:
			t.Errorf("Invalid autotype: expected float64, got %T", v)
		}
	}
}

func TestAutotypeInt(t *testing.T) {
	for _, testValue := range []string{
		`12345`,
		`-12345`,
	} {
		v := Autotype(testValue)

		switch v.(type) {
		case int64:
			continue
		default:
			t.Errorf("Invalid autotype: expected int64, got %T", v)
		}
	}
}

func TestAutotypePreserveLeadingZeroes(t *testing.T) {
	for _, testValue := range []string{
		`00`,
		`01`,
		`07753`,
		`06094`,
		`0000000010000000`,
	} {
		v := Autotype(testValue)

		switch v.(type) {
		case string:
			continue
		default:
			t.Errorf("Invalid autotype: expected string, got %T", v)
		}
	}
}

func TestAutotypeDate(t *testing.T) {
	for _, testValue := range TimeFormats {
		tvS := strings.Replace(string(testValue), `_`, ``, -1)
		tvS = strings.TrimSuffix(tvS, `07:00`)

		v := Autotype(tvS)

		switch v.(type) {
		case time.Time:
			continue
		default:
			t.Errorf("Invalid autotype %q: expected time.Time, got %T", testValue, v)
		}
	}
}

func TestAutotypeBool(t *testing.T) {
	for _, testValue := range []string{
		`true`,
		`True`,
		`false`,
		`False`,
	} {
		v := Autotype(testValue)

		switch v.(type) {
		case bool:
			continue
		default:
			t.Errorf("Invalid autotype: expected bool, got %T", v)
		}
	}

	for _, testValue := range []string{
		`trues`,
		`Falses`,
		`potato`,
	} {
		v := Autotype(testValue)

		switch v.(type) {
		case string:
			continue
		default:
			t.Errorf("Invalid autotype: expected string, got %T", v)
		}
	}
}
