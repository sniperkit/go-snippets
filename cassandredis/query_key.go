package cassandredis

import (
	"bytes"
	"strings"
)

type queryKey struct {
	namespace string
	pk        []string
}

func (q queryKey) metadataKey() string {
	var buf bytes.Buffer

	buf.WriteString(q.namespace)
	if len(q.pk) > 0 {
		buf.WriteRune('_')
		for i, v := range q.pk {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteString(v)
		}
	}

	return buf.String()
}

func newQueryKeyFromString(s string) queryKey {
	tokens := strings.Split(s, ":")
	if len(tokens) < 1 {
		return queryKey{namespace: tokens[0]}
	}

	return queryKey{
		namespace: tokens[0],
		pk:        tokens[1:],
	}
}
