// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package anlzr provides the necessary API for making source code analyzis.
package anlzr

import "github.com/DevMine/srcanlzr/src"

// An Analyzer represents a source code analyzer which extracts one or more
// metrics from the structure defined by the src package.
type Analyzer interface {
	// Analyze performs an source code analysis on a project and adds the
	// resulting metrics into r.
	Analyze(p *src.Project, r *Result) error
}

// A Result holds all source code analysis metrics output by srcanlzr.
type Result struct {
	ProgLangs      []Language        `json:"programming_languages" xml:"programming-languages"`
	AverageFuncLen float32           `json:"average_function_length" xml:"average-function-length"`
	MaxFuncLen     int64             `json:"max_function_length" xml:"max-function-length"`
	MinFuncLen     int64             `json:"min_function_length" xml:"min-function-length"`
	MedianFuncLen  int64             `json:"median_function_length" xml:"median-function-length"`
	TotalLoC       int64             `json:"total_loc" xml:"total-loc"`
	Complexity     ComplexityMetrics `json:"complexity" xml:"complexity"`
	DocCoverage    CommentRatios     `json:"documentation_coverage" xml:"documentation-coverage"`
}

// A Language represents a programming language used by the project.
type Language struct {
	src.Language
	Lines int64 `json:"lines" xml:"lines"`
}

// CyclomaticComplexity metrics, also known as McCabe metric.
type ComplexityMetrics struct {
	AveragePerFunc float32 `json:"average_per_func" xml:"average-per-func"` // Average complexity per function.
	AveragePerFile float32 `json:"average_per_file" xml:"average-per-file"` // Average complexity per file.
}

// CommentRatios contains various ratios about documentation coverage.
type CommentRatios struct {
	TypeComRatio   float32 `json:"type_comment_ratio"`
	StructComRatio float32 `json:"structure_comment_ratio"`
	ConstComRatio  float32 `json:"constant_comment_ratio"`
	VarsComRatio   float32 `json:"variable_comment_ratio"`
	FuncComRatio   float32 `json:"function_comment_ratio"`
	InterComRatio  float32 `json:"interface_comment_ratio"`
	ClassComRatio  float32 `json:"class_comment_ratio"`
	MethComRatio   float32 `json:"method_comment"`
	AttrComRatio   float32 `json:"attribute_comment_ratio"`
	EnumComRatio   float32 `json:"enumeration_comment_ratio"`
}

// RunAnalyzers runs several analyzers on a project.
func RunAnalyzers(p *src.Project, a ...Analyzer) (*Result, error) {
	r := &Result{
		ProgLangs:      []Language{},
		AverageFuncLen: -1.0,
		MaxFuncLen:     -1,
		MinFuncLen:     -1,
		MedianFuncLen:  -1,
		TotalLoC:       -1,
		Complexity:     ComplexityMetrics{AveragePerFunc: -1, AveragePerFile: -1},
		DocCoverage:    CommentRatios{},
	}

	for _, anlzr := range a {
		err := anlzr.Analyze(p, r)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}
