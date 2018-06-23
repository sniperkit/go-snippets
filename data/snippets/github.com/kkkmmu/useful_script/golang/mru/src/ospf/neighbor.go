package ospf

import (
	"net"
)

const (
	NFSMS_DOWN     = 1
	NFSMS_Attempt  = 2
	NFSMS_Init     = 3
	NFSMS_2Way     = 4
	NFSMS_Exstart  = 5
	NFSMS_Exchange = 6
	NFSMS_Loading  = 7
	NFSMS_Full     = 8
)

const (
	NFSME_HelloReceived     = 1
	NFSME_Start             = 2
	NFSME_2WayReceived      = 3
	NFSME_NegotiationDone   = 4
	NFSME_ExchangeDone      = 5
	NFSME_BadLSReq          = 6
	NFSME_LoadingDone       = 7
	NFSME_AdjOK             = 8
	NFSME_SeqNumberMismatch = 9
	NFSME_1WayReceived      = 10
	NFSME_KillNbr           = 11
	NFSME_InactivityTimer   = 12
	NFSME_LLDown            = 13
)

type Neighbor struct {
	RouterID                net.IP
	AreaID                  net.IP
	AreaType                int
	IP                      net.IP
	Options                 uint8
	State                   int
	InactivityTimer         uint32
	Master                  bool
	DDSequenceNumber        uint32
	LastReceivedDBD         DBD
	Priority                uint8
	DesignatedRouter        net.IP
	BackupDesignatedRouter  net.IP
	LinkStateRetransmitList []LSA
	DatabaseSummaryList     []LSA
	LinkStateRequestList    []LSA
	NeighborState
}
