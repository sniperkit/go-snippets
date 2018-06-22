package download

import (
	"io"
	"net/http"

	"github.com/gophergala/aeris/info"
)

func Download(i *info.Info, stream *info.Stream, output io.Writer) error {
	res, err := http.Get(stream.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// we need to determine whether we're dealing with a video stream with an
	// encrypted signature, or an already decrypted signature
	if res.Header.Get("Content-Length") == "0" {
		// signature is encrypted, we need to decrypt the signature
		res.Body.Close()

		i.DecryptSignatures()

		// issue a new HTTP request
		res, err = http.Get(stream.Url)
		if err != nil {
			return err
		}
		defer res.Body.Close()
	}

	_, err = io.Copy(output, res.Body)

	return err
}
