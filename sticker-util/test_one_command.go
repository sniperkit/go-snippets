package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TestOneCommand have flags for testOne sub-command.
type TestOneCommand struct {
	Help       bool
	Ks         common.OptionUints
	N          uint
	Restore    bool
	Ts         common.OptionUints
	TableNames common.OptionStrings

	Result map[string]interface{}

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTestOneCommand returns a new TestOneCommand.
func NewTestOneCommand(opts *Options) *TestOneCommand {
	return &TestOneCommand{
		Help:       false,
		Ks:         common.OptionUints{true, []uint{1, 3, 5}},
		N:          ^uint(0),
		Restore:    false,
		Ts:         common.OptionUints{true, []uint{0}},
		TableNames: common.OptionStrings{true, []string{"test.txt"}},
		opts:       opts,
	}
}

func (cmd *TestOneCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@testOne", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.Ks, "K", "Specify the top-K values")
	cmd.flagSet.UintVar(&cmd.N, "N", cmd.N, "Specify the maximum number of the tested entries")
	cmd.flagSet.BoolVar(&cmd.Restore, "restore", cmd.Restore, "Restore the test result if true")
	cmd.flagSet.Var(&cmd.Ts, "T", "Specify the used numbers of rounds (use all rounds if zero)")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TestOneCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run tests the .labelone model on the specified table of dataset.
func (cmd *TestOneCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TestOneCommands: %#v", cmd)
	dsname := opts.GetDatasetName()
	tblname := common.JoinTableNames(cmd.TableNames.Values)
	restoreName := fmt.Sprintf("%s.testOne.%s.%s.bin", opts.LabelOne, dsname, tblname)
	if cmd.Restore {
		opts.Logger.Printf("restroing the dumped test result in %q ...", restoreName)
		f, err := os.Open(restoreName)
		if err != nil {
			return err
		}
		defer f.Close()
		return gob.NewDecoder(f).Decode(&cmd.Result)
	}
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, cmd.N, true)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading .labelone model from %q ...", opts.LabelOne)
	model, err := common.ReadLabelOne(opts.LabelOne)
	if err != nil {
		return err
	}
	maxK := uint(0)
	for _, K := range cmd.Ks.Values {
		if maxK < K {
			maxK = K
		}
	}
	maxAvgPrecisions := make([]float32, 0, len(cmd.Ks.Values))
	for _, K := range cmd.Ks.Values {
		maxPrecisionKs := sticker.ReportMaxPrecision(ds.Y, K)
		maxSumPrecisionK := float32(0.0)
		for _, maxPrecisionKi := range maxPrecisionKs {
			maxSumPrecisionK += maxPrecisionKi
		}
		maxAvgPrecisions = append(maxAvgPrecisions, maxSumPrecisionK/float32(len(ds.Y)))
	}
	rounds := make([]interface{}, 0, len(cmd.Ts.Values))
	var avgMaxPrecisions map[uint]float32
	for _, T := range cmd.Ts.Values {
		reporter := common.NewResultsReporter(ds.Y, cmd.Ks.Values)
		if avgMaxPrecisions == nil {
			avgMaxPrecisions = reporter.AvgMaxPrecisionKs()
		}
		opts.Logger.Printf("predicting top-%d labels ...", reporter.MaxK())
		reporter.ResetTimer()
		opts.Logger.Printf("predicting top-%d labels with first %d rounds ...", reporter.MaxK(), T)
		reporter.ResetTimer()
		avgPKs, avgNKs := reporter.Report(model.PredictAll(ds.X, reporter.MaxK(), T), opts.OutputWriter)
		inferenceTime, inferenceTimePerEntry := reporter.InferenceTimes()
		rounds = append(rounds, map[string]interface{}{
			"T":                     T,
			"inferenceTime":         fmt.Sprintf("%s", inferenceTime),
			"inferenceTimePerEntry": fmt.Sprintf("%s", inferenceTimePerEntry),
			"precisions":            avgPKs,
			"nDCGs":                 avgNKs,
		})
	}
	cmd.Result = map[string]interface{}{
		"Ks":            cmd.Ks.Values,
		"maxPrecisions": maxAvgPrecisions,
		"nentries":      ds.Size(),
		"rounds":        rounds,
	}
	opts.Logger.Printf("dumping the test result to %q ...", restoreName)
	f, err := os.Create(restoreName)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(cmd.Result)
}

// ShowHelp shows the help.
func (cmd *TestOneCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @testOne [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
