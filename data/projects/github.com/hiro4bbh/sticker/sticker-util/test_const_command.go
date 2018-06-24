package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TestConstCommand have flags for testForest sub-command.
type TestConstCommand struct {
	Help       bool
	Ks         common.OptionUints
	N          uint
	TableNames common.OptionStrings

	Result map[string]interface{}

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTestConstCommand returns a new TestConstCommand.
func NewTestConstCommand(opts *Options) *TestConstCommand {
	return &TestConstCommand{
		Help:       false,
		Ks:         common.OptionUints{true, []uint{1, 3, 5}},
		N:          ^uint(0),
		TableNames: common.OptionStrings{true, []string{"test.txt"}},
		opts:       opts,
	}
}

func (cmd *TestConstCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@testConst", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.Ks, "K", "Specify the top-K values")
	cmd.flagSet.UintVar(&cmd.N, "N", cmd.N, "Specify the maximum number of the tested entries")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TestConstCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run tests the .labelforest model on the specified table of dataset.
func (cmd *TestConstCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TestConstCommands: %#v", cmd)
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, cmd.N, true)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading .labelconst model from %q ...", opts.LabelConst)
	model, err := common.ReadLabelConst(opts.LabelConst)
	if err != nil {
		return err
	}
	reporter := common.NewResultsReporter(ds.Y, cmd.Ks.Values)
	opts.Logger.Printf("predicting top-%d labels ...", reporter.MaxK())
	reporter.ResetTimer()
	reporter.Report(model.PredictAll(ds.X, reporter.MaxK()), opts.OutputWriter)
	return nil
}

// ShowHelp shows the help.
func (cmd *TestConstCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @testConst [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
