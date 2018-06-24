package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TestNearCommand have flags for testNear sub-command.
type TestNearCommand struct {
	Alpha      common.OptionFloat32
	Beta       common.OptionFloat32
	C          uint
	Help       bool
	Ks         common.OptionUints
	N          uint
	Per        uint
	S          uint
	TableNames common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTestNearCommand returns a new TestNearCommand.
func NewTestNearCommand(opts *Options) *TestNearCommand {
	return &TestNearCommand{
		Alpha:      common.OptionFloat32(1.0),
		Beta:       common.OptionFloat32(1.0),
		C:          uint(2),
		Help:       false,
		Ks:         common.OptionUints{true, []uint{1, 3, 5}},
		N:          ^uint(0),
		Per:        uint(0),
		S:          uint(1),
		TableNames: common.OptionStrings{true, []string{"test.txt"}},
		opts:       opts,
	}
}

func (cmd *TestNearCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@testNear", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.Var(&cmd.Alpha, "alpha", "Specify the smoothing parameter for weighting the voted by each neighbor")
	cmd.flagSet.Var(&cmd.Beta, "beta", "Specify the balancing parameter between the Jaccard and cosine similarity")
	cmd.flagSet.UintVar(&cmd.C, "c", cmd.C, "Specify the factor of candidate near neighbors")
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.Ks, "K", "Specify the top-K values")
	cmd.flagSet.UintVar(&cmd.N, "N", cmd.N, "Specify the maximum number of the tested entries")
	cmd.flagSet.UintVar(&cmd.Per, "per", cmd.Per, "Specify the deep-inspection timing counts (not do deep-inspection if 0)")
	cmd.flagSet.UintVar(&cmd.S, "S", cmd.S, "Specify the number of nearest neighbors")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TestNearCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run tests the .labelforest model on the specified table of dataset.
func (cmd *TestNearCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TestNearCommands: %#v", cmd)
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, cmd.N, true)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading .labelnear model from %q ...", opts.LabelNear)
	model, err := common.ReadLabelNear(opts.LabelNear)
	if err != nil {
		return err
	}
	if opts.Verbose {
		bucketUsage, bucketSizeHist := model.Hashing.Summary()
		opts.Logger.Printf("JaccardHashing(K=%d,L=%d,R=%d): bucketUsage=%d", model.Hashing.K(), model.Hashing.L(), model.Hashing.R(), bucketUsage)
		opts.Logger.Printf("JaccardHashing(K=%d,L=%d,R=%d): bucketSizeHist=%d", model.Hashing.K(), model.Hashing.L(), model.Hashing.R(), bucketSizeHist)
	}
	reporter := common.NewResultsReporter(ds.Y, cmd.Ks.Values)
	opts.Logger.Printf("predicting top-%d labels ...", reporter.MaxK())
	reporter.ResetTimer()
	if cmd.Per == 0 {
		reporter.Report(model.PredictAll(ds.X, reporter.MaxK(), cmd.C, cmd.S, float32(cmd.Alpha), float32(cmd.Beta)), opts.OutputWriter)
	} else {
		Yhat := make(sticker.LabelVectors, 0, ds.Size())
		for i, xi := range ds.X {
			yihat, labelHist, indexSimsTopS := model.Predict(xi, reporter.MaxK(), cmd.C, cmd.S, float32(cmd.Alpha), float32(cmd.Beta))
			Yhat = append(Yhat, yihat)
			if uint(i)%cmd.Per == 0 {
				if opts.DebugLogger != nil {
					InspectDataEntryWithNeighbors(opts, reporter, i, xi, ds.Y[i], yihat, labelHist, indexSimsTopS)
				}
				reporter.Report(Yhat, opts.OutputWriter)
			}
		}
		reporter.Report(Yhat, opts.OutputWriter)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *TestNearCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @testNear [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
