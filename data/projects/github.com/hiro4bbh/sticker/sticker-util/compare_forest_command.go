package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// CompareForestCommand have flags for compareForest sub-command.
type CompareForestCommand struct {
	Help        bool
	LabelForest string
	TableNames  common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewCompareForestCommand returns a new CompareForestCommand.
func NewCompareForestCommand(opts *Options) *CompareForestCommand {
	return &CompareForestCommand{
		Help:        false,
		LabelForest: "",
		TableNames:  common.OptionStrings{},
		opts:        opts,
	}
}

func (cmd *CompareForestCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@compareForest", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.StringVar(&cmd.LabelForest, "labelforest", cmd.LabelForest, "Specify the compared .labelforest file name")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *CompareForestCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run compares the test performance on the specified table of dataset between the given .labelforest model and the compared one.
func (cmd *CompareForestCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("CompareForestCommands: %#v", cmd)
	dsname := opts.GetDatasetName()
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{},
		Y: sticker.LabelVectors{},
	}
	if len(cmd.TableNames.Values) == 0 {
		return fmt.Errorf("specify the table names")
	}
	for _, tblname := range cmd.TableNames.Values {
		opts.Logger.Printf("loading table %q of dataset %q ...", tblname, dsname)
		subds, err := opts.ReadDataset(tblname)
		if err != nil {
			return err
		}
		ds.X, ds.Y = append(ds.X, subds.X...), append(ds.Y, subds.Y...)
	}
	opts.Logger.Printf("loading .labelforest model from %q ...", opts.LabelForest)
	forest1, err := common.ReadLabelForest(opts.LabelForest)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading compared .labelforest model from %q ...", cmd.LabelForest)
	forest2, err := common.ReadLabelForest(cmd.LabelForest)
	if err != nil {
		return err
	}
	opts.Logger.Printf("classifying all entries with %q ...", opts.LabelForest)
	leafIdsSlice1 := forest1.ClassifyAll(ds.X)
	opts.Logger.Printf("classifying all entries with %q ...", cmd.LabelForest)
	leafIdsSlice2 := forest2.ClassifyAll(ds.X)
	opts.Logger.Printf("calculating average total variation between %q and %q ...", opts.LabelForest, cmd.LabelForest)
	sumTV := float32(0.0)
	for i := range ds.X {
		labelFreq1, labelFreq2 := make(sticker.SparseVector), make(sticker.SparseVector)
		for treeId, leafId := range leafIdsSlice1[i] {
			labelFreq := forest1.Trees[treeId].LabelFreqSet[leafId]
			for label, freq := range labelFreq {
				labelFreq1[label] += freq
			}
		}
		for treeId, leafId := range leafIdsSlice2[i] {
			labelFreq := forest2.Trees[treeId].LabelFreqSet[leafId]
			for label, freq := range labelFreq {
				labelFreq2[label] += freq
			}
		}
		TV := sticker.AvgTotalVariationAmongSparseVectors(sticker.SparseVectors{labelFreq1, labelFreq2})
		sumTV += TV
	}
	avgTV := sumTV / float32(ds.Size())
	fmt.Printf("Average Total Variation between %q and %q: %.5g\n", opts.LabelForest, cmd.LabelForest, avgTV)
	return nil
}

// ShowHelp shows the help.
func (cmd *CompareForestCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @compareForest [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
