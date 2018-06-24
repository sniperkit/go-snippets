package testdata

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bsm/reason/core"
)

var BigClassificationModel = core.NewModel(
	core.NewCategoricalFeature("c1", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c2", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c3", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c4", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewCategoricalFeature("c5", []string{"v1", "v2", "v3", "v4", "v5"}),
	core.NewNumericalFeature("n1"),
	core.NewNumericalFeature("n2"),
	core.NewNumericalFeature("n3"),
	core.NewNumericalFeature("n4"),
	core.NewNumericalFeature("n5"),
	core.NewCategoricalFeature("target", []string{"c1", "c2"}),
)

var BigRegressionModel = core.NewModel(
	core.NewCategoricalFeatureHashBuckets("c1", 1000),
	core.NewCategoricalFeatureHashBuckets("c2", 100),
	core.NewCategoricalFeature("c3", []string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"}),
	core.NewCategoricalFeatureHashBuckets("c4", 1000),
	core.NewNumericalFeature("n1"),
	core.NewNumericalFeature("target"),
)

var (
	bigClassificationFieldIndices = map[string]int{"c1": 0, "c2": 1, "c3": 2, "c4": 3, "c5": 4, "n1": 5, "n2": 6, "n3": 7, "n4": 8, "n5": 9, "target": 10}
	bigRegressionFieldIndices     = map[string]int{"c1": 0, "c2": 1, "c3": 2, "c4": 3, "n1": 4, "target": 5}
)

type BigDataStream struct {
	file  *os.File
	recs  *csv.Reader
	model *core.Model
	fix   map[string]int

	x   core.MapExample
	err error
}

func OpenClassification(root string) (*BigDataStream, *core.Model, error) {
	return open(filepath.Join(root, "bigcls.csv"), BigClassificationModel, bigClassificationFieldIndices)
}

func OpenRegression(root string) (*BigDataStream, *core.Model, error) {
	return open(filepath.Join(root, "bigreg.csv"), BigRegressionModel, bigRegressionFieldIndices)
}

func open(fname string, model *core.Model, fix map[string]int) (*BigDataStream, *core.Model, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, nil, err
	}

	recs := csv.NewReader(f)
	recs.FieldsPerRecord = len(model.Features)

	return &BigDataStream{
		file:  f,
		recs:  recs,
		model: model,
		fix:   fix,
		x:     make(core.MapExample, recs.FieldsPerRecord),
	}, model, nil
}

func (s *BigDataStream) Next() bool {
	if s.err != nil {
		return false
	}

	fields, err := s.recs.Read()
	if err != nil {
		s.err = err
		return false
	}

	// init example
	s.x = make(core.MapExample, len(s.x))

	// read predictors
	for name, feat := range s.model.Features {
		str := fields[s.fix[name]]
		if str == "?" {
			continue
		} else if feat.Kind.IsCategorical() {
			s.x[name] = str
		} else if s.x[name], err = strconv.ParseFloat(str, 64); err != nil {
			s.err = err
			return false
		}
	}
	return true
}

func (s *BigDataStream) ReadN(n int) ([]core.Example, error) {
	res := make([]core.Example, 0, n)
	for s.Next() {
		res = append(res, s.Example())
		if len(res) == n {
			break
		}
	}
	return res, s.Err()
}

func (s *BigDataStream) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

func (s *BigDataStream) Example() core.Example { return s.x }
func (s *BigDataStream) Close() error          { return s.file.Close() }
