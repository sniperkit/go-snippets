package lan

import (
	"fmt"
	"strings"
)

//
func SrvName(srvId string) string {
	ss := strings.Split(srvId, "@")
	return ss[0]
}

func SrvId(srvName, addr string) string {
	return fmt.Sprintf("%s@%s", srvName, addr) // to do缩短长度
}

func FormatTcpAddr(addr string) string {
	if !strings.HasPrefix(addr, "tcp://") {
		addr = "tcp://" + addr
	}
	return addr
}

//
type LanCfg struct {
	Name string // server name
	Addr string // server addr -- tcp://$ip:$port
}

func NewLanCfg(name, addr string) *LanCfg {
	addr = FormatTcpAddr(addr)
	return &LanCfg{Name: name, Addr: addr}
}

func (this LanCfg) ServerId() string {
	return SrvId(this.Name, this.Addr)
}

func (this LanCfg) String() string {
	return fmt.Sprintf("%v=>%v", this.Name, this.Addr)
}
