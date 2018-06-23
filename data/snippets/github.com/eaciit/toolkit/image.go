package toolkit

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"strings"
)

func Base64ToImage(base64String, destination string) error {
	if strings.HasPrefix(base64String, "data:image") {
		parts := strings.SplitN(base64String, ",", 2)
		if len(parts) > 1 {
			base64String = parts[1]
		} else {
			return errors.New("invalid base64 string")
		}
	}

	bts := bytes.NewBufferString(base64String)
	decoded := base64.NewDecoder(base64.StdEncoding, bts)
	f, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, decoded)
	if err != nil {
		return err
	}

	return nil
}

func ImageToBase64(source string) (string, error) {
	f, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fInfo, _ := f.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(f)
	fReader.Read(buf)

	imgBase64String := base64.StdEncoding.EncodeToString(buf)
	return imgBase64String, nil
}
