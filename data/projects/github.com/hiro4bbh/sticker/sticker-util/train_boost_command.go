package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker/plugin"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TrainBoostCommand have flags for trainBoost sub-command.
type TrainBoostCommand struct {
	RankerTrainerName  string
	C, Epsilon         common.OptionFloat32
	Help               bool
	NegativeSampleSize uint
	PainterK           uint
	PainterName        string
	T                  uint
	TableNames         common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTrainBoostCommand returns a new TrainBoostCommand.
func NewTrainBoostCommand(opts *Options) *TrainBoostCommand {
	boostParams := plugin.NewLabelBoostParameters()
	return &TrainBoostCommand{
		RankerTrainerName:  boostParams.RankerTrainerName,
		C:                  common.OptionFloat32(boostParams.C),
		Epsilon:            common.OptionFloat32(boostParams.Epsilon),
		Help:               false,
		NegativeSampleSize: boostParams.NegativeSampleSize,
		PainterK:           boostParams.PainterK,
		PainterName:        boostParams.PainterName,
		T:                  boostParams.T,
		TableNames:         common.OptionStrings{true, []string{"train.txt"}},
		opts:               opts,
	}
}

func (cmd *TrainBoostCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@trainBoost", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.StringVar(&cmd.RankerTrainerName, "rankerTrainer", cmd.RankerTrainerName, "Specify the binary ranker trainer name")
	cmd.flagSet.Var(&cmd.C, "C", "Specify the inverse of the penalty parameter for each binary classifier")
	cmd.flagSet.Var(&cmd.Epsilon, "epsilon", "Specify the tolerance parameter for each binary classifier")
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.UintVar(&cmd.NegativeSampleSize, "negativeSampleSize", cmd.NegativeSampleSize, "Specify the size of each negative sample for Multi-Label Ranking Hinge Boosting (specify 0 for Multi-Label Hinge Boosting)")
	cmd.flagSet.UintVar(&cmd.PainterK, "painterK", cmd.PainterK, "Specify the maximum number of the painted target label")
	cmd.flagSet.StringVar(&cmd.PainterName, "painter", cmd.PainterName, "Specify the painter name")
	cmd.flagSet.UintVar(&cmd.T, "T", cmd.T, "Specify the maximum number of the target labels")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TrainBoostCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run trains on the specified table of dataset.
func (cmd *TrainBoostCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TrainBoostCommands: %#v", cmd)
	params := plugin.NewLabelBoostParameters()
	params.RankerTrainerName = cmd.RankerTrainerName
	params.C, params.Epsilon = float32(cmd.C), float32(cmd.Epsilon)
	params.NegativeSampleSize = cmd.NegativeSampleSize
	params.PainterK, params.PainterName = cmd.PainterK, cmd.PainterName
	params.T = cmd.T
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, ^uint(0), false)
	if err != nil {
		return err
	}
	model, err := plugin.TrainLabelBoost(ds, params, opts.DebugLogger)
	if err != nil {
		return err
	}
	filename := opts.LabelBoost
	if filename == "" {
		filename = fmt.Sprintf("./labelboost/%s.%s.T%d.labelboost", opts.GetDatasetName(), common.JoinTableNames(cmd.TableNames.Values), cmd.T)
		opts.LabelBoost = filename
	}
	opts.Logger.Printf("writing the model to %s ...", filename)
	file, err := common.CreateWithDir(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := plugin.EncodeLabelBoost(model, file); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *TrainBoostCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @trainBoost [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
