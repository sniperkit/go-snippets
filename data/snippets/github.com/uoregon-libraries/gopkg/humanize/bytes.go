package humanize

import (
	"fmt"
	"strconv"
)

// Values for different suffixes.  We only need up to exabytes, as anything
// larger can't be stored in an int64.  And, at least currently, exabytes are
// absurdly large anyway.
const (
	Kilobyte = 1024
	Megabyte = Kilobyte * 1024
	Gigabyte = Megabyte * 1024
	Terabyte = Gigabyte * 1024
	Petabyte = Terabyte * 1024
	Exabyte  = Petabyte * 1024
)

// Bytes returns a human-friendly value for filesizes
func Bytes(bytes int64) string {
	var divVal int64
	var suffix string
	var decToStr = func(val int64) string { return strconv.FormatInt(val, 10) }

	switch {
	case bytes >= Exabyte:
		divVal = Exabyte
		suffix = " EB"

	case bytes >= Petabyte:
		divVal = Petabyte
		suffix = " PB"

	case bytes >= Terabyte:
		divVal = Terabyte
		suffix = " TB"

	case bytes >= Gigabyte:
		divVal = Gigabyte
		suffix = " GB"

	case bytes >= Megabyte:
		divVal = Megabyte
		suffix = " MB"

	case bytes >= Kilobyte:
		divVal = Kilobyte
		suffix = " KB"

	default:
		// As a special case, when we have no need for division, we just build a
		// simple string inline
		return decToStr(bytes) + " B"
	}

	var whole = bytes / divVal
	var dec = ((bytes % divVal) / (divVal / 1024) * 100) / 1024
	switch {
	case dec == 0:
		return decToStr(whole) + suffix

	case dec < 10:
		return decToStr(whole) + ".0" + decToStr(dec) + suffix

	// This can happen with huge numbers that should be rounded up
	case dec > 99:
		whole++
		return decToStr(whole) + suffix

	default:
		return decToStr(whole) + "." + decToStr(dec) + suffix
	}
}

// bytesSimple isn't intended for use; it's just built to show why we don't use
// the simpler approach of using floats and fmt.Sprintf - benchmarks put this
// function taking 2.5x-7x as long as the above function, where 2.5x is only
// seen when the bytes are < 1024, which in practice is likely to be uncommon.
func bytesSimple(bytes int64) string {
	var val float64
	var size = float64(bytes)
	var suffix string

	switch {
	case bytes >= Exabyte:
		val = size / Exabyte
		suffix = " EB"

	case bytes >= Petabyte:
		val = size / Petabyte
		suffix = " PB"

	case bytes >= Terabyte:
		val = size / Terabyte
		suffix = " TB"

	case bytes >= Gigabyte:
		val = size / Gigabyte
		suffix = " GB"

	case bytes >= Megabyte:
		val = size / Megabyte
		suffix = " MB"

	case bytes >= Kilobyte:
		val = size / Kilobyte
		suffix = " KB"

	default:
		// As a special case, when we have no need for division, we just build a
		// simple string inline
		return fmt.Sprintf("%d B", bytes)
	}

	return fmt.Sprintf("%0.2f%s", val, suffix)
}
