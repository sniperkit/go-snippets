package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// InspectForestCommand have flags for inspectForest sub-command.
type InspectForestCommand struct {
	Addr string
	Help bool

	opts    *Options
	flagSet *flag.FlagSet
}

// NewInspectForestCommand returns a new InspectForestCommand.
func NewInspectForestCommand(opts *Options) *InspectForestCommand {
	return &InspectForestCommand{
		Addr: ":8080",
		Help: false,
		opts: opts,
	}
}

func (cmd *InspectForestCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@inspectForest", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.StringVar(&cmd.Addr, "addr", cmd.Addr, "Specify the HTTP server address")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *InspectForestCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run inspects the .labelforest file.
func (cmd *InspectForestCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("InspectForestCommands: %#v", cmd)
	opts.Logger.Printf("loading .labelforest model from %q ...\n", opts.LabelForest)
	forest, err := common.ReadLabelForest(opts.LabelForest)
	if err != nil {
		return err
	}
	if err := opts.RunHTTPServer(cmd.Addr, func(handleFunc func(prefix string, pattern string, handler func(writer http.ResponseWriter, req *http.Request))) error {
		handleFunc("RunInspect", "/inspectForest/tree", func(writer http.ResponseWriter, req *http.Request) {
			treeIdStrs, ok := req.URL.Query()["treeId"]
			if !ok {
				treeIdStrs = []string{"0"}
			}
			treeIdStr := treeIdStrs[0]
			treeId, err := strconv.ParseUint(treeIdStr, 10, 64)
			if err != nil {
				writer.Write([]byte(fmt.Sprintf("treeId: strconv.ParseUint: %s", err)))
				return
			}
			tmpl := template.New("")
			tmpl.Funcs(TemplateStandardFunctions)
			tmplFilename := filepath.Join(opts.HTTPResource, "inspectForestTemplate.html")
			if _, err := tmpl.ParseFiles(tmplFilename); err != nil {
				writer.Write([]byte(fmt.Sprintf("template.ParseFiles: %s: %s", tmplFilename, err)))
				return
			}
			if err := tmpl.ExecuteTemplate(writer, "inspectForestTemplate.html", map[string]interface{}{
				"featureMap": opts.featureMap,
				"labelMap":   opts.labelMap,
				"filename":   opts.LabelForest,
				"forest":     forest,
				"treeId":     treeId,
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
func (cmd *InspectForestCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @inspectForest [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
