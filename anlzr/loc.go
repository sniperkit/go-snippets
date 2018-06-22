// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anlzr

import "github.com/DevMine/srcanlzr/src"

const (
	maxInt64 = int64(^uint64(0) >> 1)
	minInt64 = -maxInt64 - 1
)

type LoC struct{}

func (lca LoC) Analyze(p *src.Project, r *Result) error {
	r.TotalLoC = p.LoC

	var totalFuncs int64
	var totalLoCFunc int64

	maxLoCFunc := minInt64
	minLoCFunc := maxInt64

	hist := make(map[int64]int64)

	for _, pkg := range p.Packages {
		for _, sf := range pkg.SrcFiles {
			for _, f := range sf.Funcs {
				totalFuncs++
				totalLoCFunc += f.LoC

				if f.LoC > maxLoCFunc {
					maxLoCFunc = f.LoC
				}

				if f.LoC < minLoCFunc {
					minLoCFunc = f.LoC
				}

				hist[f.LoC] += int64(1)
			}

			for _, cls := range sf.Classes {
				for _, m := range cls.Methods {
					totalFuncs++
					totalLoCFunc += m.LoC

					if m.LoC > maxLoCFunc {
						maxLoCFunc = m.LoC
					}

					if m.LoC < minLoCFunc {
						minLoCFunc = m.LoC
					}

					hist[m.LoC] += int64(1)

				}
			}

			for _, mod := range sf.Traits {
				for _, m := range mod.Methods {
					totalFuncs++
					totalLoCFunc += m.LoC

					if m.LoC > maxLoCFunc {
						maxLoCFunc = m.LoC
					}

					if m.LoC < minLoCFunc {
						minLoCFunc = m.LoC
					}

					hist[m.LoC] += int64(1)
				}
			}
		}
	}

	r.AverageFuncLen = float32(totalLoCFunc) / float32(totalFuncs)
	r.MaxFuncLen = maxLoCFunc
	r.MinFuncLen = minLoCFunc
	r.MedianFuncLen = lca.median(hist)

	return nil
}

// median computes the media of a given histogram.
func (lca LoC) median(h map[int64]int64) int64 {
	maxVal := minInt64
	maxKey := int64(-1)

	for k, v := range h {
		if v > maxVal {
			maxVal = v
			maxKey = k
		}
	}

	return maxKey
}
