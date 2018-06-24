package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// PruneOneCommand have flags for pruneOne sub-command.
type PruneOneCommand struct {
	Help bool
	T    uint

	opts    *Options
	flagSet *flag.FlagSet
}

// NewPruneOneCommand returns a new PruneOneCommand.
func NewPruneOneCommand(opts *Options) *PruneOneCommand {
	return &PruneOneCommand{
		Help: false,
		T:    uint(1),
		opts: opts,
	}
}

func (cmd *PruneOneCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@pruneOne", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.UintVar(&cmd.T, "T", cmd.T, "Specify the maximum number of the round in the pruned .labelone model")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *PruneOneCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run prunes .labelone model.
func (cmd *PruneOneCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("PruneOneCommands: %#v", cmd)
	opts.Logger.Printf("loading .labelone model from %q ...", opts.LabelOne)
	model, err := common.ReadLabelOne(opts.LabelOne)
	if err != nil {
		return err
	}
	opts.Logger.Printf("extracting the first %d rounds ...", cmd.T)
	newModel := model.Prune(cmd.T)
	filename := fmt.Sprintf("%s.prune%d.labelone", opts.LabelOne, cmd.T)
	opts.Logger.Printf("writing the pruned model to %s ...", filename)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	defer file.Close()
	if err := sticker.EncodeLabelOne(newModel, file); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *PruneOneCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @pruneOne [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
