package main

import (
    "errors"
    "flag"
    "fmt"
    "strings"
    "time"
)

var species = flag.String("species", "gopher", "the species we are studying")

var gopherType string

func init() {
    const (
        defaultGopher = "pocket"
        usage = "the variety of gopher"
    )
    flag.StringVar(&gopherType, "gopher_type", defaultGopher, usage)
    flag.StringVar(&gopherType, "g", defaultGopher, usage + " (shorthand) ")
}

type interval []time.Duration


func (i *interval) String() string {
    return fmt.Sprintf("The duration is %q", *i)
}

func (i *interval) Set(value string) error {
    if len(*i) > 0 {
        return errors.New("interval flag already set")
    }
    for _, dt := range strings.Split(value, ",") {
        duration, err := time.ParseDuration(dt)
        fmt.Printf("The parsed duration is %q", duration)
        if err != nil {
            return err
        }
        *i = append(*i, duration)
    }
    return nil
}

var intervalFlag interval

func init() {
    flag.Var(&intervalFlag, "deltaT", "comma-separated list of intervals to use between events")
}

func main() {
    flag.Parse()
    fmt.Printf("The interval flag is %s\n", intervalFlag)
    fmt.Printf("The duration is %s\n", time.Second)
}
