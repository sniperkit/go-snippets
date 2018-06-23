package maputil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMapPluck(t *testing.T) {
	assert := require.New(t)

	assert.Empty(Pluck(nil, nil))
	assert.Empty(Pluck(nil, []string{`name`}))
	assert.Empty(Pluck(`test`, []string{`name`}))
	assert.Empty(Pluck([]string{`test`, `values`}, []string{`name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]map[string]string{
		map[string]string{
			`name`: `Alice`,
		},
		map[string]string{
			`name`: `Bob`,
		},
		map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Mallory`,
	}, Pluck([]map[string]string{
		map[string]string{
			`name`: `Alice`,
		},
		map[string]string{
			`NAME`: `Bob`,
		},
		map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]map[string]map[string]interface{}{
		map[string]map[string]interface{}{
			`info`: map[string]interface{}{
				`name`: `Alice`,
			},
		},
		map[string]map[string]interface{}{
			`info`: map[string]interface{}{
				`name`: `Bob`,
			},
		},
		map[string]map[string]interface{}{
			`info`: map[string]interface{}{
				`name`: `Mallory`,
			},
		},
	}, []string{`info`, `name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]map[interface{}]map[interface{}]interface{}{
		map[interface{}]map[interface{}]interface{}{
			`info`: map[interface{}]interface{}{
				`name`: `Alice`,
			},
		},
		map[interface{}]map[interface{}]interface{}{
			`info`: map[interface{}]interface{}{
				`name`: `Bob`,
			},
		},
		map[interface{}]map[interface{}]interface{}{
			`info`: map[interface{}]interface{}{
				`name`: `Mallory`,
			},
		},
	}, []string{`info`, `name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]interface{}{
		map[string]string{
			`name`: `Alice`,
		},
		map[string]string{
			`name`: `Bob`,
		},
		map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck([]interface{}{
		&map[string]string{
			`name`: `Alice`,
		},
		&map[string]string{
			`name`: `Bob`,
		},
		&map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))

	assert.Equal([]interface{}{
		`Alice`,
		`Bob`,
		`Mallory`,
	}, Pluck(&[]interface{}{
		&map[string]string{
			`name`: `Alice`,
		},
		&map[string]string{
			`name`: `Bob`,
		},
		&map[string]string{
			`name`: `Mallory`,
		},
	}, []string{`name`}))
}
