package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TrainOneCommand have flags for trainOne sub-command.
type TrainOneCommand struct {
	ClassifierTrainerName string
	C, Epsilon            common.OptionFloat32
	Help                  bool
	T                     uint
	TableNames            common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTrainOneCommand returns a new TrainOneCommand.
func NewTrainOneCommand(opts *Options) *TrainOneCommand {
	boostParams := sticker.NewLabelOneParameters()
	return &TrainOneCommand{
		ClassifierTrainerName: boostParams.ClassifierTrainerName,
		C:          common.OptionFloat32(boostParams.C),
		Epsilon:    common.OptionFloat32(boostParams.Epsilon),
		Help:       false,
		T:          boostParams.T,
		TableNames: common.OptionStrings{true, []string{"train.txt"}},
		opts:       opts,
	}
}

func (cmd *TrainOneCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@trainOne", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.StringVar(&cmd.ClassifierTrainerName, "classifierTrainer", cmd.ClassifierTrainerName, "Specify the binary classifier trainer name")
	cmd.flagSet.Var(&cmd.C, "C", "Specify the inverse of the penalty parameter for each binary classifier")
	cmd.flagSet.Var(&cmd.Epsilon, "epsilon", "Specify the tolerance parameter for each binary classifier")
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.UintVar(&cmd.T, "T", cmd.T, "Specify the maximum number of the target labels")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TrainOneCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run trains on the specified table of dataset.
func (cmd *TrainOneCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TrainOneCommands: %#v", cmd)
	params := sticker.NewLabelOneParameters()
	params.ClassifierTrainerName = cmd.ClassifierTrainerName
	params.C, params.Epsilon = float32(cmd.C), float32(cmd.Epsilon)
	params.T = cmd.T
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, ^uint(0), false)
	if err != nil {
		return err
	}
	model, err := sticker.TrainLabelOne(ds, params, opts.DebugLogger)
	if err != nil {
		return err
	}
	filename := opts.LabelOne
	if filename == "" {
		filename = fmt.Sprintf("./labelone/%s.%s.T%d.labelone", opts.GetDatasetName(), common.JoinTableNames(cmd.TableNames.Values), cmd.T)
		opts.LabelOne = filename
	}
	opts.Logger.Printf("writing the model to %s ...", filename)
	file, err := common.CreateWithDir(filename)
	if err != nil {
		return err
	}
	if err := sticker.EncodeLabelOne(model, file); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *TrainOneCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @trainOne [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
