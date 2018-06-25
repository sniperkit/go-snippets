package main

import (
	"encoding/csv"
	"io"
	"strings"
)

// trimSpace trims the whitespace from each string in the passed slice.
func trimSpace(slice []string) []string {
	for i, v := range slice {
		slice[i] = strings.TrimSpace(v)
	}

	return slice
}

type ureader struct {
	r io.Reader
}

func (c *ureader) Read(buf []byte) (int, error) {
	n, err := c.r.Read(buf)

	// Replace carriage returns with newlines
	for i, b := range buf {
		if b == '\r' {
			buf[i] = '\n'
		}
	}

	return n, err
}

func isEmpty(r []string) bool {
	for _, s := range r {
		if s != "" {
			return false
		}
	}

	return true
}

// ReadRows reads rows from an io.Reader and processes each with the function.
func ReadRows(r io.Reader, f func([]string)) error {
	cr := csv.NewReader(&ureader{r})

	cr.Comment = '#'
	cr.LazyQuotes = true
	cr.FieldsPerRecord = -1
	cr.TrimLeadingSpace = true

	// Skip header.
	_, err := cr.Read()

	if err != nil {
		return err
	}

	for {
		row, err := cr.Read()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		// In-place
		trimSpace(row)

		if isEmpty(row) {
			continue
		}

		f(row)
	}

	return nil
}
