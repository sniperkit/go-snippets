package ocr

import (
	"fmt"
	"testing"

	"github.com/otiai10/gosseract"
)

// need `testeract` installed first
// brew install --with-all-languages tesseract
func TestOCR(t *testing.T) {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("friendsship.jpg")
	// client.SetLanguage("chi_sim")
	client.SetLanguage("eng")
	text, _ := client.Text()
	fmt.Println(text)
}
