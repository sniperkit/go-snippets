package n2x

import (
	"fmt"
	"strings"
	"telnet"
	"util"
)

type NSession struct {
	ID    string
	Port  string
	Label string
	Pid   string
	Conn  *telnet.Session
	Ports map[string]*Port
}

func (ns *NSession) ReservePorts(ports ...string) error {
	ns.Ports = make(map[string]*Port, 10)
	if len(ports) == 0 || !ns.isPortsValid(ports...) {
		return fmt.Errorf("You must give the port set to reserver")
	}

	for _, port := range ports {
		cmd := fmt.Sprintf("AgtPortSelector AddPorts")

		cmd += fmt.Sprintf(" %s", port)

		res, err := ns.Invoke(cmd)
		if err != nil {
			return err
		}
		ns.Ports[port] = &Port{Name: port, Handler: res, NSession: ns}
	}

	return nil
}

func (ns *NSession) isPortValid(port string) bool {
	port = strings.TrimSpace(port)
	if port != "103/1" && port != "103/2" && port != "103/3" && port != "103/4" &&
		port != "101/1" && port != "101/2" && port != "101/3" && port != "101/4" {
		return false
	}

	return true
}

func (ns *NSession) isPortsValid(ports ...string) bool {
	for _, port := range ports {
		if !ns.isPortValid(port) {
			return false
		}
	}

	return true
}

func (ns *NSession) Invoke(cmds ...string) (string, error) {
	cmd := fmt.Sprintf("%s ", "invoke")
	for _, p := range cmds {
		cmd += fmt.Sprintf(" %s", p)
	}

	util.AppendToFile("n2x_log.txt", []byte(cmd+"\n"))

	_, err := ns.Conn.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	line, err := ns.GetCommandResult()
	if err != nil {
		return "", fmt.Errorf("Cannot get result of: %s with error: %s", cmd, err.Error())
	}

	util.AppendToFile("n2x_log.txt", []byte("Result: "+line))
	res := BasicResultR.FindStringSubmatch(line)
	if len(res) != 3 {
		return "", fmt.Errorf("Run %s with invalid result: %s", cmd, line)
	}

	if res[1] == "0" {
		return strings.TrimSpace(res[2]), nil
	}

	return "", fmt.Errorf("Run %s with result: (%s: %s)", cmd, res[1], res[2])

}

func (ns *NSession) GetCommandResult() (string, error) {
	var line []byte

	b, err := ns.Conn.ReadByte()
	if err != nil {
		return "", fmt.Errorf("Cannot get command result: ", err.Error())
	}
	if b == 'i' {
		line, err = ns.Conn.ReadUntilSkip([]string{"\""}, []string{"name"})
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else if b == 'm' {
		line, err = ns.Conn.ReadUntil("brace")
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else {
		sline, err := ns.Conn.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, sline))
	}

	return string(line), nil
}

func (ns *NSession) GetReservedPorts() ([]*Port, error) {
	ps, err := ns.ListPorts()
	if err != nil {
		return nil, err
	}

	ports := make([]*Port, 0, 10)

	for _, p := range ps {
		pname, err := ns.GetPortName(p)
		if err != nil {
			return nil, fmt.Errorf("Cannot get reserved port with: %s", err)
		}

		ports = append(ports, &Port{
			Name:     pname,
			Handler:  p,
			NSession: ns,
		})
	}

	if ns.Ports == nil {
		ns.Ports = make(map[string]*Port, 10)
	}

	for _, port := range ports {
		ns.Ports[port.Name] = port
	}

	return ports, nil
}

func (ns *NSession) GetPorts() ([]*Port, error) {
	ports := make([]*Port, 0, 10)
	if ns.Ports != nil {
		for _, port := range ns.Ports {
			ports = append(ports, port)
		}
		return ports, nil
	}

	ps, err := ns.ListPorts()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		pname, err := ns.GetPortName(p)
		if err != nil {
			return nil, fmt.Errorf("Cannot get reserved port with: %s", err)
		}

		ports = append(ports, &Port{
			Name:     pname,
			Handler:  p,
			NSession: ns,
		})
	}

	if ns.Ports == nil {
		ns.Ports = make(map[string]*Port, 10)
	}

	for _, port := range ports {
		ns.Ports[port.Name] = port
	}

	return ports, nil
}

func (ns *NSession) ListModules() ([]string, error) {
	res, err := ns.Invoke("AgtPortSelector ListModules")
	if err != nil {
		return nil, err
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	modules := make([]string, 0, 10)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		modules = append(modules, field)
	}

	return modules, nil
}

func (ns *NSession) ListAvailableModules() ([]string, error) {
	res, err := ns.Invoke("AgtPortSelector ListAvailableModules")
	if err != nil {
		return nil, err
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	modules := make([]string, 0, 10)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		modules = append(modules, field)
	}

	return modules, nil
}

func (ns *NSession) ListAvailablePorts() ([]string, error) {
	res, err := ns.Invoke("AgtPortSelector ListAvailablePorts")
	if err != nil {
		return nil, err
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	ports := make([]string, 0, 10)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		ports = append(ports, field)
	}

	return ports, nil
}

func (ns *NSession) ListLockedPorts() ([]string, error) {
	res, err := ns.Invoke("AgtPortSelector ListLockedPorts")
	if err != nil {
		return nil, err
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	ports := make([]string, 0, 10)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		ports = append(ports, field)
	}

	return ports, nil
}

func (ns *NSession) ListPorts() ([]string, error) {
	res, err := ns.Invoke("AgtPortSelector ListPorts")
	if err != nil {
		return nil, err
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)
	ports := make([]string, 0, 10)

	fields := strings.Split(res, " ")
	for _, field := range fields {
		ports = append(ports, field)
	}

	return ports, nil
}

func (ns *NSession) GetPortName(handler string) (string, error) {
	handler = strings.TrimSpace(handler)
	if handler == "" {
		return "", fmt.Errorf("You must give the port handler for get name")
	}

	cmd := fmt.Sprintf("AgtPortSelector GetPortDetails %s", handler)
	res, err := ns.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.Replace(res, "{", "", -1)
	res = strings.Replace(res, "}", "", -1)

	fields := strings.Split(res, " ")
	if len(fields) != 2 {
		return "", fmt.Errorf("Cannot get port detail with: %s", res)
	}

	return strings.TrimSpace(fields[0]) + "/" + strings.TrimSpace(fields[1]), nil
}

//Use this functio to start sending traffic
func (ns *NSession) StartTest() error {
	//invoke AgtTestController StartTest
	_, err := ns.Invoke("AgtTestController StartTest")
	if err != nil {
		return fmt.Errorf("Cannot start test with: %s", err)
	}

	return nil
}

//Use this functio to stop sending traffic
func (ns *NSession) StopTest() error {
	//invoke AgtTestController StartTest
	_, err := ns.Invoke("AgtTestController StopTest")
	if err != nil {
		return fmt.Errorf("Cannot stop test with: %s", err)
	}

	return nil
}

//Use this function to start routing engine
func (ns *NSession) StartRoutingEngine() error {
	_, err := ns.Invoke("AgtRoutingEngine Start")
	if err != nil {
		return fmt.Errorf("Cannot start routing engine with: %s", err)
	}

	return nil
}

//Use this function to stop routing engine
func (ns *NSession) StopRoutingEngine() error {
	_, err := ns.Invoke("AgtRoutingEngine Stop")
	if err != nil {
		return fmt.Errorf("Cannot stop routing engine with: %s", err)
	}

	return nil
}

func (ns *NSession) Sync() error {
	ports, err := ns.GetReservedPorts()
	if err != nil {
		return fmt.Errorf("Cannot sync session %s with %s", ns.ID, err)
	}

	for _, port := range ports {
		err = port.Sync()
		if err != nil {
			return fmt.Errorf("Cannot sync session %s with %s", ns.ID, err)
		}
	}

	return nil
}
