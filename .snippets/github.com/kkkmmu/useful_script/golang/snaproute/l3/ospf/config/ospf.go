//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package config

import ()

type AreaId string
type RouterId string
type Metric int                   // 0x0 to 0xffff
type BigMetric int                // 0x0 to 0xffffff
type PositiveInteger int32        // 0x0 to 0x7fffffff
type HelloRange int               // 0x1 to 0xffff
type UpToMaxAge int               // 0x0 to 3600
type DesignatedRouterPriority int // 0x0 to 0xff
type TosType int                  // 0x0 to 30
type IpAddress string
type InterfaceIndexOrZero int

/*Architectural constants */
const (
	MaxAge        uint16 = 3600 //60 min
	LSRefreshTime uint32 = 1800 // 30min
	AllSPFRouters string = "224.0.0.5"
	AllDRouters   string = "224.0.0.6"
	McastMAC      string = "01:00:5e:00:00:05"
)

/*Event types */
const (
	ADJACENCY string = "ADJACENCY"
	INTF      string = "INTF"
	RIB       string = "RIB"
	FLOOD     string = "FLOOD"
	AREA      string = "AREA"
	SPF       string = "SPF"
	LSA       string = "LSA"
)

type Status int

const (
	Enabled   Status = 1
	Disabled  Status = 2
	Intf_Down Status = 0
	Intf_Up   Status = 1
)

type MsgType int

const (
	CREATE MsgType = 0
	UPDATE MsgType = 1
	DELETE MsgType = 2
)

type AuthType int

const (
	NoAuth         AuthType = 0
	SimplePassword AuthType = 1
	Md5            AuthType = 2
	Reserved       AuthType = 3
)

type RestartSupport int

const (
	None                RestartSupport = 1
	PlannedOnly         RestartSupport = 2
	PlannedAndUnplanned RestartSupport = 3
)

type AdvertiseAction int

const (
	DoNotAdvertise AdvertiseAction = 1
	Advertise      AdvertiseAction = 2
)

type ImportAsExtern int

const (
	ImportExternal   ImportAsExtern = 1
	ImportNoExternal ImportAsExtern = 2
	ImportNssa       ImportAsExtern = 3
)

type AreaSummary int

const (
	NoAreaSummary   AreaSummary = 1
	SendAreaSummary AreaSummary = 2
)

type NssaTranslatorRole int

const (
	Always    NssaTranslatorRole = 1
	Candidate NssaTranslatorRole = 2
)

type MetricType int

const (
	OspfMetric     MetricType = 1
	ComparableCost MetricType = 2
	NonComparable  MetricType = 3
)

type AreaRangeEffect int

const (
	AdvertiseMatching      AreaRangeEffect = 1
	DoNotAdvertiseMatching AreaRangeEffect = 2
)

type IfType int

const (
	Broadcast         IfType = 1
	Nbma              IfType = 2
	NumberedP2P       IfType = 3
	UnnumberedP2P     IfType = 4
	PointToMultipoint IfType = 5
	Stub              IfType = 6
)

type ConfFlag int

const (
	IF_IP_ADDRESS        ConfFlag = 0
	IF_ADDRESS_LESS_IF   ConfFlag = 1
	IF_ADMIN_STAT        ConfFlag = 2
	IF_TYPE              ConfFlag = 3
	IF_RTR_PRIORITY      ConfFlag = 4
	IF_TRANSIT_DELAY     ConfFlag = 5
	IF_RETRANS_INTERVAL  ConfFlag = 6
	IF_HELLO_INTERVAL    ConfFlag = 7
	IF_RTR_DEAD_INTERVAL ConfFlag = 8
	IF_POLL_INTERVAL     ConfFlag = 9
	IF_AREA_ID           ConfFlag = 10
)

var IfTypeList = []string{
	"Undefined",
	"Broadcast",
	"Nbma",
	"NumberedP2P",
	"UnnumberedP2P",
	"PointToMultipoint",
	"Stub"}

type MulticastForwarding int

const (
	Blocked   MulticastForwarding = 1
	Multicast MulticastForwarding = 2
	Unicast   MulticastForwarding = 3
)

type RestartStatus int

const (
	NotRestarting    RestartStatus = 1
	PlannedRestart   RestartStatus = 2
	UnplannedRestart RestartStatus = 3
)

type RestartExitReason int

const (
	NoAttempt       RestartExitReason = 1
	InProgress      RestartExitReason = 2
	Completed       RestartExitReason = 3
	TimeedOut       RestartExitReason = 4
	TopologyChanged RestartExitReason = 5
)

type NssaTranslatorState int

const (
	NssaTranslatorEnabled  NssaTranslatorState = 1
	NssaTranslatorElected  NssaTranslatorState = 2
	NssaTranslatorDisabled NssaTranslatorState = 3
)

type LsaType int

const (
	RouterLink       LsaType = 1
	NetworkLink      LsaType = 2
	SummaryLink      LsaType = 3
	AsSummaryLink    LsaType = 4
	AsExternalLink   LsaType = 5
	MulticastLink    LsaType = 6
	NssaExternalLink LsaType = 7
	LocalOpaqueLink  LsaType = 9
	AreaOpaqueLink   LsaType = 10
	AsOpaqueLink     LsaType = 11
)

type IfState int

const (
	Down                   IfState = 1
	Loopback               IfState = 2
	Waiting                IfState = 3
	P2P                    IfState = 4
	OtherDesignatedRouter  IfState = 5
	DesignatedRouter       IfState = 6
	BackupDesignatedRouter IfState = 7
)

type NbrState int

const (
	NbrDown          NbrState = 1
	NbrAttempt       NbrState = 2
	NbrInit          NbrState = 3
	NbrTwoWay        NbrState = 4
	NbrExchangeStart NbrState = 5
	NbrExchange      NbrState = 6
	NbrLoading       NbrState = 7
	NbrFull          NbrState = 8
)

var NbrStateList = []string{
	"Undef",
	"NbrDown",
	"NbrAttempt",
	"NbrInit",
	"NbrTwoWay",
	"NbrExchangeStart",
	"NbrExchange",
	"NbrLoading",
	"NbrFull"}

type NbrEvent int

const (
	Nbr1WayReceived    NbrEvent = 1
	Nbr2WayReceived    NbrEvent = 2
	NbrNegotiationDone NbrEvent = 3
	NbrExchangeDone    NbrEvent = 4
	NbrLoadingDone     NbrEvent = 5
)

type NbmaNbrPermanence int

const (
	DynamicNbr   NbmaNbrPermanence = 1
	PermanentNbr NbmaNbrPermanence = 2
)

type NbrRestartHelperStatus int

const (
	NotHelping NbrRestartHelperStatus = 1
	Helping    NbrRestartHelperStatus = 2
)

type GlobalConf struct {
	Vrf                string
	RouterId           RouterId
	AdminStat          Status
	ASBdrRtrStatus     bool
	TOSSupport         bool
	RestartSupport     RestartSupport
	RestartInterval    int32
	ReferenceBandwidth uint32
}

type GlobalState struct {
	RouterId          RouterId
	VersionNumber     int32
	AreaBdrRtrStatus  bool
	ExternLsaCount    int32
	ExternLsaChecksum int32
	OriginateNewLsas  int32
	RxNewLsas         int32
	OpaqueLsaSupport  bool
	RestartStatus     RestartStatus
	RestartAge        int32
	RestartExitReason RestartExitReason
	AsLsaCount        int32
	AsLsaCksumSum     int32
	StubRouterSupport bool
	//DiscontinuityTime        string
	DiscontinuityTime int32 //This should be string
}

// Indexed By AreaId
type AreaConf struct {
	AreaId                 AreaId
	AuthType               AuthType
	ImportAsExtern         ImportAsExtern
	AreaSummary            AreaSummary
	StubDefaultCost        int32
	AreaNssaTranslatorRole NssaTranslatorRole
}

type AreaState struct {
	AreaId                   AreaId
	SpfRuns                  int32
	AreaBdrRtrCount          int32
	AsBdrRtrCount            int32
	AreaLsaCount             int32
	AreaLsaCksumSum          int32
	AreaNssaTranslatorState  NssaTranslatorState
	AreaNssaTranslatorEvents int32
}

// Indexed by StubAreaId and StubTOS
type StubAreaConf struct {
	StubAreaId     AreaId
	StubTOS        TosType
	StubMetric     BigMetric
	StubMetricType MetricType
}

type StubAreaState struct {
	StubAreaId     AreaId
	StubTOS        TosType
	StubMetric     BigMetric
	StubMetricType MetricType
}

// Indexed by LsdbAreaId, LsdbType, LsdbLsid, LsdbRouterId
type LsdbState struct {
	LsdbAreaId        AreaId
	LsdbType          LsaType
	LsdbLsid          IpAddress
	LsdbRouterId      RouterId
	LsdbSequence      int
	LsdbAge           int
	LsdbCheckSum      int
	LsdbAdvertisement string
}

// Indexed By RangeAreaId, RangeNet
type AreaRangeConf struct {
	RangeAreaId     AreaId
	AreaRangeNet    IpAddress
	ArearangeMask   IpAddress
	AreaRangeEffect AreaRangeEffect
}

type AreaRangeState struct {
	RangeAreaId     AreaId
	AreaRangeNet    IpAddress
	ArearangeMask   IpAddress
	AreaRangeEffect AreaRangeEffect
}

// Indexed By HostIpAddress, HostTOS
type HostConf struct {
	HostIpAddress IpAddress
	HostTOS       TosType
	HostMetric    Metric
	HostCfgAreaID AreaId
}

type HostState struct {
	HostIpAddress IpAddress
	HostTOS       TosType
	HostMetric    Metric
	HostCfgAreaID AreaId
}

// Indexed By IfIpAddress, AddressLessIf

type InterfaceConf struct {
	IfIpAddress       IpAddress
	AddressLessIf     InterfaceIndexOrZero
	IfAreaId          AreaId
	IfType            IfType
	IfAdminStat       Status
	IfRtrPriority     DesignatedRouterPriority
	IfTransitDelay    UpToMaxAge
	IfRetransInterval UpToMaxAge
	IfHelloInterval   HelloRange
	IfRtrDeadInterval PositiveInteger
	IfPollInterval    PositiveInteger
	IfAuthKey         string
	IfAuthType        AuthType
}

type InterfaceState struct {
	IfIpAddress                IpAddress
	AddressLessIf              InterfaceIndexOrZero
	IfState                    IfState
	IfDesignatedRouter         IpAddress
	IfBackupDesignatedRouter   IpAddress
	IfEvents                   int32
	IfLsaCount                 int32
	IfLsaCksumSum              int32
	IfDesignatedRouterId       RouterId
	IfBackupDesignatedRouterId RouterId
}

// Indexed By  IfMetricIpAddress, IfMetricAddressLessIf, IfMetricTOS
type IfMetricConf struct {
	IfMetricIpAddress     IpAddress
	IfMetricAddressLessIf InterfaceIndexOrZero
	IfMetricTOS           TosType
	IfMetricValue         Metric
}

type IfMetricState struct {
	IfMetricIpAddress     IpAddress
	IfMetricAddressLessIf InterfaceIndexOrZero
	IfMetricTOS           TosType
	IfMetricValue         Metric
}

// Indexed By VirtIfAreaId, VirtIfNeighbor
type VirtIfConf struct {
	VirtIfAreaId          AreaId
	VirtIfNeighbor        RouterId
	VirtIfTransitDelay    UpToMaxAge
	VirtIfRetransInterval UpToMaxAge
	VirtIfHelloInterval   HelloRange
	VirtIfRtrDeadInterval PositiveInteger
	VirtIfAuthKey         string
	VirtIfAuthType        AuthType
}

type VirtIfState struct {
	VirtIfAreaId          AreaId
	VirtIfNeighbor        RouterId
	VirtIfTransitDelay    UpToMaxAge
	VirtIfRetransInterval UpToMaxAge
	VirtIfHelloInterval   HelloRange
	VirtIfRtrDeadInterval PositiveInteger
	VirtIfState           IfState
	VirtIfEvents          int
	VirtIfAuthKey         string
	VirtIfAuthType        AuthType
	VirtIfLsaCount        int
	VirtIfLsaCksumSum     int
}

// Indexed by NbrIpAddress, NbrAddressLessIndex
type NbrConf struct {
	NbrIpAddress        IpAddress
	NbrAddressLessIndex InterfaceIndexOrZero
	NbrPriority         DesignatedRouterPriority
}

type NeighborState struct {
	NbrIpAddress               IpAddress
	NbrAddressLessIndex        int
	NbrRtrId                   string
	NbrOptions                 int
	NbrPriority                uint8
	NbrState                   string
	NbrEvents                  int
	NbrLsRetransQLen           int
	NbmaNbrPermanence          int
	NbrHelloSuppressed         bool
	NbrRestartHelperStatus     int
	NbrRestartHelperAge        uint32
	NbrRestartHelperExitReason int
}

// Virtual Neighbor Table (Read Only)
// OSPF Virtual Neighbor Entry
// Indexed By VirtNbrArea, VirtNbrRtrId
type VirtNbrState struct {
	VirtNbrArea                    AreaId
	VirtNbrRtrId                   RouterId
	VirtNbrIpAddress               IpAddress
	VirtNbrOptions                 int
	VirtNbrState                   NbrState
	VirtNbrEvents                  int
	VirtNbrLsRetransQLen           int
	VirtNbrHelloSuppressed         bool
	VirtNbrRestartHelperStatus     NbrRestartHelperStatus
	VirtNbrRestartHelperAge        uint32
	VirtNbrRestartHelperExitReason RestartExitReason
}

// External LSA link State - Deprecated
// Indexed by ExtLsdbType, ExtLsdbLsid, ExtLsdbRouterId
type ExtLsdbState struct {
	ExtLsdbType          LsaType
	ExtLsdbLsid          IpAddress
	ExtLsdbRouterId      RouterId
	ExtLsdbSequence      int
	ExtLsdbAge           int
	ExtLsdbChecksum      int
	ExtLsdbAdvertisement string
}

// OSPF Area Aggregate Table
// Replaces OSPF Area Summary Table
// Indexed By AreaAggregateAreaId, AreaAggregateLsdbType,
// AreaAggregateNet, AreaAggregateMask
type AreaAggregateConf struct {
	AreaAggregateAreaId      AreaId
	AreaAggregateLsdbType    LsaType
	AreaAggregateNet         IpAddress
	AreaAggregateMask        IpAddress
	AreaAggregateEffect      AreaRangeEffect
	AreaAggregateExtRouteTag int
}

type AreaAggregateState struct {
	AreaAggregateAreaId      AreaId
	AreaAggregateLsdbType    LsaType
	AreaAggregateNet         IpAddress
	AreaAggregateMask        IpAddress
	AreaAggregateEffect      AreaRangeEffect
	AreaAggregateExtRouteTag uint32
}

// Link local link state database for non-virtual links
// Indexed by LocalLsdbIpAddress, LocalLsdbAddressLessIf,
// LocalLsdbType, LocalLsdbLsid, LocalLsdbRouterId
type LocalLsdbState struct {
	LocalLsdbIpAddress     IpAddress
	LocalLsdbAddressLessIf InterfaceIndexOrZero
	LocalLsdbType          LsaType
	LocalLsdbLsid          IpAddress
	LocalLsdbRouterId      RouterId
	LocalLsdbSequence      int
	LocalLsdbAge           int
	LocalLsdbChecksum      int
	LocalLsdbAdvertisement string
}

// Link State Database, link-local for Virtual Links
// Indexed By VirtLocalLsdbTransitArea, VirtLocalLsdbTransitArea,
// VirtLocalLsdbType, VirtLocalLsdbLsid, VirtLocalLsdbRouterId
type VirtLocalLsdbState struct {
	VirtLocalLsdbTransitArea   AreaId
	VirtLocalLsdbNeighbor      RouterId
	VirtLocalLsdbType          LsaType
	VirtLocalLsdbLsid          IpAddress
	VirtLocalLsdbRouterId      RouterId
	VirtLocalLsdbSequence      int
	VirtLocalLsdbAge           int
	VirtLocalLsdbChecksum      int
	VirtLocalLsdbAdvertisement string
}

// Link State Database, AS - scope
// Indexed AsLsdbType, AsLsdbLsid, AsLsdbRouterId
type AsLsdbState struct {
	AsLsdbType          LsaType
	AsLsdbLsid          IpAddress
	AsLsdbRouterId      RouterId
	AsLsdbSequence      int
	AsLsdbAge           int
	AsLsdbChecksum      int
	AsLsdbAdvertisement string
}

// Area LSA Counter Table
// Indexed By AreaLsaCountAreaId, AreaLsaCountLsaType
type OspfAreaLsaCountState struct {
	AreaLsaCountAreaId  AreaId
	AreaLsaCountLsaType LsaType
	AreaLsaCountNumber  int
}

type OspfIPv4Route struct {
	DestId          string
	AddrMask        string
	DestType        string
	OptCapabilities int32
	AreaId          string
	PathType        string
	Cost            int32
	Type2Cost       int32
	NumOfPaths      int32
	NextHops        string
	LSOrigin        string
}

type OspfEventState struct {
	TimeStamp string
	EventType string
	EventInfo string
}

type InterfaceRpcMsg struct {
	IntfConf InterfaceConf
	Op       MsgType
}
