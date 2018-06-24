package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// InspectDataEntryWithNeighbors inspects yi and yihat with estimated neighbor information
func InspectDataEntryWithNeighbors(opts *Options, reporter *common.ResultsReporter, i int, x sticker.FeatureVector, y, yhat sticker.LabelVector, labelHist map[uint32]float32, indexSimsTopS sticker.KeyValues32) {
	ymap := make(map[uint32]int)
	for _, label := range y {
		ymap[label] = 30
	}
	lenx := float32(0.0)
	for _, xpair := range x {
		lenx += xpair.Value * xpair.Value
	}
	lenx = sticker.Sqrt32(lenx)
	opts.DebugLogger.Printf("i=%d, ||x||_2=%.3g", i, lenx)
	line := "  indexSimsTopS=map["
	for k, indexSim := range indexSimsTopS {
		if k > 0 {
			line += " "
		}
		line += fmt.Sprintf("%s:%.3g", opts.LabelMap(indexSim.Key, true), indexSim.Value)
	}
	line += "]"
	opts.DebugLogger.Print(line)
	line = "  labelHist=map["
	labelHistKV := make(sticker.KeyValues32OrderedByValue, 0, len(labelHist))
	for label, freq := range labelHist {
		labelHistKV = append(labelHistKV, sticker.KeyValue32{label, freq})
	}
	sort.Sort(sort.Reverse(labelHistKV))
	for rank, labelFreq := range labelHistKV {
		if rank > 0 {
			line += " "
		}
		if _, ok := ymap[labelFreq.Key]; ok {
			if uint(rank) < reporter.MaxK() {
				ymap[labelFreq.Key] = 32
			} else {
				ymap[labelFreq.Key] = 36
			}
		}
		line += fmt.Sprintf("\033[%dm%s\033[0m:%.3g", ymap[labelFreq.Key], opts.LabelMap(labelFreq.Key, true), labelFreq.Value)
	}
	line += "]"
	opts.DebugLogger.Print(line)
	line = "  <> y=["
	for k, label := range y {
		if k > 0 {
			line += " "
		}
		line += fmt.Sprintf("\033[%dm%s\033[0m", ymap[label], opts.LabelMap(label, true))
	}
	line += "]: "
	for iK, K := range reporter.Ks() {
		if iK > 0 {
			line += ","
		}
		pK := sticker.ReportPrecision(sticker.LabelVectors{y}, K, sticker.LabelVectors{yhat})
		line += fmt.Sprintf("P@%d=%.2f%%", K, pK[0]*100)
	}
	opts.DebugLogger.Print(line)
}

// TestNearestCommand have flags for testNearest sub-command.
type TestNearestCommand struct {
	Alpha      common.OptionFloat32
	Beta       common.OptionFloat32
	Help       bool
	Ks         common.OptionUints
	N          uint
	Per        uint
	S          uint
	TableNames common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTestNearestCommand returns a new TestNearestCommand.
func NewTestNearestCommand(opts *Options) *TestNearestCommand {
	return &TestNearestCommand{
		Alpha:      common.OptionFloat32(1.0),
		Beta:       common.OptionFloat32(1.0),
		Help:       false,
		Ks:         common.OptionUints{true, []uint{1, 3, 5}},
		N:          ^uint(0),
		Per:        uint(0),
		S:          uint(1),
		TableNames: common.OptionStrings{true, []string{"test.txt"}},
		opts:       opts,
	}
}

func (cmd *TestNearestCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@testNearest", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.Var(&cmd.Alpha, "alpha", "Specify the smoothing parameter for weighting the voted by each neighbour")
	cmd.flagSet.Var(&cmd.Beta, "beta", "Specify the balancing parameter between the Jaccard and cosine similarity")
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.Ks, "K", "Specify the top-K values")
	cmd.flagSet.UintVar(&cmd.N, "N", cmd.N, "Specify the maximum number of the tested entries")
	cmd.flagSet.UintVar(&cmd.Per, "per", cmd.Per, "Specify the deep-inspection timing counts (not do deep-inspection if 0)")
	cmd.flagSet.UintVar(&cmd.S, "S", cmd.S, "Specify the number of nearest neighbours")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TestNearestCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run tests the .labelforest model on the specified table of dataset.
func (cmd *TestNearestCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TestNearestCommands: %#v", cmd)
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, cmd.N, true)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading .labelnearest model from %q ...", opts.LabelNearest)
	model, err := common.ReadLabelNearest(opts.LabelNearest)
	if err != nil {
		return err
	}
	reporter := common.NewResultsReporter(ds.Y, cmd.Ks.Values)
	opts.Logger.Printf("predicting top-%d labels ...", reporter.MaxK())
	reporter.ResetTimer()
	if cmd.Per == 0 {
		reporter.Report(model.PredictAll(ds.X, reporter.MaxK(), cmd.S, float32(cmd.Alpha), float32(cmd.Beta)), opts.OutputWriter)
	} else {
		Yhat := make(sticker.LabelVectors, 0, ds.Size())
		ctx := model.NewContext()
		for i, xi := range ds.X {
			yihat, labelHist, indexSimsTopS := model.PredictWithContext(xi, reporter.MaxK(), cmd.S, float32(cmd.Alpha), float32(cmd.Beta), ctx)
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
func (cmd *TestNearestCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @testNearest [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
