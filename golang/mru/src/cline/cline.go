package cline

import (
	"client"
	"command"
	"configuration"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"util"
)

var hostReg = regexp.MustCompile("^(?P<host>.*)#")

type Cli struct {
	client       client.Client
	conf         *configuration.Configuration
	currentMode  string
	modeToPrompt map[string]string
	promptToMode map[string]string
	logLock      sync.Mutex
}

func (c *Cli) RunCommand(cmd *command.Command) (result []byte, err error) {
	message := fmt.Sprintf("Run Command: %-40s cmode:%15s mode: %15s on %20s\n", cmd.CMD, cmd.Mode, c.CurrentMode(), c.conf.IP)
	c.Log(message)

	if cmd.Mode != c.currentMode {
		return nil, errors.New("Error: Command: " + cmd.CMD + " should be run under: " + cmd.Mode + "! But currently we are under: " + c.currentMode + " mode!")
	}

	cmd.CMD = strings.TrimSpace(cmd.CMD)
	if strings.HasPrefix(cmd.CMD, "bcm.user.proxy") {
		c.client.WriteLine(cmd.CMD) //For the stupid bcmshell
		cmd.End = "BCM.0>"
	} else if strings.HasPrefix(cmd.CMD, "diag") {
		c.client.WriteLine(cmd.CMD)
		cmd.End = "RTK.0>"
	} else if c.currentMode == "bcmshell" && cmd.CMD != "exit" && cmd.CMD != "quit" {
		c.client.WriteLine(cmd.CMD) //For the stupid bcmshell
		cmd.End = "BCM.0>"
	} else if c.currentMode == "rtkshell" && cmd.CMD != "exit" && cmd.CMD != "quit" {
		c.client.WriteLine(cmd.CMD)
		cmd.End = "RTK.0>"
	} else {
		c.client.WriteLine(cmd.CMD)
	}
	if cmd.End == "" {
		cmd.End = c.conf.Prompt
	}
	data, err := c.client.ReadUntil(cmd.End)
	if err != nil {
		fmt.Println(fmt.Sprintf("Connection to %s is broken\n", c.conf.IP))
		return nil, err
	}

	for {
		if strings.Contains(string(data), c.conf.Hostname) {
			break
		}

		pending, err := c.client.ReadUntil(cmd.End)
		if err != nil {
			fmt.Println(fmt.Sprintf("Connection to %s is broken\n", c.conf.IP))
			return nil, err
		}

		data = append(data, pending...)
	}

	util.AppendToFile("command_log.txt", []byte(fmt.Sprintf("Command: %s, Result: %s\n", cmd.CMD, string(data))))

	if c.IsErrorExist(string(data)) {
		return nil, errors.New("Cannot run command: " + cmd.CMD + " with error: <<<" + string(data) + ">>>")
	}

	old := c.currentMode
	rs := strings.Split(string(data), "\n")
	//log.Println(len(rs))
	//log.Println(c.promptToMode)
	for p, m := range c.promptToMode {
		//log.Println(p, m, rs[len(rs)-1])
		//if strings.Contains(rs[len(rs)-1], p) && m != old {
		if strings.HasPrefix(strings.TrimSpace(rs[len(rs)-1]), p) && m != old {
			c.currentMode = m
		}
	}

	if c.IsModeSwitchMustBeOccured(cmd) && old == c.currentMode {
		return nil, fmt.Errorf("Mode change must be accured after run command: %s, but there is no mode change. Result: %s", cmd.CMD, string(data))
	}

	if old != c.currentMode {
		//log.Println("After run: ", cmd.CMD, " mode switch from: ", old, " to: ", c.currentMode, "!")
		message = fmt.Sprintf("After run: %40s mode switch from : %15s to %15s. !\n", cmd.CMD, old, c.currentMode)
		c.Log(message)
	}

	return data, nil
}

func (c *Cli) WriteLine(line string) (int, error) {
	return c.client.WriteLine(line)
}

func (c *Cli) Expect(delims ...string) ([]byte, error) {
	return c.client.ReadUntil(delims...)
}

func (c *Cli) IsModeSwitchMustBeOccured(cmd *command.Command) bool {
	cmdstr := strings.TrimSpace(cmd.CMD)
	if strings.HasPrefix(cmdstr, "interface ") || strings.HasPrefix(cmdstr, "router ") ||
		strings.HasPrefix(cmdstr, "configure terminal") || strings.HasPrefix(cmdstr, "vlan database") ||
		strings.HasPrefix(cmdstr, "ip dhcp pool") || strings.HasPrefix(cmdstr, "route-map") ||
		strings.HasPrefix(cmdstr, "ip access-list") || strings.HasPrefix(cmdstr, "flow ") ||
		strings.HasPrefix(cmdstr, "policer ") || strings.HasPrefix(cmdstr, "policy") ||
		strings.HasPrefix(cmdstr, "do q sh") || strings.HasPrefix(cmdstr, "do quote sh") ||
		strings.HasPrefix(cmdstr, "address-family") || strings.HasPrefix(cmdstr, "q sh") ||
		strings.HasPrefix(cmdstr, "quote sh") || strings.HasPrefix(cmdstr, "exit") ||
		strings.HasPrefix(cmdstr, "quit") {
		return true
	}

	return false
}

func (c *Cli) GoNormalMode() ([]byte, error) {
	if c.currentMode == "config" ||
		c.currentMode == "config-vlan" ||
		c.currentMode == "config-if" ||
		c.currentMode == "config-dhcp" ||
		c.currentMode == "config-router" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "exit", End: "#"})
		if err != nil {
			return res, err
		}
	} else if c.currentMode == "shell" ||
		c.currentMode == "bcmshell" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "exit", End: "#"})
		if err != nil {
			return res, err
		}
		c.GoNormalMode()
	}

	return nil, nil
}

func (c *Cli) GoShellMode() ([]byte, error) {
	if c.currentMode == "config" ||
		c.currentMode == "config-vlan" ||
		c.currentMode == "config-if" ||
		c.currentMode == "config-dhcp" ||
		c.currentMode == "bridge" ||
		c.currentMode == "config-router" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "do q sh -l", End: "#"})
		if err != nil {
			return res, err
		}
	} else if c.currentMode == "normal" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "q sh -l", End: "#"})
		if err != nil {
			return res, err
		}
	} else if c.currentMode == "bcmshell" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "exit", End: "#"})
		if err != nil {
			return res, err
		}
	} else if c.currentMode == "rtkshell" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "exit", End: "#"})
		if err != nil {
			return res, err
		}
	}
	return nil, nil
}

func (c *Cli) GoBCMShelllMode() ([]byte, error) {
	if c.currentMode == "config" ||
		c.currentMode == "config-vlan" ||
		c.currentMode == "config-if" ||
		c.currentMode == "config-dhcp" ||
		c.currentMode == "bridge" ||
		c.currentMode == "config-router" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "do q sh -l", End: "#"})
		if err != nil {
			return res, err
		}
		c.GoBCMShelllMode()
	} else if c.currentMode == "normal" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "q sh -l", End: "#"})
		if err != nil {
			return res, err
		}
		c.GoBCMShelllMode()
	} else if c.currentMode == "shell" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "bcm.user.proxy", End: "BCM.0>"})
		if err != nil {
			return res, err
		}
	}
	return nil, nil
}

func (c *Cli) GoConfigMode() ([]byte, error) {
	if c.currentMode == "config-vlan" ||
		c.currentMode == "config-if" ||
		c.currentMode == "config-dhcp" ||
		c.currentMode == "config-flow" ||
		c.currentMode == "config-policer" ||
		c.currentMode == "config-policy" ||
		c.currentMode == "bridge" ||
		c.currentMode == "config-router" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "exit", End: "#"})
		if err != nil {
			return res, err
		}
	} else if c.currentMode == "shell" ||
		c.currentMode == "bcmshell" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "exit", End: "#"})
		if err != nil {
			return res, err
		}
		c.GoConfigMode()
	} else if c.currentMode == "normal" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "configure terminal", End: "#"})
		if err != nil {
			return res, err
		}
	} else if c.currentMode == "enable" {
		res, err := c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "enable", End: "#"})
		if err != nil {
			return res, err
		}
		res, err = c.RunCommand(&command.Command{Mode: c.CurrentMode(), CMD: "configure terminal", End: "#"})
		if err != nil {
			return res, err
		}
	}
	return nil, nil
}

func NewCli(conf *configuration.Configuration) (c *Cli, err error) {
	tc, err := client.New(conf.Username, conf.Password, conf.Protocol, conf.IP, conf.Port)
	if err != nil {
		return nil, fmt.Errorf("\n\r\tCannot connect to host: %s with error: \n\r\t\t%s!", conf.IP, err.Error())
	}

	os.Remove("command_log.txt")

	return &Cli{
		client:       tc,
		conf:         conf,
		modeToPrompt: make(map[string]string, 1),
		promptToMode: make(map[string]string, 1),
		logLock:      sync.Mutex{},
	}, nil
}

func (c *Cli) CurrentMode() string {
	return c.currentMode
}

func (c *Cli) SetModeDB(db map[string]string) {
	c.ClearModeDB()
	for mode, prompt := range db {
		err := c.AddMode(mode, prompt)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (c *Cli) Log(message string) {
	c.logLock.Lock()
	defer c.logLock.Unlock()

	logfile := c.conf.SessionID
	if logfile == "" {
		logfile = "default"
	}

	file, err := os.OpenFile("asset/log/"+logfile+"_full.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Println("cannot Open file: ", logfile+"_full.log", " ", err.Error())
		return
	}
	defer file.Close()

	file.WriteString(message)
}

func (c *Cli) ClearModeDB() {
	c.promptToMode = make(map[string]string, 1)
	c.modeToPrompt = make(map[string]string, 1)
}

func (c *Cli) Init() error {
	for mode, prompt := range c.conf.ModeDB {
		err := c.AddMode(mode, prompt)
		if err != nil {
			fmt.Println(err.Error())
			return errors.New("Register mode db error!")
		}
	}

	//ToDo: Need to handle the enable password enabled case.
	c.client.WriteLine("enable")
	data, err := c.client.ReadUntil(c.conf.Prompt)
	if err != nil {
		return fmt.Errorf("Cannot go to enable mode with error: %s", err.Error())
	}

	var rHostname string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		matches := hostReg.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}
		rHostname = matches[1]
	}

	rHostname = strings.Trim(rHostname, "["+c.conf.SFU+"]")

	if rHostname == "" || rHostname != c.conf.Hostname {
		return fmt.Errorf("Invalid hostname: %s, must be: %s", c.conf.Hostname, rHostname)
	}

	c.client.WriteLine("terminal length 0")
	_, err = c.client.ReadUntil("#")
	if err != nil {
		return fmt.Errorf("Cannot set terminal length: %s", err.Error())
	}
	c.currentMode = "normal"

	return nil
}

func (c *Cli) AddMode(mode, prompt string) error {
	if mode == "" || prompt == "" {
		return errors.New("Invalid input mode: " + mode + " prompt: " + prompt + "!")
	}

	if _, ok := c.modeToPrompt[mode]; ok {
		return errors.New("Same mode: " + mode + " already exist!")
	}

	if _, ok := c.promptToMode[prompt]; ok {
		return errors.New("Same prompt: " + prompt + " already exist!")
	}

	c.modeToPrompt[mode] = prompt
	c.promptToMode[prompt] = mode

	return nil
}

func (c *Cli) GetPromptByMode(mode string) (string, error) {
	if prompt, ok := c.modeToPrompt[mode]; ok {
		return prompt, nil
	}

	return "", errors.New("Mode: " + mode + " does not exist!")
}

func (c *Cli) GetModeByPrompt(prompt string) (string, error) {
	if mode, ok := c.modeToPrompt[prompt]; ok {
		return mode, nil
	}

	return "", errors.New("Prompt: " + prompt + " does not exist!")
}

func (c *Cli) IsErrorExist(result string) bool {
	if strings.Contains(result, "Invalid input detected at") {
		return true
	}
	return false
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
