package common

import (
	"fmt"
	"io"
	"time"

	"github.com/hiro4bbh/sticker"
)

// ResultsReporter manages the precision@K and nDCG@K results on the given LabelVectors.
type ResultsReporter struct {
	_Y                     sticker.LabelVectors
	_Ks                    []uint
	maxK                   uint
	avgMaxPKs              map[uint]float32
	pKsSet, nKsSet         map[uint][]float32
	avgPKs, avgNKs         map[uint]float32
	startTime, lastEndTime time.Time
}

// NewResultsReporter returns a new ResultsReporter.
func NewResultsReporter(Y sticker.LabelVectors, Ks []uint) *ResultsReporter {
	maxK := uint(0)
	for _, K := range Ks {
		if maxK < K {
			maxK = K
		}
	}
	pKsSet, nKsSet := make(map[uint][]float32), make(map[uint][]float32)
	for _, K := range Ks {
		pKsSet[K], nKsSet[K] = make([]float32, 0, len(Y)), make([]float32, 0, len(Y))
	}
	return &ResultsReporter{
		_Y:     Y,
		_Ks:    Ks,
		maxK:   maxK,
		pKsSet: pKsSet,
		nKsSet: nKsSet,
		avgPKs: make(map[uint]float32),
		avgNKs: make(map[uint]float32),
	}
}

// AvgMaxPrecisionKs returns the average max Precision@Ks.
func (reporter *ResultsReporter) AvgMaxPrecisionKs() map[uint]float32 {
	if reporter.avgMaxPKs == nil {
		reporter.avgMaxPKs = make(map[uint]float32)
		for _, K := range reporter._Ks {
			maxPKs := sticker.ReportMaxPrecision(reporter._Y, K)
			sumPK := float32(0.0)
			for _, maxPKi := range maxPKs {
				sumPK += maxPKi
			}
			reporter.avgMaxPKs[K] = sumPK / float32(len(reporter._Y))
		}
	}
	return reporter.avgMaxPKs
}

// InferenceTimes returns the total and average inference time between the ResetTimer time and the Report time per entry.
func (reporter *ResultsReporter) InferenceTimes() (time.Duration, time.Duration) {
	inferenceTime := reporter.lastEndTime.Sub(reporter.startTime)
	return inferenceTime, time.Duration(inferenceTime.Nanoseconds() / int64(reporter.Nprocesseds())).Round(time.Microsecond)
}

// Ks returns the Ks.
func (reporter *ResultsReporter) Ks() []uint {
	return reporter._Ks
}

// MaxK returns the maximum number of Ks.
func (reporter *ResultsReporter) MaxK() uint {
	return reporter.maxK
}

// Nprocesseds returns the number of the processed entries.
func (reporter *ResultsReporter) Nprocesseds() int {
	for _, pKs := range reporter.pKsSet {
		return len(pKs)
	}
	return 0
}

// Report calculates Precision@Ks and nDCG@Ks with the predicted label vectors Yhat, append them, and returns the average ones.
// If w is not null, this function writes each result.
func (reporter *ResultsReporter) Report(Yhat sticker.LabelVectors, w io.Writer) (avgPKs, avgNKs map[uint]float32) {
	reporter.lastEndTime = time.Now()
	startidx := reporter.Nprocesseds()
	n := len(Yhat)
	for _, K := range reporter._Ks {
		reporter.pKsSet[K] = append(reporter.pKsSet[K], sticker.ReportPrecision(reporter._Y[startidx:n], K, Yhat[startidx:])...)
		sumPK := float32(0.0)
		for _, precisionKi := range reporter.pKsSet[K] {
			sumPK += precisionKi
		}
		avgPK := sumPK / float32(n)
		reporter.nKsSet[K] = append(reporter.nKsSet[K], sticker.ReportNDCG(reporter._Y[startidx:n], K, Yhat[startidx:])...)
		sumNK := float32(0.0)
		for _, nDCGKi := range reporter.nKsSet[K] {
			sumNK += nDCGKi
		}
		avgNK := sumNK / float32(n)
		reporter.avgPKs[K], reporter.avgNKs[K] = avgPK, avgNK
	}
	inferenceTime, inferenceTimePerEntry := reporter.InferenceTimes()
	if w != nil {
		fmt.Fprintf(w, "finished inference on %d/%d entries (%-5.4g%%) in %s (about %s/entry)\n", n, len(reporter._Y), float32(n)/float32(len(reporter._Y))*100.0, inferenceTime, inferenceTimePerEntry)
	}
	if w != nil {
		for _, K := range reporter._Ks {
			fmt.Fprintf(w, "Precision@%d=%-5.4g%%/%-5.4g%%, nDCG@%d=%-5.4g%%\n", K, reporter.avgPKs[K]*100, reporter.AvgMaxPrecisionKs()[K]*100, K, reporter.avgNKs[K]*100)
		}
	}
	return reporter.avgPKs, reporter.avgNKs
}

// ResetTimer resets the start time.
func (reporter *ResultsReporter) ResetTimer() {
	reporter.startTime = time.Now()
}
