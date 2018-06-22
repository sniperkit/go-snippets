package crawler

import "strconv"

func ParseSize(size string) int64 {
	sfloat, err := strconv.ParseFloat(size[:len(size)-1], 64)
	if err != nil {
		return -1
	}
	s := int64(sfloat)
	suffix := size[len(size)-1]

	switch suffix {
	case 'G':
		return s * 1024 * 1024 * 1024
	case 'M':
		return s * 1024 * 1024
	case 'K':
		return s * 1024
	}
	return s
}
