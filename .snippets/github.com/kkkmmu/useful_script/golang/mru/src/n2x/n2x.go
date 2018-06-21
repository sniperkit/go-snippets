package n2x

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"telnet"
	"util"
)

var DEFAULTSESSIONNAME = "ATSN2X"

var ResultR = regexp.MustCompile("{(?P<result>[[:alnum:][:space:]]*)}")
var BasicResultR = regexp.MustCompile("(?P<status>[[:alnum:]-]+)[[:space:]]+(?P<result>[[:alnum:][:space:]_-{}*\".:]*)")

var ErrorNoOpenSession = errors.New("There is no open session")

type N2X struct {
	IP       string
	Port     string
	Conn     *telnet.Session
	APIs     map[string][]string
	Sessions []*NSession
	Ports    map[string]*Port
}

type APIs map[string][]string

func New(ip, port string) (*N2X, error) {
	n2x := &N2X{
		IP:    ip,
		Port:  port,
		APIs:  make(map[string][]string, 10),
		Ports: make(map[string]*Port, 10),
	}

	sess, err := telnet.New3(fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, err
	}

	n2x.Conn = sess

	os.Remove("n2x_log.txt")

	return n2x, nil
}

func (n *N2X) KillAllSessions() error {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		if err == ErrorNoOpenSession {
			return nil
		}

		return fmt.Errorf("Cannot kill all session: %s", err.Error())
	}

	for _, sess := range sesss {
		err := n.KillSession(sess.ID)
		if err != nil {
			return fmt.Errorf("Cannot session %s(%d) : %s", sess.Label, sess.ID, err.Error())
		}
	}

	return nil
}

func (n *N2X) KillSession(id string) error {
	cmd := fmt.Sprintf("AgtSessionManager CloseSession %s", id)
	_, err := n.Invoke(cmd)
	if err != nil {
		return err
	}

	return nil
}

func (n *N2X) KillSessionByName(name string) error {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		if err == ErrorNoOpenSession {
			return nil
		}
		return fmt.Errorf("Cannot kill all session: %s", err.Error())
	}

	for _, sess := range sesss {
		if sess.Label == name {
			err := n.KillSession(sess.ID)
			if err != nil {
				return fmt.Errorf("Cannot session %s(%d) : %s", sess.Label, sess.ID, err.Error())
			}
		}
	}

	return nil
}

func (n *N2X) OpenNewSession(ports ...string) (*NSession, error) {
	cmd := fmt.Sprintf("AgtSessionManager OpenSession RouterTester900 AGT_SESSION_ONLINE")
	res, err := n.Invoke(cmd)
	if err != nil {
		return nil, err
	}

	nsess := &NSession{ID: res, Ports: make(map[string]*Port, 10)}
	nport, err := n.GetSessionPort(nsess.ID)
	if err != nil {
		return nil, err
	}

	nsess.Port = nport

	nsess, err = n.ConnectToSession(nsess)
	if err != nil {
		return nil, err
	}

	err = n.SetSessionName(nsess, DEFAULTSESSIONNAME)
	if err != nil {
		return nil, err
	}

	err = nsess.ReservePorts(ports...)
	if err != nil {
		return nil, err
	}

	return nsess, nil
}

func (n *N2X) ConnectToSessionByID(id string) (*NSession, error) {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		return nil, fmt.Errorf("Cannot get all session: %s", err.Error())
	}

	for _, sess := range sesss {
		if sess.ID == id {
			sess, err = n.ConnectToSession(sess)
			if err != nil {
				return nil, fmt.Errorf("Cannot connect to session %s(%d) : %s", sess.Label, sess.ID, err.Error())
			}

			err = sess.Sync()
			if err != nil {
				return nil, err
			}

			return sess, nil
		}
	}
	return nil, fmt.Errorf("Cannot find session with id: %s", id)
}

func (n *N2X) GetSessionByID(id string) (*NSession, error) {
	return n.ConnectToSessionByID(id)
}

func (n *N2X) ConnectToSessionByName(name string) (*NSession, error) {
	sesss, err := n.GetAllOpenSessions()
	if err != nil {
		return nil, fmt.Errorf("Cannot get all session: %s", err.Error())
	}

	for _, sess := range sesss {
		if sess.Label == name {
			sess, err = n.ConnectToSessionByID(sess.ID)
			if err != nil {
				return nil, fmt.Errorf("Cannot connect to session with: %s", err.Error())
			}
			err = sess.Sync()
			if err != nil {
				return nil, err
			}

			return sess, nil
		}
	}

	return nil, fmt.Errorf("Cannot find session with name: %s", name)
}

func (n *N2X) GetSessionByName(name string) (*NSession, error) {
	return n.ConnectToSessionByName(name)
}

func (n *N2X) SetSessionName(sess *NSession, name string) error {
	cmd := fmt.Sprintf("AgtSessionManager SetSessionLabel %s %s", sess.ID, name)
	_, err := n.Invoke(cmd)
	if err != nil {
		return err
	}

	sess.Label = name

	return nil
}

func (n *N2X) GetAllOpenSessions() ([]*NSession, error) {
	res, err := n.Invoke("AgtSessionManager ListOpenSessions")
	if err != nil {
		return nil, err
	}

	var sessions = make([]*NSession, 0, 10)
	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		sess := strings.Split(matches[1], " ")
		for _, s := range sess {
			if strings.TrimSpace(s) == "" {
				continue
			}
			nsess := &NSession{ID: s, Ports: make(map[string]*Port, 10)}
			sessions = append(sessions, nsess)
		}
	}

	if len(sessions) == 0 {
		return nil, ErrorNoOpenSession
	}

	for _, sess := range sessions {
		port, err := n.GetSessionPort(sess.ID)
		if err != nil {
			return nil, fmt.Errorf("Cannot get session port of %s", sess.ID)
		}

		label, err := n.GetSessionLabel(sess.ID)
		if err != nil {
			return nil, fmt.Errorf("Cannot get session label of %s", sess.ID)
		}

		pid, err := n.GetSessionPid(sess.ID)
		if err != nil {
			return nil, fmt.Errorf("Cannot get session Pid of %s", sess.ID)
		}

		sess.Port = port
		sess.Label = strings.Trim(label, "\"")
		sess.Pid = pid
	}

	return sessions, nil
}

func (n *N2X) ConnectToSession(sess *NSession) (*NSession, error) {
	addr := fmt.Sprintf("%s:%s", n.IP, sess.Port)
	conn, err := telnet.New3(addr)
	if err != nil {
		return nil, err
	}

	sess.Conn = conn

	return sess, nil
}

func (n *N2X) GetSessionPort(id string) (string, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionPort %s", id)
	res, err := n.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.TrimSpace(res)

	return res, nil
}

func (n *N2X) GetSessionLabel(id string) (string, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionLabel %s", id)
	res, err := n.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.TrimSpace(res)

	return res, nil
}

func (n *N2X) GetSessionPid(id string) (string, error) {
	cmd := fmt.Sprintf("AgtSessionManager GetSessionPid %s", id)
	res, err := n.Invoke(cmd)
	if err != nil {
		return "", err
	}

	res = strings.TrimSpace(res)

	return res, nil
}

func (n *N2X) Invoke(cmds ...string) (string, error) {
	cmd := fmt.Sprintf("%s ", "invoke")
	for _, p := range cmds {
		cmd += fmt.Sprintf(" %s", p)
	}

	util.AppendToFile("n2x_log.txt", []byte(cmd+"\n"))

	_, err := n.Conn.WriteLine(cmd)
	if err != nil {
		return "", fmt.Errorf("Run %s with error: %s", cmd, err.Error())
	}

	line, err := n.GetCommandResult()
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

func (n *N2X) GetAllMethods() error {
	res, err := n.Invoke("AgtHelp ListObjects")
	if err != nil {
		return err
	}

	matches := ResultR.FindStringSubmatch(res)
	if len(matches) == 2 {
		res = matches[1]
	}

	objects := strings.Split(res, " ")
	for _, obj := range objects {
		if _, ok := n.APIs[obj]; !ok {
			n.APIs[obj] = nil
		} else {
			fmt.Printf("Duplicate object: %s found!\n", obj)
		}
	}

	for obj, _ := range n.APIs {
		cmd := fmt.Sprintf("AgtHelp ListMethods %s", obj)
		res, err := n.Invoke(cmd)
		if err != nil {
			n.APIs[obj] = nil
			return err
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

		n.APIs[obj] = methods
	}

	return nil
}

func (n *N2X) GetCommandResult() (string, error) {
	var line []byte

	b, err := n.Conn.ReadByte()
	if err != nil {
		return "", fmt.Errorf("Cannot get command result: ", err.Error())
	}
	if b == 'i' {
		line, err = n.Conn.ReadUntilSkip([]string{"\""}, []string{"name"})
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else if b == 'm' {
		line, err = n.Conn.ReadUntil("brace")
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, line))
	} else {
		sline, err := n.Conn.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("Cannot get result with error: %s", err.Error())
		}
		line = []byte(fmt.Sprintf("%c%s", b, sline))
	}

	return string(line), nil
}

func init() {

}
