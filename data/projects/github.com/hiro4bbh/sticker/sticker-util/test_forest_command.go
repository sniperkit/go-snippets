package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/bits"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TestForestCommand have flags for testForest sub-command.
type TestForestCommand struct {
	Help        bool
	Ks          common.OptionUints
	N           uint
	OnlyResults bool
	TableNames  common.OptionStrings
	Weighted    bool

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTestForestCommand returns a new TestForestCommand.
func NewTestForestCommand(opts *Options) *TestForestCommand {
	return &TestForestCommand{
		Help:        false,
		Ks:          common.OptionUints{true, []uint{1, 3, 5}},
		N:           ^uint(0),
		OnlyResults: false,
		TableNames:  common.OptionStrings{true, []string{"test.txt"}},
		Weighted:    false,
		opts:        opts,
	}
}

func (cmd *TestForestCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@testForest", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.Ks, "K", "Specify the top-K values")
	cmd.flagSet.UintVar(&cmd.N, "N", cmd.N, "Specify the maximum number of the tested entries")
	cmd.flagSet.BoolVar(&cmd.OnlyResults, "onlyResults", cmd.OnlyResults, "Report only the test results")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
	cmd.flagSet.BoolVar(&cmd.Weighted, "weighted", cmd.Weighted, "Use the weighted forest")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TestForestCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run tests the .labelforest model on the specified table of dataset.
func (cmd *TestForestCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TestForestCommands: %#v", cmd)
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, cmd.N, true)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading .labelforest model from %q ...", opts.LabelForest)
	forest, err := common.ReadLabelForest(opts.LabelForest)
	if err != nil {
		return err
	}
	reporter := common.NewResultsReporter(ds.Y, cmd.Ks.Values)
	reporter.ResetTimer()
	var leafIdsSlice [][]uint64
	var Yhat sticker.LabelVectors
	if cmd.Weighted {
		opts.Logger.Printf("classifying all entries with weights ...")
		var weightsSlice [][]float32
		leafIdsSlice, weightsSlice = forest.ClassifyAllWithWeight(ds.X)
		opts.Logger.Printf("predicting top-%d labels with weights ...", reporter.MaxK())
		Yhat = forest.PredictAllWithWeight(leafIdsSlice, weightsSlice, reporter.MaxK())
	} else {
		opts.Logger.Printf("classifying all entries ...")
		leafIdsSlice = forest.ClassifyAll(ds.X)
		opts.Logger.Printf("predicting top-%d labels ...", reporter.MaxK())
		Yhat = forest.PredictAll(leafIdsSlice, reporter.MaxK())
	}
	reporter.Report(Yhat, opts.OutputWriter)
	if cmd.OnlyResults {
		return nil
	}
	sumHeights := make([]int, len(forest.Trees))
	sumTV := float32(0.0)
	labelFreqSlice := make(sticker.SparseVectors, len(forest.Trees))
	for i, leafIds := range leafIdsSlice {
		for treeId, leafId := range leafIds {
			sumHeights[treeId] += bits.Len64(leafId) - 1
			labelFreqSlice[treeId] = forest.Trees[treeId].LabelFreqSet[leafId]
		}
		TV := sticker.AvgTotalVariationAmongSparseVectors(labelFreqSlice)
		if sticker.IsNaN32(TV) {
			return fmt.Errorf("#%d data point: labelFreqSlice=%#v, TV=%g", i, labelFreqSlice, TV)
		}
		sumTV += TV
	}
	avgTV := sumTV / float32(len(leafIdsSlice))
	avgHeights := make([]float32, len(forest.Trees))
	for treeId, sumHeight := range sumHeights {
		avgHeights[treeId] = float32(sumHeight) / float32(ds.Size())
	}
	avgHeight := float32(0.0)
	for _, h := range avgHeights {
		avgHeight += h
	}
	avgHeight /= float32(len(avgHeights))
	expectAvgHeight := sticker.LogBinary32(float32(ds.Size()) / float32(forest.TreeParams.MaxEntriesInLeaf))
	opts.Logger.Printf("Tree Balances: [avg(height[.])]=%.4g, log2(n/MaxEntriesInLeaf)=%.4g", avgHeights, expectAvgHeight)
	fmt.Printf("Average Tree Balance: avg([avg(height[.])])/log2(n/MaxEntriesInLeaf)=%.4g\n", avgHeight/expectAvgHeight)
	fmt.Printf("Average Total Variation among Trees: %.5g\n", avgTV)
	return nil
}

// ShowHelp shows the help.
func (cmd *TestForestCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @testForest [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
