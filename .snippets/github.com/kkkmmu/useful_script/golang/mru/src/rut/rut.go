package rut

import (
	"cline"
	"command"
	"configuration"
	"context"
	"errors"
	"fmt"
	"log"
	"logger"
	"path/filepath"
	"result"
	"script"
	"strconv"
	"strings"
	"time"
)

//RUT should be and interface
type RUT struct {
	Name       string //Name in test case
	cli        *cline.Cli
	Device     string //Device name
	Username   string
	Password   string
	Protocol   string
	IP         string
	Port       string
	BasePrompt string
	Hostname   string //hostName
	SessionID  string
	SFU        string
}

type Config struct {
	Index      int    `json:"index"`
	Device     string `json:"device"`
	Protocol   string `json:"protocol"`
	IP         string `json:"ip"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"passowrd"`
	BasePrompt string `json:"baseprompt"`
	Hostname   string `json:"hostname"`
	SessionID  string `json:"-"`
	SFU        string `json:"sfu"`
}

type DB struct {
	DB map[string]*RUT `json:"-"`
}

func NewDB() *DB {
	return &DB{
		DB: make(map[string]*RUT, 1),
	}
}

func (db *DB) AddDUT(name string, dut *RUT) {
	db.DB["DUT"+name] = dut
}

func (db *DB) GetRUTByName(name string) *RUT {
	if r, ok := db.DB[name]; ok {
		return r
	}
	return nil
}

func buildDefaultConfiguration(r *RUT) *configuration.Configuration {
	//log.Println(r.Hostname, r.Device)
	var conf configuration.Configuration
	conf.Name = r.Name
	conf.Username = r.Username
	conf.Password = r.Password
	conf.Device = r.Device
	conf.Hostname = r.Hostname
	conf.BasePrompt = r.BasePrompt
	conf.IP = r.IP
	conf.Port = r.Port
	if r.Protocol == "" {
		conf.Protocol = "telnet"
	} else {
		conf.Protocol = r.Protocol
	}

	conf.EnablePrompt = configuration.DefaultEnablePrompt
	conf.LoginPrompt = configuration.DefaultLoginPrompt
	conf.PasswordPrompt = configuration.DefaultPasswordPrompt
	conf.Prompt = configuration.PromptEnd
	conf.ModeDB = configuration.BuildModeDBFromHostNameAndBasePrompt(r.Hostname, r.BasePrompt)
	conf.SessionID = r.SessionID

	if r.SFU == "" {
		conf.SFU = configuration.DefaultSFU
	} else {
		conf.SFU = r.SFU
	}

	//log.Printf("%#v", conf)
	return &conf
}

func New(r *RUT) (*RUT, error) {
	return r, nil
}

func (d *RUT) Init() error {
	if d.Device == "V8" {
		d.BasePrompt = d.Hostname + "[" + d.SFU + "]"
	} else {
		d.BasePrompt = d.Hostname
	}

	conf := buildDefaultConfiguration(d)
	c, err := cline.NewCli(conf)
	if err != nil {
		return errors.New("Cannot create CLI instance: " + err.Error())
	}

	err = c.Init()
	if err != nil {
		return errors.New("Cannot init RUT with: " + err.Error())
	}

	d.cli = c
	return nil
}

func (d *RUT) GoInitMode() {
	d.cli.GoNormalMode()
}

func (d *RUT) SetModeDB(db map[string]string) {
	d.cli.SetModeDB(db)
}

func (d *RUT) RunCommand(ctx context.Context, cmd *command.Command) (string, error) {
	logger.Push(ctx, fmt.Sprintf("Run Command: %-40s cmode:%15s mode: %15s on %20s\n", cmd.CMD, cmd.Mode, d.cli.CurrentMode(), d.IP))
	return d.runCommand(cmd)
}

func (d *RUT) RunCommands(ctx context.Context, cmds []*command.Command) (string, error) {
	for _, c := range cmds {
		data, err := d.RunCommand(ctx, c)
		if err != nil {
			return data, err
		}
	}

	return "", nil
}

func (d *RUT) WriteLine(line string) (int, error) {
	return d.cli.WriteLine(line)
}

func (d *RUT) Expect(delims ...string) ([]byte, error) {
	return d.cli.Expect(delims...)
}

func (d RUT) runCommand(cmd *command.Command) (string, error) {
	if cmd.Delay != 0 {
		<-time.After(time.Second * time.Duration(cmd.Delay))
	}
	data, err := d.cli.RunCommand(cmd)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (r *RUT) RunScript(sc *script.Script) <-chan result.Result {
	log.Printf("Start Runing Script: %v", sc)
	res := make(chan result.Result)
	go func(chan<- result.Result) {
		for i := 0; i < sc.Count; i++ {
			for _, c := range sc.Commands {
				<-time.After(time.Second * time.Duration(c.Delay))
				log.Printf("Run command: %v", c)
				data, err := r.cli.RunCommand(c)
				log.Println(string(data), err)
				res <- result.Result{
					Command: c.CMD,
					Result:  string(data),
					Err:     err,
				}
			}
			<-time.After(time.Second * time.Duration(sc.Timer))
		}
		close(res)
	}(res)

	return res
}

func (d *RUT) CreateVlan(id int) error {
	return nil

}

func (d *RUT) DestroyVlan(id int) error {

	return nil
}

func (d *RUT) DestroyAllVlan() error {

	return nil
}

func (d *RUT) CreateVlanInterface(id int, ip string) error {

	return nil
}

func (d *RUT) DestroyVlanInterface(id int) error {

	return nil
}

func (d *RUT) AddIPAddress(ifname, ip string) error {

	return nil
}

func (d *RUT) DelIPAddress(ifname, ip string) error {

	return nil
}

func (d *RUT) AddSecondaryIPAddress(ifname, ip string) error {

	return nil
}

func (d *RUT) DelSecondaryIPAddress(ifname, ip string) error {

	return nil
}

func (d *RUT) DelAllIPAddress(ifname string) error {

	return nil
}

func (d *RUT) CreateOSPFInstance(id, tag string) error {

	return nil
}

func (d *RUT) DestroyOSPFInstance(id string) error {

	return nil
}

func isValidRUTConfig(c *Config) bool {
	if c.Index < 0 ||
		c.Device == "" ||
		c.IP == "" ||
		c.Port == "" ||
		c.Username == "" {
		return false
	}

	return true
}

func (d *RUT) IsAlive(ctx context.Context) bool {
	/*
		msg, err := d.cli.GoNormalMode()
		if err != nil {
			log.Println(err, msg)
			return false
		}
	*/

	res, err := d.RunCommand(ctx, &command.Command{
		Mode: d.cli.CurrentMode(),
		CMD:  "show running-config",
	})

	if err != nil {
		log.Println(err, res)
		return false
	}

	if strings.Contains(res, d.Hostname) {
		return true
	}
	return false
}

func (d RUT) GoShellMode() ([]byte, error) {
	return d.cli.GoShellMode()
}

func (d *RUT) FTPPut(local, ip, user, pass, dir string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("ftp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}
	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	ldir := filepath.Dir(local)
	file := filepath.Base(local)

	_, err := d.WriteLine("ftp " + ip)
	if err != nil {
		return err
	}

	data, err := d.Expect("):")
	//log.Println(string(data))
	_, err = d.WriteLine(user)
	if err != nil {
		return err
	}
	data, err = d.Expect("Password:")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("lcd " + ldir)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("cd " + dir)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("put " + file)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	_, err = d.WriteLine("exit")
	if err != nil {
		return err
	}

	data, err = d.Expect("#")
	if err != nil {
		return err
	}

	return nil
}

func (d *RUT) FTPGet(local, ip, user, pass, dir, file string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("ftp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}
	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	ldir := filepath.Dir(local)

	_, err := d.WriteLine("ftp " + ip)
	if err != nil {
		return err
	}

	data, err := d.Expect("):")
	//log.Println(string(data))
	_, err = d.WriteLine(user)
	if err != nil {
		return err
	}
	data, err = d.Expect("Password:")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("lcd " + ldir)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("cd " + dir)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	//log.Println(string(data))

	_, err = d.WriteLine("get " + file)
	if err != nil {
		return err
	}

	data, err = d.Expect("ftps>")
	if err != nil {
		return err
	}
	fmt.Println(string(data))

	_, err = d.WriteLine("exit")
	if err != nil {
		return err
	}

	data, err = d.Expect("#")
	if err != nil {
		return err
	}

	return nil
}

func (d *RUT) TCPDUMP(intf, filter, file, count string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("tcpdump should working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}

	if file == "" || !strings.HasSuffix(file, ".pcap") {
		return fmt.Errorf("Invalid file name: ", file)
	}

	dump := "tcpdump "
	if intf != "" {
		dump += " -i " + intf
	}

	if count != "" {
		dump += " -c " + count
	}

	if filter != "" {
		dump += " " + filter
	}

	dump += " -w " + file

	if count != "" {
		_, err := d.WriteLine(dump)
		if err != nil {
			return err
		}
		_, err = d.Expect("#")
		if err != nil {
			return err
		}
	} else {
		_, err := d.WriteLine(dump + "&")
		if err != nil {
			return err
		}
		_, err = d.Expect("#")
		if err != nil {
			return err
		}
		time.Sleep(time.Second * 10)
		_, err = d.WriteLine("killall tcpdump")
		if err != nil {
			return err
		}
		_, err = d.Expect("#")
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *RUT) SCPPut(local, ip, user, pass, dir string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("scp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}
	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	file := filepath.Base(local)

	_, err := d.WriteLine("scp " + local + " " + user + "@" + ip + ":" + dir + "/" + file)
	if err != nil {
		return err
	}

	known, err := d.Expect("(yes/no)?", "password:")
	fmt.Println(string(known))
	if strings.Contains(string(known), "(yes/no)?") {
		d.WriteLine("yes")
		_, err := d.Expect("password:")
		if err != nil {
			return err
		}
	}

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	data, err := d.Expect("#")
	if err != nil {
		return err
	}

	if !strings.Contains(string(data), "100%") {
		return fmt.Errorf("Upload file with error: %s", string(data))
	}
	fmt.Println(string(data))

	return nil
}

func (d *RUT) SCPGet(local, ip, user, pass, dir, file string) error {
	if d.cli.CurrentMode() != "shell" {
		_, err := d.GoShellMode()
		if err != nil {
			return fmt.Errorf("scp working under shell mode, current: %s", d.cli.CurrentMode())
		}
	}

	if !filepath.IsAbs(local) {
		return fmt.Errorf("local file must use absoluted path")
	}

	_, err := d.WriteLine("scp " + user + "@" + ip + ":" + dir + "/" + file + " " + local)
	if err != nil {
		return err
	}

	known, err := d.Expect("(yes/no)?", "password:")
	fmt.Println(string(known))
	if strings.Contains(string(known), "(yes/no)?") {
		d.WriteLine("yes")
		_, err := d.Expect("password:")
		if err != nil {
			return err
		}
	}

	_, err = d.WriteLine(pass)
	if err != nil {
		return err
	}

	data, err := d.Expect("#")
	if err != nil {
		return err
	}

	if !strings.Contains(string(data), "100%") {
		return fmt.Errorf("Upload file with error: %s", string(data))
	}
	fmt.Println(string(data))

	return nil
}

func GetRUTByConfig(c *Config) (*RUT, error) {
	if !isValidRUTConfig(c) {
		return nil, fmt.Errorf("Invalid config to create RUT: %v", c)
	}
	newrut := &RUT{
		Name:       "DUT" + strconv.Itoa(c.Index),
		Device:     c.Device,
		Username:   c.Username,
		Password:   c.Password,
		IP:         c.IP,
		Port:       c.Port,
		Hostname:   c.Hostname,
		BasePrompt: c.BasePrompt,
		SessionID:  c.SessionID,
		SFU:        c.SFU,
	}

	log.Printf("%#v", newrut)

	if err := newrut.Init(); err != nil {
		return nil, fmt.Errorf("Cannot create new DUT with config :%v. Error Message: %s", c, err.Error())
	}

	return newrut, nil
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}
