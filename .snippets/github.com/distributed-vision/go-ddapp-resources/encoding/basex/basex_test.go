package basex

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBase58EncodeEmptyString(t *testing.T) {
	actual, expected := BTCEncoder.Encode([]byte{}), ""
	if actual != expected {
		t.Errorf("EncodeEmpty Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase58EncodeNil(t *testing.T) {
	actual, expected := BTCEncoder.Encode(nil), ""
	if actual != expected {
		t.Errorf("EncodeNil Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase58Encode_test123(t *testing.T) {
	actual, expected := BTCEncoder.Encode([]byte("test123")), "5QqG6h3Xdc"
	if actual != expected {
		t.Errorf("Encode_test123 Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase58Decode_test123(t *testing.T) {
	actual, err := BTCEncoder.Decode("5QqG6h3Xdc")
	expected := []byte("test123")

	if err != nil {
		t.Errorf("Decode_test123 Failed: unexpected error: '%s", err)
	}

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("Decode_test123 Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase58EncodeDecodeArray(t *testing.T) {
	bytes1 := []byte{0x53, 0xFE, 0x92}
	var s1 = BTCEncoder.Encode(bytes1)

	actual, expected := s1, "VDLu"

	if actual != expected {
		t.Errorf("EncodeDecodeArray: encode arrray to string Failed: expected: '%s' got '%s'", expected, actual)
	}

	b := []byte{116, 32, 8, 99, 100, 232, 4, 7}
	s := BTCEncoder.Encode(b)

	actual, expected = s, "LRZSAm8VfWE"

	if actual != expected {
		t.Errorf("EncodeDecodeArray: encode arrray to string Failed: expected: '%s' got '%s'", expected, actual)
	}

	arr, _ := BTCEncoder.Decode(s)

	if bytes.Compare(arr, b) != 0 {
		t.Errorf("EncodeDecodeArray: decode string to array Failed: expected: '%d' got '%d'", b, arr)
	}
}

func TestBase58EncodeDecodeStr256(t *testing.T) {
	var str256 = [256]byte{}

	for i := 0; i <= 255; i++ {
		str256[i] = byte(i)
	}

	strB62 := BTCEncoder.Encode(str256[:])
	strData, _ := BTCEncoder.Decode(strB62)

	if bytes.Compare(strData, str256[:]) != 0 {
		t.Errorf("EncodeDecodeStr256 Failed: expected: '%d' got '%d'", str256, strData)
	}
}

func TestBase58KeyBuf(t *testing.T) {
	// this buff casued an encoding overrun beciase
	// last entry has zero length
	keybuf := []byte{24, 23, 224, 166, 164, 198, 162, 13, 94, 181, 12, 245, 108,
		24, 143, 220, 152, 181, 9, 74, 70, 81, 227, 157, 1, 41, 78, 125, 143,
		229, 88, 105, 247, 107, 128, 90, 144, 179, 55, 168, 51, 205, 190, 33,
		46, 123, 86, 123, 129, 206, 185, 206, 231, 48, 21, 76}

	key := BTCEncoder.Encode(keybuf)
	decoded, _ := BTCEncoder.Decode(key)

	if len(decoded) != len(keybuf) {
		t.Errorf("KeyBuf Failed: expected: '%d' got '%d'", len(keybuf), len(decoded))
	}

	if bytes.Compare(decoded, keybuf) != 0 {
		t.Errorf("KeyBuf Failed: expected: '%d' got '%d'", keybuf, decoded)
	}
}

func TestBase58HighByteVals(t *testing.T) {

	val63 := []byte{252}

	encoded := BTCEncoder.Encode(val63)
	decoded, _ := BTCEncoder.Decode(encoded)

	if bytes.Compare(decoded, val63) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val63, decoded)
	}

	val62 := []byte{248}

	encoded = BTCEncoder.Encode(val62)
	decoded, _ = BTCEncoder.Decode(encoded)

	if bytes.Compare(decoded, val62) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val62, decoded)
	}

	val61 := []byte{244}

	encoded = BTCEncoder.Encode(val61)
	decoded, _ = BTCEncoder.Decode(encoded)

	if bytes.Compare(decoded, val61) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val61, decoded)
	}
}

var stringTests = []struct {
	in  string
	out string
}{
	{"", ""},
	{" ", "Z"},
	{"-", "n"},
	{"0", "q"},
	{"1", "r"},
	{"-1", "4SU"},
	{"11", "4k8"},
	{"abc", "ZiCa"},
	{"1234598760", "3mJr7AoUXx2Wqd"},
	{"abcdefghijklmnopqrstuvwxyz", "3yxU3u1igY8WkgtjK92fbJQCd4BZiiT1v25f"},
	{"00000000000000000000000000000000000000000000000000000000000000", "3sN2THZeE9Eh9eYrwkvZqNstbHGvrxSAM7gXUXvyFQP8XvQLUqNCS27icwUeDT7ckHm4FUHM2mTVh1vbLmk7y"},
}

var invalidStringTests = []struct {
	in  string
	out string
}{
	{"0", ""},
	{"O", ""},
	{"I", ""},
	{"l", ""},
	{"3mJr0", ""},
	{"O3yxU", ""},
	{"3sNI", ""},
	{"4kl8", ""},
	{"0OIl", ""},
	{"!@#$%^&*()-_=+~`", ""},
}

var hexTests = []struct {
	in  string
	out string
}{
	{"61", "2g"},
	{"626262", "a3gV"},
	{"636363", "aPEr"},
	{"73696d706c792061206c6f6e6720737472696e67", "2cFupjhnEsSn59qHXstmK2ffpLv2"},
	{"00eb15231dfceb60925886b67d065299925915aeb172c06647", "1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L"},
	{"516b6fcd0f", "ABnLTmg"},
	{"bf4f89001e670274dd", "3SEo3LWLoPntC"},
	{"572e4794", "3EFU7m"},
	{"ecac89cad93923c02321", "EJDM8drfXA6uyA"},
	{"10c8511e", "Rt5zm"},
	{"00000000000000000000", "1111111111"},
}

func TestBase58Compat(t *testing.T) {
	// Base58Encode tests
	for x, test := range stringTests {
		tmp := []byte(test.in)
		if res := BTCEncoder.Encode(tmp); res != test.out {
			t.Errorf("Base58Encode test #%d failed: got: %s want: %s",
				x, res, test.out)
			continue
		}
	}

	// Base58Decode tests
	for x, test := range hexTests {
		b, err := hex.DecodeString(test.in)
		if err != nil {
			t.Errorf("hex.DecodeString failed failed #%d: got: %s", x, test.in)
			continue
		}
		if res, _ := BTCEncoder.Decode(test.out); bytes.Equal(res, b) != true {
			t.Errorf("Base58Decode test #%d failed: got: %q want: %q",
				x, res, test.in)
			continue
		}
	}

	// Base58Decode with invalid input
	for x, test := range invalidStringTests {
		if res, _ := BTCEncoder.Decode(test.in); string(res) != test.out {
			t.Errorf("Base58Decode invalidString test #%d failed: got: %q want: %q",
				x, res, test.out)
			continue
		}
	}
}

var base62 = NewEncoder("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func TestBase62EncodeEmptyString(t *testing.T) {
	actual, expected := base62.Encode([]byte{}), ""
	if actual != expected {
		t.Errorf("EncodeEmpty Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase62EncodeNil(t *testing.T) {
	actual, expected := base62.Encode(nil), ""
	if actual != expected {
		t.Errorf("EncodeNil Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase62Encode_test123(t *testing.T) {
	actual, expected := base62.Encode([]byte("test123")), "2Q3IiUJ1MJ"
	if actual != expected {
		t.Errorf("Encode_test123 Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase62Decode_test123(t *testing.T) {
	actual, err := base62.Decode("2Q3IiUJ1MJ")
	expected := []byte("test123")

	if err != nil {
		t.Errorf("Decode_test123 Failed: unexpected error: '%s", err)
	}

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("Decode_test123 Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase62EncodeDecodeArray(t *testing.T) {
	bytes1 := []byte{0x53, 0xFE, 0x92}
	var s1 = base62.Encode(bytes1)

	actual, expected := s1, "N60o"

	if actual != expected {
		t.Errorf("EncodeDecodeArray: encode arrray to string Failed: expected: '%s' got '%s'", expected, actual)
	}

	b := []byte{116, 32, 8, 99, 100, 232, 4, 7}
	s := base62.Encode(b)

	actual, expected = s, "9y88pxSoliJ"

	if actual != expected {
		t.Errorf("EncodeDecodeArray: encode arrray to string Failed: expected: '%s' got '%s'", expected, actual)
	}

	arr, _ := base62.Decode(s)

	if bytes.Compare(arr, b) != 0 {
		t.Errorf("EncodeDecodeArray: decode string to array Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestBase62EncodeDecodeStr256(t *testing.T) {
	var str256 = [256]byte{}

	for i := 0; i <= 255; i++ {
		str256[i] = byte(i)
	}

	strB62 := base62.Encode(str256[:])
	strData, _ := base62.Decode(strB62)

	if bytes.Compare(strData, str256[:]) != 0 {
		t.Errorf("EncodeDecodeStr256 Failed: expected: '%d' got '%d'", str256, strData)
	}
}

func TestBase62KeyBuf(t *testing.T) {
	// this buff casued an encoding overrun beciase
	// last entry has zero length
	keybuf := []byte{24, 23, 224, 166, 164, 198, 162, 13, 94, 181, 12, 245, 108,
		24, 143, 220, 152, 181, 9, 74, 70, 81, 227, 157, 1, 41, 78, 125, 143,
		229, 88, 105, 247, 107, 128, 90, 144, 179, 55, 168, 51, 205, 190, 33,
		46, 123, 86, 123, 129, 206, 185, 206, 231, 48, 21, 76}

	key := base62.Encode(keybuf)
	decoded, _ := base62.Decode(key)

	if len(decoded) != len(keybuf) {
		t.Errorf("KeyBuf Failed: expected: '%d' got '%d'", len(keybuf), len(decoded))
	}

	if bytes.Compare(decoded, keybuf) != 0 {
		t.Errorf("KeyBuf Failed: expected: '%d' got '%d'", keybuf, decoded)
	}
}

func TestBase62HighByteVals(t *testing.T) {

	val63 := []byte{252}

	encoded := base62.Encode(val63)
	decoded, _ := base62.Decode(encoded)

	if bytes.Compare(decoded, val63) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val63, decoded)
	}

	val62 := []byte{248}

	encoded = base62.Encode(val62)
	decoded, _ = base62.Decode(encoded)

	if bytes.Compare(decoded, val62) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val62, decoded)
	}

	val61 := []byte{244}

	encoded = base62.Encode(val61)
	decoded, _ = base62.Decode(encoded)

	if bytes.Compare(decoded, val61) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val61, decoded)
	}
}
