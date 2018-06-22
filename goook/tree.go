package main

import (
	"fmt"
	"sort"
)

type TreeBayes []*NaiveBayes

func NewTreeBayes() TreeBayes {
	return make([]*NaiveBayes, 101)
}

func (tb TreeBayes) Learn(class int, words []string) {
	curr, prev := class/10+1, class
	for check := class; check > 0 && check%10 == 0; check /= 10 {
		if curr <= 10 {
			curr, prev = 0, curr
		} else {
			curr, prev = curr/10+1, curr
		}
	}

	for prev > 0 {
		cls := tb[curr]
		if cls == nil {
			cls = NewNaiveBayes()
			tb[curr] = cls
		}

		cls.Learn(prev, words)
		if curr <= 10 {
			curr, prev = 0, curr
		} else {
			curr, prev = curr/10+1, curr
		}
	}
}

func (tb TreeBayes) Classify(words []string, threshold float64) (int, float64) {
	p := float64(1.0)
	pos := 0
	for pos < len(tb) {
		cls := tb[pos]
		if cls == nil {
			break
		}

		ps := cls.Classify(words)
		if ps[0].P < threshold {
			break
		}

		p *= ps[0].P
		pos = ps[0].Class
	}

	return pos - 1, p
}

func (tb TreeBayes) ClassifyTree(words []string) *ProbabilityTree {
	pt := &ProbabilityTree{
		Class: -1,
		P:     1.0,
	}

	for _, pR := range tb[0].Classify(words) {
		clsH := tb[pR.Class]

		ptH := &ProbabilityTree{
			Class: (pR.Class - 1) * 100,
			P:     pR.P,
		}
		pt.Children = append(pt.Children, ptH)

		if clsH == nil {
			continue
		}

		for _, pH := range clsH.Classify(words) {
			clsT := tb[pH.Class]

			ptT := &ProbabilityTree{
				Class: (pH.Class - 1) * 10,
				P:     pH.P,
			}
			ptH.Children = append(ptH.Children, ptT)

			if clsT == nil {
				continue
			}

			for _, pT := range clsT.Classify(words) {
				ptT.Children = append(ptT.Children, &ProbabilityTree{
					Class: pT.Class,
					P:     pT.P,
				})
			}
		}
	}

	pt.Sort()

	return pt
}

type ProbabilityTree struct {
	Class    int
	P        float64
	Children []*ProbabilityTree
}

func (pt *ProbabilityTree) Sort() {
	sort.Sort(ProbabilityTreeSorter(pt.Children))
	for _, child := range pt.Children {
		child.Sort()
	}
}

func (pt *ProbabilityTree) Display(p float64, prefix, indent string, ddc *DDC) {
	if pt.Class != -1 {
		fmt.Printf(
			"%3.f%%  %s%03d - %s\n",
			p*pt.P*100,
			prefix,
			pt.Class,
			ddc.Classes[pt.Class].Name,
		)

		prefix = prefix + indent
	}

	for _, child := range pt.Children {
		child.Display(p*pt.P, prefix, indent, ddc)
	}
}

type ProbabilityTreeSorter []*ProbabilityTree

func (pts ProbabilityTreeSorter) Len() int           { return len(pts) }
func (pts ProbabilityTreeSorter) Swap(i, j int)      { pts[i], pts[j] = pts[j], pts[i] }
func (pts ProbabilityTreeSorter) Less(i, j int) bool { return pts[j].P < pts[i].P }
