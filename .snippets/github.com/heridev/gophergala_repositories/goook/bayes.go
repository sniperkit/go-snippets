package main

import "sort"

const (
	Epsilon = 0.000001
)

type NaiveBayes struct {
	counts map[int]uint
	models map[int]map[string]uint
}

func NewNaiveBayes() *NaiveBayes {
	return &NaiveBayes{
		counts: make(map[int]uint),
		models: make(map[int]map[string]uint),
	}
}

func (nb *NaiveBayes) Learn(class int, words []string) {
	if len(words) == 0 {
		return
	}

	classModel, exists := nb.models[class]
	if !exists {
		classModel = make(map[string]uint)
		nb.models[class] = classModel
	}

	for _, word := range words {
		classModel[word]++
		nb.counts[class]++
	}
}

func (nb *NaiveBayes) Classify(words []string) Probabilities {
	probs := make(map[int]float64)
	var norm float64
	for id, model := range nb.models {
		var p float64
		for _, word := range words {
			p += float64(model[word])
		}

		p /= float64(nb.counts[id])

		if p < Epsilon {
			continue
		}

		probs[id] = p
		norm += p
	}

	var ps Probabilities
	for id, p := range probs {
		ps = append(ps, Probability{
			Class: id,
			P:     p / norm,
		})
	}

	sort.Sort(ps)

	return ps
}

type Probability struct {
	Class int
	P     float64
}

type Probabilities []Probability

func (ps Probabilities) Len() int           { return len(ps) }
func (ps Probabilities) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps Probabilities) Less(i, j int) bool { return ps[j].P < ps[i].P }
