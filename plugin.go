/*

Package plugin provides a Go library which helps writing monitoring plugins.

Features:

    * Setting and appending check messages
    * Aggregation of results
    * Thresholds and ranges for metrics with breaches reported
    * Exit shortcut helper methods
    * Provides extensive command line options parser

Example usage:

    package main

    // import plugin library
    import (
      "github.com/ajgb/go-plugin"
    )

    // define command line options
    var opts struct {
      Hostname string `short:"H" long:"hostname" description:"Host" required:"true"`
      Port     int    `short:"p" long:"port" description:"Port" required:"true"`
    }

    func main() {
      // initialize check
      check := plugin.New("check_service", "v1.0.0")

      // parse command line arguments
      if err := check.ParseArgs(&opts); err != nil {
        check.ExitCritical("Error parsing arguments: %s", err)
      }

      // return result on exit
      defer check.Final()

      // add service data to output
      check.AddMessage("Service %s:%d", opts.Hostname, opts.Port)

      // gather metrics - provided by some function
      serviceMetrics, err := .... // retrieve metrics
      if err != nil {
        check.ExitCritical("Connection failed: %s", err)
      }

      // add metrics to output
      for metric, value := range serviceMetrics {
        check.AddMetric(metric, value)
      }
    }

*/
package plugin

import (
	"bytes"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Plugin represents the check - its name, version and help messages. It also
// stores the check status, messages and metrics data.
type Plugin struct {
	status   Status
	messages []string
	metrics  checkMetrics
	// Plugin name
	Name string
	// Plugin version
	Version string
	// Preamble displayed in help output before flags usage
	Preamble string
	// Plugin description displayed in help after flags usage
	Description string
	// If true all metrics will be added to check message
	AllMetricsInOutput bool
	// Messages separator, default: ", "
	MessageSeparator string
}

type checkMetric struct {
	value    interface{}
	status   Status
	uom      string
	warn     string
	critical string
}

type checkMetrics map[string]*checkMetric

var pOsExit = func(code Status) { os.Exit(code.ExitCode()) }
var pOutputHandle io.Writer = os.Stdout
var pArgs = os.Args[1:]

/*
New creates a new plugin instance.

	check := plugin.New("check_service", "v1.0.0")

*/
func New(name, version string) *Plugin {
	return &Plugin{
		status:             OK,
		messages:           make([]string, 0),
		metrics:            make(checkMetrics),
		Name:               name,
		Version:            version,
		AllMetricsInOutput: false,
		MessageSeparator:   ", ",
	}
}

/*
AddMetric adds new metric to check's performance data, with name and value
parameters required. The optional string arguments include (in order):
uom (unit of measurement), warning threshold, critical threshold - for
details see Monitoring Plugins Development Guidelines.
Note: Metrics names have to be unique.

    // basic usage - add metric with value
    check.AddMetric("load5", 0.98)

    // metric with UOM
    check.AddMetric("tmp", 15789, "MB")

    // metric with warning threshold (without uom)
    check.AddMetric("rtmax", 28.723, "", 75)

    // metric with warning & critical thresholds (with uom)
    check.AddMetric("rta", 24.558, "ms", 50, 100)

*/
func (p *Plugin) AddMetric(name string, value interface{}, args ...string) error {
	argsCount := len(args)

	metric := &checkMetric{}

	if strings.ContainsRune(name, ' ') && !strings.HasPrefix(name, "'") {
		name = "'" + name + "'"
	}
	if _, ok := p.metrics[name]; ok {
		return fmt.Errorf("Duplicated metric %s", name)
	}

	metric.value = value
	if argsCount >= 1 {
		metric.uom = args[0]
	}

	val, err := i2f(value)
	if err != nil {
		return fmt.Errorf("Invalid value of %s: %v", name, value)
	}

	var alertMessage string

	if argsCount == 2 || argsCount == 3 {
		var thresholdBreached bool
		for i, a := range args[1:] {
			var thresholdName string
			var invert bool

			if len(a) == 0 {
				continue
			}

			arg := strings.TrimPrefix(a, "@")
			if a != arg {
				invert = true
			}

			thresh := strings.Split(arg, ":")

			switch i {
			case 0:
				thresholdName = "warning"
				metric.warn = a
			case 1:
				thresholdName = "critical"
				metric.critical = a
			}

			switch len(thresh) {
			case 1:
				// v < X
				tMax, err := strconv.ParseFloat(thresh[0], 64)
				if err != nil {
					return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
				}
				thresholdBreached = val < 0 || val > tMax
			case 2:
				switch {
				case thresh[0] == "~":
					tMax, err := strconv.ParseFloat(thresh[1], 64)
					if err != nil {
						return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
					}
					thresholdBreached = val > tMax
				case thresh[1] == "":
					tMin, err := strconv.ParseFloat(thresh[0], 64)
					if err != nil {
						return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
					}
					thresholdBreached = val < tMin
				default:
					tMin, err := strconv.ParseFloat(thresh[0], 64)
					if err != nil {
						return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
					}
					tMax, err := strconv.ParseFloat(thresh[1], 64)
					if err != nil {
						return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
					}
					if tMin > tMax {
						return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
					}
					thresholdBreached = val < tMin || val > tMax
				}
			default:
				return fmt.Errorf("Invalid format of %s threshold %s: %s", thresholdName, name, a)
			}

			if invert {
				thresholdBreached = !thresholdBreached
			}

			if thresholdBreached {
				metric.status = Status(i + 1) // i=0 warning, i=1 critical
				if invert {
					alertMessage = fmt.Sprintf("%s is %v%s (inside %s)", name, value, metric.uom, a)
				} else {
					alertMessage = fmt.Sprintf("%s is %v%s (outside %s)", name, value, metric.uom, a)
				}
			}

		}
	} else if argsCount > 3 {
		return fmt.Errorf("Too many arguments")
	}

	if len(alertMessage) > 0 {
		p.AddMessage(alertMessage)
	} else if p.AllMetricsInOutput {
		p.AddMessage(fmt.Sprintf("%s is %v%s", name, value, metric.uom))
	}

	p.metrics[name] = metric
	p.UpdateStatus(metric.status)
	return nil
}

/*
AddMessage appends message to check output.

    check.AddMessage("Server %s", opts.Hostname)

*/
func (p *Plugin) AddMessage(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = fmt.Sprint(format)
	}
	p.messages = append(p.messages, msg)
}

/*
AddResult aggregates results and appends message to check output - the worst
result is final.

    // would not change the final result
    check.AddResult(plugin.OK, "Server %s", opts.Hostname)

    // increases to WARNING level
    if opts.SkipSSLChecks {
        check.AddResult(plugin.WARNING, "Skiping SSL Certificate checks")
    }

*/
func (p *Plugin) AddResult(code Status, format string, args ...interface{}) {
	p.UpdateStatus(code)
	p.AddMessage(format, args...)
}

/*
Final calculates the final check output and exit status.

    check := plugin.New("check_service, "v1.0.0")
    // make sure Final() is called
    defer check.Final()

*/
func (p *Plugin) Final() {
	if r := recover(); r != nil {
		p.ExitCritical("%s panic: %v", p.Name, r)
		return // for testing only as it overrides the os.Exit
	}
	fmt.Fprintf(pOutputHandle, "%s:", p.status.String())
	if len(p.messages) > 0 {
		fmt.Fprintf(pOutputHandle, " ")
		fmt.Fprint(pOutputHandle, strings.Join(p.messages, p.MessageSeparator))
	}
	if len(p.metrics) > 0 {
		var sorted []string
		sorted = make([]string, 0, len(p.metrics))

		fmt.Fprintf(pOutputHandle, " |")
		for k := range p.metrics {
			sorted = append(sorted, k)
		}
		sort.Strings(sorted)
		for _, k := range sorted {
			fmt.Fprintf(pOutputHandle, " %s=%v%s;%s;%s;;",
				k,
				p.metrics[k].value,
				p.metrics[k].uom,
				p.metrics[k].warn,
				p.metrics[k].critical,
			)
		}
	}
	fmt.Fprintf(pOutputHandle, "\n")
	pOsExit(p.status)
}

/*
SetMessage replaces accumulated messages with new one provided.

    check.SetMessage("%s", opts.Hostname)

*/
func (p *Plugin) SetMessage(format string, args ...interface{}) {
	p.messages = []string{}
	p.AddMessage(format, args...)
}

func (p *Plugin) exit(code Status, format string, args ...interface{}) {
	p.status = code
	p.SetMessage(format, args...)
	p.metrics = make(checkMetrics)
	p.Final()
}

/*
ExitOK exits with specified message and OK exit status.
Note: existing messages and metrics are discarded.

    check.ExitOK("Test mode | metric1=1.1; metric2=2.2")

*/
func (p *Plugin) ExitOK(format string, args ...interface{}) {
	p.exit(OK, format, args...)
}

// ExitUnknown exits with specified message and UNKNOWN exit status.
// Note: existing messages and metrics are discarded.
func (p *Plugin) ExitUnknown(format string, args ...interface{}) {
	p.exit(UNKNOWN, format, args...)
}

// ExitWarning exits with specified message and WARNING exit status.
// Note: existing messages and metrics are discarded.
func (p *Plugin) ExitWarning(format string, args ...interface{}) {
	p.exit(WARNING, format, args...)
}

// ExitCritical exits with specified message and CRITICAL exit status.
// Note: existing messages and metrics are discarded.
func (p *Plugin) ExitCritical(format string, args ...interface{}) {
	p.exit(CRITICAL, format, args...)
}

/*
ParseArgs parses the command line options using flags parsing library
providing handling of short/long names, flags and lists, and default and
required options. For details please see https://godoc.org/github.com/jessevdk/go-flags.
Note: -h/--help is automatically added

	if err := check.ParseArgs(&opts); err != nil {
		check.ExitCritical("Error parsing arguments: %s", err)
	}

*/
func (p *Plugin) ParseArgs(opts interface{}) error {
	var err error

	var builtin struct {
		Help bool `short:"h" long:"help" description:"Show this help message"`
	}
	parser := flags.NewParser(opts, 0)
	_, err = parser.AddGroup("Default Options", "", &builtin)

	g := parser.Command.Group.Find("Application Options")
	if g != nil {
		g.ShortDescription = "Plugin Options"
	}

	_, err = parser.ParseArgs(pArgs)

	if builtin.Help {
		fmt.Fprintf(pOutputHandle, "%s v%s\n", p.Name, strings.TrimPrefix(p.Version, "v"))
		if len(p.Preamble) > 0 {
			fmt.Fprintln(pOutputHandle, p.Preamble)
		}
		parser.Options = flags.HelpFlag
		var b bytes.Buffer
		parser.WriteHelp(&b)
		fmt.Fprintln(pOutputHandle, b.String())

		if len(p.Description) > 0 {
			fmt.Fprintln(pOutputHandle, p.Description)
		}
		pOsExit(UNKNOWN)
	}

	return err
}

/*
UpdateStatus updates final exit status if the provided value is higher
(worse) then the current Status.

    // keep people awake
    if rand.Intn(100) % 3 == 0 {
        check.UpdateStatus(plugin.CRITICAL)
    }

*/
func (p *Plugin) UpdateStatus(status Status) {
	if int(status) > int(p.status) {
		p.status = status
	}
}

/*
Status returns current status.

    fmt.Printf("Status after first five metrics: %s\n", check.Status)

*/
func (p *Plugin) Status() Status {
	return p.status
}

func i2f(v interface{}) (float64, error) {
	var f float64
	var err error

	switch v.(type) {
	case float32:
		f = float64(v.(float32))
	case float64:
		f = v.(float64)
	default:
		f, err = strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	}
	return f, err
}
