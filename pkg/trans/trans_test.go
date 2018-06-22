package trans

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/axgle/mahonia"
	iconv "github.com/djimenez/iconv-go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var testBytes = []byte{0xC4, 0xE3, 0xBA, 0xC3, 0xA3, 0xAC, 0xCA, 0xC0, 0xBD, 0xE7, 0xA3, 0xA1}

func TestTrans_Mahonia(t *testing.T) {
	testStr := string(testBytes)
	enc := mahonia.NewDecoder("gbk")
	res := enc.ConvertString(testStr)

	log.Println(res)
}

func TestTrans_Iconv(t *testing.T) {
	var res []byte
	// NOTE: not work
	iconv.Convert(testBytes, res, "GBK", "UTF-8")
	log.Println(string(res)) // 你好，世界！
}

func TestTrans_Text(t *testing.T) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	reader := transform.NewReader(bytes.NewReader(testBytes), decoder)
	res, _ := ioutil.ReadAll(reader)
	log.Printf(string(res)) // 你好，世界！
}
