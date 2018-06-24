package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/plugin"
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// TrainForestCommand have flags for trainForest sub-command.
type TrainForestCommand struct {
	AssignerName          string
	AssignInitializerName string
	ClassifierTrainerName string
	C, Epsilon            common.OptionFloat32
	FeatureSubSamplerName string
	Help                  bool
	K                     uint
	MaxEntriesInLeaf      uint
	NtopLabels            uint
	Ntrees                uint
	SubSamplerName        string
	SubSampleSize         uint
	SuppVecK              uint
	TableNames            common.OptionStrings

	opts    *Options
	flagSet *flag.FlagSet
}

// NewTrainForestCommand returns a new TrainForestCommand.
func NewTrainForestCommand(opts *Options) *TrainForestCommand {
	treeParams := plugin.NewLabelTreeParameters()
	return &TrainForestCommand{
		AssignerName:          treeParams.AssignerName,
		AssignInitializerName: treeParams.AssignInitializerName,
		ClassifierTrainerName: treeParams.ClassifierTrainerName,
		C:                     common.OptionFloat32(treeParams.C),
		Epsilon:               common.OptionFloat32(treeParams.Epsilon),
		FeatureSubSamplerName: treeParams.FeatureSubSamplerName,
		Help:             false,
		K:                treeParams.K,
		MaxEntriesInLeaf: treeParams.MaxEntriesInLeaf,
		NtopLabels:       0,
		Ntrees:           uint(runtime.GOMAXPROCS(0)),
		SubSamplerName:   "random",
		SubSampleSize:    10000,
		SuppVecK:         treeParams.SuppVecK,
		TableNames:       common.OptionStrings{true, []string{"train.txt"}},
		opts:             opts,
	}
}

func (cmd *TrainForestCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@trainForest", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.StringVar(&cmd.AssignerName, "assigner", cmd.AssignerName, "Specify the left/right assigner name")
	cmd.flagSet.StringVar(&cmd.AssignInitializerName, "assignInitializer", cmd.AssignInitializerName, "Specify the left/right assign initializer name")
	cmd.flagSet.StringVar(&cmd.ClassifierTrainerName, "classifierTrainer", cmd.ClassifierTrainerName, "Specify the binary classifier trainer name")
	cmd.flagSet.Var(&cmd.C, "C", "Specify the inverse of the penalty for each binary classifier")
	cmd.flagSet.Var(&cmd.Epsilon, "epsilon", "Specify the tolerance parameter for each binary classifier")
	cmd.flagSet.StringVar(&cmd.FeatureSubSamplerName, "featureSubSampler", cmd.FeatureSubSamplerName, "Specify the dataset feature sub-sampler name")
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.UintVar(&cmd.K, "K", cmd.K, "Specify the maximum number of the labels in each terminal leaf")
	cmd.flagSet.UintVar(&cmd.MaxEntriesInLeaf, "maxEntriesInLeaf", cmd.MaxEntriesInLeaf, "Specify the maximum number of the entries in each leaf (best-effort)")
	cmd.flagSet.UintVar(&cmd.NtopLabels, "ntopLabels", cmd.NtopLabels, "Specify the number of the used top labels (all labels are used if 0)")
	cmd.flagSet.UintVar(&cmd.Ntrees, "ntrees", cmd.Ntrees, "Specify the number of the trained trees")
	cmd.flagSet.StringVar(&cmd.SubSamplerName, "subSampler", cmd.SubSamplerName, "Specify the dataset sub-sampler name")
	cmd.flagSet.UintVar(&cmd.SubSampleSize, "subSampleSize", cmd.SubSampleSize, "Specify each sub-sample size")
	cmd.flagSet.UintVar(&cmd.SuppVecK, "suppVecK", cmd.SuppVecK, "Specify the maximum number of the support vectors in each leaf")
	cmd.flagSet.Var(&cmd.TableNames, "table", "Specify the table names")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *TrainForestCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run trains on the specified table of dataset.
func (cmd *TrainForestCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("TrainForestCommands: %#v", cmd)
	params := plugin.NewLabelTreeParameters()
	params.AssignerName = cmd.AssignerName
	params.AssignInitializerName = cmd.AssignInitializerName
	params.ClassifierTrainerName = cmd.ClassifierTrainerName
	params.C, params.Epsilon = float32(cmd.C), float32(cmd.Epsilon)
	params.FeatureSubSamplerName = cmd.FeatureSubSamplerName
	params.K = cmd.K
	params.MaxEntriesInLeaf = cmd.MaxEntriesInLeaf
	params.SuppVecK = cmd.SuppVecK
	var subsampler plugin.DatasetEntrySubSampler
	switch cmd.SubSamplerName {
	case "deterministic":
		subsampler = plugin.NewDeterministicDatasetEntrySubSampler(cmd.SubSampleSize)
	case "random":
		subsampler = plugin.NewRandomDatasetEntrySubSampler(cmd.SubSampleSize)
	default:
		return fmt.Errorf("unknown subSampler: %s", cmd.SubSamplerName)
	}
	ds, err := opts.ReadDatasets(cmd.TableNames.Values, ^uint(0), false)
	if err != nil {
		return err
	}
	if cmd.NtopLabels > 0 {
		opts.Logger.Printf("collecting label frequencies ...")
		labelFreq := make(sticker.SparseVector)
		for _, yi := range ds.Y {
			for _, label := range yi {
				labelFreq[label] += 1
			}
		}
		opts.Logger.Printf("ranking top-%d labels ...", cmd.NtopLabels)
		labelRankK := sticker.RankTopK(labelFreq, cmd.NtopLabels)
		labelInvRankK := sticker.InvertRanks(labelRankK)
		remain := []int{}
		ndeleteds := 0
		for i, yi := range ds.Y {
			yinew := make(sticker.LabelVector, 0)
			for _, label := range yi {
				if _, ok := labelInvRankK[label]; ok {
					yinew = append(yinew, label)
				}
			}
			if len(yi) > len(yinew) {
				ndeleteds++
			}
			if len(yinew) > 0 {
				remain = append(remain, i)
			}
			ds.Y[i] = yinew
		}
		n0 := ds.Size()
		ds = ds.SubSet(remain)
		opts.Logger.Printf("deleted non-top-%d labels from %d entry(s) (total %d entry(s) -> %d entry(s) in dataset)", cmd.NtopLabels, ndeleteds, n0, ds.Size())
	}
	forest, err := plugin.TrainLabelForest(ds, cmd.Ntrees, subsampler, params, opts.DebugLogger)
	if err != nil {
		return err
	}
	filename := opts.LabelForest
	if filename == "" {
		filename = fmt.Sprintf("./labelforest/%s.%s.N%d%s%d.labelforest", opts.GetDatasetName(), common.JoinTableNames(cmd.TableNames.Values), cmd.Ntrees, strings.Title(string(cmd.SubSamplerName)), cmd.SubSampleSize)
		opts.LabelForest = filename
	}
	opts.Logger.Printf("writing the model to %s ...", filename)
	file, err := common.CreateWithDir(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := plugin.EncodeLabelForest(forest, file); err != nil {
		return fmt.Errorf("%s: %s", filename, err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *TrainForestCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @trainForest [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
