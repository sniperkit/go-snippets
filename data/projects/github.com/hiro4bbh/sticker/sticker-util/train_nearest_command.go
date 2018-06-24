package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TrainNearestCommand have flags for trainBoost sub-command.
type TrainNearestCommand struct {
	Help       bool
	TableNames common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTrainNearestCommand returns a new TrainNearestCommand.
func NewTrainNearestCommand(opts *Options) *TrainNearestCommand {
	return &TrainNearestCommand{
		Help:       false,
		TableNames: common.OptionStrings{true, []string{"train.txt"}},
		opts:       opts,
	}
}

func (cmd *TrainNearestCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@trainNearest", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TrainNearestCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run trains on the specified table of dataset.
func (cmd *TrainNearestCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TrainNearestCommands: %#v", cmd)
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, ^uint(0), false)
	if err != nil {
		return err
	}
	model, err := sticker.TrainLabelNearest(ds, opts.DebugLogger)
	if err != nil {
		return err
	}
	filename := opts.LabelNearest
	if filename == "" {
		filename = fmt.Sprintf("./labelnearest/%s.%s.labelnearest", opts.GetDatasetName(), common.JoinTableNames(cmd.TableNames.Values))
		opts.LabelNearest = filename
	}
	opts.Logger.Printf("writing the model to %s ...", filename)
	file, err := common.CreateWithDir(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := sticker.EncodeLabelNearest(model, file); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *TrainNearestCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @trainNearest [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
