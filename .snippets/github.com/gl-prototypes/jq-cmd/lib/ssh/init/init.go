package init

import (
	"github.com/gl-prototypes/jq-cmd/lib/ssh"
	"github.com/gliderlabs/com"
)

func init() {
	com.DefaultRegistry.Register(&com.Object{Value: &ssh.Component{}})
}
