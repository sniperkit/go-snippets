package n2x

import (
	"errors"
	"fmt"
	"strings"
)

//invoke AgtTestTopology AddSession 1 AGT_SESSION_BGP4_I_BGP_IPV6
//invoke AgtTestTopology AddSession 1 AGT_SESSION_BGP4_I_BGP
//invoke AgtTestTopology SetSessionPoolSize 8 200
//invoke AgtBgp4PeerPool Disable 8
//invoke AgtBgp4Ipv4PeerPool SetSutIpAddressIncrementingRange 8 123.1.1.1 24 1 1
//invoke AgtBgp4Ipv4PeerPool SetTesterIpAddressIncrementingRange 8 123.1.1.20  24 1 1
//invoke AgtBgp4Ipv6PeerPool SetSutIpAddressIncrementingRange 8 2001:db8::1 128 1 1
//invoke AgtBgp4Ipv6PeerPool SetTesterIpAddressIncrementingRange 8 2001:db8::1 128 1 1
//invoke AgtBgp4PeerPool SetSutAsNumberIncrementingRange 8 200 1 1
//invoke AgtBgp4PeerPool SetTesterAsNumberIncrementingRange 8 200 1 1
//invoke AgtBgp4PeerPool SetOpenParameter 8  AGT_BGP4_HOLDTIME 1234
//invoke AgtBgp4PeerPool SetOpenParameter 8 AGT_BGP4_KEEPALIVE_TIMER 12345
//invoke AgtBgp4PeerPool SetRoutesPerUpdate 8 1234
//invoke AgtBgp4PeerPool SetInterUpdateDelay 8 300
//invoke AgtBgp4PeerPool AddRouteProfile 8 AGT_BGP4_ROUTE_PROFILE_IPV4
//invoke AgtBgp4Ipv4RouteProfile SetRoutesIncrementingRange 215 100.0.0.1 32 1000 0
//invoke AgtBgp4Ipv6RouteProfile SetRoutesIncrementingRange 415 2001:1234:: 64 2 0
//invoke AgtBgp4RouteProfile SetFlapDistribution 415 100 0
//invoke AgtBgp4RouteProfile Withdraw 415
//invoke AgtBgp4RouteProfile Advertise 415
//invoke AgtBgp4FlapProfile StartFlap  8
//invoke AgtBgp4FlapProfile SetRoutesPerUpdate 8 100
//invoke AgtBgp4FlapProfile SetInterUpdateDelay  8 300
//invoke AgtBgp4FlapProfile SetAdvertiseToWithdrawDelay 8 50000
//invoke AgtBgp4FlapProfile SetWithdrawToAdvertiseDelay 8 100000
type BGP struct {
	Name            string
	Type            int
	Handler         string
	FlapHandler     string
	Object          string
	PeerPoolObject  string
	PeerPoolHandler string
	ASN             string
	RASN            string
	ASNType         int
	HoldTime        string
	KeepAliveTime   string
	Routes          map[string]*RouteProfile
	State           bool
	SutIP           string
	IP              string
	*Port
}

const (
	BGPASN2BYTE = iota
	BGPASN4BYTE = iota
)

const (
	IPV4_EBGP = iota
	IPV4_IBGP
	IPV6_EBGP
	IPV6_IBGP
)

var ErrorNoBGPRouteProfileExist = errors.New("No bgp route profile exist")

var BGPTypeNameMap = map[int]string{
	IPV4_EBGP: "AGT_SESSION_BGP4_E_BGP",
	IPV4_IBGP: "AGT_SESSION_BGP4_I_BGP",
	IPV6_EBGP: "AGT_SESSION_BGP4_E_BGP_IPV6",
	IPV6_IBGP: "AGT_SESSION_BGP4_I_BGP_IPV6",
}

var BGPSessionTypeTypeMap = map[string]int{
	"AGT_SESSION_BGP4_E_BGP":      IPV4_EBGP,
	"AGT_SESSION_BGP4_I_BGP":      IPV4_IBGP,
	"AGT_SESSION_BGP4_E_BGP_IPV6": IPV6_EBGP,
	"AGT_SESSION_BGP4_I_BGP_IPV6": IPV6_IBGP,
}

var BGPTypePeerPoolMap = map[int]string{
	IPV4_EBGP: "AgtBgp4PeerPool",
	IPV4_IBGP: "AgtBgp4PeerPool",
	IPV6_EBGP: "AgtBgp4Ipv6PeerPool",
	IPV6_IBGP: "AgtBgp4Ipv6PeerPool",
}

var BGPSessionTypePeerPoolMap = map[string]string{
	"AGT_SESSION_BGP4_E_BGP":      "AgtBgp4PeerPool",
	"AGT_SESSION_BGP4_I_BGP":      "AgtBgp4PeerPool",
	"AGT_SESSION_BGP4_E_BGP_IPV6": "AgtBgp4Ipv6PeerPool",
	"AGT_SESSION_BGP4_I_BGP_IPV6": "AgtBgp4Ipv6PeerPool",
}

var BGPTypeSessionMap = map[int]string{
	IPV4_EBGP: "AgtBgp4Session",
	IPV4_IBGP: "AgtBgp4Session",
	IPV6_EBGP: "AgtBgp4Session",
	IPV6_IBGP: "AgtBgp4Session",
}

var BGPTypeRouteProfileMap = map[int]string{
	IPV4_EBGP: "AgtBgp4Ipv4RouteProfile",
	IPV4_IBGP: "AgtBgp4Ipv4RouteProfile",
	IPV6_EBGP: "AgtBgp4Ipv6RouteProfile",
	IPV6_IBGP: "AgtBgp4Ipv6RouteProfile",
}

var BGPTypeRouteProfileTypeMap = map[int]string{
	IPV4_EBGP: "AGT_BGP4_ROUTE_PROFILE_IPV4",
	IPV4_IBGP: "AGT_BGP4_ROUTE_PROFILE_IPV4",
	IPV6_EBGP: "AGT_BGP4_ROUTE_PROFILE_IPV6",
	IPV6_IBGP: "AGT_BGP4_ROUTE_PROFILE_IPV6",
}

var BGPASNTypeNameMap = map[int]string{
	BGPASN2BYTE: "AGT_BGP4_2_BYTE_AS_NUMBER",
	BGPASN4BYTE: "AGT_BGP4_4_BYTE_AS_NUMBER",
}

var BGPASNNameTypeMap = map[string]int{
	"AGT_BGP4_2_BYTE_AS_NUMBER": BGPASN2BYTE,
	"AGT_BGP4_4_BYTE_AS_NUMBER": BGPASN4BYTE,
}

type RouteProfile struct {
	Handler string
	Object  string
	*BGP
	First string
	Plen  string
	Count string
	Step  string
}

func (rp *RouteProfile) Sync() error {
	return nil
}

func (rp *RouteProfile) GetRoutes() error {
	res, err := rp.GetRoutesIncrementingRange()
	if err != nil {
		return fmt.Errorf("Cannot get routes of %s with %s", rp.Handler, err)
	}

	fields := strings.Split(res, " ")
	if len(fields) != 4 {
		return fmt.Errorf("Invalid return value when get bgp routes")
	}

	rp.First = strings.TrimSpace(fields[0])
	rp.Plen = strings.TrimSpace(fields[1])

	res, err = rp.GetRouteCount()
	if err != nil {
		return fmt.Errorf("Cannot get bgp routes with: %s", err)
	}

	rp.Count = res

	return nil
}

func (rp *RouteProfile) GetRouteCount() (string, error) {
	//AgtBgp4Ipv4RouteProfile GetTotalRoutesInProfile
	cmd := fmt.Sprintf("%s GetTotalRoutesInProfile %s", rp.Object, rp.Handler)
	res, err := rp.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get bgp routes count with: %s", err)
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return "", fmt.Errorf("Cannot get bgp route count with: %s", err)
	}

	rp.Count = res

	return res, nil
}

func (rp *RouteProfile) SetName(name string) error {
	//AgtBgp4Ipv4PeerPool SetRouteProfileName
	cmd := fmt.Sprintf("%s SetRouteProfileName %s %s %s", rp.BGP.PeerPoolObject, rp.BGP.PeerPoolHandler, rp.Handler, name)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set bgp route profile name with: %s", err)
	}

	return nil
}

func (rp *RouteProfile) SetRoutesPerPeer(count string) error {
	cmd := fmt.Sprintf("%s SetRoutesPerPeer %s %s", rp.Object, rp.Handler, count)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set bgp routes per peer with: %s", err)
	}

	return nil
}

func (rp *RouteProfile) GetRoutesPerPeer() (string, error) {
	cmd := fmt.Sprintf("%s GetRoutesPerPeer %s", rp.Object, rp.Handler)
	res, err := rp.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get bgp routes per peer with: %s", err)
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return "", fmt.Errorf("Cannot get bgp route per peer with: %s", err)
	}

	return res, nil
}

func (rp *RouteProfile) SetRoutesPrefixLength(plen string) error {
	cmd := fmt.Sprintf("%s SetRoutesPrefixLength %s %s", rp.Object, rp.Handler, plen)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set bgp routes prefix length with: %s", err)
	}

	return nil
}

func (rp *RouteProfile) GetRoutesPrefixLength() (string, error) {
	cmd := fmt.Sprintf("%s GetRoutesPrefixLength %s", rp.Object, rp.Handler)
	res, err := rp.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get bgp routes prefix length with: %s", err)
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return "", fmt.Errorf("Cannot get bgp route prefix with: %s", err)
	}

	rp.Plen = res

	return res, nil
}

func (rp *RouteProfile) SetRoutesIncrementingRange(first, plen string) error {
	cmd := fmt.Sprintf("%s SetRoutesIncrementingRange %s %s %s 1 0", rp.Object, rp.Handler, first, plen)
	_, err := rp.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set bgp routes with: %s", err)
	}

	rp.First = first
	rp.Plen = plen

	return nil
}

func (rp *RouteProfile) GetRoutesIncrementingRange() (string, error) {
	cmd := fmt.Sprintf("%s GetRoutesIncrementingRange %s", rp.Object, rp.Handler)
	res, err := rp.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get bgp routes with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	if res == "" {

		return "", fmt.Errorf("Invalid return value when get increment range")
	}

	fields := strings.Split(res, " ")
	if len(fields) != 4 {
		return "", fmt.Errorf("Invalid return value when get bgp routes")
	}

	rp.First = strings.TrimSpace(fields[0])
	rp.Plen = strings.TrimSpace(fields[1])

	return res, nil
}

func (b *BGP) Enable() error {
	cmd := fmt.Sprintf("%s Enable %s", b.Object, b.Handler)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable bgp: %s with %s", b.Handler, err)
	}

	return nil
}

func (b *BGP) Disable() error {
	cmd := fmt.Sprintf("%s Disable %s", b.Object, b.Handler)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot enable bgp: %s with %s", b.Handler, err)
	}

	return nil
}

func (b *BGP) IsEnabled() bool {
	cmd := fmt.Sprintf("%s GetEnableFlag %s", b.Object, b.Handler)
	res, err := b.Invoke(cmd)
	if err != nil {
		return false
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return false
	}

	if res == "1" {
		return true
	}

	return false
}

func (b *BGP) Sync() error {
	b.State = b.IsEnabled()

	_, err := b.GetSutAsNumber()
	if err != nil {
		return fmt.Errorf("Cannot get sut asnumber %s", err)
	}

	_, err = b.GetTesterAsNumber()
	if err != nil {
		return fmt.Errorf("Cannot get sut asnumber %s", err)
	}

	_, err = b.GetTesterAsNumberType()
	if err != nil {
		return fmt.Errorf("Cannot get sut asnumber %s", err)
	}

	_, err = b.GetSutIpAddress()
	if err != nil {
		return fmt.Errorf("Cannot get sut ip address %s", err)
	}

	_, err = b.GetTesterIpAddress()
	if err != nil {
		return fmt.Errorf("Cannot get tester ip address %s", err)
	}

	_, err = b.GetAllRouteProfiles()
	if err != nil {
		if err == ErrorNoBGPRouteProfileExist {
			return nil
		}
		return fmt.Errorf("Cannot get all route profile with: %s", err)
	}

	return nil
}

func (b *BGP) GetSutAsNumber() (string, error) {
	//AgtBgp4PeerPool GetSutAsNumber
	cmd := fmt.Sprintf("%s GetSutAsNumber %s", b.Object, b.Handler)
	res, err := b.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get as number with: %s", err)
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return "", fmt.Errorf("Get sut as number failed with: empty return value")
	}

	b.RASN = res

	return res, nil
}

func (b *BGP) SetSutAsNumber(asn string) error {
	asn = strings.TrimSpace(asn)
	//AgtBgp4PeerPool SetSutAsNumber
	cmd := fmt.Sprintf("%s SetSutAsNumber %s %s", b.Object, b.Handler, asn)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get as number with: %s", err)
	}

	b.RASN = asn

	return nil
}

func (b *BGP) GetTesterAsNumber() (string, error) {
	cmd := fmt.Sprintf("%s GetTesterAsNumber %s", b.Object, b.Handler)
	res, err := b.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get tester as number with: %s", err)
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return "", fmt.Errorf("Get tester as number failed with: empty return value")
	}

	b.ASN = res

	return res, nil
}

func (b *BGP) SetTesterAsNumber(asn string) error {
	asn = strings.TrimSpace(asn)
	//AgtBgp4PeerPool SetSutAsNumber
	cmd := fmt.Sprintf("%s SetTesterAsNumber %s %s", b.Object, b.Handler, asn)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot get as number with: %s", err)
	}

	b.ASN = asn

	return nil
}

func (b *BGP) GetTesterAsNumberType() (int, error) {
	cmd := fmt.Sprintf("%s GetTesterAsNumberType %s", b.Object, b.Handler)
	res, err := b.Invoke(cmd)
	if err != nil {
		return -1, fmt.Errorf("Cannot get tester as number with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)
	if res == "" {
		return -1, fmt.Errorf("Get tester as number failed with: empty return value")
	}

	typ, ok := BGPASNNameTypeMap[res]
	if !ok {
		return -1, fmt.Errorf("Invalid asn name %s", res)
	}

	b.ASNType = typ

	return typ, nil
}

func (b *BGP) SetTesterAsNumberType(typ int) error {
	name, ok := BGPASNTypeNameMap[typ]
	if !ok {
		return fmt.Errorf("Invalid asn type: %d", typ)
	}

	cmd := fmt.Sprintf("%s SetTesterAsNumberType %s %s", b.Object, b.Handler, name)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester as number type with: %s", err)
	}

	b.ASNType = typ

	return nil
}

func (b *BGP) GetSutIpAddress() (string, error) {
	//AgtBgp4Session GetSutIpAddress
	cmd := fmt.Sprintf("%s GetSutIpAddress %s", b.Object, b.Handler)
	res, err := b.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get sut ip address with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	b.SutIP = res

	return b.SutIP, nil
}

//N2x will auto get the interface address, so this API is used only when necessary.
func (b *BGP) SetSutIpAddress(ip string) error {
	//AgtBgp4Session GetSutIpAddress
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return fmt.Errorf("Invalid ip address")
	}

	cmd := fmt.Sprintf("%s SetSutIpAddress %s %s", b.Object, b.Handler, ip)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set sut ip address with: %s", err)
	}

	b.SutIP = ip

	return nil
}

func (b *BGP) GetTesterIpAddress() (string, error) {
	cmd := fmt.Sprintf("%s GetTesterIpAddress %s", b.Object, b.Handler)
	res, err := b.Invoke(cmd)
	if err != nil {
		return "", fmt.Errorf("Cannot get tester ip address with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)

	b.IP = res

	return b.IP, nil
}

//N2x will auto get the interface address, so this API is used only when necessary.
func (b *BGP) SetTesterIpAddress(ip string) error {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return fmt.Errorf("Invalid ip address")
	}

	cmd := fmt.Sprintf("%s SetTesterIpAddress %s %s", b.Object, b.Handler, ip)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot set tester ip address with: %s", err)
	}

	b.IP = ip

	return nil
}

func (b *BGP) GetAllRouteProfiles() ([]*RouteProfile, error) {
	cmd := fmt.Sprintf("%s ListRouteProfiles %s", b.PeerPoolObject, b.PeerPoolHandler)

	res, err := b.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot get route pools  with: %s", err)
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	res = strings.Replace(res, "\"", "", -1)
	res = strings.TrimSpace(res)
	if res == "" {
		return nil, ErrorNoBGPRouteProfileExist
	}

	pools := make([]*RouteProfile, 0, 1)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		np := &RouteProfile{
			Handler: field,
			BGP:     b,
			Object:  BGPTypeRouteProfileMap[b.Type],
		}

		err = np.Sync()
		if err != nil {
			return nil, fmt.Errorf("Cannot sync route pool %s with: %s", np.Handler, err)
		}

		pools = append(pools, np)
	}

	if len(pools) == 0 {
		return nil, ErrorNoBGPRouteProfileExist
	}

	//Directly drop previous entries.
	b.Routes = make(map[string]*RouteProfile, len(pools))
	for _, pool := range pools {
		b.Routes[pool.First] = pool
	}

	return pools, nil
}

//AgtBgp4Ipv4PeerPool AddRouteProfile
func (b *BGP) AddRouteProfile() (*RouteProfile, error) {
	cmd := fmt.Sprintf("%s AddRouteProfile %s %s", b.PeerPoolObject, b.PeerPoolHandler, BGPTypeRouteProfileTypeMap[b.Type])
	res, err := b.Invoke(cmd)
	if err != nil {
		return nil, fmt.Errorf("Cannot add route profile of type %s with: %s", BGPTypeRouteProfileTypeMap[b.Type], err)
	}

	res = strings.TrimSpace(res)
	if res == "" {
		return nil, fmt.Errorf("Cannot add route profile of type %s with: %s", BGPTypeRouteProfileTypeMap[b.Type], err)
	}

	rp := &RouteProfile{Handler: res, Object: BGPTypeRouteProfileMap[b.Type], BGP: b}

	err = rp.Sync()
	if err != nil {
		return nil, fmt.Errorf("Cannot add route profile of type %s with: %s", BGPTypeRouteProfileTypeMap[b.Type], err)
	}

	if b.Routes == nil {
		b.Routes = make(map[string]*RouteProfile, 1)
	}

	b.Routes[rp.First] = rp

	return rp, nil
}

func (b *BGP) RemoveRouteProfile(rp *RouteProfile) error {
	if rp.Handler == "" {
		return fmt.Errorf("Cannot remove invalid route profile")
	}

	cmd := fmt.Sprintf("%s RemoveRouteProfile %s", b.PeerPoolHandler, rp.Handler)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot delete route profile %s with %s", rp.Handler, err)
	}

	if orp, ok := b.Routes[rp.First]; ok {
		delete(b.Routes, orp.First)
	}

	return nil
}

func (b *BGP) isRouteProfileExist(handler string) bool {
	if b.Routes == nil || len(b.Routes) == 0 {
		return false
	}

	for _, rp := range b.Routes {
		if rp.Handler == handler {
			return true
		}
	}

	return false
}

func (b *BGP) GetRouteProfileByID(handler string) (*RouteProfile, error) {
	if b.Routes == nil || len(b.Routes) == 0 {
		return nil, ErrorNoBGPRouteProfileExist
	}

	for _, rp := range b.Routes {
		if rp.Handler == handler {
			return rp, nil
		}
	}

	return nil, fmt.Errorf("Cannot find any route profile with id: %s", handler)
}

func (b *BGP) RemoveRouteProfileByID(handler string) error {
	rp, err := b.GetRouteProfileByID(handler)
	if err != nil {
		return fmt.Errorf("Cannot remove route profile %s with: %s", handler, err)
	}

	return b.RemoveRouteProfile(rp)
}

func (b *BGP) RemoveRouteProfileByIP(ip string) error {
	if b.Routes == nil || len(b.Routes) == 0 {
		return fmt.Errorf("Cannot remove route profile by ip %s with %s", ip, ErrorNoBGPRouteProfileExist)
	}

	var del *RouteProfile

	for first, rp := range b.Routes {
		if first == ip {
			del = rp
			break
		}
	}

	if del == nil {
		return fmt.Errorf("Cannot remove route profile by ip %s with cannot find the profiel", ip)
	}

	return b.RemoveRouteProfile(del)
}

func (b *BGP) RemoveAllRouteProfiles() error {
	cmd := fmt.Sprintf("%s ClearRouteProfiles", b.PeerPoolHandler)
	_, err := b.Invoke(cmd)
	if err != nil {
		return fmt.Errorf("Cannot remove all route profile  with: %s", err)
	}

	b.Routes = nil

	return nil
}

func (b *BGP) AddRoutes(first, plen, count, step string) (*RouteProfile, error) {
	if _, ok := b.Routes[first]; ok {
		return nil, fmt.Errorf("Same routes already exist")
	}

	rp, err := b.AddRouteProfile()
	if err != nil {
		return nil, fmt.Errorf("Cannot add routes with %s", err)
	}

	err = rp.SetRoutesPrefixLength(plen)
	if err != nil {
		return nil, fmt.Errorf("Cannot set prefix length with: %s", err)
	}

	err = rp.SetRoutesPerPeer(count)
	if err != nil {
		return nil, fmt.Errorf("Cannot set routes count: %s", err)
	}

	err = rp.SetRoutesIncrementingRange(first, plen)
	if err != nil {
		return nil, fmt.Errorf("Cannot set routes with: %s", err)
	}

	err = rp.SetName(first)
	if err != nil {
		return nil, fmt.Errorf("Cannot set routes with: %s", err)
	}

	err = b.Sync()
	if err != nil {
		return nil, fmt.Errorf("Cannot sync bgp configurations")
	}

	return rp, err
}

func (b *BGP) RemoveRoutes(first string) error {
	return b.RemoveRouteProfileByIP(first)
}

func (b *BGP) GetTester4ByteAsNumberFormat() (string, error) {
	return "", nil
}

func (b *BGP) SetTester4ByteAsNumberFormat(typ string) error {
	return nil
}

func (b *BGP) GetTester4ByteAsNumber() (string, error) {
	return "", nil
}

func (b *BGP) SetTester4ByteAsNumber(asn string) error {
	return nil
}

func (b *BGP) GetOpenParameter(param string) (string, error) {
	return "", nil
}

func (b *BGP) SetOpenParameter(param, value string) error {
	return nil

}
