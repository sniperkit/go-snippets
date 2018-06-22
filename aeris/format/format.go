package format

import (
	"errors"
	"strings"
)

type Format struct {
	Container string
	Audio     AudioFormat
	Video     VideoFormat
}

type VideoFormat struct {
	Encoding   string
	Resolution string
}

type AudioFormat struct {
	Encoding string
	Bitrate  int
}

func (f *Format) Extension() (string, error) {
	container := strings.ToLower(f.Container)

	switch container {
	case "mp4", "mp3":
		return "." + container, nil
	}

	return "", errors.New("Cannot determine extension for container format " + f.Container)
}
