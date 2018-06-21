package ospf

import (
	"errors"
	"net"
)

const (
	IFSMS_Down         = 1
	IFSMS_Loopback     = 2
	IFSMS_Waiting      = 3
	IFSMS_PointToPoint = 4
	IFSMS_DROther      = 5
	IFSMS_Backup       = 6
	IFSMS_DR           = 7
)

const (
	IFSME_InterfaceUp    = 0
	IFSME_WaitTimer      = 1
	IFSME_BackupSeen     = 2
	IFSME_NeighborChange = 3
	IFSME_LoopInd        = 4
	IFSME_UnloopInd      = 5
	IFSME_InterfaceDown  = 6
)

type Interface struct {
	Name                   string
	Type                   int
	State                  int
	AreaID                 net.IP
	RouterPriority         uint8
	HelloInterval          uint16
	RouterDeadInterval     uint32
	RxmtInterval           int
	InfTransDelay          int
	NetworkType            int
	Options                uint8
	Cost                   uint32
	AuType                 uint32
	AuthenticationKey      uint64
	Interface              *net.Interface
	IP                     net.IP
	NetworkMask            net.IPMask
	HelloTimer             uint32
	WaitTimer              uint32
	DesignatedRouter       net.IP
	BackupDesignatedRouter net.IP
	Neighbors              map[string]*Neighbor
	InterfaceState
}

func NewOSPFInterface(ifp *net.Interface) (*Interface, error) {
	newIfp := Interface{
		Name:               ifp.Name,
		RouterPriority:     1,
		HelloInterval:      10,
		RouterDeadInterval: 40,
		RxmtInterval:       10,
		InfTransDelay:      1,
		NetworkType:        1,
		Interface:          ifp,
		Options:            0x2,
		Neighbors:          make(map[string]*Neighbor, 1),
	}

	ip, mask, err := GetIPv4Address(ifp)
	if err != nil {
		panic(err)
	}

	newIfp.IP = *ip
	newIfp.NetworkMask = *mask

	return &newIfp, nil
}

func GetIPv4Address(ifp *net.Interface) (*net.IP, *net.IPMask, error) {
	var ip *net.IP
	var mask *net.IPMask
	addrs, _ := ifp.Addrs()
	for _, a := range addrs {
		r, n, err := net.ParseCIDR(a.String())
		if err != nil {
			continue
		}

		if r.To4() != nil {
			ip = &r
			mask = &n.Mask

			return ip, mask, nil
		}
	}

	return nil, nil, errors.New("Interface has no ipv4 addresss")
}

func (oi *Interface) GetAttachedNeighbors() []net.IP {
	if len(oi.Neighbors) == 0 {
		return nil
	}

	ns := make([]net.IP, 0, len(oi.Neighbors))
	for _, n := range oi.Neighbors {
		ns = append(ns, n.RouterID)
	}

	return ns
}
