package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hiro4bbh/sticker"
	"github.com/hiro4bbh/sticker/plugin"
)

// OptionFloat32 is the data structure for using float32 in flag.
// This implements interface flag.Value.
type OptionFloat32 float32

// String is for interface flag.Value.
func (opt *OptionFloat32) String() string {
	return fmt.Sprintf("%v", float32(*opt))
}

// Set is for interface flag.Value.
func (opt *OptionFloat32) Set(value string) error {
	valueF64, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return err
	}
	*opt = OptionFloat32(float32(valueF64))
	return nil
}

// OptionStrings is the data structure for using string slice in flag.
// This implements interface flag.Value.
type OptionStrings struct {
	// WithDefault indicates whether Values are the default values or not.
	// If true, Values are cleared when calling Set.
	WithDefault bool
	// Values is the list of the values.
	Values []string
}

// String is for interface flag.Value.
func (opt *OptionStrings) String() string {
	return fmt.Sprintf("%v", []string(opt.Values))
}

// Set is for interface flag.Value.
func (opt *OptionStrings) Set(value string) error {
	if opt.WithDefault {
		opt.Values = []string{value}
		opt.WithDefault = false
	} else {
		opt.Values = append(opt.Values, value)
	}
	return nil
}

// OptionUints is the data structure for using uint slice in flag.
// This implements interface flag.Value.
type OptionUints struct {
	// WithDefault indicates whether Values are the default values or not.
	// If true, Values are cleared when calling Set.
	WithDefault bool
	// Values is the list of the values.
	Values []uint
}

// String is for interface flag.Value.
func (opt *OptionUints) String() string {
	return fmt.Sprintf("%v", []uint(opt.Values))
}

// Set is for interface flag.Value.
func (opt *OptionUints) Set(value string) error {
	valueUint64, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	if opt.WithDefault {
		opt.Values = []uint{uint(valueUint64)}
		opt.WithDefault = false
	} else {
		opt.Values = append(opt.Values, uint(valueUint64))
	}
	return nil
}

// Options is an interface for sticker-util.Options.
type Options interface {
	FeatureMap(feature uint32, quote bool) string
	GetDatasetName() string
	GetDebugLogger() *log.Logger
	GetErrorWriter() io.Writer
	GetLabelNext() string
	GetLogger() *log.Logger
	GetOutputWriter() io.Writer
	LabelMap(label uint32, quote bool) string
	ReadDataset(tblname string) (*sticker.Dataset, error)
	SetLabelNext(labelNext string)
}

// CreateWithDir creates the parent directories if needed, the new file with the given filename.
//
// This function returns an error in creating the directories or file.
func CreateWithDir(filename string) (*os.File, error) {
	dirpath := filepath.Dir(filename)
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("%s: %s", dirpath, err)
	}
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", filename, err)
	}
	return file, nil
}

// JoinTableNames returns the joined table name from the given table names cut its file extension.
func JoinTableNames(tblnames []string) string {
	name := ""
	for i, tblname := range tblnames {
		if i > 0 {
			name += "_"
		}
		name += strings.TrimSuffix(tblname, filepath.Ext(tblname))
	}
	return name
}

// ReadLabelBoost reads the .labelboost model file.
func ReadLabelBoost(filename string) (*plugin.LabelBoost, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadLabelBoost: %s: %s", filename, err)
	}
	defer file.Close()
	var model plugin.LabelBoost
	if err := plugin.DecodeLabelBoost(&model, file); err != nil {
		return nil, fmt.Errorf("ReadLabelBoost: %s: %s", filename, err)
	}
	return &model, nil
}

// ReadLabelConst reads the .labelconst model file.
func ReadLabelConst(filename string) (*sticker.LabelConst, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadLabelConst: %s: %s", filename, err)
	}
	defer file.Close()
	var model sticker.LabelConst
	if err := sticker.DecodeLabelConst(&model, file); err != nil {
		return nil, fmt.Errorf("ReadLabelConst: %s: %s", filename, err)
	}
	return &model, nil
}

// ReadLabelForest reads the .labelforest model file.
func ReadLabelForest(filename string) (*plugin.LabelForest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadLabelForest: %s: %s", filename, err)
	}
	defer file.Close()
	var forest plugin.LabelForest
	if err := plugin.DecodeLabelForest(&forest, file); err != nil {
		return nil, fmt.Errorf("ReadLabelForest: %s: %s", filename, err)
	}
	return &forest, nil
}

// ReadLabelNear reads the .labelnear model file.
func ReadLabelNear(filename string) (*sticker.LabelNear, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadLabelNear: %s: %s", filename, err)
	}
	defer file.Close()
	var model sticker.LabelNear
	if err := sticker.DecodeLabelNear(&model, file); err != nil {
		return nil, fmt.Errorf("ReadLabelNear: %s: %s", filename, err)
	}
	return &model, nil
}

// ReadLabelNearest reads the .labelnearest model file.
func ReadLabelNearest(filename string) (*sticker.LabelNearest, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadLabelNearest: %s: %s", filename, err)
	}
	defer file.Close()
	var model sticker.LabelNearest
	if err := sticker.DecodeLabelNearest(&model, file); err != nil {
		return nil, fmt.Errorf("ReadLabelNearest: %s: %s", filename, err)
	}
	return &model, nil
}

// ReadLabelOne reads the .labelone model file.
func ReadLabelOne(filename string) (*sticker.LabelOne, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ReadLabelOne: %s: %s", filename, err)
	}
	defer file.Close()
	var model sticker.LabelOne
	if err := sticker.DecodeLabelOne(&model, file); err != nil {
		return nil, fmt.Errorf("ReadLabelOne: %s: %s", filename, err)
	}
	return &model, nil
}
