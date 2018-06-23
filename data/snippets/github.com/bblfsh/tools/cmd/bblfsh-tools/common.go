package main

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/bblfsh/tools"

	"github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/bblfsh/sdk.v1/protocol"
	"gopkg.in/bblfsh/sdk.v1/uast"
	"gopkg.in/src-d/go-errors.v0"
)

var (
	ErrParserFatal = errors.NewKind("Fatal response from parser: %s")
	ErrParserError = errors.NewKind("Error response from parser: %s")
)

type Common struct {
	Address  string `long:"address" description:"server adress to connect to" default:"localhost:9432"`
	Language string `long:"language" description:"language of the input" default:""`
	Args     struct {
		File string `positional-arg-name:"file" required:"true"`
	} `positional-args:"yes"`
}

func (c *Common) execute(args []string, tool tools.Tooler) error {
	logrus.Debugf("executing command")

	request, err := c.buildRequest()
	if err != nil {
		return err
	}

	uast, err := c.parseRequest(request)
	if err != nil {
		return err
	}

	return tool.Exec(uast)
}

func (c *Common) buildRequest() (*protocol.ParseRequest, error) {
	logrus.Debugf("reading file %s", c.Args.File)
	content, err := ioutil.ReadFile(c.Args.File)
	if err != nil {
		return nil, err
	}

	request := &protocol.ParseRequest{
		Filename: filepath.Base(c.Args.File),
		Language: c.Language,
		Content:  string(content),
	}
	return request, nil
}

func (c *Common) parseRequest(request *protocol.ParseRequest) (*uast.Node, error) {
	logrus.Debugf("dialing request at %s", c.Address)
	connection, err := grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := protocol.NewProtocolServiceClient(connection)
	response, err := client.Parse(context.TODO(), request)
	if err != nil {
		return nil, err
	}
	switch response.Status {
	case protocol.Fatal:
		return nil, ErrParserFatal.New(strings.Join(response.Errors, "\n"))
	case protocol.Error:
		return nil, ErrParserError.New(strings.Join(response.Errors, "\n"))
	default:
		return response.UAST, nil
	}
}
