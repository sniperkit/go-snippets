// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"errors"
	"os"
	"path/filepath"
)

func openOrCreateFile(fileName string) (*os.File, error) {
	// TODO maybe should set the flag and FileMode by user.
	file, err := os.OpenFile(fileName, os.O_APPEND, 0644)
	if err != nil {
		dirPath := filepath.Dir(fileName)
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return nil, err
		}
		file, err = os.Create(fileName)
		if err != nil {
			return nil, errors.New("can not create the file.")
		}
	}
	return file, nil
}
