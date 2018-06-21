package n2x

import (
	"errors"
	"fmt"
	"strings"
)

type OSPF struct {
	Name               string
	Version            string
	Handler            string
	RID                string
	SutRID             string
	AreaID             string
	RouterHandler      string
	HelloInterval      string
	DeadInterval       string
	PollInterval       string
	TransmitDelay      string
	Priority           string
	ExternalRoutePools map[string]*RoutePool
	SummaryRoutePools  map[string]*RoutePool
	Networks           map[string]*Network
	Routers            map[string]*Router
	Cost               string
	MTU                string
	Pool               *RouterPool
	LSDB               *LSDB
	Neighbour          *OSPFNeighbour
	*Port
}

type Router struct {
	Handler string
}

type Network struct {
	Handler string
}

const (
	Summary = iota
	External
)

type RoutePool struct {
	Type    int
	First   string
	Plen    string
	Count   string
	Handler string
	Step    string
	Default bool
	Object  string
	LSAPool *LSAPool
	*OSPF
}

type LSAPool struct {
	Object  string
	Handler string
}

var ErrorNoExternalRoutePoolExist = errors.New("No external route pool exist")
var ErrorNoObjectConnections = errors.New("No external route pool exist")

func (rp *RoutePool) Sync() error {
	if rp.Object == "" {
		return fmt.Errorf("You cannot get routes of pool without object")
	}

	err := rp.GetRoutes()
	if err != nil {
		return fmt.Errorf("Cannot sync routepool %s with: %s", rp.Handler, err)
	}

	if !rp.OSPF.isObjectAlreadyConnected(rp.Handler) {
		err = rp.OSPF.ConnectObjects(rp.Handler)
		if err != nil {
			return fmt.Errorf("Cannot sync route pool with: %s", err)
		}
	}

	_, err = rp.GetExternalLsaPool()
	if err != nil {
		return fmt.Errorf("Cannot sync route pool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) GetRoutes() error {
	if rp.Handler == "" {
		return fmt.Errorf("You cannot get routes to empty route pool")
	}

	cmd := fmt.Sprintf("%s GetRoutes %s", rp.Object, rp.Handler)
	res, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get routes of routepool with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	fields := strings.Split(res, " ")
	if len(fields) != 4 {
		return fmt.Errorf("Cannot invaid return value when get routes: %s", res)
	}

	//pFirstRoute pPrefixLength pNumRoutes pModifier
	rp.First = strings.TrimSpace(fields[0])
	rp.Plen = strings.TrimSpace(fields[1])
	rp.Count = strings.TrimSpace(fields[2])
	rp.Step = strings.TrimSpace(fields[3])

	return nil
}

func (rp *RoutePool) Advertise() error {
	//AgtOspfExternalLsaPool Advertise
	cmd := fmt.Sprintf("%s Advertise %s", rp.LSAPool.Object, rp.LSAPool.Handler)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot advertise lsa pool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) Withdraw() error {
	//AgtOspfExternalLsaPool Withdraw
	cmd := fmt.Sprintf("%s Withdraw %s", rp.LSAPool.Object, rp.LSAPool.Handler)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot withdraw lsa pool with: %s", err)
	}

	return nil
}

//AgtOspfExternalRoutePool SetRoutes
func (rp *RoutePool) SetRoutes(ip, plen, count, step string) error {
	if rp.Handler == "" {
		return fmt.Errorf("You cannot set routes to empty route pool")
	}

	ip = strings.TrimSpace(ip)
	plen = strings.TrimSpace(plen)
	count = strings.TrimSpace(count)
	step = strings.TrimSpace(step)

	cmd := fmt.Sprintf("%s SetRoutes %s %s %s %s %s", rp.Object, rp.Handler, ip, plen, count, step)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set routes to routepool with: %s", err)
	}

	rp.First = ip
	rp.Plen = plen
	rp.Count = count
	rp.Step = step

	return nil
}

func (rp *RoutePool) SetMetric(metric string) error {
	cmd := fmt.Sprintf("%s SetMetric %s %s", rp.Object, rp.Handler, metric)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set metric to routepool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) SetMetricType(typ string) error {
	typ = strings.TrimSpace(typ)
	if typ != "1" && typ != "2" {
		return fmt.Errorf("Invalid OSPF external metric type: %s", typ)
	}

	mtypMap := map[string]string{
		"1": "AGT_OSPF_EXTERNAL_METRIC_TYPE_1",
		"2": "AGT_OSPF_EXTERNAL_METRIC_TYPE_2",
	}

	cmd := fmt.Sprintf("%s SetMetricType %s %s", rp.Object, rp.Handler, mtypMap[typ])
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set metrictype to routepool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) SetForwardingAddress(fwdr string) error {
	cmd := fmt.Sprintf("%s SetForwardingAddress %s %s", rp.Object, rp.Handler, fwdr)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set forwarding address to routepool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) SetNssaFlag() error {
	cmd := fmt.Sprintf("%s SetNssaFlag %s %s", rp.Object, rp.Handler, "1")
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set nssa flags to routepool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) UnSetNssaFlag() error {
	cmd := fmt.Sprintf("%s SetNssaFlag %s %s", rp.Object, rp.Handler, "0")
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot unset nssa flags to routepool with: %s", err)
	}

	return nil
}

func (rp *RoutePool) GetExternalLsaPool() (*LSAPool, error) {
	//AgtOspfExternalRoutePool GetExternalLsaPool
	cmd := fmt.Sprintf("%s GetExternalLsaPool %s", rp.Object, rp.Handler)
	res, err := rp.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get lsa pool of %s with: %s", rp.Handler, err)
	}

	res = strings.TrimSpace(res)

	rp.LSAPool = &LSAPool{Handler: res, Object: "AgtOspfExternalLsaPool"}

	return rp.LSAPool, nil
}

type OSPFNeighbour struct {
	Handler string
	AreaID  string
	RID     string
	State   string
	*Port
}

type RouterPool struct {
	Handler string
	*Port
}

type LSDB struct {
	Handler string
	*Port
	RouterLSAs         map[string]*LSA
	NetworkLSAs        map[string]*LSA
	NetworkSummaryLSAs map[string]*LSA
	ASBRSummaryLSAs    map[string]*LSA
	ASExternalLSAs     map[string]*LSA
	NSSAExternalLSAs   map[string]*LSA
}

const (
	OSPF_ROUTER_LSA = iota
	OSPF_NETWORK_LSA
	OSPF_SUMMARY_LSAPOOL
	OSPF_ASBR_SUMMARY_LSA
	OSPF_EXTERNAL_LSAPOOL
	OSPF_NSSA_EXTERNAL_LSAPOOL
	OSPF_OPAQUE_LSA
	OSPF_GRACEFUL_RESTART_LSA
	OSPF_TE_ROUTER_LSA
	OSPF_TE_LINK_LSA
	OSPFV3_ROUTER_LSA
	OSPFV3_NETWORK_LSA
	OSPFV3_INTER_AREA_PREFIX_LSAPOOL
	OSPFV3_INTER_AREA_ROUTER_LSA
	OSPFV3_EXTERNAL_LSAPOOL
	OSPFV3_NSSA_EXTERNAL_LSAPOOL
	OSPFV3_LINK_LSA
	OSPFV3_INTRA_AREA_PREFIX_LSA
	OSPFV3_GRACEFUL_RESTART_LSA
)

var LsaTypeNameMap = map[int]string{
	OSPF_ROUTER_LSA:                  "AGT_OSPF_ROUTER_LSA",
	OSPF_NETWORK_LSA:                 "AGT_OSPF_NETWORK_LSA",
	OSPF_SUMMARY_LSAPOOL:             "AGT_OSPF_SUMMARY_LSAPOOL",
	OSPF_ASBR_SUMMARY_LSA:            "AGT_OSPF_ASBR_SUMMARY_LSA",
	OSPF_EXTERNAL_LSAPOOL:            "AGT_OSPF_EXTERNAL_LSAPOOL",
	OSPF_NSSA_EXTERNAL_LSAPOOL:       "AGT_OSPF_NSSA_EXTERNAL_LSAPOOL",
	OSPF_OPAQUE_LSA:                  "AGT_OSPF_OPAQUE_LSA",
	OSPF_GRACEFUL_RESTART_LSA:        "AGT_OSPF_GRACEFUL_RESTART_LSA",
	OSPF_TE_ROUTER_LSA:               "AGT_OSPF_TE_ROUTER_LSA",
	OSPF_TE_LINK_LSA:                 "AGT_OSPF_TE_LINK_LSA",
	OSPFV3_ROUTER_LSA:                "AGT_OSPFV3_ROUTER_LSA",
	OSPFV3_NETWORK_LSA:               "AGT_OSPFV3_NETWORK_LSA",
	OSPFV3_INTER_AREA_PREFIX_LSAPOOL: "AGT_OSPFV3_INTER_AREA_PREFIX_LSAPOOL",
	OSPFV3_INTER_AREA_ROUTER_LSA:     "AGT_OSPFV3_INTER_AREA_ROUTER_LSA",
	OSPFV3_EXTERNAL_LSAPOOL:          "AGT_OSPFV3_EXTERNAL_LSAPOOL",
	OSPFV3_NSSA_EXTERNAL_LSAPOOL:     "AGT_OSPFV3_NSSA_EXTERNAL_LSAPOOL",
	OSPFV3_LINK_LSA:                  "AGT_OSPFV3_LINK_LSA",
	OSPFV3_INTRA_AREA_PREFIX_LSA:     "AGT_OSPFV3_INTRA_AREA_PREFIX_LSA",
	OSPFV3_GRACEFUL_RESTART_LSA:      "AGT_OSPFV3_GRACEFUL_RESTART_LSA",
}

type LSA struct {
	Handler string
	*Port
}

//AgtOspfLsaDatabase ListLsas
func (lsdb *LSDB) GetAllLSAs() ([]*LSA, error) {
	cmd := fmt.Sprintf("AgtOspfLsaDatabase ListLsas %s", lsdb.Handler)
	_, err := lsdb.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get ospf neighbor rid %s : %s", lsdb.Handler, err.Error())
	}

	return nil, nil
}

//AgtOspfLsaDatabase AddLsa
func (lsdb *LSDB) AddLSA(lsaType int) ([]*LSA, error) {
	name, ok := LsaTypeNameMap[lsaType]
	if !ok {
		return nil, fmt.Errorf("Unknown lsa type: %d", lsaType)
	}

	cmd := fmt.Sprintf("AgtOspfLsaDatabase AddLsa %s %s", lsdb.Handler, name)
	_, err := lsdb.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot Add Lsa to rid %s : %s", lsdb.Handler, err.Error())
	}

	return nil, nil
}

func (on *OSPFNeighbour) Sync() error {
	if on.Handler == "" {
		return fmt.Errorf("You must set the handler for this neighbor")
	}

	if on.Port == nil {
		return fmt.Errorf("You must set the port for this neighbor")
	}

	_, err := on.GetRouterID()
	if err != nil {
		return fmt.Errorf("Cannot sync neighbor with: %s", err)
	}

	_, err = on.GetIPAddress()
	if err != nil {
		return fmt.Errorf("Cannot sync neighbor with: %s", err)
	}

	return nil
}

func (on *OSPFNeighbour) GetRouterID() (string, error) {
	cmd := fmt.Sprintf("AgtOspfNeighbor GetRouterId %s", on.Handler)
	res, err := on.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get ospf neighbor rid %s : %s", on.Handler, err.Error())
	}

	on.RID = strings.TrimSpace(res)

	return on.RID, nil
}

func (on *OSPFNeighbour) SetRouterID(rid string) error {
	//AgtOspfNeighbor SetRouterId
	cmd := fmt.Sprintf("AgtOspfNeighbor SetRouterId %s %s", on.Handler, rid)
	_, err := on.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set ospf neighbor rid %s : %s", on.Handler, err.Error())
	}

	return nil
}

func (on *OSPFNeighbour) GetIPAddress() (string, error) {
	cmd := fmt.Sprintf("AgtOspfNeighbor GetIpAddress %s", on.Handler)
	res, err := on.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get ospf neighbor ip %s : %s", on.Handler, err.Error())
	}

	return strings.TrimSpace(res), nil
}

func (on *OSPFNeighbour) SetIPAddress(ip string) error {
	cmd := fmt.Sprintf("AgtOspfNeighbor SetIpAddress %s %s", on.Handler, ip)
	_, err := on.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set ospf neighbor ip %s : %s", on.Handler, err.Error())
	}

	return nil
}

func (o *OSPF) Sync() error {
	if o.Handler == "" {
		return fmt.Errorf("You must set the OSPF session handler first")
	}

	_, err := o.GetRouterHandler()
	if err != nil {
		return fmt.Errorf("Cannot sync ospf with: %s", err)
	}

	name, err := o.GetName()
	if err != nil {
		return fmt.Errorf("Cannot sync ospf with: %s", err)
	}

	area, err := o.GetAreaID()
	if err != nil {
		return fmt.Errorf("Cannot sync ospf with: %s", err)
	}

	rid, err := o.GetRouterID()
	if err != nil {
		return fmt.Errorf("Cannot sync ospf with: %s", err)
	}

	n, err := o.GetNeighbour()
	if err != nil {
		return fmt.Errorf("Cannot sync ospf with: %s", err)
	}

	lsdb, err := o.GetLSDB()
	if err != nil {
		return fmt.Errorf("Cannot sync ospf with: %s", err)
	}

	o.Name = name
	o.AreaID = area
	o.RID = rid
	o.Neighbour = n
	o.LSDB = lsdb

	return nil
}

func (o *OSPF) SetName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("OSPF Session name must not be empty")
	}

	cmd := fmt.Sprintf("AgtTestTopology SetSessionName %s %s", o.Handler, name)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set ospf session name for %s : %s", o.Handler, err.Error())
	}

	o.Name = name

	return nil
}

func (o *OSPF) GetName() (string, error) {
	cmd := fmt.Sprintf("AgtTestTopology GetSessionName %s", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get ospf session name for %s : %s", o.Handler, err.Error())
	}

	res = strings.Replace(res, "\"", "", -1)

	o.Name = strings.TrimSpace(res)

	return o.Name, nil
}

func (o *OSPF) Enable() error {
	cmd := fmt.Sprintf("AgtTestTopology EnableSession %s", o.Handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable ospf session for %s : %s", o.Handler, err.Error())
	}

	return nil
}

func (o *OSPF) Disable() error {
	cmd := fmt.Sprintf("AgtTestTopology DisableSession %s", o.Handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable ospf session for %s : %s", o.Handler, err.Error())
	}

	return nil
}

//AgtOspfSession GetSessionRouter
func (o *OSPF) GetRouterHandler() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetSessionRouter %s", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get ospf session router for %s : %s", o.Handler, err.Error())
	}

	o.RouterHandler = strings.TrimSpace(res)

	return o.RouterHandler, nil
}

func (o *OSPF) SetRouterID(id string) error {
	cmd := fmt.Sprintf("AgtOspfRouter SetRouterId %s %s", o.RouterHandler, id)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set ospf router id for %s(%s) : %s", o.Handler, o.RouterHandler, err.Error())
	}

	o.RID = id

	return nil
}

func (o *OSPF) GetRouterID() (string, error) {
	cmd := fmt.Sprintf("AgtOspfRouter GetRouterId %s", o.RouterHandler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get ospf router id for %s : %s", o.Handler, err.Error())
	}

	o.RID = strings.TrimSpace(res)

	return o.RID, nil
}

func (o *OSPF) SetAreaID(id string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetAreaId %s %s", o.Handler, id)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set ospf area id for %s : %s", o.Handler, err.Error())
	}

	o.AreaID = id

	return nil
}

func (o *OSPF) GetAreaID() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetAreaId %s", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get ospf area id for %s : %s", o.Handler, err.Error())
	}

	o.AreaID = strings.TrimSpace(res)

	return o.RID, nil
}

func (o *OSPF) GetNeighbour() (*OSPFNeighbour, error) {
	cmd := fmt.Sprintf("AgtOspfSession ListNeighbors %s", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get ospf neighbors for %s : %s", o.Handler, err.Error())
	}

	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "{", "", -1)
	res = strings.TrimSpace(res)

	o.Neighbour = nil

	fields := strings.Split(res, " ")
	for i, field := range fields {
		if i == 0 {
			if strings.TrimSpace(field) == "" {
				continue
			}

			n := &OSPFNeighbour{
				Handler: strings.TrimSpace(field),
				Port:    o.Port,
			}

			err := n.Sync()
			if err != nil {
				return nil, fmt.Errorf("Cannot get neighbor with: %s", err)
			}

			//Currently we only use the first one
			o.Neighbour = n
			break
		}
	}

	if o.Neighbour == nil {
		return nil, fmt.Errorf("Cannot find any ospf neighbor on: %s", o.Handler)
	}

	return o.Neighbour, nil
}

//AgtOspfSession AddNeighbor
func (o *OSPF) AddNeighbor() (*OSPFNeighbour, error) {
	cmd := fmt.Sprintf("AgtOspfSession AddNeighbor %s", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get add neighbor for port %s : %s", o.Handler, err.Error())
	}

	o.Neighbour.Handler = strings.TrimSpace(res)

	return o.Neighbour, nil
}

func (o *OSPF) SetSutIPAddress(ip string) error {
	if o.Neighbour == nil {
		return fmt.Errorf("No neighbor exist for: %s", o.Handler)
	}

	return o.Neighbour.SetIPAddress(ip)
}

func (o *OSPF) SetSutRouterID(rid string) error {
	if o.Neighbour == nil {
		return fmt.Errorf("No neighbor exist for: %s", o.Handler)
	}

	return o.Neighbour.SetRouterID(rid)
}

//AGT_OSPF_ROUTER_PRIORITY,
//AGT_OSPF_HELLO_INTERVAL,
//AGT_OSPF_DEAD_INTERVAL,
//AGT_OSPF_POLL_INTERVAL,
//AGT_OSPF_LSA_RETRANSMIT_DELAY,
//AGT_OSPF_INTERFACE_COST,
//AGT_OSPF_TRANSIT_DELAY,
//AGT_OSPF_OPTIONS,
//AGT_OSPF_MTU,
//AGT_OSPF_INSTANCE_ID,
//AGT_OSPF_PASS_EVTS_TO_CLIENT,
//AGT_OSPF_GRACE_PERIOD,
//AGT_OSPF_RESTART_REASON,
//AGT_OSPF_RESTART_INTERVAL,
//AGT_OSPF_RESTART_MODE

func (o *OSPF) SetPriority(pri string) error {
	cmd := fmt.Sprintf("AgtOspfSession EnableNeighborDiscovery %s", o.Handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set priorty for %s : %s", o.Handler, err.Error())
	}

	cmd = fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_ROUTER_PRIORITY %s", o.Handler, pri)
	_, err = o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set priorty for %s : %s", o.Handler, err.Error())
	}

	o.Priority = strings.TrimSpace(pri)

	return nil
}

func (o *OSPF) EnableNeighrorDiscovery() error {
	cmd := fmt.Sprintf("AgtOspfSession EnableNeighborDiscovery %s", o.Handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable ospf neighbor discovery for %s : %s", o.Handler, err.Error())
	}

	return nil
}

func (o *OSPF) DisableNeighrorDiscovery() error {
	cmd := fmt.Sprintf("AgtOspfSession DisableNeighborDiscovery %s", o.Handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disable ospf neighbor discovery for %s : %s", o.Handler, err.Error())
	}

	return nil
}

func (o *OSPF) SetHelloInterval(inter string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_HELLO_INTERVAL %s", o.Handler, inter)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set Hello interval for %s : %s", o.Handler, err.Error())
	}

	o.HelloInterval = strings.TrimSpace(inter)

	return nil
}

func (o *OSPF) SetDeadInterval(inter string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_DEAD_INTERVAL %s", o.Handler, inter)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set dead interval for %s : %s", o.Handler, err.Error())
	}

	o.DeadInterval = strings.TrimSpace(inter)

	return nil
}

func (o *OSPF) SetPollInterval(inter string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_POLL_INTERVAL %s", o.Handler, inter)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set poll interval for %s : %s", o.Handler, err.Error())
	}

	o.PollInterval = strings.TrimSpace(inter)

	return nil
}

func (o *OSPF) SetTransmitDelay(delay string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_TRANSIT_DELAY %s", o.Handler, delay)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set transmit delay for %s : %s", o.Handler, err.Error())
	}

	o.TransmitDelay = strings.TrimSpace(delay)

	return nil
}

func (o *OSPF) SetMTU(mtu string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_MTU %s", o.Handler, mtu)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set mtu for %s : %s", o.Handler, err.Error())
	}

	o.MTU = strings.TrimSpace(mtu)

	return nil
}

func (o *OSPF) SetInterfaceCost(cost string) error {
	cmd := fmt.Sprintf("AgtOspfSession SetInterfaceParameter %s AGT_OSPF_INTERFACE_COST %s", o.Handler, cost)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set cost for %s : %s", o.Handler, err.Error())
	}

	o.Cost = strings.TrimSpace(cost)

	return nil
}

func (o *OSPF) GetPriority() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_ROUTER_PRIORITY", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get priorty for %s : %s", o.Handler, err.Error())
	}

	o.Priority = strings.TrimSpace(res)

	return o.Priority, nil
}

func (o *OSPF) GetHelloInterval() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_HELLO_INTERVAL", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get Hello interval for %s : %s", o.Handler, err.Error())
	}

	o.HelloInterval = strings.TrimSpace(res)

	return o.HelloInterval, nil
}

func (o *OSPF) GetDeadInterval() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_DEAD_INTERVAL", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get dead interval for %s : %s", o.Handler, err.Error())
	}

	o.DeadInterval = strings.TrimSpace(res)

	return o.DeadInterval, nil
}

func (o *OSPF) GetPollInterval() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_POLL_INTERVAL", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get poll interval for %s : %s", o.Handler, err.Error())
	}

	o.PollInterval = strings.TrimSpace(res)

	return o.PollInterval, nil
}

func (o *OSPF) GetTransmitDelay() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_TRANSIT_DELAY", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get transmit delay for %s : %s", o.Handler, err.Error())
	}

	o.TransmitDelay = strings.TrimSpace(res)

	return o.TransmitDelay, nil
}

func (o *OSPF) GetMTU() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_MTU", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get mtu for %s : %s", o.Handler, err.Error())
	}

	o.MTU = strings.TrimSpace(res)

	return o.MTU, nil
}

func (o *OSPF) GetInterfaceCost() (string, error) {
	cmd := fmt.Sprintf("AgtOspfSession GetInterfaceParameter %s AGT_OSPF_INTERFACE_COST", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get cost for %s : %s", o.Handler, err.Error())
	}

	o.Cost = strings.TrimSpace(res)

	return o.Cost, nil
}

//invoke AgtTestTopology AddSession 1 ospfv2RouterPoolDevice
func (o *OSPF) AddRouterPool() (*RouterPool, error) {
	cmd := fmt.Sprintf("AgtTestTopology AddSession %s ospfv2RouterPoolDevice", o.Port.Handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot create router pool on port %s : %s", o.Name, err.Error())
	}

	return nil, nil
}

func (o *OSPF) GetLSDB() (*LSDB, error) {
	//    AgtOspfSession GetLsaDatabase
	cmd := fmt.Sprintf("AgtOspfSession GetLsaDatabase %s", o.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get lsdb on port %s : %s", o.Name, err.Error())
	}

	handler := strings.TrimSpace(res)
	if handler == "" {
		return nil, fmt.Errorf("Cannot get lsdb on port %s : %s", o.Name, err.Error())
	}

	lsdb := &LSDB{
		Handler: handler,
		Port:    o.Port,
	}

	o.LSDB = lsdb

	return o.LSDB, nil
}

/*
AGT_OSPF_ROUTER: an individual OSPF router
AGT_OSPFV2_ROUTER_PATTERN: a grid, ring, star, tree, or full mesh of OSPFv2 routers
AGT_OSPFV3_ROUTER_PATTERN: a grid, ring, star, tree, or full mesh of OSPFv3 routers
AGT_OSPF_ROUTER_GRID: a grid of routers (obsoleted by the above patterns)
AGT_OSPF_NETWORK: a subnet behind the selected session router
AGT_OSPF_SUMMARY_ROUTEPOOL: a pool of summary route addresses
AGT_OSPF_EXTERNAL_ROUTEPOOL: a pool of external route addresses
AGT_OSPFV3_INTER_AREA_PREFIX_ROUTEPOOL: a pool of OSPFv3 inter area prefix LSA, prefixes.
AGT_OSPFV3_EXTERNAL_ROUTEPOOL: a pool of external route addresses
*/
func (o *OSPF) AddRouter() (*Router, error) {
	cmd := fmt.Sprintf("AgtOspfTopology Add %s AGT_OSPF_ROUTER", o.LSDB.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf router to %s with: %s", o.Handler, err)
	}

	return &Router{Handler: strings.TrimSpace(res)}, nil
}

func (o *OSPF) AddNetwork() (*Network, error) {
	cmd := fmt.Sprintf("AgtOspfTopology Add %s AGT_OSPF_NETWORK", o.LSDB.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf network to %s with: %s", o.Handler, err)
	}

	return &Network{Handler: strings.TrimSpace(res)}, nil
}

func (o *OSPF) AddSummaryRoutePool() (*RoutePool, error) {
	cmd := fmt.Sprintf("AgtOspfTopology Add %s AGT_OSPF_SUMMARY_ROUTEPOOL", o.LSDB.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf summary route pool to %s with: %s", o.Handler, err)
	}

	return &RoutePool{Type: Summary, Handler: strings.TrimSpace(res)}, nil
}

func (o *OSPF) getDefaultExternalRoutePool() *RoutePool {
	for _, ep := range o.ExternalRoutePools {
		if ep.Default {
			return ep
		}
	}

	return nil
}

func (o *OSPF) GetDefaultExternalRoutePool() (*RoutePool, error) {
	if o.IsDefaultExternalRoutePoolExist() {
		return o.getDefaultExternalRoutePool(), nil
	}

	rp, err := o.AddExternalRoutePool()
	if err != nil {
		return nil, fmt.Errorf("Cannot get default external route pool with: %s", err)
	}

	return rp, nil
}

func (o *OSPF) IsDefaultExternalRoutePoolExist() bool {
	if o.ExternalRoutePools == nil || len(o.ExternalRoutePools) == 0 {
		return false
	}

	for _, ep := range o.ExternalRoutePools {
		if ep.Default {
			return true
		}
	}

	return false
}

func (o *OSPF) AddExternalPrefixes(start, plen, count, step, metric, mtyp, fwdr string) error {
	np, err := o.AddExternalRoutePool()
	if err != nil {
		return fmt.Errorf("Cannot add external prefix %s with %s", start, err)
	}

	err = np.SetRoutes(start, plen, count, step)
	if err != nil {
		return fmt.Errorf("Cannot add external prefix %s with %s", start, err)
	}

	err = np.SetMetricType(mtyp)
	if err != nil {
		return fmt.Errorf("Cannot add external prefix %s with %s", start, err)
	}

	err = np.SetMetric(metric)
	if err != nil {
		return fmt.Errorf("Cannot add external prefix %s with %s", start, err)
	}

	err = np.SetForwardingAddress(fwdr)
	if err != nil {
		return fmt.Errorf("Cannot add external prefix %s with %s", start, err)
	}

	return nil
}

func (o *OSPF) AddNSSAExternalPrefixes(start, plen, count, step, metric, mtyp, fwdr string) error {
	err := o.AddExternalPrefixes(start, plen, count, step, metric, mtyp, fwdr)
	if err != nil {
		return fmt.Errorf("Cannot add nssa external prefix with: %s", err)
	}

	rp := o.getExternalRoutePoolByIP(start)
	if rp == nil {
		return fmt.Errorf("Cannot add nssa external with cannot find route pool:%s", start)
	}

	err = rp.SetNssaFlag()
	if err != nil {
		return fmt.Errorf("Cannot add nssa external prefix with: %s", err)
	}

	return nil
}

func (o *OSPF) RemoveExternalPrefixes(first string) error {
	first = strings.TrimSpace(first)
	rp := o.getExternalRoutePoolByIP(first)
	if rp == nil {
		return fmt.Errorf("Cannot find any external pool by: %s", first)
	}

	return o.RemoveExternalRoutePool(rp)
}

func (o *OSPF) AdvertiseExternalPrefixes(first string) error {
	rp := o.getExternalRoutePoolByIP(first)
	if rp == nil {
		return fmt.Errorf("Cannot advertise prefixes: %s with cannot find the pool", first)
	}

	err := o.ConnectObjects(rp.Handler)
	if err != nil {
		return fmt.Errorf("Cannot advertise prefixes: %s with %s", first, err)
	}

	return nil
}

func (o *OSPF) WithdrawExternalPrefixes(first string) error {
	rp := o.getExternalRoutePoolByIP(first)
	if rp == nil {
		return fmt.Errorf("Cannot withdraw prefixes: %s with cannot find the pool", first)
	}

	//I temparary use this method.*/
	err := o.DisconnectObjects(rp.Handler)
	if err != nil {
		return fmt.Errorf("Cannot withdraw prefixes: %s with %s", first, err)
	}

	return nil
}

func (o *OSPF) AdvertiseAllExternalPrefixes() error {
	if o.ExternalRoutePools == nil || len(o.ExternalRoutePools) == 0 {
		return ErrorNoExternalRoutePoolExist
	}

	for ip, _ := range o.ExternalRoutePools {
		err := o.AdvertiseExternalPrefixes(ip)
		if err != nil {
			return fmt.Errorf("Cannot withdraw prefix %s with: %s", ip, err)
		}
	}
	return nil
}

func (o *OSPF) WithdrawAllExternalPrefixes() error {
	if o.ExternalRoutePools == nil || len(o.ExternalRoutePools) == 0 {
		return ErrorNoExternalRoutePoolExist
	}

	for ip, _ := range o.ExternalRoutePools {
		err := o.WithdrawExternalPrefixes(ip)
		if err != nil {
			return fmt.Errorf("Cannot withdraw prefix %s with: %s", ip, err)
		}
	}

	return nil
}

func (o *OSPF) AddExternalRoutePool() (*RoutePool, error) {
	cmd := fmt.Sprintf("AgtOspfTopology Add %s AGT_OSPF_EXTERNAL_ROUTEPOOL", o.LSDB.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf external route pool to %s with: %s", o.Handler, err)
	}

	if o.ExternalRoutePools == nil || len(o.ExternalRoutePools) == 0 {
		o.ExternalRoutePools = make(map[string]*RoutePool, 1)
	}

	isDefault := o.IsDefaultExternalRoutePoolExist()

	pool := &RoutePool{
		Object:  "AgtOspfExternalRoutePool",
		Type:    External,
		Handler: strings.TrimSpace(res),
		OSPF:    o,
		Default: !isDefault,
	}

	err = o.ConnectObjects(pool.Handler)
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf external route pool with: %s", err)
	}

	err = pool.Sync()
	if err != nil {
		return nil, fmt.Errorf("Cannot add ospf external route pool to %s with: %s", o.Handler, err)
	}

	o.ExternalRoutePools[pool.First] = pool

	return pool, nil
}

func (o *OSPF) isObjectAlreadyConnected(handler string) bool {
	handler = strings.TrimSpace(handler)

	conns, err := o.GetObjectConnections()
	if err != nil {
		return false
	}

	for _, conn := range conns {
		if conn == handler {
			return true
		}
	}

	return false
}

func (o *OSPF) GetObjectConnections() ([]string, error) {
	cmd := fmt.Sprintf("AgtOspfRouter ListObjectConnections %s", o.RouterHandler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get ospf object connections to %s with: %s", o.RouterHandler, err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	conns := make([]string, 0, 2)
	fields := strings.Split(res, " ")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		conns = append(conns, field)
	}

	if len(conns) == 0 {
		return nil, ErrorNoObjectConnections
	}

	return conns, nil
}

func (o *OSPF) RemoveExternalRoutePool(pool *RoutePool) error {
	if pool == nil || pool.Handler == "" || pool.Type != External {
		return fmt.Errorf("Inavlid route pool to delete")
	}

	return o.RemoveExternalRoutePoolByID(pool.Handler)
}

func (o *OSPF) RemoveExternalRoutePoolByID(handler string) error {
	cmd := fmt.Sprintf("AgtOspfTopology Remove %s ", handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot delete ospf external route pool  %s with: %s", handler, err)
	}

	pool := o.getExternalRoutePoolByID(handler)
	if pool != nil {
		delete(o.ExternalRoutePools, pool.First)
	}

	return nil
}

func (o *OSPF) getExternalRoutePoolByID(handler string) *RoutePool {
	handler = strings.TrimSpace(handler)
	for _, pool := range o.ExternalRoutePools {
		if pool.Handler == handler {
			return pool
		}
	}

	return nil
}

func (o *OSPF) getExternalRoutePoolByIP(ip string) *RoutePool {
	ip = strings.TrimSpace(ip)
	for first, pool := range o.ExternalRoutePools {
		if first == ip {
			return pool
		}
	}

	return nil
}

func (o *OSPF) RemoveExternalRoutePoolByIP(ip string) error {
	if o.ExternalRoutePools == nil || len(o.ExternalRoutePools) == 0 {
		return fmt.Errorf("Cannot remove external route pool: %s with: no pool exist")
	}

	ip = strings.TrimSpace(ip)
	var p *RoutePool
	for first, pool := range o.ExternalRoutePools {
		if first == ip {
			p = pool
		}
	}

	err := o.RemoveExternalRoutePool(p)
	if err != nil {
		return fmt.Errorf("Cannot delete ospf external route pool  %s with: %s", p.Handler, err)
	}

	return nil
}

func (o *OSPF) GetAllExternalRoutePools() ([]*RoutePool, error) {
	cmd := fmt.Sprintf("AgtOspfTopology ListHandlesByType %s AGT_OSPF_EXTERNAL_ROUTEPOOL", o.LSDB.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get all ospf external route pool of:%s with: %s", o.Handler, err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	//Clear all cached entry.
	o.ExternalRoutePools = make(map[string]*RoutePool, 1)

	fields := strings.Split(res, " ")
	pools := make([]*RoutePool, 0, 1)

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		np := &RoutePool{
			Handler: field,
			Type:    External,
			Object:  "AgtOspfExternalRoutePool",
			OSPF:    o,
		}

		err = np.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot get all external route pool with: %s", err)
		}

		pools = append(pools, np)
	}

	if len(pools) == 0 {
		return nil, ErrorNoExternalRoutePoolExist
	}

	for i, pool := range pools {
		if i == 0 {
			pool.Default = true
		}

		o.ExternalRoutePools[pool.First] = pool
	}

	return pools, nil
}

func (o *OSPF) RemoveAllExternalRoutePools() error {
	pools, err := o.GetAllExternalRoutePools()
	if err != nil {
		if err == ErrorNoExternalRoutePoolExist {
			return nil
		}

		return fmt.Errorf("Cannot remove all external route pool with: %s", err)
	}

	for _, pool := range pools {
		if err = o.RemoveExternalRoutePool(pool); err != nil {
			return fmt.Errorf("Cannot remove external route pool %s with: %s", pool.Handler, err)
		}

		delete(o.ExternalRoutePools, pool.First)
	}

	return nil
}

func (o *OSPF) RemoveAllPools() error {
	//AgtOspfTopology ListHandles
	cmd := fmt.Sprintf("AgtOspfTopology ListHandles %s", o.LSDB.Handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get all route pool of:%s with: %s", o.Handler, err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		pt, err := o.GetPoolType(field)
		if err != nil {
			return fmt.Errorf("Cannot remove pool: %s", field)
		}

		if field == o.RouterHandler {
			continue
		}

		if pt == "AGT_OSPF_NEIGHBOR" {
			continue
		}

		err = o.RemovePoolByID(field)
		if err != nil {
			return fmt.Errorf("Cannot remove pool: %s", field)
		}
	}

	o.ExternalRoutePools = nil
	o.SummaryRoutePools = nil
	o.Networks = nil
	o.Routers = nil

	return nil
}

func (o *OSPF) GetPoolType(handler string) (string, error) {
	cmd := fmt.Sprintf("AgtOspfTopology GetType %s ", handler)
	res, err := o.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot delete  route pool  %s with: %s", handler, err)
	}

	res = strings.Replace(res, "\"", "", -1)

	return strings.TrimSpace(res), nil
}

func (o *OSPF) RemovePoolByID(handler string) error {
	cmd := fmt.Sprintf("AgtOspfTopology Remove %s ", handler)
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot delete  route pool  %s with: %s", handler, err)
	}

	return nil
}

func (o *OSPF) ConnectObjects(handler string) error {
	cmd := fmt.Sprintf("AgtOspfRouter ConnectObjects %s %s", o.RouterHandler, strings.TrimSpace(handler))
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot connect %s to  %s with: %s", handler, o.RouterHandler, err)
	}

	return nil
}

func (o *OSPF) DisconnectObjects(handler string) error {
	cmd := fmt.Sprintf("AgtOspfRouter DisconnectObjects %s %s", o.RouterHandler, strings.TrimSpace(handler))
	_, err := o.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot disconnect %s to  %s with: %s", handler, o.RouterHandler, err)
	}

	return nil
}
