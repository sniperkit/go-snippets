package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// InspectOneCommand have flags for inspectForest sub-command.
type InspectOneCommand struct {
	Addr string
	Help bool

	opts    *Options
	flagSet *flag.FlagSet
}

// NewInspectOneCommand returns a new InspectOneCommand.
func NewInspectOneCommand(opts *Options) *InspectOneCommand {
	return &InspectOneCommand{
		Addr: ":8080",
		Help: false,
		opts: opts,
	}
}

func (cmd *InspectOneCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@inspectOne", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.StringVar(&cmd.Addr, "addr", cmd.Addr, "Specify the HTTP server address")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *InspectOneCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run inspects the .labelboost file.
func (cmd *InspectOneCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("InspectOneCommands: %#v", cmd)
	opts.Logger.Printf("loading .labelone model from %q ...\n", opts.LabelOne)
	model, err := common.ReadLabelOne(opts.LabelOne)
	if err != nil {
		return err
	}
	results := make([]interface{}, 0, len(opts.TestOnes))
	for _, testOne := range opts.TestOnes {
		results = append(results, map[string]interface{}{
			"tableName": common.JoinTableNames(testOne.TableNames.Values),
			"result":    testOne.Result,
		})
	}
	if err := opts.RunHTTPServer(cmd.Addr, func(handleFunc func(prefix string, pattern string, handler func(writer http.ResponseWriter, req *http.Request))) error {
		handleFunc("RunInspect", "/inspectOne", func(writer http.ResponseWriter, req *http.Request) {
			tmpl := template.New("")
			tmpl.Funcs(TemplateStandardFunctions)
			tmplFilename := filepath.Join(opts.HTTPResource, "inspectOneTemplate.html")
			if _, err := tmpl.ParseFiles(tmplFilename); err != nil {
				writer.Write([]byte(fmt.Sprintf("template.ParseFiles: %s: %s", tmplFilename, err)))
				return
			}
			if err := tmpl.ExecuteTemplate(writer, "inspectOneTemplate.html", map[string]interface{}{
				"featureMap":  opts.featureMap,
				"labelMap":    opts.labelMap,
				"filename":    opts.LabelOne,
				"model":       model,
				"testResults": results,
			}); err != nil {
				writer.Write([]byte(fmt.Sprintf("template.ExecuteTemplate: %s", err)))
				return
			}
		})
		return nil
	}); err != nil {
		return fmt.Errorf("RunInspect: RunHTTPServer: %s", err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *InspectOneCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @inspectOne [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
