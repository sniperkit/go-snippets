package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TrainNearCommand have flags for trainBoost sub-command.
type TrainNearCommand struct {
	Help       bool
	K          uint
	L          uint
	R          uint
	TableNames common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTrainNearCommand returns a new TrainNearCommand.
func NewTrainNearCommand(opts *Options) *TrainNearCommand {
	params := sticker.NewLabelNearParameters()
	return &TrainNearCommand{
		Help:       false,
		K:          params.K,
		L:          params.L,
		R:          params.R,
		TableNames: common.OptionStrings{true, []string{"train.txt"}},
		opts:       opts,
	}
}

func (cmd *TrainNearCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@trainNear", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.UintVar(&cmd.K, "K", cmd.K, "Show the number of the hash tables")
	cmd.flagSet.UintVar(&cmd.L, "L", cmd.L, "Show the bit-width of bucket indices in each hash table")
	cmd.flagSet.UintVar(&cmd.R, "R", cmd.R, "Show the size of a reservoir of each backet")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TrainNearCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run trains on the specified table of dataset.
func (cmd *TrainNearCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TrainNearCommands: %#v", cmd)
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, ^uint(0), false)
	if err != nil {
		return err
	}
	params := sticker.NewLabelNearParameters()
	params.K, params.L, params.R = cmd.K, cmd.L, cmd.R
	model, err := sticker.TrainLabelNear(ds, params, opts.DebugLogger)
	if err != nil {
		return err
	}
	filename := opts.LabelNear
	if filename == "" {
		filename = fmt.Sprintf("./labelnear/%s.%s.labelnear", opts.GetDatasetName(), common.JoinTableNames(cmd.TableNames.Values))
		opts.LabelNear = filename
	}
	opts.Logger.Printf("writing the model to %s ...", filename)
	file, err := common.CreateWithDir(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := sticker.EncodeLabelNear(model, file); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *TrainNearCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @trainNear [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
