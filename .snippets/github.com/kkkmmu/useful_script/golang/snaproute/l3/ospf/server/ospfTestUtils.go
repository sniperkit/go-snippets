//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"fmt"
	"infra/sysd/sysdCommonDefs"
	"l3/ospf/config"
	"log/syslog"
	"net"
	"ospfd"
	"ribdInt"
	"sync"
	"time"
	"utils/commonDefs"
	"utils/logging"
)

const (
	SUCCESS = 0
	FAIL    = 1
)

/* Pkt types */
var hello []byte
var lsaupd []byte
var lsareq []byte
var lsaack []byte
var header []byte
var lsa_network []byte
var lsa_router []byte
var lsa_summary []byte
var lsa_asExt []byte
var lsa_update []byte
var lsa_fake []byte

var ospfHdrMd OspfHdrMetadata
var ipHdrMd IpHdrMetadata
var ethHdrMd EthHdrMetadata
var key IntfConfKey
var srcMAC net.HardwareAddr
var ipIntfProp IPIntfProperty
var ifType int
var nbrConf OspfNeighborEntry
var nbrKey NeighborConfKey
var intConf IntfConf
var dstMAC net.HardwareAddr
var dstIP net.IP
var netmask uint32
var lsid uint32
var ospf *OSPFServer
var eventMsg DbEventMsg
var conf config.IfMetricConf

/* Infra Structs */
var portProperty PortProperty
var vlanProperty VlanProperty
var ipv4Msg IPv4IntfNotifyMsg
var ipProperty IpProperty
var ipIntfProperty IPIntfProperty

/* Global conf */
var gConf config.GlobalConf

/* Area conf */
var areaConf config.AreaConf
var areaConfKey AreaConfKey
var ospfArea *ospfd.OspfAreaEntry

/* Intf FSM */
var msg NbrStateChangeMsg
var msgNbrFull NbrFullStateMsg
var intf IntfConf
var hellodata OSPFHelloData
var ifConf config.InterfaceConf
var ospfIf *ospfd.OspfIfEntry

/* Nbr FSM */
var ospfNbrEntry OspfNeighborEntry
var nbrConfMsg ospfNeighborConfMsg
var nbrDbPkt ospfDatabaseDescriptionData
var nbrIntfMsg IntfToNeighMsg
var nbrDbdMsg ospfNeighborDBDMsg
var nbrLsaReqMsg ospfNeighborLSAreqMsg
var db_list []*ospfNeighborDBSummary
var req ospfLSAReq
var req1 ospfLSAReq

var nbr_req *ospfNeighborReq
var nbr_req_list []*ospfNeighborReq
var lsa_header1 ospfLSAHeader
var lsa_header2 ospfLSAHeader
var ack_msg *ospfNeighborLSAAckMsg

/* Lsdb and flooding */
var areaId uint32
var lsdbKey LsdbKey
var summaryKey LsaKey
var summaryLsa SummaryLsa

var routerKey LsaKey
var routerLsa RouterLsa
var link1 LinkDetail
var lin []LinkDetail
var networkKey LsaKey
var networkLsa NetworkLsa
var val1 LsdbSliceEnt
var val2 LsdbSliceEnt
var val3 LsdbSliceEnt
var val4 LsdbSliceEnt
var val5 LsdbSliceEnt

var lsa_reqs []ospfLSAReq
var val LsdbSliceEnt
var lsdbMsg DbLsdbMsg
var ackTxMsg ospfNeighborAckTxMsg
var floodMsg ospfFloodMsg

/* Routing table */
var rKey RoutingTblEntryKey
var rKeyAbr RoutingTblEntryKey
var rKeyAsbr RoutingTblEntryKey
var rKeyAsabr RoutingTblEntryKey
var rEntry GlobalRoutingTblEntry
var entry RoutingTblEntry
var nextHop NextHop
var nhmap map[NextHop]bool
var areaRoutingTable AreaRoutingTbl
var areaidkey AreaIdKey
var vKeyR VertexKey
var vKeyN VertexKey
var vKeyT VertexKey

var vertexR Vertex
var vertexN Vertex
var vertexT Vertex
var treeVertex TreeVertex
var sVertex1 StubVertex
var sVertex2 StubVertex
var sVertex3 StubVertex

var route ribdInt.Routes

func OSPFNewLogger(name string, tag string, listenToConfig bool) (*logging.Writer, error) {
	var err error
	srLogger := new(logging.Writer)
	srLogger.MyComponentName = name

	srLogger.SysLogger, err = syslog.New(syslog.LOG_DEBUG|syslog.LOG_DAEMON, tag)
	if err != nil {
		fmt.Println("Failed to initialize syslog - ", err)
		return srLogger, err
	}

	srLogger.MyLogLevel = sysdCommonDefs.DEBUG
	return srLogger, err
}

func initAttr() {
	ospf.initOspfGlobalConfDefault()
	ospf.InitDBChannels()

	initPacketData()

	ifType = int(config.Broadcast)
	srcMAC = net.HardwareAddr{0x01, 0x00, 0x50, 0x00, 0x00, 0x07}
	dstMAC = net.HardwareAddr{0x24, 00, 0x50, 0x00, 0x00, 0x05}
	dstIP = net.IP{10, 1, 1, 2}

	areaConfKey = AreaConfKey{
		AreaId: config.AreaId("10.0.0.0"),
	}
	areaConf = config.AreaConf{
		AreaId:          config.AreaId("10.0.0.0"),
		AuthType:        config.AuthType(1),
		ImportAsExtern:  config.ImportAsExtern(1),
		AreaSummary:     config.AreaSummary(2),
		StubDefaultCost: int32(20),
	}

	ospfArea = &ospfd.OspfAreaEntry{
		AuthType:       int32(areaConf.AuthType),
		ImportAsExtern: int32(areaConf.ImportAsExtern),
		AreaSummary:    int32(areaConf.AreaSummary),
	}

	gConf.RouterId = "20.0.1.1"
	gConf.AdminStat = config.Disabled
	gConf.ASBdrRtrStatus = true
	gConf.TOSSupport = false
	gConf.RestartSupport = config.None
	gConf.RestartInterval = 40
	gConf.ReferenceBandwidth = 100

	ifConf = config.InterfaceConf{
		IfIpAddress:       config.IpAddress("10.1.1.2"),
		AddressLessIf:     config.InterfaceIndexOrZero(0),
		IfAreaId:          config.AreaId("10.0.0.0"),
		IfType:            config.IfType(0),
		IfAdminStat:       config.Status(0),
		IfRtrPriority:     config.DesignatedRouterPriority(1),
		IfTransitDelay:    config.UpToMaxAge(10),
		IfRetransInterval: config.UpToMaxAge(40),
		IfHelloInterval:   config.HelloRange(10),
		IfRtrDeadInterval: config.PositiveInteger(10),
		IfPollInterval:    config.PositiveInteger(10),
	}

	ospfIf = &ospfd.OspfIfEntry{
		IfIpAddress:       string(ifConf.IfIpAddress),
		AddressLessIf:     int32(ifConf.AddressLessIf),
		IfAreaId:          string(ifConf.IfAreaId),
		IfRtrPriority:     int32(ifConf.IfRtrPriority),
		IfTransitDelay:    int32(ifConf.IfTransitDelay),
		IfRetransInterval: int32(ifConf.IfRetransInterval),
		IfHelloInterval:   int32(ifConf.IfHelloInterval),
		IfRtrDeadInterval: int32(ifConf.IfRtrDeadInterval),
		IfPollInterval:    int32(ifConf.IfPollInterval),
	}

	hellodata = OSPFHelloData{
		netmask:             []byte{10, 0, 0, 0},
		helloInterval:       uint16(10),
		options:             uint8(0),
		rtrPrio:             uint8(1),
		rtrDeadInterval:     uint32(40),
		designatedRtr:       []byte{10, 1, 1, 2},
		backupDesignatedRtr: []byte{10, 1, 1, 7},
		neighbor:            []byte{10, 1, 1, 2},
	}

	ospfHdrMd = OspfHdrMetadata{
		pktType:  HelloType,
		pktlen:   OSPF_HELLO_MIN_SIZE,
		backbone: false,
		routerId: []byte{10, 1, 1, 10},
		areaId:   101010,
	}

	ipHdrMd = IpHdrMetadata{
		srcIP:     []byte{10, 1, 1, 2},
		dstIP:     []byte{10, 1, 1, 2},
		dstIPType: Normal,
	}
	ethHdrMd = EthHdrMetadata{
		srcMAC: srcMAC,
	}

	key = IntfConfKey{
		IPAddr:  config.IpAddress(net.IP{10, 1, 1, 2}),
		IntfIdx: config.InterfaceIndexOrZero(2),
	}

	ipIntfProp = IPIntfProperty{
		IfName:  "fpPort1",
		IpAddr:  net.IP{10, 1, 1, 1},
		MacAddr: srcMAC,
		NetMask: []byte{10, 1, 0, 0},
		Mtu:     8124,
		Cost:    10,
	}
	nbrConf = OspfNeighborEntry{
		OspfNbrRtrId:          20,
		OspfNbrIPAddr:         net.IP{10, 1, 1, 2},
		OspfRtrPrio:           17,
		intfConfKey:           key,
		OspfNbrOptions:        0,
		OspfNbrState:          config.NbrInit,
		isStateUpdate:         false,
		isDRBDR:               true,
		ospfNbrSeqNum:         1223,
		req_list_mutex:        &sync.Mutex{},
		db_summary_list_mutex: &sync.Mutex{},
		retx_list_mutex:       &sync.Mutex{},
	}

	nbrKey = NeighborConfKey{
		IPAddr:  config.IpAddress(net.IP{10, 1, 1, 2}),
		IntfIdx: config.InterfaceIndexOrZero(2),
	}

	msg = NbrStateChangeMsg{
		nbrKey: nbrKey,
	}

	msgNbrFull = NbrFullStateMsg{
		FullState: true,
		NbrRtrId:  10,
		nbrKey:    nbrKey,
	}

	intf = IntfConf{
		IfAreaId:          []byte{0, 0, 0, 1},
		IfType:            config.Broadcast,
		IfAdminStat:       config.Enabled,
		IfRtrPriority:     uint8(2),
		IfTransitDelay:    config.UpToMaxAge(3600),
		IfRetransInterval: config.UpToMaxAge(3600),
		IfHelloInterval:   uint16(10),
		IfRtrDeadInterval: uint32(40),
		IfPollInterval:    config.PositiveInteger(50),
		IfDemand:          false,

		/* IntefaceState: Start */
		IfDRIp:     []byte{10, 1, 1, 2},
		IfBDRIp:    []byte{10, 1, 1, 10},
		IfFSMState: config.Down,
		IfDRtrId:   uint32(10),
		IfBDRtrId:  uint32(20),
		/* IntefaceState: End */
		IfName:    "fpPort1",
		IfIpAddr:  net.IP{10, 1, 1, 2},
		IfMacAddr: net.HardwareAddr{0x01, 0x00, 0x50, 0x00, 0x00, 0x07},
		IfNetmask: []byte{10, 0, 0, 0},
		IfMtu:     8124,
	}

	ospfNbrEntry = OspfNeighborEntry{
		OspfNbrRtrId:           20,
		OspfNbrIPAddr:          net.IP{10, 1, 1, 2},
		OspfRtrPrio:            2,
		intfConfKey:            key,
		OspfNbrOptions:         0,
		OspfNbrState:           config.NbrInit,
		isStateUpdate:          true,
		OspfNbrInactivityTimer: time.Now(),
		OspfNbrDeadTimer:       40,
		ospfNbrSeqNum:          2002,
		isSeqNumUpdate:         true,
		isMaster:               true,
		isMasterUpdate:         true,
		ospfNbrLsaIndex:        0,
	}

	nbrConfMsg = ospfNeighborConfMsg{
		ospfNbrConfKey: nbrKey,
		ospfNbrEntry:   ospfNbrEntry,
		nbrMsgType:     NBRADD,
	}

	nbrDbPkt = ospfDatabaseDescriptionData{
		options:            0,
		interface_mtu:      1500,
		dd_sequence_number: 2000,
		ibit:               true,
		mbit:               true,
		msbit:              false,
	}

	nbrIntfMsg = IntfToNeighMsg{
		IntfConfKey:  key,
		RouterId:     12,
		RtrPrio:      1,
		NeighborIP:   net.IP{10, 1, 1, 2},
		nbrDeadTimer: 20,
		TwoWayStatus: true,
		nbrDR:        []byte{10, 0, 0, 1},
		nbrBDR:       []byte{10, 0, 0, 2},
		nbrMAC:       net.HardwareAddr{0x02, 0x00, 0x50, 0x00, 0x00, 0x08},
	}

	nbrDbdMsg = ospfNeighborDBDMsg{
		ospfNbrConfKey: nbrKey,
		nbrFull:        false,
		ospfNbrDBDData: nbrDbPkt,
	}

	req = ospfLSAReq{
		ls_type:       uint32(1),
		link_state_id: uint32(2001),
		adv_router_id: uint32(4001),
	}

	req1 = ospfLSAReq{
		ls_type:       uint32(2),
		link_state_id: uint32(2001),
		adv_router_id: uint32(4001),
	}

	lsa_reqs = []ospfLSAReq{}
	lsa_reqs = append(lsa_reqs, req)
	lsa_reqs = append(lsa_reqs, req1)
	nbrLsaReqMsg = ospfNeighborLSAreqMsg{
		lsa_slice: []ospfLSAReq{
			req,
		},
		nbrKey: nbrKey,
	}

	ack_msg = newospfNeighborLSAAckMsg()

	ackTxMsg.lsa_headers_byte = lsaack
	ackTxMsg.nbrKey = nbrKey

	eventMsg = DbEventMsg{
		eventType: config.ADJACENCY,
	}
	eventMsg.eventInfo = "SeqNumberMismatch. Nbr should be master "

	initInfra()
	initLsdbData()
	initRoutingTable()
	//populateNbrLists()

}

func initInfra() {
	conf = config.IfMetricConf{
		IfMetricIpAddress:     config.IpAddress("10.1.1.10"),
		IfMetricAddressLessIf: config.InterfaceIndexOrZero(0),
		IfMetricTOS:           config.TosType(2),
		IfMetricValue:         config.Metric(20),
	}
	portProperty = PortProperty{
		Name:  "fpPort1",
		Mtu:   int32(1500),
		Speed: uint32(100),
	}
	vlanProperty = VlanProperty{
		Name:       "fpPort1",
		UntagPorts: []int32{12, 1, 2},
	}

	ipv4Msg = IPv4IntfNotifyMsg{
		IpAddr: string("10.1.1.2/24"),
		IfId:   uint16(0),
		IfType: uint8(commonDefs.IfTypePort),
	}

	ipProperty = IpProperty{
		IfId:   uint16(0),
		IfType: uint8(commonDefs.IfTypePort),
		IpAddr: "10.1.1.2",
	}
	ipIntfProperty = IPIntfProperty{
		IfName:  "fpPort1",
		IpAddr:  dstIP,
		MacAddr: dstMAC,
		NetMask: []byte{0x0a, 0x0, 0x0, 0x0},
		Mtu:     1500,
		Cost:    uint32(20),
	}

}
func initLsdbData() {
	areaid := convertAreaOrRouterIdUint32("10.0.0.0")
	netmask = convertAreaOrRouterIdUint32("255.0.0.0")
	lsid = convertAreaOrRouterIdUint32("10.1.1.1")
	lsdbKey = LsdbKey{
		AreaId: areaid,
	}

	summaryKey = LsaKey{
		LSType:    uint8(Summary3LSA),
		LSId:      lsid,
		AdvRouter: lsid,
	}

	lsamdata := LsaMetadata{
		LSAge:         uint16(1),
		Options:       uint8(0),
		LSSequenceNum: int(1800),
		LSChecksum:    uint16(12),
		LSLen:         uint16(28),
	}

	summaryLsa = SummaryLsa{
		LsaMd:   lsamdata,
		Netmask: netmask,
		Metric:  uint32(20),
	}
	val1.AreaId = lsdbKey.AreaId
	val1.LSType = RouterLSA
	val1.LSId = lsid
	val1.AdvRtr = lsid

	val2.AreaId = lsdbKey.AreaId
	val2.LSType = NetworkLSA
	val2.LSId = lsid
	val2.AdvRtr = lsid

	val3.AreaId = lsdbKey.AreaId
	val3.LSType = Summary3LSA
	val3.LSId = lsid
	val3.AdvRtr = lsid

	val4.AreaId = lsdbKey.AreaId
	val4.LSType = Summary4LSA
	val4.LSId = lsid
	val4.AdvRtr = lsid

	val5.AreaId = lsdbKey.AreaId
	val5.LSType = ASExternalLSA
	val5.LSId = lsid
	val5.AdvRtr = lsid

	routerKey = LsaKey{
		LSType:    uint8(RouterLSA),
		LSId:      lsid,
		AdvRouter: lsid,
	}

	link := make([]LinkDetail, 2)

	link1 = LinkDetail{
		LinkId:     uint32(1234), /* Link ID */
		LinkData:   uint32(1),    /* Link Data */
		LinkType:   TransitLink,  /* Link Type */
		NumOfTOS:   uint8(0),     /* # TOS Metrics */
		LinkMetric: uint16(20),   /* Metric */
	}

	link = append(link, link1)
	link1.LinkId = uint32(lsid)
	link1.LinkData = 1
	link1.LinkType = StubLink
	link = append(link, link1)

	link1.LinkId = uint32(lsid)
	link1.LinkData = 1
	link1.LinkType = P2PLink
	link = append(link, link1)

	routerLsa = RouterLsa{
		LsaMd:       lsamdata,
		BitV:        true,
		BitE:        true,
		BitB:        true,
		NumOfLinks:  3,
		LinkDetails: link,
	}

	networkKey = LsaKey{
		LSType:    uint8(NetworkLSA),
		LSId:      lsid,
		AdvRouter: lsid,
	}

	att := []uint32{12, 11, 10}

	networkLsa = NetworkLsa{
		LsaMd:       lsamdata,
		Netmask:     uint32(10), /* Network Mask */
		AttachedRtr: att,
	}

	val.AreaId = lsdbKey.AreaId
	val.LSType = summaryKey.LSType
	val.LSId = summaryKey.LSId
	val.AdvRtr = summaryKey.AdvRouter

	lsdbMsg = DbLsdbMsg{
		entry: val,
		op:    true,
	}
	vKeyR = VertexKey{
		Type:   RouterVertex,
		ID:     routerKey.LSId,
		AdvRtr: lsid,
	}

	vKeyN = VertexKey{
		Type:   SNetworkVertex,
		ID:     networkKey.LSId,
		AdvRtr: lsid,
	}

	vKeyT = VertexKey{
		Type:   SNetworkVertex,
		ID:     networkKey.LSId,
		AdvRtr: lsid,
	}

	vertexR = Vertex{
		NbrVertexKey:  []VertexKey{vKeyN},
		NbrVertexCost: []uint16{10},
		LsaKey:        routerKey,
		AreaId:        lsdbKey.AreaId,
		Visited:       false,
		LinkStateId:   lsid,
		NetMask:       uint32(0),
	}
	vertexR.LinkData = make(map[VertexKey]uint32)
	vertexR.LinkData[vKeyR] = link1.LinkData

	p := []VertexKey{vKeyR, vKeyN, vKeyT}

	treeVertex = TreeVertex{
		Paths:      []Path{p},
		Distance:   uint16(20),
		NumOfPaths: 3,
	}
	floodMsg = ospfFloodMsg{
		nbrKey:  nbrKey,
		intfKey: key,
		areaId:  lsdbKey.AreaId,
		lsType:  RouterLSA,
		linkid:  routerKey.LSId,
		lsaKey:  routerKey,
		lsOp:    LSAFLOOD,
		pkt:     lsa_router,
	}
	db_list = []*ospfNeighborDBSummary{}
	db_summary1 := newospfNeighborDBSummary()
	db_summary2 := newospfNeighborDBSummary()
	db_summary1.valid = true
	db_summary2.valid = true
	lsaHeader := getLsaHeaderFromLsa(routerLsa.LsaMd.LSAge, routerLsa.LsaMd.Options, routerKey.LSType,
		routerKey.LSId, routerKey.AdvRouter, uint32(routerLsa.LsaMd.LSSequenceNum),
		routerLsa.LsaMd.LSChecksum, routerLsa.LsaMd.LSLen)

	db_summary1.lsa_headers = lsaHeader
	db_summary2.lsa_headers = lsaHeader
	db_list = append(db_list, db_summary1)
	db_list = append(db_list, db_summary2)

}

func initRoutingTable() {
	rKey = RoutingTblEntryKey{
		DestType: Network,
		AddrMask: 0,
		DestId:   0,
	}
	rKeyAbr = RoutingTblEntryKey{
		DestType: AreaBdrRouter,
		AddrMask: 0,
		DestId:   0,
	}

	rKeyAsbr = RoutingTblEntryKey{
		DestType: ASBdrRouter,
		AddrMask: 0,
		DestId:   0,
	}

	rKeyAsabr = RoutingTblEntryKey{
		DestType: ASAreaBdrRouter,
		AddrMask: 0,
		DestId:   0,
	}

	nhmap = make(map[NextHop]bool)
	nextHop = NextHop{
		IfIPAddr:  uint32(2345),
		IfIdx:     uint32(10),
		NextHopIP: uint32(222),
		AdvRtr:    uint32(120),
	}

	nhmap[nextHop] = true

	entry = RoutingTblEntry{
		OptCapabilities: uint8(1),  // Optional Capabilities
		PathType:        IntraArea, // Path Type
		Cost:            uint16(20),
		Type2Cost:       uint16(10),
		LSOrigin:        summaryKey,
		NumOfPaths:      10,
		NextHops:        nhmap,
	}

	rEntry = GlobalRoutingTblEntry{
		AreaId:        lsdbKey.AreaId,
		RoutingTblEnt: entry,
	}

	ospf.GlobalRoutingTbl[rKey] = rEntry
	ospf.GlobalRoutingTbl[rKeyAbr] = rEntry
	ospf.GlobalRoutingTbl[rKeyAsbr] = rEntry
	ospf.GlobalRoutingTbl[rKeyAsabr] = rEntry

	ospf.OldGlobalRoutingTbl = ospf.GlobalRoutingTbl
	ospf.TempGlobalRoutingTbl = ospf.GlobalRoutingTbl
	areaRoutingTable.RoutingTblMap = make(map[RoutingTblEntryKey]RoutingTblEntry)
	areaRoutingTable.RoutingTblMap[rKey] = entry

	areaidkey = AreaIdKey{
		AreaId: lsdbKey.AreaId,
	}
	ospf.TempAreaRoutingTbl = make(map[AreaIdKey]AreaRoutingTbl)
	ospf.TempAreaRoutingTbl[areaidkey] = areaRoutingTable
	ospf.TempGlobalRoutingTbl[rKey] = rEntry
	ospf.OldGlobalRoutingTbl[rKey] = rEntry
	sVertex1 = StubVertex{
		NbrVertexKey:  vKeyR,
		NbrVertexCost: uint16(20),
		LinkData:      uint32(11),
		LsaKey:        routerKey,
		AreaId:        lsdbKey.AreaId,
		LinkStateId:   lsid,
	}

	sVertex2 = StubVertex{
		NbrVertexKey:  vKeyN,
		NbrVertexCost: uint16(20),
		LinkData:      uint32(11),
		LsaKey:        routerKey,
		AreaId:        lsdbKey.AreaId,
		LinkStateId:   lsid,
	}

	sVertex3 = StubVertex{
		NbrVertexKey:  vKeyT,
		NbrVertexCost: uint16(20),
		LinkData:      uint32(11),
		LsaKey:        routerKey,
		AreaId:        lsdbKey.AreaId,
		LinkStateId:   lsid,
	}

	route = ribdInt.Routes{
		Ipaddr: "10.1.1.2",
		Mask:   "10.0.0.0",
		Metric: 10,
	}

}

func initPacketData() {
	header = []byte{0x02, 0x05, 0x00, 0x40, 0x04, 0x04, 0x04, 0x04, 0x00, 0x00, 0x00, 0x14, 0x09, 0x83, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	hello = []byte{0x01, 0x00, 0x5e, 0x00, 0x00, 0x05, 0xca, 0x11, 0x09, 0xb3,
		0x00, 0x1c, 0x08, 0x00, 0x45, 0xc0, 0x00, 0x50, 0x8d, 0xed, 0x00, 0x00,
		0x01, 0x59, 0x3f, 0x5a, 0x0a, 0x4b, 0x00, 0xfe, 0xe0, 0x00, 0x00, 0x05,
		0x02, 0x01, 0x00, 0x30, 0x4b, 0x01, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00,
		0x3e, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0x00, 0x00, 0x0a, 0x12, 0x01, 0x00, 0x00, 0x00, 0x28,
		0x0a, 0x4b, 0x00, 0xfe, 0x0a, 0x4b, 0x00, 0x01, 0x4b, 0x01, 0x00, 0x01,
		0xff, 0xf6, 0x00, 0x03, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01}

	lsaupd = []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x22, 0x01, 0x05, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00, 0x00, 0x05, 0x0a, 0x40, 0x00, 0x30, 0x00,
		0x00, 0x00, 0x02, 0xc0, 0xa8, 0x14, 0x00, 0xff, 0xff, 0xff, 0x00, 0x03, 0x00, 0x00,
		0x0a, 0x0a, 0x00, 0x14, 0x00, 0xff, 0xff, 0xff, 0xfc, 0x03, 0x00, 0x00, 0x0a}

	lsareq = []byte{0x00, 0x00, 0x00, 0x01, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05}

	lsaack = []byte{0x00, 0x01, 0x22, 0x01, 0x05, 0x05,
		0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00, 0x00, 0x06, 0x78, 0xac, 0x00, 0x30}

	lsa_network = []byte{0x01, 0xbe, 0x22, 0x02, 0x0a, 0x00, 0x14, 0x02, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00,
		0x00, 0x01, 0xf6, 0xed, 0x00, 0x20, 0xff, 0xff, 0xff, 0xfc, 0x05, 0x05, 0x05, 0x05, 0x04, 0x04,
		0x04, 0x04}

	lsa_router = []byte{0x01, 0xbe, 0x22, 0x01, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00, 0x00, 0x04, 0x7c, 0xaa, 0x00, 0x30, 0x00, 0x00, 0x00, 0x02, 0xc0, 0xa8, 0x14, 0x00, 0xff, 0xff, 0xff, 0x00, 0x03, 0x00, 0x00, 0x0a, 0x0a, 0x00, 0x14, 0x02, 0x0a, 0x00, 0x14, 0x02, 0x02, 0x00, 0x00, 0x0a}

	lsa_summary = []byte{0x00, 0x0b, 0x22, 0x03, 0xc0, 0xa8, 0x0a, 0x00, 0x04, 0x04, 0x04, 0x04, 0x80, 0x00, 0x00, 0x01, 0x1e, 0x7d, 0x00, 0x1c, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x1e}

	lsa_asExt = []byte{0x00, 0xc5, 0x20, 0x05, 0xac, 0x10, 0x02, 0x00, 0x02, 0x02,
		0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x33, 0x56, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00,
		0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	lsa_update = []byte{0x00, 0x00, 0x00, 0x0b, 0x01, 0xbe,
		0x22, 0x01, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00, 0x00, 0x04, 0x7c, 0xaa,
		0x00, 0x30, 0x00, 0x00, 0x00, 0x02, 0xc0, 0xa8, 0x14, 0x00, 0xff, 0xff, 0xff, 0x00, 0x03, 0x00,
		0x00, 0x0a, 0x0a, 0x00, 0x14, 0x02, 0x0a, 0x00, 0x14, 0x02, 0x02, 0x00, 0x00, 0x0a, 0x00, 0x0a,
		0x22, 0x01, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x80, 0x00, 0x00, 0x06, 0x36, 0xb1,
		0x00, 0x24, 0x01, 0x00, 0x00, 0x01, 0x0a, 0x00, 0x14, 0x00, 0xff, 0xff, 0xff, 0xfc, 0x03, 0x00,
		0x00, 0x0a, 0x01, 0xbe, 0x22, 0x02, 0x0a, 0x00, 0x14, 0x02, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00,
		0x00, 0x01, 0xf6, 0xed, 0x00, 0x20, 0xff, 0xff, 0xff, 0xfc, 0x05, 0x05, 0x05, 0x05, 0x04, 0x04,
		0x04, 0x04, 0x00, 0x0b, 0x22, 0x03, 0xc0, 0xa8, 0x0a, 0x00, 0x04, 0x04, 0x04, 0x04, 0x80, 0x00,
		0x00, 0x01, 0x1e, 0x7d, 0x00, 0x1c, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x1e, 0x00, 0x0b,
		0x22, 0x03, 0x0a, 0x00, 0x0a, 0x00, 0x04, 0x04, 0x04, 0x04, 0x80, 0x00, 0x00, 0x01, 0xd6, 0x31,
		0x00, 0x1c, 0xff, 0xff, 0xff, 0xfc, 0x00, 0x00, 0x00, 0x14, 0x00, 0x0b, 0x22, 0x03, 0x0a, 0x00,
		0x00, 0x00, 0x04, 0x04, 0x04, 0x04, 0x80, 0x00, 0x00, 0x01, 0xe0, 0x3b, 0x00, 0x1c, 0xff, 0xff,
		0xff, 0xfc, 0x00, 0x00, 0x00, 0x0a, 0x00, 0x0b, 0x22, 0x04, 0x02, 0x02, 0x02, 0x02, 0x04, 0x04,
		0x04, 0x04, 0x80, 0x00, 0x00, 0x01, 0x6f, 0xa0, 0x00, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x14, 0x00, 0xc5, 0x20, 0x05, 0xac, 0x10, 0x03, 0x00, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00,
		0x00, 0x01, 0x28, 0x60, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00, 0x00, 0x64, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc5, 0x20, 0x05, 0xac, 0x10, 0x02, 0x00, 0x02, 0x02,
		0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x33, 0x56, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00,
		0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc5, 0x20, 0x05, 0xac, 0x10,
		0x01, 0x00, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x3e, 0x4c, 0x00, 0x24, 0xff, 0xff,
		0xff, 0x00, 0x80, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xc5,
		0x20, 0x05, 0xac, 0x10, 0x00, 0x00, 0x02, 0x02, 0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x37, 0x57,
		0x00, 0x24, 0xff, 0xff, 0xff, 0xfc, 0x80, 0x00, 0x00, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00}

	/*	lsa_update = []byte{0x00, 0x00, 0x00, 0x01, 0xe, 0x10,
		0x22, 0x02, 0x0a, 0x00, 0x14, 0x02, 0x05, 0x05, 0x05, 0x05, 0x80, 0x00, 0x00, 0x02, 0xf4, 0xee,
		0x00, 0x20, 0xff, 0xff, 0xff, 0xfc, 0x05, 0x05, 0x05, 0x05, 0x04, 0x04, 0x04, 0x04} */

	lsa_fake = []byte{0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x02, 0x02,
		0x02, 0x02, 0x80, 0x00, 0x00, 0x01, 0x33, 0x56, 0x00, 0x24, 0xff, 0xff, 0xff, 0x00, 0x80, 0x00,
		0x00, 0x64, 0x00, 0x00, 0x00, 0x00}

}

func startDummyIntfChannels(key IntfConfKey) {
	ent, _ := ospf.IntfConfMap[key]
	ent.NeighborMap = make(map[NeighborConfKey]NeighborData)
	ent.NeighChangeCh = make(chan NeighChangeMsg)
	ent.NeighCreateCh = make(chan NeighCreateMsg)
	ent.BackupSeenCh = make(chan BackupSeenMsg)
	for {
		select {
		case data := <-ent.NeighChangeCh:
			fmt.Println("Received data from NeighChangeCh ", data)
		case data := <-ent.NeighCreateCh:
			fmt.Println("Received data from NeighCreateCh ", data)
		case data := <-ent.BackupSeenCh:
			fmt.Println("Received data from BackupSeenCh ", data)
		}
	}
}

func startDummyChannels(server *OSPFServer) {

	for {
		select {
		case data := <-server.neighborHelloEventCh:
			fmt.Println("Receieved data from  neighborHelloEventCh : ", data)

		case data := <-server.neighborDBDEventCh:
			fmt.Println("Receieved data from neighbor DBD : ", data)

		case data := <-server.NetworkDRChangeCh:
			fmt.Println("Received data from NetworkDRChangeCh ", data)

		case data := <-server.neighborConfCh:
			fmt.Println("Next state for nbr  ", data.ospfNbrEntry.OspfNbrState)

		case data := <-server.ospfNbrDBDSendCh:
			fmt.Println("Received data from ospfNbrDBDSendCh", data)

		case data := <-server.DbEventOp:
			fmt.Println("Received data from DbEventOp", data)

		case data := <-server.DbLsdbOp:
			fmt.Println("Received data on DbLsdbOp", data)

		case data := <-server.StartCalcSPFCh:
			fmt.Println("Recieved data on StartCalcSPFCh ", data)
			server.DoneCalcSPFCh <- true

		case data := <-server.ospfNbrLsaUpdSendCh:
			fmt.Println("Received data on ospfNbrLsaUpdSendCh ", data)

		case data := <-server.neighborIntfEventCh:
			fmt.Println("Received data on neighborIntfEventCh ", data)
			//case data := <-server.StartCalcSPFCh:
			//	fmt.Println("Received data on StartCalcSPFCh ", data)
		case data := <-server.ospfNbrLsaReqSendCh:
			fmt.Println("Received data on ospfNbrLsaReqSendCh ", data)
		case data := <-server.neighborLSAUpdEventCh:
			fmt.Println("Received data on neighborLSAUpdEventCh ", data)
		case data := <-server.LsdbUpdateCh:
			fmt.Println("Received data on LsdbUpdateCh ", data)
		case data := <-server.ospfNbrLsaAckSendCh:
			fmt.Println("Received data on ospfNbrLsaAckSendCh ", data)
		case data := <-server.neighborLSAACKEventCh:
			fmt.Println("Received data on neighborLSAACKEventCh ", data)
		case data := <-server.neighborLSAReqEventCh:
			fmt.Println("Received data on neighborLSAReqEventCh ", data)
		case data := <-server.IntfSliceRefreshCh:
			fmt.Println("Received data on IntfSliceRefreshCh ", data)
			server.IntfSliceRefreshDoneCh <- true
		case data := <-server.neighborSliceStartCh:
			fmt.Println("Received data on neighborSliceStartCh ", data)
		case data := <-server.ExternalRouteNotif:
			fmt.Println("Received data on ExternalRouteNotif ", data)
		}
	}

}

func getServerObject() *OSPFServer {
	logger, err := OSPFNewLogger("ospfd", "OSPFTEST", true)
	if err != nil {
		fmt.Println("ospftest: creating logger failed")
	}
	ospfServer := NewOSPFServer(logger)
	if ospfServer == nil {
		fmt.Sprintln("ospf server object is null ")
	}
	ospf = ospfServer
	return ospfServer
}
