package goconfigure

import (
	"github.com/coreos/go-etcd/etcd"
	"strings"
)

func nestedEtcdToMap(node *etcd.Node) (v map[string]interface{}) {
	v = make(map[string]interface{})

	for _, n := range node.Nodes {
		keys := strings.Split(n.Key, "/")
		key := keys[len(keys)-1]

		if n.Dir {
			v[key] = nestedEtcdToMap(n)
		} else {
			v[key] = n.Value
		}
	}

	return
}
