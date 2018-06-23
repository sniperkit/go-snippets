package adorn

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyTestCases(t *testing.T) {
	filenames, err := filepath.Glob("./testcases/*.json")
	require.NoError(t, err)
	for _, filename := range filenames {
		base := filepath.Base(filename)
		t.Run(base, func(t *testing.T) {
			cfg, err := LoadConfigFromFile(filename)
			require.NoError(t, err)
			src, err := GenerateString(cfg)
			require.NoError(t, err)
			noext := strings.TrimSuffix(filename, filepath.Ext(filename))
			expect, err := ioutil.ReadFile(noext + ".go")
			require.NoError(t, err)
			assert.Equal(t, string(expect), src)
		})
	}
}
