package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var (
	gopherGalaFlag  = flag.Bool("gophergala", false, "-gophergala")
	fromSecondsFlag = flag.Bool("s", false, "-s")
	fromMinutesFlag = flag.Bool("m", false, "-m")
	fromHoursFlag   = flag.Bool("h", false, "-h")
	fromDaysFlag    = flag.Bool("d", false, "-d")
)

func main() {
	flag.Parse()
	var deadline string
	if *gopherGalaFlag {
		deadline = "2015-01-26T01:00:00+00:00"
	} else {
		deadline = flag.Arg(0)
	}
	date, err := time.Parse(time.RFC3339, deadline)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: %T{%v}\n", err, err)
		fmt.Fprintf(os.Stdout, "Deadline has to be in RFC3339 format, i.e.: '%s'\n", time.RFC3339)
		os.Exit(1)
	}
	dformat := getFormat()
	tick(os.Stdout, date, dformat)
}

func tick(writer io.Writer, date time.Time, dformat durationFormat) {
	ticker := time.Tick(1 * time.Second)
	for t := range ticker {
		remaining := date.Sub(t)
		if remaining.Seconds() < 0.0 {
			fmt.Fprintf(writer, "\n>>> Deadline reached. Exiting.\n")
			os.Exit(0)
		}
		timeString := format(dformat, remaining)
		template := fmt.Sprintf("%%%ds\r>>> T minus %%s", len(timeString))
		fmt.Fprintf(writer, template, "", timeString)
	}
}

func getFormat() durationFormat {
	if *fromSecondsFlag {
		return fromSeconds
	}
	if *fromMinutesFlag {
		return fromMinutes
	}
	if *fromHoursFlag {
		return fromHours
	}
	if *fromDaysFlag {
		return fromDays
	}
	return Full
}

type durationFormat int

const (
	fromSeconds durationFormat = iota << 1
	fromMinutes
	fromHours
	fromDays
	Full
)

const (
	Day  = 24 * time.Hour
	Year = 365 * Day
)

func format(typ durationFormat, d time.Duration) string {
	format := func(x time.Duration, typ string) string {
		plural := "s"
		if x == 1 {
			plural = ""
		}
		return fmt.Sprintf("%d %s%s", x, typ, plural)
	}
	var segments []string
	seconds := (d / time.Second) % 60
	if typ == fromSeconds {
		segments = append(segments, format((d/time.Second), "second"))
		return reverseJoin(segments, ", ")
	} else {
		segments = append(segments, format(seconds, "second"))
	}
	minutes := (d / time.Minute) % 60
	if typ == fromMinutes {
		segments = append(segments, format((d/time.Minute), "minute"))
		return reverseJoin(segments, ", ")
	} else {
		segments = append(segments, format(minutes, "minute"))
	}
	if minutes != (d / time.Minute) {
		hours := (d / time.Hour) % 24
		if typ == fromHours {
			segments = append(segments, format((d/time.Hour), "hour"))
			return reverseJoin(segments, ", ")
		} else {
			segments = append(segments, format(hours, "hour"))
		}
		if hours != (d / time.Hour) {
			days := (d / Day) % 365
			if typ == fromDays {
				segments = append(segments, format((d/Day), "day"))
				return reverseJoin(segments, ", ")
			} else {
				segments = append(segments, format(days, "day"))
			}
			if days != (d / Day) {
				years := d / Year
				segments = append(segments, format(years, "year"))
			}
		}
	}
	return reverseJoin(segments, ", ")
}

func reverseJoin(array []string, separator string) string {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}
	return strings.Join(array, ", ")
}
