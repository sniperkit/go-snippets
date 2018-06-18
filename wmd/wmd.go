package wmd

import (
	//"fmt"
	"math"
	"strings"

	"github.com/yizha/go/tp"
	"github.com/yizha/go/w2v"
)

type nbdoc struct {
	nbow []float64
	wvec []w2v.Vector
}

func toNbDoc(words []string, m *w2v.Model) *nbdoc {
	wvs := make([]w2v.Vector, 0, len(words))
	wmap := make(map[string][]int) // word --> [id, cnt]
	wcnt := 0
	for _, w := range words {
		w = strings.ToLower(w)
		wmeta, ok := wmap[w]
		if ok {
			wmeta[1] = wmeta[1] + 1
		} else {
			wv := m.GetVectorByWord(w)
			if wv == nil {
				continue
			}
			wcnt += 1
			id := len(wvs)
			wvs = append(wvs, wv)
			wmap[w] = []int{id, 1}
		}
		//fmt.Printf("word: %v, wmeta: %v,wmap: %v\n", w, wmeta, wmap)
	}

	//fmt.Printf("wmap: %v\n", wmap)
	//fmt.Println()

	wvsLen := len(wvs)
	if wvsLen > 0 {
		nbow := make([]float64, wvsLen)
		for _, meta := range wmap {
			nbow[meta[0]] = float64(meta[1]) / float64(wcnt)
		}
		return &nbdoc{
			nbow: nbow,
			wvec: wvs,
		}
	} else {
		return nil
	}
}

func calculateDistanceMatrix(d1, d2 *nbdoc, vsize int) [][]float64 {
	wv1, wv2 := d1.wvec, d2.wvec
	wv1Len, wv2Len := len(wv1), len(wv2)
	dm := make([][]float64, wv1Len)
	for i := 0; i < wv1Len; i++ {
		v1 := wv1[i]
		dm[i] = make([]float64, wv2Len)
		for j := 0; j < wv2Len; j++ {
			v2 := wv2[j]
			sum := float64(0)
			for k := 0; k < vsize; k++ {
				diff := v1[k] - v2[k]
				sum += diff * diff
			}
			dm[i][j] = math.Sqrt(sum)
		}
	}
	return dm
}

// returns the word-move-distance between the given two words slice
// if one of the words slice doesn't have any word in the model then
// this function returns math.Inf(1).
func Wmd(d1, d2 []string, m *w2v.Model) (float64, error) {
	nbd1 := toNbDoc(d1, m)
	nbd2 := toNbDoc(d2, m)

	if nbd1 == nil || nbd2 == nil {
		return math.Inf(1), nil
	}

	dm := calculateDistanceMatrix(nbd1, nbd2, m.FeatureSize)

	//fmt.Printf("supply: %v\n", nbd1.nbow)
	//fmt.Printf("demand: %v\n", nbd2.nbow)
	//fmt.Printf("costs:  %v\n", dm)
	p, err := tp.CreateProblem(nbd1.nbow, nbd2.nbow, dm)
	if err != nil {
		return -1, err
	}
	err = p.Solve()
	if err != nil {
		return -1, err
	}
	return p.GetCost(), nil
}
