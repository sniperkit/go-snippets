package base58

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestEncodeEmptyString(t *testing.T) {
	actual, expected := Encode([]byte{}), ""
	if actual != expected {
		t.Errorf("EncodeEmpty Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestEncodeNil(t *testing.T) {
	actual, expected := Encode(nil), ""
	if actual != expected {
		t.Errorf("EncodeNil Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestEncode_test123(t *testing.T) {
	actual, expected := Encode([]byte("test123")), "5QqG6h3Xdc"
	if actual != expected {
		t.Errorf("Encode_test123 Failed: expected: '%s' got '%s'", expected, actual)
	}
}

func TestDecode_test123(t *testing.T) {
	actual, err := Decode("5QqG6h3Xdc")
	expected := []byte("test123")

	if err != nil {
		t.Errorf("Decode_test123 failed with error: '%s'", err)
	}

	if bytes.Compare(actual, expected) != 0 {
		t.Errorf("Decode_test123 Failed: expected: '%d' got '%d'", expected, actual)
	}
}

func TestEncodeDecodeArray(t *testing.T) {
	bytes1 := []byte{0x53, 0xFE, 0x92}
	var s1 = Encode(bytes1)

	actual, expected := s1, "VDLu"

	if actual != expected {
		t.Errorf("EncodeDecodeArray: encode arrray to string Failed: expected: '%s' got '%s'", expected, actual)
	}

	b := []byte{116, 32, 8, 99, 100, 232, 4, 7}
	s := Encode(b)

	actual, expected = s, "LRZSAm8VfWE"

	if actual != expected {
		t.Errorf("EncodeDecodeArray: encode arrray to string Failed: expected: '%s' got '%s'", expected, actual)
	}

	arr, _ := Decode(s)

	if bytes.Compare(arr, b) != 0 {
		t.Errorf("EncodeDecodeArray: decode string to array Failed: expected: '%d' got '%d'", b, arr)
	}
}

func TestEncodeDecodeStr256(t *testing.T) {
	var str256 = [256]byte{}

	for i := 0; i <= 255; i++ {
		str256[i] = byte(i)
	}

	strB62 := Encode(str256[:])
	strData, _ := Decode(strB62)

	if bytes.Compare(strData, str256[:]) != 0 {
		t.Errorf("EncodeDecodeStr256 Failed: expected: '%d' got '%d'", str256, strData)
	}
}

func TestKeyBuf(t *testing.T) {
	// this buff casued an encoding overrun beciase
	// last entry has zero length
	keybuf := []byte{24, 23, 224, 166, 164, 198, 162, 13, 94, 181, 12, 245, 108,
		24, 143, 220, 152, 181, 9, 74, 70, 81, 227, 157, 1, 41, 78, 125, 143,
		229, 88, 105, 247, 107, 128, 90, 144, 179, 55, 168, 51, 205, 190, 33,
		46, 123, 86, 123, 129, 206, 185, 206, 231, 48, 21, 76}

	key := Encode(keybuf)
	decoded, _ := Decode(key)

	if len(decoded) != len(keybuf) {
		t.Errorf("KeyBuf Failed: expected: '%d' got '%d'", len(keybuf), len(decoded))
	}

	if bytes.Compare(decoded, keybuf) != 0 {
		t.Errorf("KeyBuf Failed: expected: '%d' got '%d'", keybuf, decoded)
	}
}

func TestHighByteVals(t *testing.T) {

	val63 := []byte{252}

	encoded := Encode(val63)
	decoded, _ := Decode(encoded)

	if bytes.Compare(decoded, val63) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val63, decoded)
	}

	val62 := []byte{248}

	encoded = Encode(val62)
	decoded, _ = Decode(encoded)

	if bytes.Compare(decoded, val62) != 0 {
		t.Errorf("HighByteVals Failed: expected: '%d' got '%d'", val62, decoded)
	}

	val61 := []byte{244}

	encoded = Encode(val61)
	decoded, _ = Decode(encoded)

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
		if res := Encode(tmp); res != test.out {
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
		if res, _ := Decode(test.out); bytes.Equal(res, b) != true {
			t.Errorf("Base58Decode test #%d failed: got: %q want: %q",
				x, res, test.in)
			continue
		}
	}

	// Base58Decode with invalid input
	for x, test := range invalidStringTests {
		if res, _ := Decode(test.in); string(res) != test.out {
			t.Errorf("Base58Decode invalidString test #%d failed: got: %q want: %q",
				x, res, test.out)
			continue
		}
	}
}
