package pdf

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/uoregon-libraries/gopkg/logger"
)

// ImageDPI holds x and y dpis gathered from PDF images
type ImageDPI struct {
	X float64
	Y float64
}

// ImageDPIs returns an array of ImageDPIs by reading the images in the
// given PDF with "pdfimages -list"
func ImageDPIs(path string) []ImageDPI {
	var cmdParts = []string{"pdfimages", "-list", path}
	var cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	var output, err = cmd.CombinedOutput()
	if err != nil {
		logger.Errorf(`Failed to run "%s": %s`, strings.Join(cmdParts, " "), err)
		for _, line := range bytes.Split(output, []byte("\n")) {
			logger.Debugf("--> %s", line)
		}

		return nil
	}

	var dpis []ImageDPI
	for i, line := range bytes.Split(output, []byte("\n")) {
		// The first two lines don't give us any information
		if bytes.HasPrefix(line, []byte("page")) || bytes.HasPrefix(line, []byte("--------")) {
			continue
		}

		// The last line appears to always be blank
		if len(line) == 0 {
			continue
		}

		var parts = bytes.Fields(line)
		if len(parts) < 15 {
			logger.Errorf("Too few fields in line %d of %q output: %q", i+1, strings.Join(cmdParts, " "), line)
			return nil
		}

		var xdpiString, ydpiString, sizeString = string(parts[12]), string(parts[13]), string(parts[14])

		// In rare cases, we have embedded images with no DPI information
		// somehow...  we have to ignore these situations....
		if xdpiString == "inf" || ydpiString == "inf" {
			continue
		}

		// Then there are the cases where the size string is empty or so low we
		// just don't care about the image
		if sizeString == "-" || sizeString[len(sizeString)-1] == 'B' {
			logger.Debugf("Invalid DPI information in line %d of %q output: %q (skipping: small image)",
				i+1, strings.Join(cmdParts, " "), line)
			continue
		}

		var xdpi, _ = strconv.ParseFloat(xdpiString, 64)
		var ydpi, _ = strconv.ParseFloat(ydpiString, 64)

		if xdpi == 0 || ydpi == 0 {
			logger.Errorf("Invalid DPI information in line %d of %q output: %q", i+1, strings.Join(cmdParts, " "), line)
			return nil
		}

		dpis = append(dpis, ImageDPI{X: xdpi, Y: ydpi})
	}

	return dpis
}
