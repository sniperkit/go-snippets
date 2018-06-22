package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var reYamlInclude = regexp.MustCompile(`(.*)?!include\s+(.*)`)

func loadYaml(path, baseDir string) ([]byte, error) {
	yamlData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err

	}
	yamlString := string(yamlData)

	matches := reYamlInclude.FindAllStringSubmatch(yamlString, -1)
	if len(matches) == 0 {
		return yamlData, nil
	}

	for _, match := range matches {
		var (
			rawLine   = match[0]
			rawIndent = match[1]
			pattern   = match[2]
		)

		if !strings.HasPrefix(pattern, "/") {
			pattern = filepath.Join(baseDir, pattern)
		}

		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf(
				"can't glob by pattern '%s': %s",
				pattern, err,
			)
		}

		var including string
		for _, file := range files {
			subYamlData, err := loadYaml(file, baseDir)
			if err != nil {
				return nil, fmt.Errorf(
					"can't include by pattern '%s': %s", pattern, err,
				)
			}

			including = string(subYamlData) + "\n"
		}

		if len(rawIndent) > 0 {
			level := 0
			nextLevel := false
			if string(rawIndent[0]) == " " {
				trimmed := strings.TrimLeft(rawIndent, " ")
				spaces := len(rawIndent) - len(trimmed)
				level = spaces / 4

				if trimmed != "" {
					level = level + 1
					nextLevel = true
				}
			}

			lines := strings.Split(including, "\n")
			for i, line := range lines {
				if !nextLevel && i == 0 {
					continue
				}

				lines[i] = strings.Repeat(" ", level*4) + line
			}

			including = strings.Join(lines, "\n")
			if nextLevel {
				including = "\n" + including
			}

			including = rawIndent + including
		}

		yamlString = strings.Replace(yamlString, rawLine, including, -1)
	}

	return []byte(yamlString), nil
}
