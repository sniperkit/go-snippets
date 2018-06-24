package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// ShuffleCommand have flags for shuffle sub-command.
type ShuffleCommand struct {
	Help       bool
	Ks         common.OptionUints
	ReportOnly bool
	S          uint
	Seed       int
	TableNames common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewShuffleCommand returns a new ShuffleCommand.
func NewShuffleCommand(opts *Options) *ShuffleCommand {
	return &ShuffleCommand{
		Help:       false,
		Ks:         common.OptionUints{true, []uint{1, 3, 5}},
		ReportOnly: false,
		S:          5,
		Seed:       0,
		TableNames: common.OptionStrings{true, []string{"train.txt", "test.txt"}},
		opts:       opts,
	}
}

func (cmd *ShuffleCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@shuffle", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.Var(&cmd.Ks, "K", "Specify the K values for reporting the attainable precision@K")
	cmd.flagSet.BoolVar(&cmd.ReportOnly, "reportOnly", cmd.ReportOnly, "Report only, and do not write the splitted tables")
	cmd.flagSet.UintVar(&cmd.S, "S", cmd.S, "Specify the number of splitted tables")
	cmd.flagSet.IntVar(&cmd.Seed, "seed", cmd.Seed, "Specify the seed in shuffling")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *ShuffleCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run shuffles the specified tables of dataset, and split into K tables.
func (cmd *ShuffleCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("ShuffleCommands: %#v", cmd)
	S := cmd.S
	n := 0
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{},
		Y: sticker.LabelVectors{},
	}
	dsname := opts.GetDatasetName()
	if len(cmd.TableNames.Values) == 0 {
		return fmt.Errorf("specify the table names")
	}
	starts := make([]int, 0, len(cmd.TableNames.Values))
	for _, tblname := range cmd.TableNames.Values {
		opts.Logger.Printf("loading table %q from dataset %q ...", tblname, dsname)
		subds, err := opts.ReadDataset(tblname)
		if err != nil {
			return err
		}
		starts = append(starts, ds.Size())
		ds.X, ds.Y = append(ds.X, subds.X...), append(ds.Y, subds.Y...)
		n += subds.Size()
	}
	starts = append(starts, ds.Size())
	opts.Logger.Printf("Counting the unique labels on each dataset splits in order ...")
	report := func(ds *sticker.Dataset, starts []int, tblnames []string) {
		for t := 0; t < len(tblnames); t++ {
			trainLabelSet, testLabelSet := make(map[uint32]bool), make(map[uint32]bool)
			testStart, testEnd := starts[t], starts[t+1]
			for i := 0; i < testStart; i++ {
				for _, label := range ds.Y[i] {
					trainLabelSet[label] = true
				}
			}
			for i := testStart; i < testEnd; i++ {
				for _, label := range ds.Y[i] {
					testLabelSet[label] = true
				}
			}
			for i := testEnd; i < ds.Size(); i++ {
				for _, label := range ds.Y[i] {
					trainLabelSet[label] = true
				}
			}
			countTrain, countBoth, countTest := 0, 0, 0
			for label := range trainLabelSet {
				if _, ok := testLabelSet[label]; ok {
					countBoth++
				} else {
					countTrain++
				}
			}
			for label := range testLabelSet {
				if _, ok := trainLabelSet[label]; !ok {
					countTest++
				}
			}
			trainTblnames := make([]string, len(tblnames)-1)
			copy(trainTblnames, tblnames[:t])
			copy(trainTblnames[t:], tblnames[t+1:])
			testTblnames := []string{tblnames[t]}
			fmt.Fprintf(opts.OutputWriter, "train=%#v, test=%#v: label counts: (train,test)=(%d,%d); (trainOnly,both,testOnly)=(%d,%d,%d)\n", trainTblnames, testTblnames, countTrain+countBoth, countTest+countBoth, countTrain, countBoth, countTest)
			for _, K := range cmd.Ks.Values {
				avgPrecisionK := float32(0.0)
				for i := testStart; i < testEnd; i++ {
					precisionKi := float32(0.0)
					for _, label := range ds.Y[i] {
						if trainLabelSet[label] {
							precisionKi++
						}
					}
					if precisionKi > float32(K) {
						precisionKi = float32(K)
					}
					avgPrecisionK += precisionKi / float32(K)
				}
				avgPrecisionK /= float32(testEnd - testStart)
				fmt.Fprintf(opts.OutputWriter, "train=%#v, test=%#v: attainable Precision@%d=%-5.4g%%\n", trainTblnames, testTblnames, K, avgPrecisionK*100)
			}
		}
	}
	report(ds, starts, cmd.TableNames.Values)
	opts.Logger.Printf("shuffling %d entries into %d tables with the random number generator (seed=%d) ...", n, S, cmd.Seed)
	rng := rand.New(rand.NewSource(int64(cmd.Seed)))
	for i := 0; i < n; i++ {
		j := i + rng.Intn(n-i)
		ds.X[i], ds.Y[i], ds.X[j], ds.Y[j] = ds.X[j], ds.Y[j], ds.X[i], ds.Y[i]
	}
	joinedTblname := common.JoinTableNames(cmd.TableNames.Values)
	unit := (uint(n) + S - 1) / S
	starts, ends, tblnames := make([]int, 0, S), make([]int, 0, S), make([]string, 0, S)
	for s := uint(0); s < S; s++ {
		start, end := s*unit, (s+1)*unit
		if end > uint(n) {
			end = uint(n)
		}
		starts, ends, tblnames = append(starts, int(start)), append(ends, int(end)), append(tblnames, fmt.Sprintf("%s.%d.txt", joinedTblname, s))
	}
	starts = append(starts, ds.Size())
	opts.Logger.Printf("Counting the unique labels on %d-fold datasets ...", S)
	report(ds, starts, tblnames)
	if !cmd.ReportOnly {
		for s := uint(0); s < S; s++ {
			start, end, tblname := starts[s], ends[s], tblnames[s]
			tblpath := filepath.Join(opts.DatasetPath, tblname)
			opts.Logger.Printf("writing %d entries (%d-%d) into #%d table %q ...", end-start, start, end-1, s, tblpath)
			file, err := os.Create(tblpath)
			if err != nil {
				return fmt.Errorf("os.Create: %s: %s", tblpath, err)
			}
			subds := &sticker.Dataset{
				X: ds.X[start:end],
				Y: ds.Y[start:end],
			}
			writer := bufio.NewWriter(file)
			if err := subds.WriteTextDataset(writer); err != nil {
				file.Close()
				return fmt.Errorf("Dataset.WriteToText: %s: %s", tblpath, err)
			}
			writer.Flush()
			file.Close()
		}
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *ShuffleCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @shuffle [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
