// Copyright 2014-2015 The DevMine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anlzr

import "github.com/DevMine/srcanlzr/src"

type LocPerLang struct{}

func (lpl LocPerLang) Analyze(p *src.Project, r *Result) error {

	// if the project has only one programming language,
	// which is mostly the case
	if len(p.Langs) == 1 {
		r.ProgLangs = append(r.ProgLangs,
			Language{Language: *p.Langs[0], Lines: p.LoC})
		return nil
	}

	m := make(map[string]Language)

	for _, pkg := range p.Packages {
		for _, srf := range pkg.SrcFiles {
			var lang Language
			var ok bool
			if lang, ok = m[srf.Lang.Lang]; !ok {
				lang = Language{Language: *srf.Lang}
			}

			lang.Lines += srf.LoC
			m[srf.Lang.Lang] = lang
		}
	}

	for _, v := range m {
		r.ProgLangs = append(r.ProgLangs, v)
	}

	return nil
}
