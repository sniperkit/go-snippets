package main

import (
	"command"
	"context"
	"flag"
	"fmt"
	/*
		"github.com/google/gopacket"
		"github.com/google/gopacket/pcap"
	*/
	"github.com/fatih/color"
	"io/ioutil"
	"net"
	"path/filepath"
	"regexp"
	"rut"
	"strconv"
	"strings"
	"sync"
	"telnet"
	"time"
	"util"
)

var (
/*
	handle *pcap.Handle
	err    error
*/
)

type bcmTable struct {
	Name   string
	Size   int64
	Fields []string
}

type Interface struct {
	Name            string
	VRF             string
	Type            string
	AdminStatus     string
	OperationStatus string
	Index           string
	MTU             string
	Bandwidth       string
	Mode            string
	Hardware        string
	HardwareAddress string
	Loopback        bool
	Multicast       bool
	Broadcast       bool
	IP              []string
	IP6             []string
	RxPPS           string
	RxBPS           string
	TxPPS           string
	TxBPS           string
	RxUcastPacket   uint64
	RxMcastPacket   uint64
	RxBcastPacket   uint64
	RxTotalPacket   uint64
	RxTotalByte     uint64
	TxUcastPacket   uint64
	TxMcastPacket   uint64
	TxBcastPacket   uint64
	TxTotalPacket   uint64
	TxTotalByte     uint64
	RxDiscard       uint64
	RxErrors        uint64
	TxDiscard       uint64
	TxErrors        uint64
}

func (ifp *Interface) String() string {
	var ifs string
	var Green = color.New(color.FgGreen)
	if ifp.AdminStatus == "UP" && ifp.OperationStatus == "UP" {
		ifs += Green.Sprintf("%12s: ", ifp.Name)
		ifs += Green.Sprintf("%10s ", ifp.VRF)
		ifs += Green.Sprintf("%6s ", ifp.Mode)
		ifs += Green.Sprintf("%8s ", ifp.Index)
		ifs += Green.Sprintf("%5s/%5s ", ifp.AdminStatus, ifp.OperationStatus)
		ifs += Green.Sprintf("%5s ", ifp.MTU)
		ifs += Green.Sprintf("%5s ", ifp.Bandwidth)
		//ifs += Green.Sprintf("%5s ", ifp.Hardware)
		ifs += Green.Sprintf("%15s ", ifp.HardwareAddress)
		if *Monitor == "status" {
			for _, ip := range ifp.IP {
				ifs += Green.Sprintf(" %18s ", ip)
			}

			for _, ip6 := range ifp.IP6 {
				ifs += Green.Sprintf(" %42s ", ip6)
			}
		} else if *Monitor == "statistics" {
			ifs += Green.Sprintf("RXPPs(%5s) ", ifp.RxPPS)
			ifs += Green.Sprintf("RXBPs(%5s) ", ifp.RxBPS)
			//ifs += Green.Sprintf("RXUcast(%8d) ", ifp.RxUcastPacket)
			ifs += Green.Sprintf("RXTotal(%5d) ", ifp.RxTotalPacket)
			if ifp.RxDiscard != 0 {
				ifs += Green.Sprintf("RXDiscard(%5d) ", ifp.RxDiscard)
			}
			if ifp.RxErrors != 0 {
				ifs += Green.Sprintf("RXErrors(%5d) ", ifp.RxErrors)
			}
			ifs += Green.Sprintf("TXPPs(%5s) ", ifp.TxPPS)
			ifs += Green.Sprintf("TXBPs(%5s) ", ifp.TxBPS)
			//ifs += Green.Sprintf("TXUcast(%8d) ", ifp.TxUcastPacket)
			ifs += Green.Sprintf("TXTotal(%5d) ", ifp.TxTotalPacket)
			if ifp.TxDiscard != 0 {
				ifs += Green.Sprintf("TXDiscard(%5d) ", ifp.TxDiscard)
			}

			if ifp.TxErrors != 0 {
				ifs += Green.Sprintf("TXErrors(%5d) ", ifp.TxErrors)
			}
		}
	} else {
		ifs += fmt.Sprintf("%12s: ", ifp.Name)
		ifs += fmt.Sprintf("%10s ", ifp.VRF)
		ifs += fmt.Sprintf("%6s ", ifp.Mode)
		ifs += fmt.Sprintf("%8s ", ifp.Index)
		ifs += fmt.Sprintf("%5s/%5s ", ifp.AdminStatus, ifp.OperationStatus)
		ifs += fmt.Sprintf("%5s ", ifp.MTU)
		ifs += fmt.Sprintf("%5s ", ifp.Bandwidth)
		//ifs += fmt.Sprintf("%5s ", ifp.Hardware)
		ifs += fmt.Sprintf("%15s ", ifp.HardwareAddress)
		if *Monitor == "status" {
			for _, ip := range ifp.IP {
				ifs += fmt.Sprintf(" %18s ", ip)
			}

			for _, ip6 := range ifp.IP6 {
				ifs += fmt.Sprintf(" %42s ", ip6)
			}
		}
	}

	return ifs
}

var InterfaceList = make([]*Interface, 0, 10)
var InterfaceLock sync.Mutex

//Entries: 32768 with indices
var tf = regexp.MustCompile("(?P<fname>[[:alnum:]_:]+)<")
var ts = regexp.MustCompile("Entries:[[:space:]]+(?P<size>[[:alnum:]]+)[[:space:]]+with indices")
var ifr = regexp.MustCompile("Interface[[:space:]]+(?P<ifname>[[:alnum:]/.]+)[[:space:]]+is")
var ifstatus = regexp.MustCompile("(?P<status><[[:alnum:],]+>)")
var ifhw = regexp.MustCompile("Hardware is (?P<hw>[[:alnum:]]+)[[:space:]]+Current HW addr: (?P<hwaddr>[[:alnum:].]+)")
var ipr = regexp.MustCompile("inet[[:space:]](?P<ip>[[:alnum:]./]+)[[:space:]]+")
var ip6r = regexp.MustCompile("inet[[:space:]](?P<ip>[[:alnum:]:/]+)[[:space:]]+")
var ifindexr = regexp.MustCompile("index[[:space:]]+(?P<index>[[:digit:]]+)[[:space:]]")
var ifmtur = regexp.MustCompile("mtu[[:space:]]+(?P<mtu>[[:digit:]]+)[[:space:]]")
var ifbandwidthr = regexp.MustCompile("Bandwidth[[:space:]]+(?P<bw>[[:alnum:]]+)[[:space:]]")
var ifvrfr = regexp.MustCompile("VRF Binding:[[:space:]]+Associated with[[:space:]]+(?P<vrf>[[:alnum:]_-]+)")
var ifrater = regexp.MustCompile("5 sec :[[:space:]]+(?P<tpps>[[:alnum:]]+)[[:space:]]+(?P<tbps>[[:alnum:]]+)[[:space:]]+(?P<rpps>[[:alnum:]]+)[[:space:]]+(?P<rbps>[[:alnum:]]+)")
var ifstatupr = regexp.MustCompile("IfInUcastPkts[[:space:]]+:[[:space:]]+(?P<rptk>[[:alnum:]]+)[[:space:]]+IfOutUcastPkts[[:space:]]+:[[:space:]]+(?P<tpkt>[[:alnum:]]+)")
var ifstatmpr = regexp.MustCompile("IfInMcastPkts[[:space:]]+:[[:space:]]+(?P<rptk>[[:alnum:]]+)[[:space:]]+IfOutMcastPkts[[:space:]]+:[[:space:]]+(?P<tpkt>[[:alnum:]]+)")
var ifstatbpr = regexp.MustCompile("IfInBcastPkts[[:space:]]+:[[:space:]]+(?P<rptk>[[:alnum:]]+)[[:space:]]+IfOutBcastPkts[[:space:]]+:[[:space:]]+(?P<tpkt>[[:alnum:]]+)")
var ifstatDiscardr = regexp.MustCompile("IfInDiscards[[:space:]]+:[[:space:]]+(?P<rptk>[[:alnum:]]+)[[:space:]]+IfOutDiscards[[:space:]]+:[[:space:]]+(?P<tpkt>[[:alnum:]]+)")
var ifstatErrorr = regexp.MustCompile("IfInErrors[[:space:]]+:[[:space:]]+(?P<rptk>[[:alnum:]]+)[[:space:]]+IfOutErrors[[:space:]]+:[[:space:]]+(?P<tpkt>[[:alnum:]]+)")

var UseLessFilter = regexp.MustCompile("unit[[:space:]]+[[:digit:]]+[[:space:]]+[a-zA-Z?]+[[:space:]]+registers")

var IP = flag.String("ip", "", "IP address of the remote device")
var Port = flag.String("port", "23", "Port to connect")
var Host = flag.String("hostname", "SWITCH", "Host name of the remote device")
var User = flag.String("username", "admin", "Username of the remote device")
var Password = flag.String("password", "", "Passwrod of the remote device")
var Table = flag.String("table", "", "Table name to dump")
var Register = flag.String("register", "", "Register name to dump")
var Command = flag.String("command", "", "Bcm shell command to run")
var Shell = flag.String("shell", "", "shell command to run")
var Config = flag.String("config", "", "Configure file name, your configuration should start from priviledge more")
var Protocol = flag.String("protocol", "telnet", "Login protocol(ssh/telnet)")
var Upload = flag.String("upload", "", "upload protocol (ftp/scp)")
var Download = flag.String("download", "", "download protocol (ftp/scp)")
var Server = flag.String("server", "", "server ip")
var ServerUser = flag.String("suser", "", "server user name")
var ServerPassword = flag.String("spass", "", "server password")
var Local = flag.String("local", "", "Local file name(Must use abslute path)")
var Remote = flag.String("remote", "", "Remote file name(Must use abslute path)")
var Tcpdump = flag.String("tcpdump", "", "Run tcpdump on remote device interface")
var TcpdumpCount = flag.String("count", "100", "Packet count to dump")
var TcpdumpFile = flag.String("dfile", "any.pcap", "File name of tcpdump")
var TcpdumpFilter = flag.String("filter", "", "tcpdump packet filter")
var PCAPDecode = flag.String("pcapdecode", "", "Decode a pcap file.")
var Check = flag.String("check", "", "Check related registers/tables, Use option to assign the dump format(raw/all/chg)")
var Set = flag.String("set", "", "Set related registers/tables of a given name")
var Name = flag.String("name", "", "Name of various things")
var Option = flag.String("option", "", "Options for another command")
var Value = flag.String("value", "", "values set to register/tables")
var Field = flag.String("field", "", "Field name of a table/register")
var Monitor = flag.String("monitor", "", "Monitor the cpu/memory/status/statitics usage of device")

func main() {
	flag.Parse()

	if *Table == "" && *Register == "" && *Command == "" && *Config == "" && *Shell == "" &&
		*Upload == "" && *Download == "" && *Tcpdump == "" && *PCAPDecode == "" && *Check == "" &&
		*Set == "" && *Monitor == "" {
		fmt.Println("You must give an operation to run(Dump table/Register, bcmshell command, shell command, config, upload/download, tcpdump")
		fmt.Println("\t Use -h to get help information")
		return
	}

	ip := net.ParseIP(*IP)
	if ip == nil {
		fmt.Printf("Invalid IP address: %s\n", *IP)
		return
	}

	if *Port == "" {
		fmt.Printf("Invalid port: %s\n", *IP)
		return
	}

	if *Host == "" {
		fmt.Println("Invalid Host name")
		return
	}

	if *User == "" {
		fmt.Println("Invalid user name")
		return
	}

	if *Protocol == "ssh" || *Upload == "ssh" || *Download == "ssh" {
		if *ServerUser == "" || *ServerPassword == "" {
			fmt.Println("Server User/Password is necessary for ssh")
			return
		}
	}

	if *Set != "" && *Value == "" {
		fmt.Println("You must give the value to set with value option")
		return
	}

	if *Upload != "" || *Download != "" {
		if *Server == "" {
			fmt.Println("Server IP is necessary for file download/upload!")
			return
		}

		if *ServerUser == "" {
			fmt.Println("Server user name is necessary for file download/upload!")
			return
		}

		if *Local == "" || *Remote == "" {
			fmt.Println("Local/Remote file name is necessary for file Up/Download")
			return
		}

		if !filepath.IsAbs(*Local) || !filepath.IsAbs(*Remote) {
			fmt.Println("Local/Remote file name should use absolute path.")
			return
		}
	}

	dev, err := rut.New(&rut.RUT{
		Name:     "M3000_210",
		Device:   "V5",
		Protocol: *Protocol,
		IP:       *IP,
		Port:     *Port,
		Username: *User,
		Hostname: *Host,
		Password: *Password,
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = dev.Init()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()

	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "q sh -l",
	})

	if *Table != "" {
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "bcm d all " + *Table,
		})
		fmt.Println(data)
		util.SaveToFile(filenameNormalize(*Table)+".txt", []byte(data))
	}

	if *Register != "" {
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "bcm g " + *Register,
		})
		fmt.Println(data)
		util.SaveToFile(filenameNormalize(*Register)+".txt", []byte(data))
	}

	if *Command != "" {
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "bcm " + *Command,
		})
		fmt.Println(data)
		util.SaveToFile(filenameNormalize(*Command)+".txt", []byte(data))
	}

	if *Shell != "" {
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  *Shell,
		})
		fmt.Println(data)
		util.SaveToFile(filenameNormalize(*Shell)+".txt", []byte(data))
	}

	if *Config != "" {
		err = Configure(*User, *Password, *IP, *Port, *Config)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if *Tcpdump != "" {
		err = dev.TCPDUMP(*Tcpdump, *TcpdumpFilter, *TcpdumpFile, *TcpdumpCount)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	if *Check != "" {
		if *Check == "register" {
			if *Name != "" {
				regs := strings.Split(*Name, ",")
				for _, reg := range regs {
					data, err = dev.RunCommand(ctx, &command.Command{
						Mode: "shell",
						CMD:  fmt.Sprintf("bcm l %s", reg),
					})

					if strings.Contains(string(data), "No matching register found") ||
						strings.Contains(string(data), "Unknown register address") {
						fmt.Printf("No register find for name/address: %s\n", reg)
						return
					}
				}
				for _, reg := range regs {
					CheckRegister(ctx, dev, reg)
				}
			} else {
				//Dump all registers
				CheckRegister(ctx, dev, "-t")
			}
		} else if *Check == "table" {
			if *Option != "raw" && *Option != "all" && *Option != "chg" && *Option != "" {
				fmt.Printf("Invalid option: \"%s\"\n", *Option)
				fmt.Println("\tThe valid option for table dump is raw/all/chg, by default is chg")
				return
			}

			op := *Option
			if op == "" {
				op = "chg"
			}

			if *Name != "" {
				tbls := strings.Split(*Name, ",")
				for _, tbl := range tbls {
					data, err := dev.RunCommand(ctx, &command.Command{
						Mode: "shell",
						CMD:  fmt.Sprintf("bcm list %s", tbl),
					})
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					if strings.Contains(string(data), "No memory found with the substring") {
						fmt.Printf("No table find for: %s\n", tbl)
						return
					}
				}

				for _, tbl := range tbls {
					CheckTable(ctx, dev, tbl, op)
				}
			} else {
				//Dump all memory
				CheckTable(ctx, dev, "", op)
			}
		} else {
			fmt.Println("Currently only support register/table check")
			return
		}
	}

	if *Set != "" {
		var fname string
		if *Field == "" {
			fname = "all"
		}

		if *Set == "register" {
			if *Name != "" {
				regs := strings.Split(*Name, ",")
				for _, reg := range regs {
					data, err = dev.RunCommand(ctx, &command.Command{
						Mode: "shell",
						CMD:  fmt.Sprintf("bcm l %s", reg),
					})

					if strings.Contains(string(data), "No matching register found") ||
						strings.Contains(string(data), "Unknown register address") {
						fmt.Printf("No register find for name/address: %s\n", reg)
						return
					}
				}
				for _, reg := range regs {
					SetRegister(ctx, dev, reg, fname, *Value)
				}
			} else {
				//Dump all registers
				SetRegister(ctx, dev, "-t", fname, *Value)
			}
		} else if *Set == "table" {
			if *Name != "" {
				tbls := strings.Split(*Name, ",")
				for _, tbl := range tbls {
					data, err := dev.RunCommand(ctx, &command.Command{
						Mode: "shell",
						CMD:  fmt.Sprintf("bcm list %s", tbl),
					})
					if err != nil {
						fmt.Println(err.Error())
						return
					}

					if strings.Contains(string(data), "No memory found with the substring") {
						fmt.Printf("No table find for: %s\n", tbl)
						return
					}
				}

				for _, tbl := range tbls {
					SetTable(ctx, dev, tbl, fname, *Value)
				}
			} else {
				//Dump all memory
				SetTable(ctx, dev, "", fname, *Value)
			}
		} else {
			fmt.Println("Currently only support register/table set")
			return
		}
	}

	if *Monitor != "" {
		if *Monitor == "memory" {
			memoryMonitor(ctx, dev)
		} else if *Monitor == "cpuload" {
			cpuloadMonitor(ctx, dev)
		} else if *Monitor == "status" {
			statisticsMonitor(ctx, dev)
		} else if *Monitor == "statistics" {
			statisticsMonitor(ctx, dev)
		} else if *Monitor == "cpustat" {
			cpuStatisticsMonitor(ctx, dev)
		}
	}

	if *Upload != "" {
		if *Upload == "ssh" {
			if err = dev.SCPPut(*Local, *Server, *ServerUser, *ServerPassword, filepath.Dir(*Remote)); err != nil {
				fmt.Printf("Cannot upload file: %s to %s with: %s\n", *Local, *IP, err.Error())
				return
			}
		} else {
			if err = dev.FTPPut(*Local, *Server, *ServerUser, *ServerPassword, filepath.Dir(*Remote)); err != nil {
				fmt.Printf("Cannot upload file: %s to %s with: %s\n", *Local, *IP, err.Error())
				return
			}
		}
	}

	if *Download != "" {
		if *Download == "ssh" {
			if err = dev.SCPGet(*Local, *Server, *ServerUser, *ServerPassword, filepath.Dir(*Remote), filepath.Base(*Remote)); err != nil {
				fmt.Printf("Cannot download file: %s from %s with: %s\n", *Remote, *IP, err.Error())
				return
			}
		} else {

			if err = dev.FTPGet(*Local, *Server, *ServerUser, *ServerPassword, filepath.Dir(*Remote), filepath.Base(*Remote)); err != nil {
				fmt.Printf("Cannot download file: %s from %s with: %s\n", *Remote, *IP, err.Error())
				return
			}
		}
	}

	/* Current compile error on 3.14
	if *PCAPDecode != "" {
		if !strings.HasSuffix(*PCAPDecode, ".pcap") {
			fmt.Println("You must give a pcap file for decode")
			return
		}

		handle, err = pcap.OpenOffline(*PCAPDecode)
		if err != nil {
			log.Fatal(err)
		}
		defer handle.Close()

		// Loop through packets in file
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			fmt.Println(packet)
		}
	}
	*/
}

func Configure(user, pass, ip, port, config string) error {
	c, err := telnet.New(user, pass, ip, port)
	if err != nil {
		return err
	}

	c.WriteLine("enable")
	_, err = c.ReadUntil("#")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	c.WriteLine("terminal length 0")
	_, err = c.ReadUntil("#")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	f, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	lines := strings.Split(string(f), "\n")
	for _, l := range lines {
		if strings.TrimSpace(string(l)) == "" {
			continue
		}

		c.WriteLine(l)
		result, err := c.ReadUntil("#")
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		if strings.Contains(string(result), "Invalid input detected at") ||
			strings.Contains(string(result), "VTY configuration is locked by other VTY") {
			fmt.Printf("Run command \"%s\" with error: \n%s\n", string(l), string(result))
			return fmt.Errorf("Run command \"%s\" with error: \n%s\n", string(l), string(result))
		}
		fmt.Println(string(result))
	}

	return nil
}

func CheckTable(ctx context.Context, dev *rut.RUT, name, op string) {
	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("bcm list %s", name),
	})

	if strings.Contains(string(data), "No memory found with the substring") {
		fmt.Printf("No table find for: %s\n", name)
		return
	}
	fmt.Println(string(data))
	//Table name exact match
	if strings.Contains(string(data), "Description:") {
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  fmt.Sprintf("bcm d %s %s", op, name),
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	} else { //Get a list of tabless with give name
		result := strings.Split(string(data), "Entry/Copy Description")
		tbls := result[1][:strings.Index(result[1], "Flags:")]
		lines := strings.Split(tbls, "\r\x00\r\n")
		for _, l := range lines {
			if strings.TrimSpace(l) == "" {
				continue
			}
			l = strings.Replace(l, "*", " ", -1)
			fields := strings.Split(l, " ")
			for _, field := range fields {
				if !strings.Contains(strings.ToUpper(field), strings.ToUpper(name)) {
					continue
				}

				//*Exact match is handled by exact match logic. If goes here, we must be not exact match.
				//This is used to skip some table description.
				if strings.TrimSpace(field) == name {
					continue
				}

				if strings.TrimSpace(field) == "" {
					continue
				}

				if strings.ToUpper(strings.TrimSpace(field)) != strings.TrimSpace(field) {
					continue
				}

				if strings.Contains(field, "-") {
					continue
				}

				data, err = dev.RunCommand(ctx, &command.Command{
					Mode: "shell",
					CMD:  fmt.Sprintf("bcm d %s %s", op, field),
				})
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				fmt.Println(string(data))
				//Each line only contains one table name.
				break
			}
		}
	}
}

func getTableInfo(ctx context.Context, dev *rut.RUT, name string) (*bcmTable, error) {
	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("bcm list %s", name),
	})
	if err != nil {
		return nil, err
	}

	if strings.Contains(string(data), "No memory found with the substring") {
		fmt.Printf("No table find for: %s\n", name)
		return nil, fmt.Errorf("No table find for: %s\n", name)
	}

	if !strings.Contains(string(data), "Description:") {
		return nil, fmt.Errorf("Invalid dump result: %s\n", string(data))

	}

	var table bcmTable

	table.Name = name
	fields := make([]string, 0, 1)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {

		if strings.TrimSpace(line) == "" {
			continue
		}
		line := strings.TrimSpace(line)
		mt := tf.FindAllStringSubmatch(line, -1)
		if len(mt) == 1 {
			fields = append(fields, string(mt[0][1]))
		}
	}

	table.Fields = fields

	mt := ts.FindStringSubmatch(string(data))
	if len(mt) == 2 {
		size, err := strconv.ParseInt(mt[1], 10, 64)
		if err != nil {
			fmt.Println("Cannot get table size for: %s", name)
			return nil, fmt.Errorf("Cannot get table size for: %s", name)
		}
		table.Size = size
	}

	return &table, nil
}

func SetTable(ctx context.Context, dev *rut.RUT, name, fname, value string) {
	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("bcm list %s", name),
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if strings.Contains(string(data), "No memory found with the substring") {
		fmt.Printf("No table find for: %s\n", name)
		return
	}
	fmt.Println(string(data))
	//Table name exact match
	if strings.Contains(string(data), "Description:") {
		table, err := getTableInfo(ctx, dev, name)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if fname == "all" {
			for _, f := range table.Fields {
				data, err = dev.RunCommand(ctx, &command.Command{
					Mode: "shell",
					CMD:  fmt.Sprintf("bcm mod %s %d %d %s=%s", table.Name, 0, table.Size-1, f, value),
				})
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(string(data))
			}
		} else {
			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  fmt.Sprintf("bcm mod %s %d %d %s=%s", table.Name, 0, table.Size-1, fname, value),
			})
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println(string(data))
		}
	} else { //Get a list of tabless with give name
		result := strings.Split(string(data), "Entry/Copy Description")
		tbls := result[1][:strings.Index(result[1], "Flags:")]
		lines := strings.Split(tbls, "\r\x00\r\n")
		for _, l := range lines {
			if strings.TrimSpace(l) == "" {
				continue
			}
			l = strings.Replace(l, "*", " ", -1)
			fields := strings.Split(l, " ")
			for _, field := range fields {
				if !strings.Contains(strings.ToUpper(field), strings.ToUpper(name)) {
					continue
				}

				//*Exact match is handled by exact match logic. If goes here, we must be not exact match.
				//This is used to skip some table description.
				if strings.TrimSpace(field) == name {
					continue
				}

				if strings.TrimSpace(field) == "" {
					continue
				}

				if strings.ToUpper(strings.TrimSpace(field)) != strings.TrimSpace(field) {
					continue
				}

				if strings.Contains(field, "-") {
					continue
				}

				table, err := getTableInfo(ctx, dev, field)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Printf("%s=%q\n", fields, table)
				if fname == "all" {
					for _, f := range table.Fields {
						data, err = dev.RunCommand(ctx, &command.Command{
							Mode: "shell",
							CMD:  fmt.Sprintf("bcm mod %s %d %d %s=%s", table.Name, 0, table.Size-1, f, value),
						})
						if err != nil {
							fmt.Println(err.Error())
							return
						}
						fmt.Println(string(data))
					}
				} else {
					data, err = dev.RunCommand(ctx, &command.Command{
						Mode: "shell",
						CMD:  fmt.Sprintf("bcm mod %s %d %d %s=%s", table.Name, 0, table.Size-1, fname, value),
					})
					if err != nil {
						fmt.Println(err.Error())
						return
					}
					fmt.Println(string(data))
				}

				//Each line only contains one table name.
				break
			}
		}
	}
}

func CheckRegister(ctx context.Context, dev *rut.RUT, name string) {
	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("bcm l %s", name),
	})

	if strings.Contains(string(data), "No matching register found") ||
		strings.Contains(string(data), "Unknown register address") {
		fmt.Printf("No register find for name/address: %s\n", name)
		return
	}
	//Register name exact match
	if strings.Contains(string(data), "Description:") {
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  fmt.Sprintf("bcm g %s", name),
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	} else { //Get a list of registers with give name
		var regs string
		var lines []string
		if name == "-t" {
			data := UseLessFilter.ReplaceAllString(data, "")
			lines = strings.Split(string(data), "cmd: l -t")
		} else {
			lines = strings.Split(string(data), "possible matches are:")
		}
		regs = lines[1][:strings.Index(lines[1], "BCM.")]
		regs = strings.TrimSpace(regs)
		regs = strings.Replace(regs, "\r\x00\r\n", "", -1)
		regs = strings.Replace(regs, "\x00\r\n", "", -1)
		regs = strings.Replace(regs, "\r\x00", "", -1)
		all := strings.Split(regs, " ")
		for _, r := range all {
			if strings.TrimSpace(r) == "" || strings.TrimSpace(r) == " " {
				continue
			}
			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  fmt.Sprintf("bcm g %s", r),
			})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println(string(data))
		}

	}
}

func SetRegister(ctx context.Context, dev *rut.RUT, name, fname, value string) {
	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("bcm l %s", name),
	})

	if strings.Contains(string(data), "No matching register found") ||
		strings.Contains(string(data), "Unknown register address") {
		fmt.Printf("No register find for name/address: %s\n", name)
		return
	}

	var cmd string
	//Register name exact match
	if strings.Contains(string(data), "Description:") {
		if fname == "all" {
			cmd = fmt.Sprintf("bcm s %s", name, value)
		} else {
			cmd = fmt.Sprintf("bcm m %s %s=%s", name, fname, value)
		}
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  cmd,
		})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(data))
	} else { //Get a list of registers with give name
		var regs string
		var lines []string
		if name == "-t" {
			data := UseLessFilter.ReplaceAllString(data, "")
			lines = strings.Split(string(data), "cmd: l -t")
		} else {
			lines = strings.Split(string(data), "possible matches are:")
		}
		regs = lines[1][:strings.Index(lines[1], "BCM.")]
		regs = strings.TrimSpace(regs)
		regs = strings.Replace(regs, "\r\x00\r\n", "", -1)
		regs = strings.Replace(regs, "\x00\r\n", "", -1)
		regs = strings.Replace(regs, "\r\x00", "", -1)
		all := strings.Split(regs, " ")
		for _, r := range all {
			if strings.TrimSpace(r) == "" || strings.TrimSpace(r) == " " {
				continue
			}

			if fname == "all" {
				cmd = fmt.Sprintf("bcm s %s", r, value)
			} else {
				cmd = fmt.Sprintf("bcm m %s %s=%s", r, fname, value)
			}

			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "shell",
				CMD:  cmd,
			})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println(string(data))
		}

	}
}

func filenameNormalize(command string) string {
	return strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.TrimSpace(command), "/", "", -1), " ", "", -1), ":", "", -1), ";", "", -1)
}

func memoryMonitor(ctx context.Context, dev *rut.RUT) {
	//var pResult []byte
	var iResult []byte
	var cResult []byte
	for range time.Tick(time.Duration(time.Second * 5)) {
		/*Keep track to last time result */
		/*
			if cResult != nil {
				pResult = cResult
			}
		*/

		result := make([]byte, 0, 1024)
		data, err := dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "cat /proc/meminfo",
		})
		if err != nil {
			panic(err)
		}

		result = append(result, data...)

		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "cat /proc/slabinfo",
		})
		if err != nil {
			panic(err)
		}

		result = append(result, data...)
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "ps aux",
		})
		if err != nil {
			panic(err)
		}
		result = append(result, data...)

		/*Keep track to the first time result*/
		if iResult == nil {
			iResult = result
		}

		cResult = result

		if iResult != nil && cResult != nil {
			util.Diff2(string(cResult), string(iResult))
		}

		util.SaveToFile("memory_monitor.txt", []byte(result))
	}
}

func cpuloadMonitor(ctx context.Context, dev *rut.RUT) {
	//var pResult []byte
	var iResult []byte
	var cResult []byte
	for range time.Tick(time.Duration(time.Second * 5)) {
		/*Keep track to last time result */
		/*
			if cResult != nil {
				pResult = cResult
			}
		*/

		data, err := dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "ps aux",
		})
		if err != nil {
			panic(err)
		}

		/*Keep track to the first time result*/
		if iResult == nil {
			iResult = []byte(data)
		}

		cResult = []byte(data)

		if iResult != nil && cResult != nil {
			util.Diff2(string(cResult), string(iResult))
		}

		util.SaveToFile("cpu_monitor.txt", []byte(cResult))
	}
}

func statisticsMonitor(ctx context.Context, dev *rut.RUT) {
	err := getAllInterface(ctx, dev)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(InterfaceList))
	go updateInterfaceStatus(ctx, dev)
	go updateInterfaceStatistics(ctx, dev)

	for range time.Tick(time.Duration(time.Second * 5)) {
		for _, ifp := range InterfaceList {
			if *Monitor == "statistics" && ifp.Hardware == "ETH" && ifp.AdminStatus == "UP" && ifp.OperationStatus == "UP" {
				fmt.Printf("%s\n", ifp)
			} else {
				fmt.Printf("%s\n", ifp)
			}
		}
	}
}

func cpuStatisticsMonitor(ctx context.Context, dev *rut.RUT) {

}

func getAllInterface(ctx context.Context, dev *rut.RUT) error {
	/* First go to command mode */
	data, err := dev.RunCommand(ctx, &command.Command{
		Mode: "shell",
		CMD:  "exit",
	})
	if err != nil {
		panic(err)
	}

	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "show interface",
	})
	if err != nil {
		panic(err)
	}

	InterfaceLock.Lock()
	ifs := ifr.FindAllStringSubmatch(string(data), -1)
	for _, ift := range ifs {
		var newIf = Interface{Name: ift[1]}
		InterfaceList = append(InterfaceList, &newIf)
	}
	InterfaceLock.Unlock()

	/* Go back to shell mode */
	data, err = dev.RunCommand(ctx, &command.Command{
		Mode: "normal",
		CMD:  "q sh -l",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
	if len(InterfaceList) == 0 {
		return fmt.Errorf("No interface found.")
	}

	return nil
}

func updateInterfaceStatus(ctx context.Context, dev *rut.RUT) {
	var tick time.Duration
	if len(InterfaceList) == 0 {
		tick = time.Duration(time.Second)
	} else {
		tick = time.Duration(time.Second * 10)
	}
	for range time.Tick(tick) {
		IfMap := make(map[string]*Interface, len(InterfaceList))
		InterfaceLock.Lock()
		for _, ifp := range InterfaceList {
			IfMap[ifp.Name] = ifp
		}
		data, err := dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "exit",
		})
		if err != nil {
			panic(err)
		}

		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "normal",
			CMD:  "show interface",
		})
		if err != nil {
			panic(err)
		}

		ifns := ifr.FindAllStringSubmatch(string(data), -1)
		ifis := ifr.Split(string(data), -1)
		for i, ifn := range ifns {
			if _, ok := IfMap[ifn[1]]; !ok {
				var newIf = Interface{Name: ifn[1]}
				InterfaceList = append(InterfaceList, &newIf)
				IfMap[ifn[1]] = &newIf
			}

			ifp := IfMap[ifn[1]]
			ifi := ifis[i+1]

			//Get interface staus information
			ifs := ifstatus.FindStringSubmatch(ifi)
			if strings.Contains(ifs[1], "UP") {
				ifp.AdminStatus = "UP"
			} else {
				ifp.AdminStatus = "Down"
			}

			if strings.Contains(ifs[1], "LOOPBACK") {
				ifp.Loopback = true
			} else {
				ifp.Loopback = false
			}

			if strings.Contains(ifs[1], "RUNNING") {
				ifp.OperationStatus = "UP"
			} else {
				ifp.OperationStatus = "DOWN"
			}

			if strings.Contains(ifs[1], "MULTICAST") {
				ifp.Multicast = true
			} else {
				ifp.Multicast = false
			}

			if strings.Contains(ifs[1], "BROADCAST") {
				ifp.Broadcast = true
			} else {
				ifp.Broadcast = false
			}

			if !ifp.Loopback {
				ifh := ifhw.FindStringSubmatch(ifi)
				ifp.Hardware = ifh[1]
				ifp.HardwareAddress = ifh[2]
			} else {
				ifp.Hardware = "LB"
			}

			ifindex := ifindexr.FindStringSubmatch(ifi)
			ifp.Index = ifindex[1]

			ifmtu := ifmtur.FindStringSubmatch(ifi)
			ifp.MTU = ifmtu[1]

			ifbw := ifbandwidthr.FindStringSubmatch(ifi)
			ifp.Bandwidth = ifbw[1]

			if strings.Contains(ifi, "no switchport") {
				ifp.Mode = "Router"
			} else {
				ifp.Mode = "Switch"
			}

			ips := ipr.FindAllStringSubmatch(ifi, -1)
			if len(ips) > 0 {
				ifp.IP = make([]string, 0, 1)
				for _, ip := range ips {
					ifp.IP = append(ifp.IP, ip[1])
				}
			}

			ip6s := ip6r.FindAllStringSubmatch(ifi, -1)
			if len(ip6s) > 0 {
				ifp.IP6 = make([]string, 0, 1)
				for _, ip6 := range ip6s {
					ifp.IP6 = append(ifp.IP6, ip6[1])
				}
			}

			if strings.Contains(ifi, "VRF Binding: Not bound") {
				ifp.VRF = "Global"
			} else {
				ifvrf := ifvrfr.FindStringSubmatch(ifi)
				ifp.VRF = ifvrf[1]
			}
		}

		/* Go back to shell mode */
		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "normal",
			CMD:  "q sh -l",
		})
		if err != nil {
			panic(err)
		}

		InterfaceLock.Unlock()
	}
}

func updateInterfaceStatistics(ctx context.Context, dev *rut.RUT) {
	for range time.Tick(time.Duration(time.Second * 4)) {
		InterfaceLock.Lock()
		data, err := dev.RunCommand(ctx, &command.Command{
			Mode: "shell",
			CMD:  "exit",
		})
		if err != nil {
			panic(err)
		}

		IfMap := make(map[string]*Interface, len(InterfaceList))

		for _, ifp := range InterfaceList {
			if ifp.AdminStatus == "UP" && ifp.OperationStatus == "UP" && ifp.Hardware == "ETH" {
				IfMap[ifp.Name] = ifp
			}
		}

		for _, ifp := range IfMap {
			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "normal",
				CMD:  fmt.Sprintf("show interface statistics avg %s", ifp.Name),
			})

			if err != nil {
				panic(err)
			}
			rates := ifrater.FindStringSubmatch(string(data))
			ifp.RxPPS = rates[1]
			ifp.RxBPS = rates[2]
			ifp.TxPPS = rates[3]
			ifp.TxBPS = rates[4]

			data, err = dev.RunCommand(ctx, &command.Command{
				Mode: "normal",
				CMD:  fmt.Sprintf("show interface statistics interface %s", ifp.Name),
			})

			if err != nil {
				panic(err)
			}

			upkts := ifstatupr.FindStringSubmatch(string(data))
			ifp.RxUcastPacket, err = strconv.ParseUint(upkts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			ifp.TxUcastPacket, err = strconv.ParseUint(upkts[2], 10, 64)
			if err != nil {
				panic(err)
			}

			mpkts := ifstatmpr.FindStringSubmatch(string(data))
			ifp.RxMcastPacket, err = strconv.ParseUint(mpkts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			ifp.TxMcastPacket, err = strconv.ParseUint(mpkts[2], 10, 64)
			if err != nil {
				panic(err)
			}

			bpkts := ifstatbpr.FindStringSubmatch(string(data))
			ifp.RxBcastPacket, err = strconv.ParseUint(bpkts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			ifp.TxBcastPacket, err = strconv.ParseUint(bpkts[2], 10, 64)
			if err != nil {
				panic(err)
			}

			ifp.RxTotalPacket = ifp.RxUcastPacket + ifp.RxBcastPacket + ifp.RxMcastPacket
			ifp.TxTotalPacket = ifp.TxUcastPacket + ifp.TxBcastPacket + ifp.TxMcastPacket

			dpkts := ifstatDiscardr.FindStringSubmatch(string(data))
			ifp.RxDiscard, err = strconv.ParseUint(dpkts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			ifp.TxDiscard, err = strconv.ParseUint(dpkts[2], 10, 64)
			if err != nil {
				panic(err)
			}

			epkts := ifstatErrorr.FindStringSubmatch(string(data))
			ifp.RxErrors, err = strconv.ParseUint(epkts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			ifp.TxErrors, err = strconv.ParseUint(epkts[2], 10, 64)
			if err != nil {
				panic(err)
			}
		}

		data, err = dev.RunCommand(ctx, &command.Command{
			Mode: "normal",
			CMD:  "q sh -l",
		})
		if err != nil {
			panic(err)
		}
		InterfaceLock.Unlock()
	}
}
