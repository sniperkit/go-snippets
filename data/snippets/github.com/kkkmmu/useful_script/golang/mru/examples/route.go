package main

import (
	"command"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"net"
	"regexp"
	"rut"
	"strconv"
	"strings"
)

const (
	IPV4 = 1
	IPV6 = 2
)

type Threshold struct {
	MaxNexthopIndex     int64
	MaxIPv4ALPMIndex    int64
	MaxIPv6ALPM64Index  int64
	MaxIPv6ALPM128Index int64
}

var threshold = Threshold{
	MaxNexthopIndex:     49152,
	MaxIPv4ALPMIndex:    393216,
	MaxIPv6ALPM64Index:  262144,
	MaxIPv6ALPM128Index: 131072,
}

type RouteEntry struct {
	Index           int64
	Valid           string
	RPE             string
	ReservedECMPPTR string
	Reserved0       string
	PRI             string
	NextHopIndex    int64
	Length          int64
	Key             net.IP
	Hit             bool
	EvenParity      string
	EntryOnly       string
	ECMPPTR         int64
	ECMP            bool
	DstDiscard      bool
	SrcDiscard      bool
	DefaultRoute    bool
	Data            string
	Class_ID        string
	NH              Nexthop
	EG              ECMPGroup
	AF              int
	IsValid         bool
}

type HostEntry struct {
	Index        int64
	Valid        string
	RPE          string
	PRI          string
	NextHopIndex int64
	IP           net.IP
	Hit          bool
	Length       int64
	ECMPPTR      int64
	ECMP         bool
	DstDiscard   bool
	Data         string
	KeyType      int64
	NH           Nexthop
	EG           ECMPGroup
	AF           int
	IsValid      bool
}

var Dev *rut.RUT

type IF struct {
	Index int64
	Vid   int64
	MAC   string
}

type Nexthop struct {
	Index     int64
	OIF       IF
	DstPort   int64
	DstMac    string
	Drop      bool
	CopyToCPU bool
	MTU       int64
	TGID      int64
	VLAN      int64
}

type ECMPGroup struct {
	Index       int64
	MemberCount int64
	Member      []*Nexthop
	ECMPBasePTR int64
}

var CTX = context.Background()

func (eg *ECMPGroup) ParseNexthop() {
	nhs := make([]*Nexthop, 0, eg.MemberCount)

	var i int64
	for i = 0; i < eg.MemberCount; i++ {
		nh, err := Dev.RunCommand(CTX, &command.Command{
			Mode: "shell",
			CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_ECMP %d %d", eg.ECMPBasePTR+i, eg.ECMPBasePTR+i),
		})

		if err != nil {
			panic(err)
		}

		if res, err := match(nh, getNextHopIndex); err != nil {
			panic(err)
		} else {
			nhi, err := strconv.ParseInt(res, 0, 64)
			if err != nil {
				panic(err)
			} else {
				nh, err := ParseNexthopByIndex(nhi)
				if err != nil {
					panic(err)
				}

				nhs = append(nhs, nh)
			}
		}
	}

	eg.Member = nhs
}

func ParseNexthopByIndex(index int64) (*Nexthop, error) {
	if index < 0 || index > threshold.MaxNexthopIndex {
		panic("Invalid nexthop index")
	}

	nh := &Nexthop{
		Index: index,
	}
	nh.GetINGNexthopInfo()
	nh.GetEGRNexthopInfo()
	nh.ParseOIF()

	return nh, nil
}

func (nh *Nexthop) GetINGNexthopInfo() {
	inh, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 ING_L3_NEXT_HOP %d %d", nh.Index, nh.Index),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(inh, getINGNexthopCopyToCPU); err != nil {
		panic(err)
	} else {
		if res == "0" {
			nh.CopyToCPU = false
		} else {
			nh.CopyToCPU = true
		}
	}

	if res, err := match(inh, getINGNexthopDrop); err != nil {
		panic(err)
	} else {
		if res == "0" {
			nh.Drop = false
		} else {
			nh.Drop = true
		}
	}

	if res, err := match(inh, getINGNexthopTGID); err != nil {
		panic(err)
	} else {
		tgid, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			nh.TGID = tgid
		}
	}

	if res, err := match(inh, getINGNexthopPNUM); err != nil {
		panic(err)
	} else {
		pnum, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			nh.DstPort = pnum
		}
	}

	if res, err := match(inh, getINGNexthopVLAN); err != nil {
		panic(err)
	} else {
		vlan, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			nh.VLAN = vlan
		}
	}

	/*
		if res, err := match(inh, getINGNexthopOIF); err != nil {
			panic(err)
		} else {
			oif, err := strconv.ParseInt(res, 0, 32)
			if err != nil {
				panic(err)
			} else {
				nh.OIF = oif
			}
		}
	*/

	//fmt.Println(inh, err)
}

func (nh *Nexthop) GetEGRNexthopInfo() {
	enh, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 EGR_L3_NEXT_HOP %d %d", nh.Index, nh.Index),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(enh, getEGRNexthopMACAddress); err != nil {
		panic(err)
	} else {
		nh.DstMac = FixMACAddress(res)
	}

	if res, err := match(enh, getEGRIfNum); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		}
		nh.OIF.Index = index
	}
}

func (re *RouteEntry) String() string {
	var Out *color.Color
	if re.Hit {
		Out = color.New(color.FgGreen)
	} else {
		Out = color.New(color.FgWhite)
	}
	if re.AF == IPV4 {
		if re.Length > 32 || re.NextHopIndex > threshold.MaxNexthopIndex {
			if re.Length > 32 {
				return Out.Sprintf("[%6d]: %39s/%-3d >> %20s", re.Index, re.Key, re.Length, "is not a valid IPv4 Address")
			} else {
				return Out.Sprintf("[%6d]: %39s/%-3d >> has a invalid nexthop index: %d", re.Index, re.Key, re.Length, re.NextHopIndex)
			}
		} else if !re.ECMP {
			return Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> NH[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", re.Index, re.Key, re.Length, re.DstDiscard, re.SrcDiscard, re.NextHopIndex, re.NH.OIF.MAC, re.NH.DstMac, re.NH.OIF.Index, re.NH.OIF.Vid, re.NH.DstPort)
		} else {
			base := Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> is  ECMP Route, ECMP_PTR: %5d, BASE_PTR: %5d, MemberCount: %2d", re.Index, re.Key, re.Length, re.DstDiscard, re.SrcDiscard, re.ECMPPTR, re.EG.ECMPBasePTR, re.EG.MemberCount)
			for i := 0; i < len(re.EG.Member); i++ {
				base += "\n"
				base += Out.Sprintf("%73s[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", "NH", re.EG.Member[i].Index, re.EG.Member[i].OIF.MAC, re.EG.Member[i].DstMac, re.EG.Member[i].OIF.Index, re.EG.Member[i].OIF.Vid, re.EG.Member[i].DstPort)
			}

			return base
		}
	} else if re.AF == IPV6 {
		if re.Length > 128 || re.NextHopIndex > threshold.MaxNexthopIndex {
			if re.Length > 128 {
				return Out.Sprintf("[%6d]: %39s/%-3d >> %20s", re.Index, re.Key, re.Length, "is not a valid IPv6 Address")
			} else {
				return Out.Sprintf("[%6d]: %39s/%-3d >> has a invalid nexthop index: %d", re.Index, re.Key, re.Length, re.NextHopIndex)
			}
		} else if !re.ECMP {
			return Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> NH[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", re.Index, re.Key, re.Length, re.DstDiscard, re.SrcDiscard, re.NextHopIndex, re.NH.OIF.MAC, re.NH.DstMac, re.NH.OIF.Index, re.NH.OIF.Vid, re.NH.DstPort)
		} else {
			base := Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> is  ECMP Route, ECMP_PTR: %5d, BASE_PTR: %5d, MemberCount: %2d", re.Index, re.Key, re.Length, re.DstDiscard, re.SrcDiscard, re.ECMPPTR, re.EG.ECMPBasePTR, re.EG.MemberCount)
			for i := 0; i < len(re.EG.Member); i++ {
				base += "\n"
				base += Out.Sprintf("%73s[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", "NH", re.EG.Member[i].Index, re.EG.Member[i].OIF.MAC, re.EG.Member[i].DstMac, re.EG.Member[i].OIF.Index, re.EG.Member[i].OIF.Vid, re.EG.Member[i].DstPort)
			}

			return base
		}
	}

	return fmt.Sprintf("Invalid route entry: %s/%d\n", re.Key, re.Length)
}

func (he *HostEntry) String() string {
	var Out *color.Color
	if he.Hit {
		Out = color.New(color.FgGreen)
	} else {
		Out = color.New(color.FgWhite)
	}
	if he.AF == IPV4 && he.KeyType == 0 {
		if he.NextHopIndex > threshold.MaxNexthopIndex {
			return Out.Sprintf("[%6d]: %39s/%-3d >> has a invalid nexthop index: %d", he.Index, he.IP, he.Length, he.NextHopIndex)
		} else if !he.ECMP {
			return Out.Sprintf("[%6d]: %39s/%-3d (%5t|%2d) >> NH[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", he.Index, he.IP, he.Length, he.DstDiscard, he.KeyType, he.NextHopIndex, he.NH.OIF.MAC, he.NH.DstMac, he.NH.OIF.Index, he.NH.OIF.Vid, he.NH.DstPort)
		} else {
			base := Out.Sprintf("[%6d]: %39s/%-3d (%5t|%2d) >> is  ECMP Route, ECMP_PTR: %5d, BASE_PTR: %5d, MemberCount: %2d", he.Index, he.IP, he.Length, he.DstDiscard, he.KeyType, he.ECMPPTR, he.EG.ECMPBasePTR, he.EG.MemberCount)
			for i := 0; i < len(he.EG.Member); i++ {
				base += "\n"
				base += Out.Sprintf("%73s[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", "NH", he.EG.Member[i].Index, he.EG.Member[i].OIF.MAC, he.EG.Member[i].DstMac, he.EG.Member[i].OIF.Index, he.EG.Member[i].OIF.Vid, he.EG.Member[i].DstPort)
			}

			return base
		}
	} else if he.AF == IPV6 && he.KeyType == 2 {
		if he.NextHopIndex > threshold.MaxNexthopIndex {
			return Out.Sprintf("[%6d]: %39s/%-3d >> has a invalid nexthop index: %d", he.Index, he.IP, he.Length, he.NextHopIndex)
		} else if !he.ECMP {
			return Out.Sprintf("[%6d]: %39s/%-3d (%5t|%2d) >> NH[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", he.Index, he.IP, he.Length, he.DstDiscard, he.KeyType, he.NextHopIndex, he.NH.OIF.MAC, he.NH.DstMac, he.NH.OIF.Index, he.NH.OIF.Vid, he.NH.DstPort)
		} else {
			base := Out.Sprintf("[%6d]: %39s/%-3d (%5t|%2d) >> is  ECMP Route, ECMP_PTR: %5d, BASE_PTR: %5d, MemberCount: %2d", he.Index, he.IP, he.Length, he.DstDiscard, he.KeyType, he.ECMPPTR, he.EG.ECMPBasePTR, he.EG.MemberCount)
			for i := 0; i < len(he.EG.Member); i++ {
				base += "\n"
				base += Out.Sprintf("%73s[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", "NH", he.EG.Member[i].Index, he.EG.Member[i].OIF.MAC, he.EG.Member[i].DstMac, he.EG.Member[i].OIF.Index, he.EG.Member[i].OIF.Vid, he.EG.Member[i].DstPort)
			}

			return base
		}
	}

	//Currently just skip the invalid entry.
	return ""
}

var getEntryIndex = regexp.MustCompile(`[[:word:]_]+\.\*\[(?P<index>[0-9]+)\]`)
var getValidBit = regexp.MustCompile(`VALID=(?P<valid>[0-9]+)`)
var getNextHopIndex = regexp.MustCompile(`,NEXT_HOP_INDEX=(?P<nhi>[0x]?[[:alnum:]]+)`)
var getLength = regexp.MustCompile(`,LENGTH=(?P<len>[0x]?[[:alnum:]]+)`)
var getKey = regexp.MustCompile(`,KEY=(?P<key>[0x]?[[:alnum:]]+)`)
var getHitBit = regexp.MustCompile(`,HIT=(?P<hit>[0-9]+)`)
var getECMPPTR = regexp.MustCompile(`,ECMP_PTR=(?P<ecmpptr>[0x]?[[:alnum:]]+)`)
var getECMP = regexp.MustCompile(`,ECMP=(?P<ecmpptr>[0x]?[[:alnum:]]+)`)
var getSrcDiscard = regexp.MustCompile(`,SRC_DISCARD=(?P<srcdis>[0-9]+)`)
var getDstDiscard = regexp.MustCompile(`,DST_DISCARD=(?P<dstdis>[0-9]+)`)
var getDefaultRoute = regexp.MustCompile(`,DEFAULTROUTE=(?P<default>[0-9]+)`)

//ING_L3_NEXT_HOP
var getEGRNexthopMACAddress = regexp.MustCompile(`L3:MAC_ADDRESS=(?P<nmac>[0x]?[[:alnum:]]+)`)
var getEGRIfNum = regexp.MustCompile(`L3:INTF_NUM=(?P<ifn>[0x]?[[:alnum:]]+)`)
var getEGRIVID = regexp.MustCompile(`L3:IVID=(?P<ivid>[0x]?[[:alnum:]]+)`)
var getEGROVID = regexp.MustCompile(`L3:OVID=(?P<ovid>[0x]?[[:alnum:]]+)`)
var getEGRNexthopEntryType = regexp.MustCompile(`,ENTRY_TYPE=(?P<type>[0x]?[[:alnum:]]+)`)

//EGR_L3_NEXT_HOP
var getINGNexthopCopyToCPU = regexp.MustCompile(`COPY_TO_CPU=(?P<ctc>[0x]?[[:alnum:]]+)`)
var getINGNexthopDrop = regexp.MustCompile(`DROP=(?P<ctc>[0x]?[[:alnum:]]+)`)
var getINGNexthopTGID = regexp.MustCompile(`TGID=(?P<tgid>[0x]?[[:alnum:]]+)`)
var getINGNexthopPNUM = regexp.MustCompile(`,PORT_NUM=(?P<pnum>[0x]?[[:alnum:]]+)`)
var getINGNexthopVLAN = regexp.MustCompile(`VLAN_ID=(?P<vid>[0x]?[[:alnum:]]+)`)
var getINGNexthopOIF = regexp.MustCompile(`,L3_OIF=(?P<oif>[0x]?[[:alnum:]]+)`)
var getINGNexthopEntryType = regexp.MustCompile(`,ENTRY_TYPE=(?P<type>[0x]?[[:alnum:]]+)`)
var getINGNexthopMTU = regexp.MustCompile(`MTU_SIZE=(?P<mtu>[0x]?[[:alnum:]]+)`)

func (re *RouteEntry) Validate() {
	if re.AF == IPV4 {
		if re.Index >= threshold.MaxIPv4ALPMIndex {
			re.IsValid = false
		}
	}

	if re.AF == IPV6 {
		if re.Index >= threshold.MaxIPv6ALPM64Index {
			re.IsValid = false
		}
	}

	if re.NextHopIndex >= threshold.MaxNexthopIndex {
		re.IsValid = false
	}

	re.IsValid = true
}

func (re *RouteEntry) ParseINGNexthopInfo() {
	inh, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 ING_L3_NEXT_HOP %d %d", re.NextHopIndex, re.NextHopIndex),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(inh, getINGNexthopCopyToCPU); err != nil {
		panic(err)
	} else {
		if res == "0" {
			re.NH.CopyToCPU = false
		} else {
			re.NH.CopyToCPU = true
		}
	}

	if res, err := match(inh, getINGNexthopDrop); err != nil {
		panic(err)
	} else {
		if res == "0" {
			re.NH.Drop = false
		} else {
			re.NH.Drop = true
		}
	}

	if res, err := match(inh, getINGNexthopTGID); err != nil {
		panic(err)
	} else {
		tgid, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			re.NH.TGID = tgid
		}
	}

	if res, err := match(inh, getINGNexthopPNUM); err != nil {
		panic(err)
	} else {
		pnum, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			re.NH.DstPort = pnum
		}
	}

	if res, err := match(inh, getINGNexthopVLAN); err != nil {
		panic(err)
	} else {
		vlan, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			re.NH.VLAN = vlan
		}
	}

	/*
		if res, err := match(inh, getINGNexthopOIF); err != nil {
			panic(err)
		} else {
			oif, err := strconv.ParseInt(res, 0, 32)
			if err != nil {
				panic(err)
			} else {
				re.NH.OIF = oif
			}
		}
	*/

	//fmt.Println(inh, err)
}

func (he *HostEntry) ParseINGNexthopInfo() {
	inh, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 ING_L3_NEXT_HOP %d %d", he.NextHopIndex, he.NextHopIndex),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(inh, getINGNexthopCopyToCPU); err != nil {
		panic(err)
	} else {
		if res == "0" {
			he.NH.CopyToCPU = false
		} else {
			he.NH.CopyToCPU = true
		}
	}

	if res, err := match(inh, getINGNexthopDrop); err != nil {
		panic(err)
	} else {
		if res == "0" {
			he.NH.Drop = false
		} else {
			he.NH.Drop = true
		}
	}

	if res, err := match(inh, getINGNexthopTGID); err != nil {
		panic(err)
	} else {
		tgid, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			he.NH.TGID = tgid
		}
	}

	if res, err := match(inh, getINGNexthopPNUM); err != nil {
		panic(err)
	} else {
		pnum, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			he.NH.DstPort = pnum
		}
	}

	if res, err := match(inh, getINGNexthopVLAN); err != nil {
		panic(err)
	} else {
		vlan, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			he.NH.VLAN = vlan
		}
	}

	/*
		if res, err := match(inh, getINGNexthopOIF); err != nil {
			panic(err)
		} else {
			oif, err := strconv.ParseInt(res, 0, 32)
			if err != nil {
				panic(err)
			} else {
				re.NH.OIF = oif
			}
		}
	*/

	//fmt.Println(inh, err)
}

func (re *RouteEntry) ParseEgrNexthopInfo() {
	enh, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 EGR_L3_NEXT_HOP %d %d", re.NextHopIndex, re.NextHopIndex),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(enh, getEGRNexthopMACAddress); err != nil {
		panic(err)
	} else {
		re.NH.DstMac = FixMACAddress(res)
	}

	if res, err := match(enh, getEGRIfNum); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		}
		re.NH.OIF.Index = index
	}

	//fmt.Println(enh, err)
}

func (he *HostEntry) ParseEgrNexthopInfo() {
	enh, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 EGR_L3_NEXT_HOP %d %d", he.NextHopIndex, he.NextHopIndex),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(enh, getEGRNexthopMACAddress); err != nil {
		panic(err)
	} else {
		he.NH.DstMac = FixMACAddress(res)
	}

	if res, err := match(enh, getEGRIfNum); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		}
		he.NH.OIF.Index = index
	}

	//fmt.Println(enh, err)
}

var getOIFMACAddress = regexp.MustCompile(`,MAC_ADDRESS=(?P<nmac>[0x]?[[:alnum:]]+)`)
var getOIFVID = regexp.MustCompile(`\<VID=(?P<vid>[0x]?[[:alnum:]]+)`)

func (nh *Nexthop) ParseOIF() {
	oif, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 EGR_L3_INTF  %d, %d", nh.OIF.Index, nh.OIF.Index),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(oif, getOIFMACAddress); err != nil {
		panic(err)
	} else {
		nh.OIF.MAC = FixMACAddress(res)
	}

	if res, err := match(oif, getOIFVID); err != nil {
		panic(err)
	} else {
		vlan, err := strconv.ParseInt(res, 0, 32)
		if err != nil {
			panic(err)
		} else {
			nh.OIF.Vid = vlan
		}
	}
}

func (re *RouteEntry) ParseNexthopInfo() {
	if re.NextHopIndex >= threshold.MaxNexthopIndex {
		//fmt.Printf("Skip the nexthop parse for %s due to Invalid nexthop index!", re)
		return
	}

	if !re.ECMP {
		re.NH.Index = re.NextHopIndex
		re.ParseINGNexthopInfo()
		re.ParseEgrNexthopInfo()
		re.NH.ParseOIF()
	} else {
		re.EG.Index = re.ECMPPTR
		re.ParseECMPGroup()
	}
}

func (he *HostEntry) ParseNexthopInfo() {
	if he.NextHopIndex >= threshold.MaxNexthopIndex {
		//fmt.Printf("Skip the nexthop parse for %s due to Invalid nexthop index!", re)
		return
	}

	if !he.ECMP {
		he.NH.Index = he.NextHopIndex
		he.ParseINGNexthopInfo()
		he.ParseEgrNexthopInfo()
		he.NH.ParseOIF()
	} else {
		he.EG.Index = he.ECMPPTR
		he.ParseECMPGroup()
	}
}

//L3_ECMP_GROUP.*[256]: <URPF_COUNT=1,RSVD_COUNT=0,RSVD_BASE_PTR=0,RH_FLOW_SET_SIZE=0,RH_FLOW_SET_BASE=0,RESERVED_0=0,L3_OIF_7_TYPE=0,L3_OIF_7=0,L3_OIF_6_TYPE=0,L3_OIF_6=0,L3_OIF_5_TYPE=0,L3_OIF_5=0,L3_OIF_4_TYPE=0,L3_OIF_4=0,L3_OIF_3_TYPE=0,L3_OIF_3=0,L3_OIF_2_TYPE=0,L3_OIF_2=0,L3_OIF_1_TYPE=0,L3_OIF_1=0x46,L3_OIF_0_TYPE=0,L3_OIF_0=0x50,EVEN_PARITY_1=0,EVEN_PARITY_0=0,ENHANCED_HASHING_ENABLE=0,ECMP_GT8=0,COUNT=1,BASE_PTR=0x1000>

var getECMPMemberCount = regexp.MustCompile(`\,COUNT=(?P<count>[0x]?[[:alnum:]]+)`)
var getECMPBasePTR = regexp.MustCompile(`,BASE_PTR=(?P<baseptr>[0x]?[[:alnum:]]+)`)

func (re *RouteEntry) ParseECMPGroup() {
	if re.ECMP == false {
		log.Println("Cannot parse ECMP for none ECMP entry")
		return
	}

	ecmpg, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_ECMP_GROUP %d %d", re.ECMPPTR, re.ECMPPTR),
	})

	if err != nil {
		panic(err)
	}

	//fmt.Println(ecmpg)

	if res, err := match(ecmpg, getECMPMemberCount); err != nil {
		panic(err)
	} else {
		mcount, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		re.EG.MemberCount = mcount + 1
	}

	if res, err := match(ecmpg, getECMPBasePTR); err != nil {
		panic(err)
	} else {
		base, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		re.EG.ECMPBasePTR = base
	}

	re.EG.ParseNexthop()
}

func (he *HostEntry) ParseECMPGroup() {
	if he.ECMP == false {
		log.Println("Cannot parse ECMP for none ECMP entry")
		return
	}

	ecmpg, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_ECMP_GROUP %d %d", he.ECMPPTR, he.ECMPPTR),
	})

	if err != nil {
		panic(err)
	}

	//fmt.Println(ecmpg)

	if res, err := match(ecmpg, getECMPMemberCount); err != nil {
		panic(err)
	} else {
		mcount, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		he.EG.MemberCount = mcount + 1
	}

	if res, err := match(ecmpg, getECMPBasePTR); err != nil {
		panic(err)
	} else {
		base, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		he.EG.ECMPBasePTR = base
	}

	he.EG.ParseNexthop()
}

type FIB struct {
	DB map[string]*RouteEntry
}

func IsValid(es string) bool {
	if es == "" {
		return false
	}

	if !strings.Contains(es, "VALID") ||
		!strings.Contains(es, "DST_DISCARD") ||
		!strings.Contains(es, "NEXT_HOP_INDEX") ||
		!strings.Contains(es, "KEY") ||
		!strings.Contains(es, "ECMP") ||
		!strings.Contains(es, "ECMP_PTR") {
		return false
	}

	return true
}

func FixMACAddress(s string) string {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	if len(s) == 0 {
		s = "000000000000"
	} else if len(s) == 1 {
		s = "00000000000" + s
	} else if len(s) == 2 {
		s = "0000000000" + s
	} else if len(s) == 3 {
		s = "000000000" + s
	} else if len(s) == 4 {
		s = "00000000" + s
	} else if len(s) == 5 {
		s = "0000000" + s
	} else if len(s) == 6 {
		s = "000000" + s
	} else if len(s) == 7 {
		s = "00000" + s
	} else if len(s) == 8 {
		s = "0000" + s
	} else if len(s) == 9 {
		s = "000" + s
	} else if len(s) == 10 {
		s = "00" + s
	} else if len(s) == 11 {
		s = "0" + s
	}

	f1, _ := strconv.ParseInt("0x"+s[:2], 0, 32)
	f2, _ := strconv.ParseInt("0x"+s[2:4], 0, 32)
	f3, _ := strconv.ParseInt("0x"+s[4:6], 0, 32)
	f4, _ := strconv.ParseInt("0x"+s[6:8], 0, 32)
	f5, _ := strconv.ParseInt("0x"+s[8:10], 0, 32)
	f6, _ := strconv.ParseInt("0x"+s[10:12], 0, 32)

	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", f1, f2, f3, f4, f5, f6)
	return mac
}

func FixIPv4Address(s string) net.IP {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	if len(s) == 0 {
		s = "00000000"
	} else if len(s) == 1 {
		s = "0000000" + s
	} else if len(s) == 2 {
		s = "000000" + s
	} else if len(s) == 3 {
		s = "00000" + s
	} else if len(s) == 4 {
		s = "0000" + s
	} else if len(s) == 5 {
		s = "000" + s
	} else if len(s) == 6 {
		s = "00" + s
	} else if len(s) == 7 {
		s = "0" + s
	}

	f1, _ := strconv.ParseInt("0x"+s[:2], 0, 32)
	f2, _ := strconv.ParseInt("0x"+s[2:4], 0, 32)
	f3, _ := strconv.ParseInt("0x"+s[4:6], 0, 32)
	f4, _ := strconv.ParseInt("0x"+s[6:8], 0, 32)

	return net.IPv4(byte(f1), byte(f2), byte(f3), byte(f4))
}

func FixIPv6Address(s string) net.IP {
	if strings.HasPrefix(s, "0x") {
		s = s[2:]
	}

	if len(s) == 17 {
		s = "000000000000000" + s
	} else if len(s) == 18 {
		s = "00000000000000" + s
	} else if len(s) == 19 {
		s = "0000000000000" + s
	} else if len(s) == 20 {
		s = "000000000000" + s
	} else if len(s) == 21 {
		s = "00000000000" + s
	} else if len(s) == 22 {
		s = "0000000000" + s
	} else if len(s) == 23 {
		s = "000000000" + s
	} else if len(s) == 24 {
		s = "00000000" + s
	} else if len(s) == 25 {
		s = "0000000" + s
	} else if len(s) == 26 {
		s = "000000" + s
	} else if len(s) == 27 {
		s = "00000" + s
	} else if len(s) == 28 {
		s = "0000" + s
	} else if len(s) == 29 {
		s = "000" + s
	} else if len(s) == 30 {
		s = "00" + s
	} else if len(s) == 31 {
		s = "0" + s
	} else if len(s) == 16 {
		s = s + "0000000000000000"
	} else if len(s) == 15 {
		s = "0" + s + "0000000000000000"
	} else if len(s) == 14 {
		s = "00" + s + "0000000000000000"
	} else if len(s) == 13 {
		s = "000" + s + "0000000000000000"
	} else if len(s) == 12 {
		s = "0000" + s + "0000000000000000"
	} else if len(s) == 11 {
		s = "00000" + s + "0000000000000000"
	} else if len(s) == 10 {
		s = "000000" + s + "0000000000000000"
	} else if len(s) == 9 {
		s = "0000000" + s + "0000000000000000"
	} else if len(s) == 8 {
		s = "00000000" + s + "0000000000000000"
	} else if len(s) == 7 {
		s = "000000000" + s + "0000000000000000"
	} else if len(s) == 6 {
		s = "0000000000" + s + "0000000000000000"
	} else if len(s) == 5 {
		s = "00000000000" + s + "0000000000000000"
	} else if len(s) == 4 {
		s = "000000000000" + s + "0000000000000000"
	} else if len(s) == 3 {
		s = "0000000000000" + s + "0000000000000000"
	} else if len(s) == 2 {
		s = "00000000000000" + s + "0000000000000000"
	} else if len(s) == 1 {
		s = "000000000000000" + s + "0000000000000000"
	}

	if len(s) != 32 {
		panic("Invalid IPv6 address to parse")
	}
	return net.ParseIP(s[:4] + ":" + s[4:8] + ":" + s[8:12] + ":" + s[12:16] + ":" + s[16:20] + ":" + s[20:24] + ":" + s[24:28] + ":" + s[28:32])
}

//L3_DEFIP_ALPM_IPV4.*[360456]: <VALID=1,SRC_DISCARD=0,RPE=0,RESERVED_ECMP_PTR=0,RESERVED_0=0,PRI=0,NEXT_HOP_INDEX=1,LENGTH=0x18,KEY=0x46000000,HIT=0,EVEN_PARITY=0,ENTRY_ONLY=0x26118000002,ECMP_PTR=1,ECMP=0,DST_DISCARD=0,DEFAULTROUTE=0,DATA=2,CLASS_ID=0>
func ParseRouteEntryString(es string, af int) (*RouteEntry, error) {
	if !IsValid(es) {
		return nil, errors.New("Invalid input string: " + es)
	}

	var Entry RouteEntry
	if res, err := match(es, getEntryIndex); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.Index = index
	}
	if res, err := match(es, getValidBit); err != nil {
		panic(err)
	} else {
		Entry.Valid = res
	}

	if res, err := match(es, getNextHopIndex); err != nil {
		panic(err)
	} else {
		nhi, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		} else {
			Entry.NextHopIndex = nhi
		}
	}

	if res, err := match(es, getLength); err != nil {
		panic(err)
	} else {
		length, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.Length = length
	}

	if res, err := match(es, getKey); err != nil {
		panic(err)
	} else {
		if af == IPV4 {
			Entry.AF = af
			Entry.Key = FixIPv4Address(res)
		} else if af == IPV6 {
			Entry.AF = af
			Entry.Key = FixIPv6Address(res)
		} else {
			panic("Unknown Address family")
		}
	}

	if res, err := match(es, getHitBit); err != nil {
		panic(err)
	} else {
		if res == "1" {
			Entry.Hit = true
		} else {
			Entry.Hit = false
		}
	}

	if res, err := match(es, getECMPPTR); err != nil {
		panic(err)
	} else {
		ptr, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.ECMPPTR = ptr
	}

	if res, err := match(es, getECMP); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.ECMP = false
		} else {
			Entry.ECMP = true
		}
	}

	if res, err := match(es, getSrcDiscard); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.SrcDiscard = false
		} else {
			Entry.SrcDiscard = true
		}
	}

	if res, err := match(es, getDstDiscard); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.DstDiscard = false
		} else {
			Entry.DstDiscard = true
		}
	}

	if res, err := match(es, getDefaultRoute); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.DefaultRoute = false
		} else {
			Entry.DefaultRoute = true
		}
	}

	Entry.Validate()

	return &Entry, nil
}

func match(s string, r *regexp.Regexp) (string, error) {
	matches := r.FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1], nil
	}

	return "", errors.New("Cannot match for string: " + s + " Re: " + r.String())
}

var IP = flag.String("ip", "10.71.20.115", "IP address of the remote device")
var Host = flag.String("hostname", "V8500", "Host name of the remote device")
var User = flag.String("username", "admin", "Username of the remote device")
var Password = flag.String("password", "", "Passwrod of the remote device")
var SFU = flag.String("sfu", "A", "SFU (A/B)")
var Table = flag.String("table", "all", "Address family to dump (IPv4/IPv6/all/arp/nd)")

func main() {
	flag.Parse()

	ip := net.ParseIP(*IP)
	if ip == nil {
		fmt.Printf("Invalid IP address: %s\n", *IP)
		return
	}

	if *Host == "" {
		fmt.Println("Invalid Host name")
		return
	}

	if *User == "" {
		fmt.Println("Invalidusername")
		return
	}

	if *SFU != "A" && *SFU != "B" {
		fmt.Printf("Invalid SFU: %s\n", *SFU)
		return
	}

	if *Table != "ipv4" && *Table != "ipv6" && *Table != "all" && *Table != "arp" && *Table != "nd" {
		fmt.Println("Invalid Address family to dump")
		return
	}

	dev, err := prepare(&rut.RUT{
		Name:     "DUT1",
		Device:   "V8",
		IP:       *IP,
		Port:     "23",
		Username: *User,
		Hostname: *Host,
		Password: *Password,
		SFU:      *SFU,
	})

	if err != nil {
		fmt.Printf("Cannot connect to :%s with error: %s", *IP, err.Error())
		return
	}

	if *Table == "ipv4" {
		DumpIPv4Entry(dev)
	} else if *Table == "ipv6" {
		DumpIPv6Entry(dev)
	} else if *Table == "all" {
		DumpIPv4Entry(dev)
		DumpIPv6Entry(dev)
	} else if *Table == "arp" {
		DumpIPv4HostEntry(dev)
	} else if *Table == "nd" {
		DumpIPv6HostEntry(dev)
	}
}

//L3_ENTRY_IPV4_UNICAST.*[3456]: <VRF_ID=0,VALID=1,TRILL:TREE_ID=0,TRILL:RESERVED_104_43=0x800000,TRILL:KEY=0x20,TRILL:INGRESS_RBRIDGE_NICKNAME=1,TRILL:HASH_LSB=1,TRILL:EXPECTED_TGID=0x60,TRILL:EXPECTED_T=0,TRILL:EXPECTED_PORT_NUM=0x60,TRILL:EXPECTED_MODULE_ID=8,TRILL:DATA=0x460,RSVD_VRF_ID=0,RSVD_NEXT_HOP_INDEX=0,RPE=0,RMEP:RESERVED_104_83=0,RMEP:HASH_LSB=0,RESERVED_104_82=0,PRI=0,NEXT_HOP_INDEX=8,LOCAL_ADDRESS=0,LMEP:HASH_LSB=0,KEY_TYPE=0,KEY=0x8c0000020,IP_ADDR=0x46000001,IPV4UC:VRF_ID=0,IPV4UC:RSVD_VRF_ID=0,IPV4UC:RSVD_NEXT_HOP_INDEX=0,IPV4UC:RPE=0,IPV4UC:RESERVED_104_82=0,IPV4UC:PRI=0,IPV4UC:NEXT_HOP_INDEX=8,IPV4UC:LOCAL_ADDRESS=0,IPV4UC:KEY=0x8c0000020,IPV4UC:IP_ADDR=0x46000001,IPV4UC:HASH_LSB=1,IPV4UC:ECMP_PTR=8,IPV4UC:ECMP=0,IPV4UC:DST_DISCARD=0,IPV4UC:DATA=0x10000,IPV4UC:CLASS_ID=0,IPV4UC:BFD_ENABLE=0,HIT=1,HASH_LSB=1,FCOE:VRF_ID=0x46,FCOE:S_ID=1,FCOE:RSVD_VRF_ID=0,FCOE:RSVD_NEXT_HOP_INDEX=0,FCOE:RPE=0,FCOE:RESERVED_ECMP_PTR0=2,FCOE:RESERVED_104_73=0,FCOE:PRI=0,FCOE:NEXT_HOP_INDEX=0x800,FCOE:MASKED_D_ID=1,FCOE:LOCAL_ADDRESS=0,FCOE:KEY=0x8c0000020,FCOE:HASH_LSB=1,FCOE:ECMP_PTR0=0,FCOE:ECMP=0,FCOE:D_ID=1,FCOE:DST_DISCARD=0,FCOE:DATA=0x1000000,FCOE:CLASS_ID=0,EVEN_PARITY=0,ECMP_PTR=8,ECMP=0,DUMMY_V6=0,DUMMY_IPMC=0,DUMMY=0,DST_DISCARD=0,DATA=0x10000,CLASS_ID=0,BFD_ENABLE=0>

var getHostEntryNextHopIndex = regexp.MustCompile(`,IPV4UC:NEXT_HOP_INDEX=(?P<nhi>[0x]?[[:alnum:]]+)`)
var getIPv6HostEntryNextHopIndex = regexp.MustCompile(`,IPV6UC:NEXT_HOP_INDEX=(?P<nhi>[0x]?[[:alnum:]]+)`)
var getHostEntryIPv4Address = regexp.MustCompile(`,IPV4UC:IP_ADDR=(?P<key>[0x]?[[:alnum:]]+)`)
var getIPv6HostEntryAddressUpper = regexp.MustCompile(`,IPV6UC:IP_ADDR_UPR_64=(?P<key>[0x]?[[:alnum:]]+)`)
var getIPv6HostEntryAddressLower = regexp.MustCompile(`,IPV6UC:IP_ADDR_LWR_64=(?P<key>[0x]?[[:alnum:]]+)`)
var getHostEntryECMPPTR = regexp.MustCompile(`,IPV4UC:ECMP_PTR=(?P<ecmpptr>[0x]?[[:alnum:]]+)`)
var getIPv6HostEntryECMPPTR = regexp.MustCompile(`,IPV6UC:ECMP_PTR=(?P<ecmpptr>[0x]?[[:alnum:]]+)`)
var getHostEntryECMP = regexp.MustCompile(`,IPV4UC:ECMP=(?P<ecmpptr>[0x]?[[:alnum:]]+)`)
var getIPv6HostEntryECMP = regexp.MustCompile(`,IPV6UC:ECMP=(?P<ecmpptr>[0x]?[[:alnum:]]+)`)
var getHostEntryDstDiscard = regexp.MustCompile(`,IPV4UC:DST_DISCARD=(?P<dstdis>[0-9]+)`)
var getIPv6HostEntryDstDiscard = regexp.MustCompile(`,IPV6UC:DST_DISCARD=(?P<dstdis>[0-9]+)`)
var getHostEntryKeyType = regexp.MustCompile(`,KEY_TYPE=(?P<keytype>[0-9]+)`)
var getIPv6HostEntryKeyType = regexp.MustCompile(`,KEY_TYPE_1=(?P<keytype>[0-9]+)`)
var getIPv6HostEntryValidBit = regexp.MustCompile(`VALID_1=(?P<valid>[0-9]+)`)
var getIPv6HostEntryHitBit = regexp.MustCompile(`,HIT_1=(?P<hit>[0-9]+)`)

func ParseHostEntryString(es string, af int) (*HostEntry, error) {
	if !IsValid(es) {
		return nil, errors.New("Invalid input string: " + es)
	}

	var Entry HostEntry
	if af == IPV4 {
		Entry.Length = 32
	} else {
		Entry.Length = 128
	}

	if res, err := match(es, getEntryIndex); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.Index = index
	}
	if res, err := match(es, getValidBit); err != nil {
		panic(err)
	} else {
		Entry.Valid = res
	}

	if res, err := match(es, getHostEntryNextHopIndex); err != nil {
		panic(err)
	} else {
		nhi, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		} else {
			Entry.NextHopIndex = nhi
		}
	}

	if res, err := match(es, getHostEntryKeyType); err != nil {
		panic(err)
	} else {
		kt, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		} else {
			Entry.KeyType = kt
		}
	}

	if res, err := match(es, getHostEntryIPv4Address); err != nil {
		panic(err)
	} else {
		if af == IPV4 {
			Entry.AF = af
			Entry.IP = FixIPv4Address(res)
		} else if af == IPV6 {
			Entry.AF = af
			Entry.IP = FixIPv6Address(res)
		} else {
			panic("Unknown Address family")
		}
	}

	if res, err := match(es, getHitBit); err != nil {
		panic(err)
	} else {
		if res == "1" {
			Entry.Hit = true
		} else {
			Entry.Hit = false
		}
	}

	if res, err := match(es, getHostEntryECMPPTR); err != nil {
		panic(err)
	} else {
		ptr, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.ECMPPTR = ptr
	}

	if res, err := match(es, getHostEntryECMP); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.ECMP = false
		} else {
			Entry.ECMP = true
		}
	}

	if res, err := match(es, getHostEntryDstDiscard); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.DstDiscard = false
		} else {
			Entry.DstDiscard = true
		}
	}

	return &Entry, nil
}

func ParseIPv6HostEntryString(es string, af int) (*HostEntry, error) {
	if !IsValid(es) {
		return nil, errors.New("Invalid input string: " + es)
	}

	var Entry HostEntry
	if af == IPV4 {
		Entry.Length = 32
	} else {
		Entry.Length = 128
	}

	if res, err := match(es, getEntryIndex); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.Index = index
	}
	if res, err := match(es, getIPv6HostEntryValidBit); err != nil {
		panic(err)
	} else {
		Entry.Valid = res
	}

	if res, err := match(es, getIPv6HostEntryNextHopIndex); err != nil {
		panic(err)
	} else {
		nhi, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		} else {
			Entry.NextHopIndex = nhi
		}
	}

	if res, err := match(es, getIPv6HostEntryKeyType); err != nil {
		panic(err)
	} else {
		kt, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		} else {
			Entry.KeyType = kt
		}
	}

	if upper, err := match(es, getIPv6HostEntryAddressUpper); err != nil {
		panic(err)
	} else {
		if lower, err := match(es, getIPv6HostEntryAddressLower); err != nil {
			panic(err)
		} else {
			if af == IPV6 {
				Entry.AF = af
				Entry.IP = FixIPv6Address(MakeIPv6Address(upper, lower))
			} else {
				panic("Unknown Address family")
			}
		}
	}

	if res, err := match(es, getIPv6HostEntryHitBit); err != nil {
		panic(err)
	} else {
		if res == "1" {
			Entry.Hit = true
		} else {
			Entry.Hit = false
		}
	}

	if res, err := match(es, getIPv6HostEntryECMPPTR); err != nil {
		panic(err)
	} else {
		ptr, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		Entry.ECMPPTR = ptr
	}

	if res, err := match(es, getIPv6HostEntryECMP); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.ECMP = false
		} else {
			Entry.ECMP = true
		}
	}

	if res, err := match(es, getIPv6HostEntryDstDiscard); err != nil {
		panic(err)
	} else {
		if res == "0" {
			Entry.DstDiscard = false
		} else {
			Entry.DstDiscard = true
		}
	}

	return &Entry, nil
}

func MakeIPv6Address(upper, lower string) string {
	if strings.HasPrefix(lower, "0x") {
		lower = lower[2:]
	}

	if len(lower) == 15 {
		lower = "0" + lower
	} else if len(lower) == 14 {
		lower = "00" + lower
	} else if len(lower) == 13 {
		lower = "000" + lower
	} else if len(lower) == 12 {
		lower = "0000" + lower
	} else if len(lower) == 11 {
		lower = "00000" + lower
	} else if len(lower) == 10 {
		lower = "00000" + lower
	} else if len(lower) == 9 {
		lower = "000000" + lower
	} else if len(lower) == 8 {
		lower = "0000000" + lower
	} else if len(lower) == 8 {
		lower = "00000000" + lower
	} else if len(lower) == 7 {
		lower = "000000000" + lower
	} else if len(lower) == 6 {
		lower = "0000000000" + lower
	} else if len(lower) == 5 {
		lower = "00000000000" + lower
	} else if len(lower) == 4 {
		lower = "000000000000" + lower
	} else if len(lower) == 3 {
		lower = "0000000000000" + lower
	} else if len(lower) == 2 {
		lower = "00000000000000" + lower
	} else if len(lower) == 1 {
		lower = "000000000000000" + lower
	} else if len(lower) == 0 {
		lower = "0000000000000000" + lower
	}

	return upper + lower
}

func DumpIPv4HostEntry(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_ENTRY_IPV4_UNICAST  0 16383  | grep VALID=1",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		h, _ := ParseHostEntryString(l, IPV4)
		if h != nil {
			h.ParseNexthopInfo()
		}
		if h != nil && h.AF == IPV4 && h.KeyType == 0 {
			fmt.Println(h)
		}
	}
}

func DumpIPv6HostEntry(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_ENTRY_IPV6_UNICAST 0 8191  | grep VALID_1=1",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		h, _ := ParseIPv6HostEntryString(l, IPV6)
		if h != nil {
			h.ParseNexthopInfo()
		}
		fmt.Println(h)
	}
}

func DumpIPv4Entry(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP_ALPM_IPV4 0 393215  | grep VALID=1",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		r, _ := ParseRouteEntryString(l, IPV4)
		if r != nil {
			r.ParseNexthopInfo()
		}
		fmt.Println(r)
	}
}

func DumpIPv6Entry(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP_ALPM_IPV6_64 0 262143 | grep VALID=1",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		r, _ := ParseRouteEntryString(l, IPV6)
		if r != nil {
			r.ParseNexthopInfo()
		}
		fmt.Println(r)
	}

	res, err = dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP_ALPM_IPV6_128 0 131071 | grep VALID=1",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		r, _ := ParseRouteEntryString(l, IPV6)
		if r != nil {
			r.ParseNexthopInfo()
		}
		fmt.Println(r)
	}
}

func prepare(r *rut.RUT) (*rut.RUT, error) {
	dev, err := rut.New(r)
	if err != nil {
		return nil, err
	}

	err = dev.Init()
	if err != nil {
		return nil, err
	}

	Dev = dev

	_, err = dev.RunCommand(CTX, &command.Command{
		Mode: "normal",
		CMD:  "terminal length 0",
	})

	if err != nil {
		return nil, err
	}

	_, err = dev.RunCommand(CTX, &command.Command{
		Mode: "normal",
		CMD:  "configure terminal",
	})
	if err != nil {
		return nil, err
	}

	_, err = dev.RunCommand(CTX, &command.Command{
		Mode: "config",
		CMD:  "do q sh -l",
	})
	if err != nil {
		return nil, err
	}

	_, err = dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  "ls -al",
	})

	if err != nil {
		return nil, err
	}

	return dev, nil
}
