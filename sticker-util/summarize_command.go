package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hiro4bbh/sticker"
)

// SummarizeCommand have flags for summarize sub-command.
type SummarizeCommand struct {
	Addr      string
	Help      bool
	TableName string

	opts    *Options
	flagSet *flag.FlagSet
}

// NewSummarizeCommand returns a new SummarizeCommand.
func NewSummarizeCommand(opts *Options) *SummarizeCommand {
	return &SummarizeCommand{
		Addr:      ":8080",
		Help:      false,
		TableName: "",
		opts:      opts,
	}
}

func (cmd *SummarizeCommand) initializeFlagSet() {
	cmd.flagSet = flag.NewFlagSet("@summarize", flag.ContinueOnError)
	cmd.flagSet.Usage = func() {}
	cmd.flagSet.SetOutput(ioutil.Discard)
	cmd.flagSet.StringVar(&cmd.Addr, "addr", cmd.Addr, "Specify the HTTP server address")
	cmd.flagSet.BoolVar(&cmd.Help, "h", cmd.Help, "Show the help and exit")
	cmd.flagSet.BoolVar(&cmd.Help, "help", cmd.Help, "Show the help and exit")
	cmd.flagSet.StringVar(&cmd.TableName, "table", cmd.TableName, "Specify the table name")
}

// Parse parses the flags in args, and returns the remain parts of args.
//
// This function returns an error in parsing.
func (cmd *SummarizeCommand) Parse(args []string) ([]string, error) {
	cmd.initializeFlagSet()
	if err := cmd.flagSet.Parse(args); err != nil {
		return nil, err
	}
	return cmd.flagSet.Args(), nil
}

// Run prints the summary of the specified table of the dataset.
func (cmd *SummarizeCommand) Run() error {
	if cmd.Help {
		cmd.ShowHelp()
		return nil
	}
	opts := cmd.opts
	opts.Logger.Printf("SummarizeCommands: %#v", cmd)
	dsname := opts.GetDatasetName()
	opts.Logger.Printf("loading table %q from dataset %q ...", cmd.TableName, dsname)
	ds, err := opts.ReadDataset(cmd.TableName)
	if err != nil {
		return err
	}
	nentries := ds.Size()
	opts.Logger.Printf("collecting feature and label frequencies ...")
	labelInvMap := make(map[string]uint32)
	for id, label := range opts.labelMap {
		labelInvMap[label] = uint32(id)
	}
	nfeatures, nlabels := ds.X.Dim(), ds.Y.Dim()
	featureActVec, labelCountVec := make([]float32, 0, nentries), make([]float32, 0, nentries)
	featureFreqs, labelFreqs := make(sticker.SparseVector), make(sticker.SparseVector)
	for i, xi := range ds.X {
		yi := ds.Y[i]
		featureActVec, labelCountVec = append(featureActVec, float32(len(xi))), append(labelCountVec, float32(len(yi)))
		for _, xipair := range xi {
			featureFreqs[xipair.Key]++
		}
		for _, label := range yi {
			labelFreqs[label]++
		}
	}
	featureActsMin, featureActsQ25, featureActsMed, featureActsQ75, featureActsMax, featureActsAvg := sticker.SummarizeFloat32Slice(featureActVec)
	featureActs := map[string]float32{
		"min": featureActsMin,
		"q25": featureActsQ25,
		"med": featureActsMed,
		"q75": featureActsQ75,
		"max": featureActsMax,
		"avg": featureActsAvg,
	}
	labelCountsMin, labelCountsQ25, labelCountsMed, labelCountsQ75, labelCountsMax, labelCountsAvg := sticker.SummarizeFloat32Slice(labelCountVec)
	labelCounts := map[string]float32{
		"min": labelCountsMin,
		"q25": labelCountsQ25,
		"med": labelCountsMed,
		"q75": labelCountsQ75,
		"max": labelCountsMax,
		"avg": labelCountsAvg,
	}
	K := 20
	leastRanksSet := make([][]int, K)
	opts.Logger.Printf("estimating the least ranks for completely predicting top-%d labels ...", K)
	labelRankL := sticker.RankTopK(labelFreqs, uint(len(labelFreqs)))
	labelInvRankL := sticker.InvertRanks(labelRankL)
	for k := range leastRanksSet {
		leastRanksSet[k] = []int{}
	}
	for _, yi := range ds.Y {
		labelRanks := make(sticker.KeyValues32OrderedByValue, 0, len(yi))
		for _, label := range yi {
			labelRanks = append(labelRanks, sticker.KeyValue32{label, float32(labelInvRankL[label])})
		}
		sort.Sort(labelRanks)
		for k := 0; k < K; k++ {
			kmax := k
			if kmax >= len(labelRanks) {
				kmax = len(labelRanks) - 1
			}
			leastRanksSet[k] = append(leastRanksSet[k], int(labelRanks[kmax].Value))
		}
	}
	for _, leastRanks := range leastRanksSet {
		sort.Ints(leastRanks)
	}
	if err := opts.RunHTTPServer(cmd.Addr, func(handleFunc func(prefix string, pattern string, handler func(writer http.ResponseWriter, req *http.Request))) error {
		handleFunc("RunSummarize", "/featureFreq.csv", func(writer http.ResponseWriter, req *http.Request) {
			_, withName := req.URL.Query()["withName"]
			featureFreqsKV := sticker.KeyValues32OrderedByValue{}
			for feature, freq := range featureFreqs {
				featureFreqsKV = append(featureFreqsKV, sticker.KeyValue32{feature, float32(freq)})
			}
			sort.Sort(sort.Reverse(featureFreqsKV))
			if withName {
				fmt.Fprintf(writer, "id,name,freq\n")
				for _, featureFreq := range featureFreqsKV {
					fmt.Fprintf(writer, "%d,%s,%g\n", featureFreq.Key, opts.FeatureMap(featureFreq.Key, true), featureFreq.Value)
				}
			} else {
				fmt.Fprintf(writer, "id,freq\n")
				for _, featureFreq := range featureFreqsKV {
					fmt.Fprintf(writer, "%d,%g\n", featureFreq.Key, featureFreq.Value)
				}
			}
		})
		handleFunc("RunSummarize", "/search.json", func(writer http.ResponseWriter, req *http.Request) {
			countStrs, ok := req.URL.Query()["count"]
			if !ok {
				countStrs = []string{"50"}
			}
			countStr := countStrs[0]
			count, err := strconv.ParseUint(countStr, 10, 64)
			if err != nil {
				writer.Write([]byte(fmt.Sprintf("{\"status\": \"error\", \"message\": \"illegal count\"}")))
				return
			}
			labelsStrs, ok := req.URL.Query()["labels"]
			if !ok {
				labelsStrs = []string{""}
			}
			labelsStr := strings.TrimSpace(labelsStrs[0])
			labels := regexp.MustCompile("\".+\"|\\S+").FindAllString(labelsStr, -1)
			if labels == nil {
				result := map[string]interface{}{
					"status":      "success",
					"ncandidates": nentries,
					"nhits":       0,
					"entries":     []map[string]interface{}{},
				}
				if b, err := json.Marshal(result); err == nil {
					writer.Write(b)
				} else {
					fmt.Fprintf(writer, "{\"status\": \"error\", \"message\": %q}", err)
				}
				return
			}
			for i, label := range labels {
				switch label[0] {
				case '"':
					label, err := strconv.Unquote(label)
					if err != nil {
						writer.Write([]byte(fmt.Sprintf("{\"status\": \"error\", \"message\": %q}", fmt.Sprintf("#%d label %q: %s", i, label, err))))
						return
					}
					labels[i] = label
				}
			}
			labelsMap := make(map[uint32]struct{})
			for _, labelName := range labels {
				label, ok := labelInvMap[labelName]
				if !ok {
					label64, err := strconv.ParseUint(labelName, 10, 64)
					if err != nil {
						writer.Write([]byte(fmt.Sprintf("{\"status\": \"error\", \"message\": %q}", fmt.Sprintf("unknown label %q", labelName))))
						return
					}
					label = uint32(label64)
				}
				labelsMap[label] = struct{}{}
			}
			entries := make([]map[string]interface{}, 0, count)
			nhits := 0
			for i, yi := range ds.Y {
				if len(labelsMap) > 0 {
					nmatcheds := 0
					for _, label := range yi {
						if _, ok := labelsMap[label]; ok {
							nmatcheds++
						}
					}
					if nmatcheds < len(labelsMap) {
						continue
					}
				}
				nhits++
				if uint64(len(entries)) >= count {
					continue
				}
				xi := ds.X[i]
				features := make(map[string]float32, len(xi))
				for _, xipair := range xi {
					features[opts.FeatureMap(xipair.Key, false)] = xipair.Value
				}
				labels := make(map[string]string, len(yi))
				for _, label := range yi {
					labelType := "none"
					if _, ok := labelsMap[label]; ok {
						labelType = "highlighted"
					}
					labels[opts.LabelMap(label, false)] = labelType
				}
				entry := map[string]interface{}{
					"i": i,
					"x": features,
					"y": labels,
				}
				entries = append(entries, entry)
			}
			result := map[string]interface{}{
				"status":      "success",
				"ncandidates": nentries,
				"nhits":       nhits,
				"entries":     entries,
			}
			if b, err := json.Marshal(result); err == nil {
				writer.Write(b)
			} else {
				fmt.Fprintf(writer, "{\"status\": \"error\", \"message\": %q}", err)
			}
		})
		handleFunc("RunSummarize", "/labelFreq.csv", func(writer http.ResponseWriter, req *http.Request) {
			_, withName := req.URL.Query()["withName"]
			labelFreqsKV := sticker.KeyValues32OrderedByValue{}
			for label, freq := range labelFreqs {
				labelFreqsKV = append(labelFreqsKV, sticker.KeyValue32{label, float32(freq)})
			}
			sort.Sort(sort.Reverse(labelFreqsKV))
			if withName {
				fmt.Fprintf(writer, "id,name,freq\n")
				for _, labelFreq := range labelFreqsKV {
					fmt.Fprintf(writer, "%d,%s,%g\n", labelFreq.Key, opts.LabelMap(labelFreq.Key, true), labelFreq.Value)
				}
			} else {
				fmt.Fprintf(writer, "id,freq\n")
				for _, labelFreq := range labelFreqsKV {
					fmt.Fprintf(writer, "%d,%g\n", labelFreq.Key, labelFreq.Value)
				}
			}
		})
		handleFunc("RunSummarize", "/summary", func(writer http.ResponseWriter, req *http.Request) {
			tmpl := template.New("")
			tmpl.Funcs(TemplateStandardFunctions)
			tmplFilename := filepath.Join(opts.HTTPResource, "summaryTemplate.html")
			if _, err := tmpl.ParseFiles(tmplFilename); err != nil {
				writer.Write([]byte(fmt.Sprintf("template.ParseFiles: %s: %s", tmplFilename, err)))
				return
			}
			if err := tmpl.ExecuteTemplate(writer, "summaryTemplate.html", map[string]interface{}{
				"dsname":        dsname,
				"nfeatures":     nfeatures,
				"nlabels":       nlabels,
				"nentries":      nentries,
				"tblname":       cmd.TableName,
				"featureActs":   featureActs,
				"labelCounts":   labelCounts,
				"leastRanksSet": leastRanksSet,
			}); err != nil {
				writer.Write([]byte(fmt.Sprintf("template.ExecuteTemplate: %s", err)))
				return
			}
		})
		return nil
	}); err != nil {
		return fmt.Errorf("RunInspect: RunHTTPServer: %s", err)
	}
	return nil
}

// ShowHelp shows the help.
func (cmd *SummarizeCommand) ShowHelp() {
	fmt.Fprintf(cmd.opts.ErrorWriter, "sticker-util\nCopyright 2017- Tatsuhiro Aoshima (hiro4bbh@gmail.com).\n\nUsage: @summarize [subCommandOptions]\n")
	if cmd.flagSet == nil {
		cmd.initializeFlagSet()
	}
	cmd.flagSet.SetOutput(cmd.opts.ErrorWriter)
	cmd.flagSet.PrintDefaults()
	cmd.flagSet.SetOutput(ioutil.Discard)
}
