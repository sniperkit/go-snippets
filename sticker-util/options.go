package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"runtime/pprof"
	"strings"
	"time"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/plugin/next"
)

// Options have options for common flags and flags for each sub-command.
// See the help for details.
type Options struct {
	// The following members are the common flags.
	CPUProfile     string
	Debug          bool
	FeatureMapName string
	Help           bool
	HTTPResource   string
	LabelBoost     string
	LabelConst     string
	LabelForest    string
	LabelNear      string
	LabelNearest   string
	LabelNext      string
	LabelOne       string
	LabelMapName   string
	Verbose        bool
	DatasetPath    string
	// The following members are for each sub-commands.
	CompareForest *CompareForestCommand
	InspectForest *InspectForestCommand
	InspectOne    *InspectOneCommand
	PruneOne      *PruneOneCommand
	Shuffle       *ShuffleCommand
	Summarize     *SummarizeCommand
	TestBoosts    []*TestBoostCommand
	TestConsts    []*TestConstCommand
	TestForests   []*TestForestCommand
	TestNears     []*TestNearCommand
	TestNearests  []*TestNearestCommand
	TestNexts     []next.TestCommand
	TestOnes      []*TestOneCommand
	TrainBoost    *TrainBoostCommand
	TrainConst    *TrainConstCommand
	TrainForest   *TrainForestCommand
	TrainNear     *TrainNearCommand
	TrainNearest  *TrainNearestCommand
	TrainNext     next.TrainCommand
	TrainOne      *TrainOneCommand

	// The following members are for logging or debugging use.
	ErrorWriter, OutputWriter io.Writer
	Logger                    *log.Logger // -verbose
	DebugLogger               *log.Logger // -debug

	execpath             string
	flagSet              *flag.FlagSet
	featureMap, labelMap []string
}

// NewOptions returns a new Options with default values.
func NewOptions(execpath string, outputWriter, errorWriter io.Writer) *Options {
	return &Options{
		CPUProfile:     "",
		Debug:          false,
		FeatureMapName: "feature_map.txt",
		Help:           false,
		HTTPResource:   filepath.Join(build.Default.GOPATH, "src/github.com/hiro4bbh/sticker/sticker-util/res"),
		LabelBoost:     "",
		LabelConst:     "",
		LabelForest:    "",
		LabelNear:      "",
		LabelNearest:   "",
		LabelNext:      "",
		LabelOne:       "",
		LabelMapName:   "label_map.txt",
		Verbose:        false,
		DatasetPath:    "",

		CompareForest: nil,
		InspectForest: nil,
		InspectOne:    nil,
		PruneOne:      nil,
		Shuffle:       nil,
		Summarize:     nil,
		TestBoosts:    nil,
		TestConsts:    nil,
		TestForests:   nil,
		TestNearests:  nil,
		TestNexts:     nil,
		TestOnes:      nil,
		TrainBoost:    nil,
		TrainConst:    nil,
		TrainForest:   nil,
		TrainNearest:  nil,
		TrainNext:     nil,
		TrainOne:      nil,

		OutputWriter: outputWriter,
		ErrorWriter:  errorWriter,

		execpath: execpath,
	}
}

// FeatureMap returns the feature name.
func (opts *Options) FeatureMap(feature uint32, quote bool) string {
	if feature < uint32(len(opts.featureMap)) {
		if quote {
			return fmt.Sprintf("%q", opts.featureMap[feature])
		}
		return opts.featureMap[feature]
	}
	return fmt.Sprintf("%d", feature)
}

// GetDebugLogger returns opt.DebugLogger.
func (opts *Options) GetDebugLogger() *log.Logger {
	return opts.DebugLogger
}

// GetErrorWriter returns opts.ErrorWriter.
func (opts *Options) GetErrorWriter() io.Writer {
	return opts.ErrorWriter
}

// GetDatasetName returns the name of dataset from the dataset path.
func (opts *Options) GetDatasetName() string {
	return filepath.Base(filepath.Dir(opts.DatasetPath))
}

// GetLabelNext returns opts.LabelNext.
func (opts *Options) GetLabelNext() string {
	return opts.LabelNext
}

// GetLogger returns opts.Logger.
func (opts *Options) GetLogger() *log.Logger {
	return opts.Logger
}

// GetOutputWriter returns opts.OutputWriter.
func (opts *Options) GetOutputWriter() io.Writer {
	return opts.OutputWriter
}

func (opts *Options) initializeFlagSet() {
	opts.flagSet = flag.NewFlagSet("sticker-util", flag.ContinueOnError)
	opts.flagSet.SetOutput(ioutil.Discard)
	opts.flagSet.StringVar(&opts.CPUProfile, "cpuprofile", opts.CPUProfile, "Specify the CPU profile filename")
	opts.flagSet.BoolVar(&opts.Debug, "debug", opts.Debug, "Turn on debug logging")
	opts.flagSet.StringVar(&opts.FeatureMapName, "featureMap", opts.FeatureMapName, "Specify the feature map filename")
	opts.flagSet.BoolVar(&opts.Help, "h", opts.Help, "Show the help and exit")
	opts.flagSet.BoolVar(&opts.Help, "help", opts.Help, "Show the help and exit")
	opts.flagSet.StringVar(&opts.HTTPResource, "httpResource", opts.HTTPResource, "Specify the HTTP server resource root path")
	opts.flagSet.StringVar(&opts.LabelBoost, "labelboost", opts.LabelBoost, "Specify the .labelboost filename")
	opts.flagSet.StringVar(&opts.LabelConst, "labelconst", opts.LabelConst, "Specify the .labelconst filename")
	opts.flagSet.StringVar(&opts.LabelForest, "labelforest", opts.LabelForest, "Specify the .labelforest filename")
	opts.flagSet.StringVar(&opts.LabelNear, "labelnear", opts.LabelNear, "Specify the .labelnear filename")
	opts.flagSet.StringVar(&opts.LabelNearest, "labelnearest", opts.LabelNearest, "Specify the .labelnearest filename")
	opts.flagSet.StringVar(&opts.LabelNext, "labelnext", opts.LabelNext, "Specify the .labelnext filename")
	opts.flagSet.StringVar(&opts.LabelOne, "labelone", opts.LabelOne, "Specify the .labelone filename")
	opts.flagSet.StringVar(&opts.LabelMapName, "labelMap", opts.LabelMapName, "Specify the label map filename")
	opts.flagSet.BoolVar(&opts.Verbose, "verbose", opts.Verbose, "Log verbosely")
}

// LabelMap returns the label name.
func (opts *Options) LabelMap(label uint32, quote bool) string {
	if label < uint32(len(opts.labelMap)) {
		if quote {
			return fmt.Sprintf("%q", opts.labelMap[label])
		}
		return opts.labelMap[label]
	}
	return fmt.Sprintf("%d", label)
}

// Parse parses the flags in args.
//
// This function returns an error in parsing.
func (opts *Options) Parse(args []string) error {
	opts.initializeFlagSet()
	if err := opts.flagSet.Parse(args); err != nil {
		return err
	}
	if opts.Help {
		return nil
	}
	args = opts.flagSet.Args()
	if len(args) == 0 {
		return fmt.Errorf("specify the path of dataset")
	}
	opts.DatasetPath = args[0]
	if !strings.HasSuffix(opts.DatasetPath, "/") {
		opts.DatasetPath = filepath.FromSlash(opts.DatasetPath + "/")
	}
	opts.InspectForest, opts.TestForests, opts.TrainForest = nil, []*TestForestCommand{}, nil
	args = args[1:]
	var err error
	for len(args) > 0 {
		cmd := args[0]
		args = args[1:]
		switch cmd {
		case "@compareForest":
			if opts.CompareForest != nil {
				return fmt.Errorf("cannot specify multiple @compareForest commands")
			}
			opts.CompareForest = NewCompareForestCommand(opts)
			if args, err = opts.CompareForest.Parse(args); err != nil {
				return fmt.Errorf("@compareForest: %s", err)
			}
		case "@inspectForest":
			if opts.InspectForest != nil {
				return fmt.Errorf("cannot specify multiple @inspectForest commands")
			}
			opts.InspectForest = NewInspectForestCommand(opts)
			if args, err = opts.InspectForest.Parse(args); err != nil {
				return fmt.Errorf("@inspectForest: %s", err)
			}
		case "@inspectOne":
			if opts.InspectOne != nil {
				return fmt.Errorf("cannot specify multiple @inspectOne commands")
			}
			opts.InspectOne = NewInspectOneCommand(opts)
			if args, err = opts.InspectOne.Parse(args); err != nil {
				return fmt.Errorf("@inspectOne: %s", err)
			}
		case "@pruneOne":
			if opts.PruneOne != nil {
				return fmt.Errorf("cannot specify multiple @pruneOne commands")
			}
			opts.PruneOne = NewPruneOneCommand(opts)
			if args, err = opts.PruneOne.Parse(args); err != nil {
				return fmt.Errorf("@pruneOne: %s", err)
			}
		case "@shuffle":
			if opts.Shuffle != nil {
				return fmt.Errorf("cannot specify multiple @shuffle commands")
			}
			opts.Shuffle = NewShuffleCommand(opts)
			if args, err = opts.Shuffle.Parse(args); err != nil {
				return fmt.Errorf("@shuffle: %s", err)
			}
		case "@summarize":
			if opts.Summarize != nil {
				return fmt.Errorf("cannot specify multiple @summarize commands")
			}
			opts.Summarize = NewSummarizeCommand(opts)
			if args, err = opts.Summarize.Parse(args); err != nil {
				return fmt.Errorf("@summarize: %s", err)
			}
		case "@testBoost":
			cmd := NewTestBoostCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testBoost: %s", err)
			}
			opts.TestBoosts = append(opts.TestBoosts, cmd)
		case "@testConst":
			cmd := NewTestConstCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testConst: %s", err)
			}
			opts.TestConsts = append(opts.TestConsts, cmd)
		case "@testForest":
			cmd := NewTestForestCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testForest: %s", err)
			}
			opts.TestForests = append(opts.TestForests, cmd)
		case "@testNear":
			cmd := NewTestNearCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testNear: %s", err)
			}
			opts.TestNears = append(opts.TestNears, cmd)
		case "@testNearest":
			cmd := NewTestNearestCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testNearest: %s", err)
			}
			opts.TestNearests = append(opts.TestNearests, cmd)
		case "@testNext":
			cmd := next.NewTestCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testNext: %s", err)
			}
			opts.TestNexts = append(opts.TestNexts, cmd)
		case "@testOne":
			cmd := NewTestOneCommand(opts)
			if args, err = cmd.Parse(args); err != nil {
				return fmt.Errorf("@testOne: %s", err)
			}
			opts.TestOnes = append(opts.TestOnes, cmd)
		case "@trainBoost":
			if opts.TrainBoost != nil {
				return fmt.Errorf("cannot specify multiple @trainBoost commands")
			}
			opts.TrainBoost = NewTrainBoostCommand(opts)
			if args, err = opts.TrainBoost.Parse(args); err != nil {
				return fmt.Errorf("@trainBoost: %s", err)
			}
		case "@trainConst":
			if opts.TrainConst != nil {
				return fmt.Errorf("cannot specify multiple @trainConst commands")
			}
			opts.TrainConst = NewTrainConstCommand(opts)
			if args, err = opts.TrainConst.Parse(args); err != nil {
				return fmt.Errorf("@trainConst: %s", err)
			}
		case "@trainForest":
			if opts.TrainForest != nil {
				return fmt.Errorf("cannot specify multiple @trainForest commands")
			}
			opts.TrainForest = NewTrainForestCommand(opts)
			if args, err = opts.TrainForest.Parse(args); err != nil {
				return fmt.Errorf("@trainForest: %s", err)
			}
		case "@trainNear":
			if opts.TrainNear != nil {
				return fmt.Errorf("cannot specify multiple @trainNear commands")
			}
			opts.TrainNear = NewTrainNearCommand(opts)
			if args, err = opts.TrainNear.Parse(args); err != nil {
				return fmt.Errorf("@trainNear: %s", err)
			}
		case "@trainNearest":
			if opts.TrainNearest != nil {
				return fmt.Errorf("cannot specify multiple @trainNearest commands")
			}
			opts.TrainNearest = NewTrainNearestCommand(opts)
			if args, err = opts.TrainNearest.Parse(args); err != nil {
				return fmt.Errorf("@trainNearest: %s", err)
			}
		case "@trainNext":
			if opts.TrainNext != nil {
				return fmt.Errorf("cannot specify multiple @trainNext commands")
			}
			opts.TrainNext = next.NewTrainCommand(opts)
			if args, err = opts.TrainNext.Parse(args); err != nil {
				return fmt.Errorf("@trainNext: %s", err)
			}
		case "@trainOne":
			if opts.TrainOne != nil {
				return fmt.Errorf("cannot specify multiple @trainOne commands")
			}
			opts.TrainOne = NewTrainOneCommand(opts)
			if args, err = opts.TrainOne.Parse(args); err != nil {
				return fmt.Errorf("@trainOne: %s", err)
			}
		default:
			return fmt.Errorf("unknown command: %s", cmd)
		}
	}
	return nil
}

// ReadDataset reads the dataset of the given table name in the dataset path.
//
// This function returns an error in reading the dataset.
func (opts *Options) ReadDataset(tblname string) (*sticker.Dataset, error) {
	filename := filepath.Join(opts.DatasetPath, tblname)
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadDataset: %s: %s", filename, err)
	}
	defer file.Close()
	ds, err := sticker.ReadTextDataset(file)
	if err != nil {
		return nil, fmt.Errorf("ReadDataset: %s: %s", filename, err)
	}
	return ds, nil
}

// ReadDatasets reads the multiple datasets with at most maxentries data entries.
// If sampling is true, then the data entries are randomly sampled without replacement.
//
// This function returns an error in reading the datasets, or tblnames is empty.
func (opts *Options) ReadDatasets(tblnames []string, maxentries uint, sampling bool) (*sticker.Dataset, error) {
	dsname := opts.GetDatasetName()
	ds := &sticker.Dataset{
		X: sticker.FeatureVectors{},
		Y: sticker.LabelVectors{},
	}
	if len(tblnames) == 0 {
		return nil, fmt.Errorf("specify the table names")
	}
	for _, tblname := range tblnames {
		opts.Logger.Printf("loading table %q of dataset %q ...", tblname, dsname)
		subds, err := opts.ReadDataset(tblname)
		if err != nil {
			return nil, err
		}
		ds.X, ds.Y = append(ds.X, subds.X...), append(ds.Y, subds.Y...)
	}
	if sampling {
		rng := rand.New(rand.NewSource(0))
		for i := 0; i < int(maxentries); i++ {
			j := i + rng.Intn(ds.Size()-i)
			ds.X[i], ds.X[j], ds.Y[i], ds.Y[j] = ds.X[j], ds.X[i], ds.Y[j], ds.Y[i]
		}
	}
	if maxentries < uint(ds.Size()) {
		ds.X, ds.Y = ds.X[:maxentries], ds.Y[:maxentries]
	}
	return ds, nil
}

// RunHTTPServer runs a HTTP server.
// initializer registers handlers to http.ServeMux via handleFunc(), and can return an error.
// handleFunc takes the function name used in log as prefix, URL pattern, and handler.
func (opts *Options) RunHTTPServer(addr string, initializer func(handleFunc func(prefix string, pattern string, handler func(writer http.ResponseWriter, req *http.Request))) error) error {
	mux := http.NewServeMux()
	handleFunc := func(prefix string, pattern string, handler func(writer http.ResponseWriter, req *http.Request)) {
		mux.HandleFunc(pattern, func(writer http.ResponseWriter, req *http.Request) {
			opts.Logger.Printf("%s: %s %s", prefix, req.Method, req.URL)
			handler(writer, req)
		})
	}
	if err := initializer(handleFunc); err != nil {
		return err
	}
	handleFunc("ResourceHandler", "/res/", func(writer http.ResponseWriter, req *http.Request) {
		http.StripPrefix("/res/", http.FileServer(http.Dir(opts.HTTPResource))).ServeHTTP(writer, req)
	})
	handleFunc("NotFoundHandler", "/", func(writer http.ResponseWriter, req *http.Request) {
		http.NotFoundHandler().ServeHTTP(writer, req)
	})
	opts.Logger.Printf("server is listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

// ReadMap reads the map file in the dataset path.
// If mapname is "", then the map is empty.
//
// This function returns an error in reading the file.
func (opts *Options) ReadMap(mapname string) ([]string, error) {
	if mapname == "" {
		return make([]string, 0), nil
	}
	filename := filepath.Join(opts.DatasetPath, mapname)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadMap: %s: %s", filename, err)
	}
	m := strings.Split(string(b), "\n")
	for i, s := range m {
		m[i] = strings.TrimSpace(s)
	}
	return m, nil
}

// Run runs the specified sub-commands.
func (opts *Options) Run() error {
	if opts.Help {
		opts.ShowHelp()
		return nil
	}
	if opts.Verbose {
		opts.Logger = log.New(opts.OutputWriter, "!LOG   ", log.LstdFlags)
		opts.Logger.Printf("turned on verbose logging")
	} else {
		opts.Logger = log.New(ioutil.Discard, "", 0)
	}
	opts.Logger.Printf("Options: %#v", opts)
	if opts.Debug {
		opts.DebugLogger = log.New(opts.OutputWriter, "!DEBUG ", log.LstdFlags)
		opts.DebugLogger.Printf("turned on debug logging")
	} else {
		opts.DebugLogger = log.New(ioutil.Discard, "", 0)
	}
	if opts.CPUProfile != "" {
		opts.Logger.Printf("starting CPU profiling on %s ...", opts.CPUProfile)
		f, err := os.Create(opts.CPUProfile)
		if err != nil {
			return fmt.Errorf("%s: %s", opts.CPUProfile, err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	opts.Logger.Printf("loading feature map from %q ...", opts.FeatureMapName)
	var err error
	opts.featureMap, err = opts.ReadMap(opts.FeatureMapName)
	if err != nil {
		return err
	}
	opts.Logger.Printf("loading label map from %q ...", opts.LabelMapName)
	opts.labelMap, err = opts.ReadMap(opts.LabelMapName)
	if err != nil {
		return err
	}
	if opts.Shuffle != nil {
		startTime := time.Now()
		if err := opts.Shuffle.Run(); err != nil {
			return fmt.Errorf("@shuffle: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @shuffle in %s", finishTime.Sub(startTime))
	}
	if opts.Summarize != nil {
		if err := opts.Summarize.Run(); err != nil {
			return err
		}
	}
	if opts.TrainBoost != nil {
		startTime := time.Now()
		if err := opts.TrainBoost.Run(); err != nil {
			return fmt.Errorf("@trainBoost: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainBoost in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if opts.TrainConst != nil {
		startTime := time.Now()
		if err := opts.TrainConst.Run(); err != nil {
			return fmt.Errorf("@trainConst: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainConst in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if opts.TrainForest != nil {
		startTime := time.Now()
		if err := opts.TrainForest.Run(); err != nil {
			return fmt.Errorf("@trainForest: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainForest in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if opts.TrainNear != nil {
		startTime := time.Now()
		if err := opts.TrainNear.Run(); err != nil {
			return fmt.Errorf("@trainNear: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainNear in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if opts.TrainNearest != nil {
		startTime := time.Now()
		if err := opts.TrainNearest.Run(); err != nil {
			return fmt.Errorf("@trainNearest: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainNearest in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if opts.TrainNext != nil {
		startTime := time.Now()
		if err := opts.TrainNext.Run(); err != nil {
			return fmt.Errorf("@trainNext: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainNext in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if opts.TrainOne != nil {
		startTime := time.Now()
		if err := opts.TrainOne.Run(); err != nil {
			return fmt.Errorf("@trainOne: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @trainOne in %s", finishTime.Sub(startTime))
		debug.FreeOSMemory()
	}
	if len(opts.TestBoosts) > 0 {
		for i, cmd := range opts.TestBoosts {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testBoost: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testBoost in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if len(opts.TestConsts) > 0 {
		for i, cmd := range opts.TestConsts {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testConst: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testConst in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if len(opts.TestForests) > 0 {
		for i, cmd := range opts.TestForests {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testForest: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testForest in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if len(opts.TestNears) > 0 {
		for i, cmd := range opts.TestNears {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testNear: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testNear in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if len(opts.TestNearests) > 0 {
		for i, cmd := range opts.TestNearests {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testNearest: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testNearest in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if len(opts.TestNexts) > 0 {
		for i, cmd := range opts.TestNexts {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testNext: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testNext in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if len(opts.TestOnes) > 0 {
		for i, cmd := range opts.TestOnes {
			startTime := time.Now()
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("@testOne: %s", err)
			}
			finishTime := time.Now()
			opts.Logger.Printf("finished #%d @testOne in %s", i+1, finishTime.Sub(startTime))
			debug.FreeOSMemory()
		}
	}
	if opts.PruneOne != nil {
		startTime := time.Now()
		if err := opts.PruneOne.Run(); err != nil {
			return fmt.Errorf("@pruneOne: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @pruneOne in %s", finishTime.Sub(startTime))
	}
	if opts.CompareForest != nil {
		startTime := time.Now()
		if err := opts.CompareForest.Run(); err != nil {
			return fmt.Errorf("@compareForest: %s", err)
		}
		finishTime := time.Now()
		opts.Logger.Printf("finished @compareForest in %s", finishTime.Sub(startTime))
	}
	if opts.InspectForest != nil {
		if err := opts.InspectForest.Run(); err != nil {
			return fmt.Errorf("@inspectForest: %s", err)
		}
	}
	if opts.InspectOne != nil {
		if err := opts.InspectOne.Run(); err != nil {
			return fmt.Errorf("@inspectOne: %s", err)
		}
	}
	return nil
}

// SetLabelNext sets the given string to opts.LabelNext.
func (opts *Options) SetLabelNext(value string) {
	opts.LabelNext = value
}

// ShowHelp shows the help.
func (opts *Options) ShowHelp() {
	fmt.Fprintf(opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: %s [commonOptions] datasetPath (@{compareForest|inspectForest|inspectOne|pruneOne|shuffle|summarize|trainBoost|trainConst|trainForest|trainNear|trainNearest|trainNew|trainOne|testBoost|testConst|testForest|testNear|testNearest|testNext|testOne} [subCommandOptions])*\n", opts.execpath)
	if opts.flagSet == nil {
		opts.initializeFlagSet()
	}
	opts.flagSet.SetOutput(opts.ErrorWriter)
	opts.flagSet.PrintDefaults()
	opts.flagSet.SetOutput(ioutil.Discard)
}
