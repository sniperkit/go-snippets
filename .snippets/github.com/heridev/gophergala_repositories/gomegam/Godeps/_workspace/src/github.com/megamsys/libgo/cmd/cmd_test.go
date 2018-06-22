package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"
	"gopkg.in/check.v1"
)

type recordingExiter int

func (e *recordingExiter) Exit(code int) {
	*e = recordingExiter(code)
}

func (e recordingExiter) value() int {
	return int(e)
}

type TestCommand struct{}

func (c *TestCommand) Info() *Info {
	return &Info{
		Name:  "foo",
		Desc:  "Foo do anything or nothing.",
		Usage: "foo",
	}
}

func (c *TestCommand) Run(context *Context, client *Client) error {
	io.WriteString(context.Stdout, "Running TestCommand")
	return nil
}

type ErrorCommand struct {
	msg string
}

func (c *ErrorCommand) Info() *Info {
	return &Info{Name: "error"}
}

func (c *ErrorCommand) Run(context *Context, client *Client) error {
	return errors.New(c.msg)
}

func (s *S) TestRegister(c *check.C) {
	manager.Register(&TestCommand{})
	badCall := func() { manager.Register(&TestCommand{}) }
	c.Assert(badCall, check.PanicMatches, "command already registered: foo")
}

func (s *S) TestManagerRunShouldWriteErrorsOnStderr(c *check.C) {
	manager.Register(&ErrorCommand{msg: "You are wrong\n"})
	manager.Run([]string{"error"})
	c.Assert(manager.stderr.(*bytes.Buffer).String(), check.Equals, "Error: You are wrong\n")
}

func (s *S) TestRun(c *check.C) {
	manager.Register(&TestCommand{})
	manager.Run([]string{"foo"})
	c.Assert(manager.stdout.(*bytes.Buffer).String(), check.Equals, "Running TestCommand")
}

/*func (s *S) TestFileSystem(c *check.C) {
	fsystem = &testing.RecordingFs{}
	c.Assert(filesystem(), check.DeepEquals, fsystem)
	fsystem = nil
	c.Assert(filesystem(), check.DeepEquals, fs.OsFs{})
}*/

func (s *S) TestHelpCommandShouldBeRegisteredByDefault(c *check.C) {
	var stdout, stderr bytes.Buffer
	m := NewManager("gulpd", "0.1", "", &stdout, &stderr, os.Stdin)
	_, exists := m.Commands["help"]
	c.Assert(exists, check.Equals, true)
}

func (s *S) TestHelpReturnErrorIfTheGivenCommandDoesNotExist(c *check.C) {
	command := help{manager: manager}
	context := Context{[]string{"someone-create"}, manager.stdout, manager.stderr, manager.stdin}
	err := command.Run(&context,nil)
	c.Assert(err, check.NotNil)
	c.Assert(err, check.ErrorMatches, `^command "someone-create" does not exist.$`)
}

func (s *S) TestVersion(c *check.C) {
	var stdout, stderr bytes.Buffer
	manager := NewManager("gulpd", "0.1", "", &stdout, &stderr, os.Stdin)
	command := version{manager: manager}
	context := Context{[]string{}, manager.stdout, manager.stderr, manager.stdin}
	err := command.Run(&context,nil)
	c.Assert(err, check.IsNil)
	c.Assert(manager.stdout.(*bytes.Buffer).String(), check.Equals, "gulpd version 0.1.\n")
}

func (s *S) TestVersionInfo(c *check.C) {
	expected := &Info{
		Name:    "version",
		MinArgs: 0,
		Usage:   "version",
		Desc:    "display the current version",
	}
	c.Assert((&version{}).Info(), check.DeepEquals, expected)
}

type ArgCmd struct{}

func (c *ArgCmd) Info() *Info {
	return &Info{
		Name:    "arg",
		MinArgs: 1,
		MaxArgs: 2,
		Usage:   "arg [args]",
		Desc:    "some desc",
	}
}

func (cmd *ArgCmd) Run(ctx *Context, client *Client) error {
	return nil
}

func (s *S) TestRunWrongArgsNumberShouldRunsHelpAndReturnStatus1(c *check.C) {
	expected := `gulpd version 0.1.

ERROR: wrong number of arguments.

Usage: gulpd arg [args]

some desc

Minimum # of arguments: 1
Maximum # of arguments: 2
`
	manager.Register(&ArgCmd{})
	manager.Run([]string{"arg"})
	c.Assert(manager.stdout.(*bytes.Buffer).String(), check.Equals, expected)
	c.Assert(manager.e.(*recordingExiter).value(), check.Equals, 1)
}

func (s *S) TestHelpShouldReturnUsageWithTheCommandName(c *check.C) {
	expected := `gulpd version 0.1.

Usage: gulpd foo

Foo do anything or nothing.

`
	var stdout, stderr bytes.Buffer
	manager := NewManager("gulpd", "0.1", "", &stdout, &stderr, os.Stdin)
	manager.Register(&TestCommand{})
	context := Context{[]string{"foo"}, manager.stdout, manager.stderr, manager.stdin}
	command := help{manager: manager}
	err := command.Run(&context,nil)
	c.Assert(err, check.IsNil)
	c.Assert(manager.stdout.(*bytes.Buffer).String(), check.Equals, expected)
}

func (s *S) TestExtractProgramNameWithAbsolutePath(c *check.C) {
	got := ExtractProgramName("/home/ram/bin/gulpd")
	c.Assert(got, check.Equals, "gulpd")
}

func (s *S) TestExtractProgramNameWithRelativePath(c *check.C) {
	got := ExtractProgramName("./gulpd")
	c.Assert(got, check.Equals, "gulpd")
}

func (s *S) TestExtractProgramNameWithinThePATH(c *check.C) {
	got := ExtractProgramName("gulpd")
	c.Assert(got, check.Equals, "gulpd")
}

func (s *S) TestFinisherReturnsOsExiterIfNotDefined(c *check.C) {
	m := Manager{}
	c.Assert(m.finisher(), check.FitsTypeOf, osExiter{})
}

func (s *S) TestFinisherReturnTheDefinedE(c *check.C) {
	var exiter recordingExiter
	m := Manager{e: &exiter}
	c.Assert(m.finisher(), check.FitsTypeOf, &exiter)
}
