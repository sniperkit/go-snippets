package main

import (
	//"bufio"
	"fmt"
	//"os"
	"regexp"
	"strings"
	"telnet"
	//"time"
)

//invoke AgtDevice ListAvailableEmulations

//Legacy address configuration
//invoke AgtEthernetAddresses  ListSutIpAddresses 2
//invoke AgtEthernetAddresses IsVlanSupported 1
//invoke AgtEthernetAddresses IsMacAddressAccessSupported 2
//invoke AgtEthernetAddresses AddSutIpAddress 1 100.1.1.1/24
//invoke AgtEthernetAddresses AddSutIpAddresses 2 1000 123.1.1.1 32 1
//invoke AgtEthernetAddresses GetSutMacAddress 2 123.1.8.199
//invoke AgtEthernetAddresses SetSutMacAddress 2 123.1.8.199 00:20:11:22:33:44
//invoke AgtEthernetAddresses ListVlans 2
//invoke AgtAddressesEthernet GetLocalMacAddress 2 1
//invoke AgtEthernetIpv6Addresses2 ListSutIpv6AddressesByVlan 1 0

//invoke  AgtOpticalInterface ListAvailablePortModes 2
//invoke  AgtOpticalInterface GetPortMode 2
//invoke AgtOpticalInterface IsPluginMediaSupported 2
//invoke AgtOpticalInterface IsPluginMediaInserted 2
//invoke AgtOpticalInterface GetPluginMediaType 2
//invoke AgtOpticalInterface ListAvailableMediaTypes 2
//invoke AgtOpticalInterface SetMediaType 2 AGT_MEDIA_SFP
//invoke AgtOpticalInterface GetMediaType 1
//invoke AgtOpticalInterface IsMediaTypeSelectionSupported 2 AGT_MEDIA_SFP
//invoke AgtOpticalInterface AllLasersOn
//invoke  AgtOpticalInterface AllLasersOff
//invoke AgtOpticalInterface LaserOn 2
//invoke AgtOpticalInterface LaserOff 2

//invoke IntLiveCapture Initialize

//invoke AgtIpStatus GetStatus 2
//invoke AgtIpStatus GetStatusSummary

//Use this command to get the device handles under a port
//invoke AgtDeviceList ListHandles 2
//invoke AgtDeviceList Add 2
//invoke AgtDeviceList Remove 2
//invoke AgtDeviceList GetName 2
//invoke  AgtDeviceList SetName 2 T1
//invoke AgtDeviceList GetLinkDefaults
//invoke AgtProfileList ListHandles
//invoke AgtProfileList ListNames

//invoke  AgtAddressPool SetNetworkLayer 2  AGT_NETWORK_LAYER_IPV4_IPV6
//Set device ip address
//invoke  AgtAddressesIp SetLocalIpAddress 2 AGT_ADDRESS_FAMILY_IPV4 10.1.1.111
//invoke AgtAddressesIp GetLocalIpAddress 2 AGT_ADDRESS_FAMILY_IPV4
//Set device ipv6 address
//invoke AgtAddressesIp SetLocalIpAddress 2 AGT_ADDRESS_FAMILY_IPV6 2001:db8::1000
//invoke AgtAddressesIp GetLocalIpAddress 2 AGT_ADDRESS_FAMILY_IPV6

//Add OSPF emulation
//invoke AgtDevice AddEmulation 2 AGT_EMULATION_OSPFV2_ROUTER

//BasiC ospf
//This command will return the session ID
//invoke AgtTestTopology ListSessions 2
//invoke AgtTestTopology ListSutIpAddresses 2
//invoke AgtTestTopology GetSutIpAddress 2
//invoke AgtTestTopology ListSutIpv6Addresses 2
//invoke AgtTestTopology ListTesterIpv6Addresses 2
//invoke AgtTestTopology GetSutIpv6Address 2
//invoke AgtTestTopology GetTesterIpv6Address 2
//invoke AgtTestTopology GetSessionInterface 2
//invoke AgtTestTopology AddSession 2 AGT_SESSION_OSPF
//invoke AgtTestTopology GetSessionType 3
//invoke AgtTestTopology EnableSession 3
//invoke AgtTestTopology DisableSession 3
//invoke AgtTestTopology RemoveSession 3
//invoke AgtTestTopology SetSessionName 3 Test
//invoke AgtTestTopology ListSessionEmulations 3
//invoke AgtTestTopology GetEmulationState 4
//invoke AgtTestTopology ListEmulationCustomStates ospfv2Router
//invoke AgtOspfSession SetInterfaceAddress $SESSION_ID 111.1.1.101 30
//invoke AgtOspfSession SetAreaId 3 0.0.0.1
//invoke AgtOspfSession GetLastError 3
//invoke AgtOspfSession GetNeighborState 3
//invoke AgtOspfSession GetInterfaceState 3
//invoke AgtOspfSession GetStateSummary
//invoke AgtOspfSession GetOspfv3StateSummary
//invoke AgtOspfSession GetOspfv2StateSummary
//invoke AgtOspfSession ListNeighbors 3
//invoke AgtOspfSession RemoveAllNeighbors 3
//invoke AgtOspfSession GetSessionRouter 3
//here we use the router id from last step to set the tester router ID
//invoke AgtOspfRouter SetRouterId 7 1.2.3.4
//invoke AgtOspfRouter GetRouterId 7
//invoke AgtOspfSession AddNeighbor 3
//invoke AgtOspfRouter IsTeEnabled 7
//invoke AgtOspfRouter GetRouterLsa 7
//invoke AgtOspfRouter ListObjectConnections 7
//invoke AgtOspfRouter Advertise 7
//invoke AgtOspfRouter GetVersion 7
//invoke AgtOspfRouter GetLinkDetails 7 9
//invoke AgtOspfNeighbor SetIpAddress 9 111.1.1.102
//invoke AgtOspfNeighbor SetRouterId 9 4.3.2.1
//invoke AgtOspfRouter GetRouterDetails 7

//This command return the database handler
//invoke AgtOspfSession GetLsaDatabase 3
//This command returns the lsa handler
//invoke AgtOspfLsaDatabase ListLsas 1
//invoke AgtOspfLsaDatabase GetLsaType 1
//invoke AgtOspfRouterLsa GetLinkDetails 1 0
//invoke AgtOspfRouterLsa GetMaximumLinks
//invoke AgtOspfRouterLsa GetRouterType 1

//invoke AgtDeviceList ListHandles 2
//invoke AgtDevice ListEmulations 4
//invoke AgtEmulation Enable 4 ospfv2Router
//Start the routing engine
//invoke AgtRoutingEngine Start
//invoke AgtRoutingEngine Stop
//invoke AgtRoutingEngine GetState

//Statistics
//invoke AgtStatistics ListAvailableStatistics

//Add traffic group
//This Command add traffic group in constant profile
//invoke AgtStreamGroupList AddStreamGroups 2 AGT_PACKET_STREAM_GROUP 4
//invoke AgtStreamGroupList ListHandles
//invoke AgtStreamGroupList ListName
//invoke AgtConstantProfile ListStreamGroups 2
//Traffic profile management
//invoke AgtProfileList ListHandles
//invoke AgtProfileList AddProfile 2 AGT_CONSTANT_PROFILE
//invoke AgtProfileList AddProfile 2 AGT_BURST_PROFILE
//invoke AgtProfileList AddProfile 2 AGT_CUSTOM_PROFILE
//invoke AgtProfileList ListNames
//We get all traffic profile type by this command.
//We need create profile and stream group, and then add the stream group into profile.
//invoke AgtProfileList ListTypes
//invoke AgtProfileList ListProfilesOnPort 2
//invoke AgtCustomProfile GetType 6
//invoke AgtCustomProfile GetSourcePort 6
//invoke AgtCustomProfile ListStreamGroups 6
//invoke AgtCustomProfile AddStreamGroups 6 17

//Use this command to start sending traffic and statistics tracking
//invoke AgtTestController StartTest

var ResultR = regexp.MustCompile("{(?P<result>[[:alnum:][:space:]]*)}")
var BasicResultR = regexp.MustCompile("(?P<status>[[:alnum:]-]+)[[:space:]]+(?P<result>[[:alnum:][:space:]_-{}*]*)")

var Session *telnet.Session
var MGRSession *telnet.Session

var APIs = make(map[string][]string, 10)
var MGRAPIs = make(map[string][]string, 10)

func GetAllObjects() error {
	res, err := Invoke("AgtHelp ListObjects")
	if err != nil {
		return err
	}

	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		res = matches[1]
	}

	objects := strings.Split(res, " ")
	for _, obj := range objects {
		if _, ok := APIs[obj]; !ok {
			APIs[obj] = nil
		} else {
			fmt.Printf("Duplicate object: %s found!\n", obj)
		}
	}

	return nil
}

func GetMethodParameters(object, method string) error {
	cmd := fmt.Sprintf("AgtHelp ListMethodParameters %s %s", object, method)
	res, err := Invoke(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("%20s %30s:   %q\n", object, method, strings.TrimSpace(res))

	return nil
}

func GetAllMethods() error {
	err := GetAllObjects()
	if err != nil {
		return err
	}

	for obj, _ := range APIs {
		cmd := fmt.Sprintf("AgtHelp ListMethods %s", obj)
		res, err := Invoke(cmd)
		if err != nil {
			APIs[obj] = nil
		}

		matches := ResultR.FindStringSubmatch(res)
		if len(matches) == 2 {
			res = matches[1]
		}

		fields := strings.Split(res, " ")
		methods := make([]string, 0, 10)
		for _, field := range fields {
			methods = append(methods, field)
		}

		APIs[obj] = methods
	}

	return nil
}

func MGRGetAllObjects() error {
	res, err := MGRInvoke("AgtHelp ListObjects")
	if err != nil {
		return err
	}

	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		res = matches[1]
	}

	objects := strings.Split(res, " ")
	for _, obj := range objects {
		if _, ok := MGRAPIs[obj]; !ok {
			MGRAPIs[obj] = nil
		} else {
			fmt.Printf("Duplicate object: %s found!\n", obj)
		}
	}

	return nil
}

func MGRGetAllMethods() error {
	err := MGRGetAllObjects()
	if err != nil {
		return err
	}

	for obj, _ := range MGRAPIs {
		cmd := fmt.Sprintf("AgtHelp ListMethods %s", obj)
		res, err := MGRInvoke(cmd)
		if err != nil {
			MGRAPIs[obj] = nil
		}

		matches := ResultR.FindStringSubmatch(res)
		if len(matches) == 2 {
			res = matches[1]
		}

		fields := strings.Split(res, " ")
		methods := make([]string, 0, 10)
		for _, field := range fields {
			methods = append(methods, field)
		}

		MGRAPIs[obj] = methods
	}

	return nil
}

func MGRGetMethodParameters(object, method string) error {
	cmd := fmt.Sprintf("AgtHelp ListMethodParameters %s %s", object, method)
	fmt.Println(cmd)
	res, err := MGRInvoke(cmd)
	if err != nil {
		return err
	}

	fmt.Printf("%20s %30s:   %q\n", object, method, strings.TrimSpace(res))

	return nil
}

func GetAllOpenSessions() ([]string, error) {
	res, err := Invoke("AgtSessionManager ListOpenSessions")
	if err != nil {
		return nil, err
	}

	var sessions = make([]string, 0, 10)
	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		sess := strings.Split(matches[1], " ")
		for _, s := range sess {
			sessions = append(sessions, s)
		}
	}

	return sessions, nil
}

func ConnectToSession(sess string) (*telnet.Session, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionPort %s", sess)
	res, err := Invoke(cmd)
	if err != nil {
		return nil, err
	}

	fmt.Println(res)
	res = strings.TrimSpace(res)

	mgraddr := fmt.Sprintf("%s:%s", "10.71.20.231", res)
	mgrsess, err := telnet.New3(mgraddr)
	if err != nil {
		return nil, err
	}

	return mgrsess, nil
}

func CloseSession(sess string) error {
	cmd := fmt.Sprintf("AgtSessionManager CloseSession %s", sess)
	_, err := Invoke(cmd)
	if err != nil {
		return err
	}

	return nil
}

func OpenNewSession() (string, error) {
	//AgtOpenSession RouterTester900
	cmd := fmt.Sprintf("AgtSessionManager OpenSession RouterTester900 AGT_SESSION_ONLINE")
	res, err := Invoke(cmd)
	if err != nil {
		return "", err
	}

	return res, nil
}

func main() {
	err := GetAllMethods()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%q\n", APIs)

	for obj, methods := range APIs {
		for _, method := range methods {
			GetMethodParameters(obj, method)
		}
	}

	id, err := OpenNewSession()
	if err != nil {
		panic(err)
	}

	sess, err := GetAllOpenSessions()
	if err != nil {
		panic(err)
	}

	for _, s := range sess {
		fmt.Println(s)
	}
	MGRSession, err = ConnectToSession(id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%q\n", sess)

	res, err := MGRReservePorts("101/1")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	CloseSession(id)
}

func Run(method ...string) (string, error) {
	if len(method) == 0 {
		return "", fmt.Errorf("You must give the method to run")
	}

	var cmd string
	for _, p := range method {
		cmd += fmt.Sprintf(" %s", p)
	}
	_, err := Session.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	return GetCommandResult()
}

func GetCommandResult() (string, error) {
	var line []byte

	b, err := Session.ReadByte()
	if err != nil {
		return "", fmt.Errorf("Cannot get command result: ", err.Error())
	}
	if b == 'i' {
		line, err = Session.ReadUntilSkip([]string{"\""}, []string{"name"})
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else {
		sline, err := Session.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, sline))
	}

	return string(line), nil
}

func MGRGetCommandResult() (string, error) {
	var line []byte

	b, err := MGRSession.ReadByte()
	if err != nil {
		return "", fmt.Errorf("Cannot get command result: ", err.Error())
	}
	if b == 'i' {
		line, err = MGRSession.ReadUntilSkip([]string{"\""}, []string{"name"})
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else {
		sline, err := MGRSession.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, sline))
	}

	return string(line), nil
}

func isPortValid(port string) bool {
	port = strings.TrimSpace(port)
	if port != "103/1" && port != "103/2" && port != "103/3" && port != "103/4" &&
		port != "101/1" && port != "101/2" && port != "101/3" && port != "101/4" {
		return false
	}

	return true
}

func isPortsValid(ports ...string) bool {
	for _, port := range ports {
		if !isPortValid(port) {
			return false
		}
	}

	return true
}

func MGRReservePorts(ports ...string) (string, error) {
	if len(ports) == 0 || !isPortsValid(ports...) {
		return "", fmt.Errorf("You must give the port set to reserver")
	}

	cmd := fmt.Sprintf("AgtPortSelector AddPorts")

	for _, port := range ports {
		cmd += fmt.Sprintf(" %s", port)
	}

	res, err := MGRInvoke(cmd)
	if err != nil {
		return "", err
	}

	return res, nil
}

func MGRInvoke(method ...string) (string, error) {
	cmd := fmt.Sprintf("%s ", "invoke")
	for _, p := range method {
		cmd += fmt.Sprintf(" %s", p)
	}
	_, err := MGRSession.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	line, err := MGRGetCommandResult()
	if err != nil {
		return "", fmt.Errorf("Cannot get result of: %s with error: %s", cmd, err.Error())
	}

	res := BasicResultR.FindStringSubmatch(line)
	if len(res) != 3 {
		return "", fmt.Errorf("Run %s with invalid result: %s", cmd, line)
	}

	if res[1] == "0" {
		return res[2], nil
	}

	return "", fmt.Errorf("Run %s with result: (%s: %s)", cmd, res[1], res[2])
}

func Invoke(method ...string) (string, error) {
	cmd := fmt.Sprintf("%s ", "invoke")
	for _, p := range method {
		cmd += fmt.Sprintf(" %s", p)
	}
	_, err := Session.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	line, err := GetCommandResult()
	if err != nil {
		return "", fmt.Errorf("Cannot get result of: %s with error: %s", cmd, err.Error())
	}

	res := BasicResultR.FindStringSubmatch(line)
	if len(res) != 3 {
		return "", fmt.Errorf("Run %s with invalid result: %s", cmd, line)
	}

	if res[1] == "0" {
		return res[2], nil
	}

	return "", fmt.Errorf("Run %s with result: (%s: %s)", cmd, res[1], res[2])
}

func init() {
	sess, err := telnet.New3("10.71.20.231:9001")
	if err != nil {
		panic(err)
	}

	Session = sess

	/*
		go func(r *telnet.Session) {
			for {
				line, err := r.ReadString('\n')
				if err != nil {
					panic(err)
				}
				fmt.Printf("\r> %s", line)
				fmt.Printf("\r>")
			}
		}(sess)

		go func(r *telnet.Session) {
			fmt.Println(">")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				if scanner.Text() == "quit" {
					os.Exit(0)
				}

				line := scanner.Text()
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "run ") {
					line = strings.TrimPrefix(line, "run ")
					if strings.TrimSpace(line) != "" {
						_, err := r.WriteLine(line)
						if err != nil {
							panic(err)
						}
					}
				}
				fmt.Printf(">")
			}
		}(sess)
	*/
}
