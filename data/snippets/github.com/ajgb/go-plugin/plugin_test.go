package plugin

import (
	"bytes"
	"math"
	"testing"
)

type ExitHandler struct {
	code   Status
	args   []string
	output bytes.Buffer
	length int
}

func (w *ExitHandler) Write(p []byte) (int, error) {
	w.output.Write(p)
	w.length += len(p)
	return len(p), nil
}

var exitHandler *ExitHandler

func TestFloatConversion(t *testing.T) {
	tests := []struct {
		in        interface{}
		precision float64

		value float64
		err   error
	}{
		{1, 0.01, 1.0, nil},
		{float32(1.0001), 0.00001, 1.0001, nil},
		{"1", 0.01, 1.0, nil},
		{"1.00000001", 0.000000001, 1.00000001, nil},
		{"1487801591.176291383", 0.000000001, 1487801591.176291383, nil},
	}

	for _, test := range tests {
		got, err := i2f(test.in)
		if test.err != err {
			t.Errorf("Got err: %s, expected %s", err, test.err)
		}
		if !floatNearlyEqual(got, test.value, test.precision) {
			t.Errorf("Got value: %v, expected: %v", got, test.value)
		}
	}
}

type FormatArgs struct {
	format string
	params []interface{}
}

type AddMessageOutputTest struct {
	name             string
	version          string
	messages         []FormatArgs
	separator        string
	expectedExitCode Status
	expectedOutput   string
}

func TestAddMessage(t *testing.T) {
	tests := []AddMessageOutputTest{
		{
			"check_plugin", "v1.0",
			nil, "",
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All ok", nil},
			}, "",
			OK, "OK: All ok\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All ok", nil},
			}, "=",
			OK, "OK: All ok\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"% All ok %", nil},
			}, "=",
			OK, "OK: % All ok %\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All ok", nil},
				{"Nothing to see", nil},
			}, "",
			OK, "OK: All ok, Nothing to see\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All ok", nil},
				{"Nothing to see", nil},
			}, ":",
			OK, "OK: All ok:Nothing to see\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All %s => %d", []interface{}{"ok", 123}},
			}, "",
			OK, "OK: All ok => 123\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All %s => %d", []interface{}{"ok", 123}},
			}, "=",
			OK, "OK: All ok => 123\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All %s => %d", []interface{}{"ok", 123}},
				{"Nothing to %s => %d", []interface{}{"see", 456}},
			}, "",
			OK, "OK: All ok => 123, Nothing to see => 456\n",
		},
		{
			"check_plugin", "v1.0",
			[]FormatArgs{
				{"All %s => %d", []interface{}{"ok", 123}},
				{"Nothing to %s => %d", []interface{}{"see", 456}},
			}, ":",
			OK, "OK: All ok => 123:Nothing to see => 456\n",
		},
	}

	for _, test := range tests {
		exitHandler := initExitHandler()

		func() {
			check := New(test.name, test.version)
			defer check.Final()
			for _, m := range test.messages {
				check.AddMessage(m.format, m.params...)
			}
			if len(test.separator) > 0 {
				check.MessageSeparator = test.separator
			}
		}()

		gotOutput := exitHandler.output.String()
		if gotOutput != test.expectedOutput {
			t.Errorf("Got output: '%s', expected: '%s'", gotOutput, test.expectedOutput)
		}

		if exitHandler.code != test.expectedExitCode {
			t.Errorf("Got code: %d, expected: %d", exitHandler.code, test.expectedExitCode)
		}
	}
}

type ResultFormatArgs struct {
	result    Status
	newStatus Status
	format    string
	params    []interface{}
}

type AddResultOutputTest struct {
	name             string
	version          string
	results          []ResultFormatArgs
	expectedExitCode Status
	expectedOutput   string
}

func TestAddResult(t *testing.T) {
	tests := []AddResultOutputTest{
		{
			"check_plugin", "v1.0",
			nil,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]ResultFormatArgs{
				{OK, OK, "All ok", nil},
			},
			OK, "OK: All ok\n",
		},
		{
			"check_plugin", "v1.0",
			[]ResultFormatArgs{
				{OK, OK, "% All ok %", nil},
			},
			OK, "OK: % All ok %\n",
		},
		{
			"check_plugin", "v1.0",
			[]ResultFormatArgs{
				{OK, OK, "All ok", nil},
				{OK, OK, "Still ok", nil},
			},
			OK, "OK: All ok, Still ok\n",
		},
		{
			"check_plugin", "v1.0",
			[]ResultFormatArgs{
				{OK, OK, "All ok", nil},
				{WARNING, WARNING, "Warning! fault detected", nil},
				{OK, WARNING, "Still ok?", nil},
			},
			WARNING, "WARNING: All ok, Warning! fault detected, Still ok?\n",
		},
		{
			"check_plugin", "v1.0",
			[]ResultFormatArgs{
				{OK, OK, "All ok", nil},
				{OK, OK, "Still ok...", nil},
				{CRITICAL, CRITICAL, "Serious fault detected!", nil},
			},
			CRITICAL, "CRITICAL: All ok, Still ok..., Serious fault detected!\n",
		},
	}

	for _, test := range tests {
		exitHandler := initExitHandler()

		check := New(test.name, test.version)
		for _, r := range test.results {
			check.AddResult(r.result, r.format, r.params...)
			if r.newStatus != check.Status() {
				t.Errorf("Got updated status: '%s', expected: '%s'", check.Status(), r.newStatus)
			}
		}
		check.Final()

		gotOutput := exitHandler.output.String()
		if gotOutput != test.expectedOutput {
			t.Errorf("Got output: '%s', expected: '%s'", gotOutput, test.expectedOutput)
		}

		if exitHandler.code != test.expectedExitCode {
			t.Errorf("Got code: %d, expected: %d", exitHandler.code, test.expectedExitCode)
		}
	}
}

type ExitHelpersTest struct {
	name             string
	version          string
	method           func(*Plugin, string, ...interface{})
	message          FormatArgs
	expectedExitCode Status
	expectedOutput   string
}

func TestExitHelpers(t *testing.T) {
	tests := []ExitHelpersTest{
		{
			"check_plugin", "v1.0",
			(*Plugin).ExitOK,
			FormatArgs{"with percent: %", nil},
			OK, "OK: with percent: %\n",
		},
		{
			"check_plugin", "v1.0",
			(*Plugin).ExitOK,
			FormatArgs{"with exit code: %d", []interface{}{0}},
			OK, "OK: with exit code: 0\n",
		},
		{
			"check_plugin", "v1.0",
			(*Plugin).ExitWarning,
			FormatArgs{"with exit code: %d", []interface{}{1}},
			WARNING, "WARNING: with exit code: 1\n",
		},
		{
			"check_plugin", "v1.0",
			(*Plugin).ExitCritical,
			FormatArgs{"with exit code: %d", []interface{}{2}},
			CRITICAL, "CRITICAL: with exit code: 2\n",
		},
		{
			"check_plugin", "v1.0",
			(*Plugin).ExitUnknown,
			FormatArgs{"with exit code: %d", []interface{}{3}},
			3, "UNKNOWN: with exit code: 3\n",
		},
	}

	for _, test := range tests {
		exitHandler := initExitHandler()

		check := New(test.name, test.version)
		test.method(check, test.message.format, test.message.params...)

		gotOutput := exitHandler.output.String()
		if gotOutput != test.expectedOutput {
			t.Errorf("Got output: '%s', expected: '%s'", gotOutput, test.expectedOutput)
		}

		if exitHandler.code != test.expectedExitCode {
			t.Errorf("Got code: %d, expected: %d", exitHandler.code, test.expectedExitCode)
		}
	}
}

type MetricArgs struct {
	name             string
	value            interface{}
	uomAndThresholds []string
	err              string
}

type AddMetricOutputTest struct {
	name             string
	version          string
	metrics          []MetricArgs
	includeAll       bool
	expectedExitCode Status
	expectedOutput   string
}

func TestAddMetric(t *testing.T) {
	tests := []AddMetricOutputTest{
		{
			"check_plugin", "v1.0",
			nil, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, nil, ""},
			}, false,
			OK, "OK: | m1=123.456;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"white space", 123.456, nil, ""},
			}, false,
			OK, "OK: | 'white space'=123.456;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", "abc", nil, "Invalid value of m1: abc"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, nil, ""},
				{"m2", 456.789, nil, ""},
			}, false,
			OK, "OK: | m1=123.456;;;; m2=456.789;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, nil, ""},
				{"m2", 456.789, nil, ""},
				{"m1", 0.0, nil, "Duplicated metric m1"},
			}, false,
			OK, "OK: | m1=123.456;;;; m2=456.789;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"MB"}, ""},
			}, true,
			OK, "OK: m1 is 123.456MB | m1=123.456MB;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"percent", 23.45, []string{"%"}, ""},
			}, true,
			OK, "OK: percent is 23.45% | percent=23.45%;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"MB", ""}, ""},
			}, true,
			OK, "OK: m1 is 123.456MB | m1=123.456MB;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"MB", "1", "2", "3"}, "Too many arguments"},
			}, true,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "100"}, ""},
			}, false,
			WARNING, "WARNING: m1 is 123.456 (outside 100) | m1=123.456;100;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "c"}, "Invalid format of warning threshold m1: c"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"TB", "100", "123"}, ""},
			}, false,
			CRITICAL, "CRITICAL: m1 is 123.456TB (outside 123) | m1=123.456TB;100;123;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "0:100"}, ""},
			}, false,
			WARNING, "WARNING: m1 is 123.456 (outside 0:100) | m1=123.456;0:100;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "2000:100"}, "Invalid format of warning threshold m1: 2000:100"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "20:100:200"}, "Invalid format of warning threshold m1: 20:100:200"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "200:"}, ""},
			}, false,
			WARNING, "WARNING: m1 is 123.456 (outside 200:) | m1=123.456;200:;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "f:"}, "Invalid format of warning threshold m1: f:"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "100:g"}, "Invalid format of warning threshold m1: 100:g"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"", "b:100"}, "Invalid format of warning threshold m1: b:100"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"TB", "80:100", "~:123"}, ""},
			}, false,
			CRITICAL, "CRITICAL: m1 is 123.456TB (outside ~:123) | m1=123.456TB;80:100;~:123;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"TB", "80:100", "~:d"}, "Invalid format of critical threshold m1: ~:d"},
			}, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"TB", "100", "@123"}, ""},
			}, false,
			WARNING, "WARNING: m1 is 123.456TB (outside 100) | m1=123.456TB;100;@123;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, []string{"TB", "100", "@200"}, ""},
			}, false,
			CRITICAL, "CRITICAL: m1 is 123.456TB (inside @200) | m1=123.456TB;100;@200;;\n",
		},
	}

	for _, test := range tests {
		exitHandler := initExitHandler()

		func() {
			check := New(test.name, test.version)
			defer check.Final()
			if test.includeAll {
				check.AllMetricsInOutput = true
			}
			for _, m := range test.metrics {
				metricErr := check.AddMetric(m.name, m.value, m.uomAndThresholds...)
				if m.err != "" && metricErr.Error() != m.err {
					t.Errorf("Got error: '%s', expected: '%s'", metricErr, m.err)
				}
			}
		}()

		gotOutput := exitHandler.output.String()
		if gotOutput != test.expectedOutput {
			t.Errorf("Got output: '%s', expected: '%s'", gotOutput, test.expectedOutput)
		}

		if exitHandler.code != test.expectedExitCode {
			t.Errorf("Got code: %d, expected: %d", exitHandler.code, test.expectedExitCode)
		}
	}
}

type FinalOutputTest struct {
	name             string
	version          string
	metrics          []MetricArgs
	throwException   bool
	expectedExitCode Status
	expectedOutput   string
}

func TestFinal(t *testing.T) {
	tests := []FinalOutputTest{
		{
			"check_plugin", "v1.0",
			nil, false,
			OK, "OK:\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, nil, ""},
			},
			false,
			OK, "OK: | m1=123.456;;;;\n",
		},
		{
			"check_plugin", "v1.0",
			[]MetricArgs{
				{"m1", 123.456, nil, ""},
			},
			true,
			CRITICAL, "CRITICAL: check_plugin panic: Forced exception\n",
		},
	}

	for _, test := range tests {
		exitHandler := initExitHandler()

		func() {
			check := New(test.name, test.version)
			defer check.Final()
			for _, m := range test.metrics {
				metricErr := check.AddMetric(m.name, m.value, m.uomAndThresholds...)
				if m.err != "" && metricErr.Error() != m.err {
					t.Errorf("Got error: '%s', expected: '%s'", metricErr, m.err)
				}
			}
			if test.throwException {
				panic("Forced exception")
			}
		}()

		gotOutput := exitHandler.output.String()
		if gotOutput != test.expectedOutput {
			t.Errorf("Got output: '%s', expected: '%s'", gotOutput, test.expectedOutput)
		}

		if exitHandler.code != test.expectedExitCode {
			t.Errorf("Got code: %d, expected: %d", exitHandler.code, test.expectedExitCode)
		}
	}
}

var OptionParseTest struct {
	Hostname string `short:"H" long:"hostname" description:"Hostname"`
}

type ParseArgsTest struct {
	name             string
	version          string
	args             []string
	preamble         string
	description      string
	expectedHostname string
	expectedExitCode Status
	expectedOutput   string
}

func TestParseArgs(t *testing.T) {
	tests := []ParseArgsTest{
		{
			"check_plugin", "v1.0",
			[]string{"-H", "localhost", "-h"},
			"",
			"",
			"localhost",
			UNKNOWN, `check_plugin v1.0
Usage:
  go-plugin.test [OPTIONS]

Plugin Options:
  -H, --hostname= Hostname

Default Options:
  -h, --help      Show this help message

`,
		},
		{
			"check_plugin", "v1.0",
			[]string{"-H", "localhost", "-h"},
			"Test output",
			"Description:\n123",
			"localhost",
			UNKNOWN, `check_plugin v1.0
Test output
Usage:
  go-plugin.test [OPTIONS]

Plugin Options:
  -H, --hostname= Hostname (default: localhost)

Default Options:
  -h, --help      Show this help message

Description:
123
`,
		},
	}

	for _, test := range tests {
		exitHandler := initExitHandler(test.args)

		check := New(test.name, test.version)
		if len(test.preamble) > 0 {
			check.Preamble = test.preamble
		}
		if len(test.description) > 0 {
			check.Description = test.description
		}
		check.ParseArgs(&OptionParseTest)

		if OptionParseTest.Hostname != test.expectedHostname {
			t.Errorf("Got hostname: '%s', expected: '%s'", OptionParseTest.Hostname, test.expectedHostname)
		}

		gotOutput := exitHandler.output.String()
		if gotOutput != test.expectedOutput {
			t.Errorf("Got output: '%s', expected: '%s'", gotOutput, test.expectedOutput)
		}

		if exitHandler.code != test.expectedExitCode {
			t.Errorf("Got code: %d, expected: %d", exitHandler.code, test.expectedExitCode)
		}
	}
}

func initExitHandler(args ...[]string) (exitHandler *ExitHandler) {
	exitHandler = &ExitHandler{}
	pOsExit = func(code Status) { exitHandler.code = code }
	pOutputHandle = exitHandler
	if len(args) > 0 {
		pArgs = args[0]
	}

	return
}

func floatNearlyEqual(x, y, epsilon float64) bool {
	absX := math.Abs(x)
	absY := math.Abs(y)
	diff := math.Abs(x - y)

	if x == y {
		return true
	} else if x == 0 || y == 0 || diff < math.SmallestNonzeroFloat64 {
		return diff < (epsilon * math.SmallestNonzeroFloat64)
	}

	return diff/math.Min((absX+absY), math.MaxFloat64) < epsilon
}
