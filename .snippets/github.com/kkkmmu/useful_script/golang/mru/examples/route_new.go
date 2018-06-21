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

type ALPMBucket struct {
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

func DumpECMPGroupByIndex(index int64) (ECMPGroup, error) {
	var EG = ECMPGroup{}
	ecmpg, err := Dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf("scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_ECMP_GROUP %d %d", index, index),
	})

	if err != nil {
		panic(err)
	}

	if res, err := match(ecmpg, getECMPMemberCount); err != nil {
		panic(err)
	} else {
		mcount, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		EG.MemberCount = mcount + 1
	}

	if res, err := match(ecmpg, getECMPBasePTR); err != nil {
		panic(err)
	} else {
		base, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		EG.ECMPBasePTR = base
	}

	EG.ParseNexthop()

	return EG, nil
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
	if es == "" || strings.Contains(es, "scontrol") {
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

func FixIPv4NetMask(s string) net.IPMask {
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

	return net.IPv4Mask(byte(f1), byte(f2), byte(f3), byte(f4))
}

func FixIPv6AddressDEFIP(s string) net.IP {
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
		s = s + "0" + "0000000000000000"
	} else if len(s) == 14 {
		s = s + "00" + "0000000000000000"
	} else if len(s) == 13 {
		s = s + "000" + "0000000000000000"
	} else if len(s) == 12 {
		s = s + "0000" + "0000000000000000"
	} else if len(s) == 11 {
		s = s + "00000" + "0000000000000000"
	} else if len(s) == 10 {
		s = s + "000000" + "0000000000000000"
	} else if len(s) == 9 {
		s = s + "0000000" + "0000000000000000"
	} else if len(s) == 8 {
		s = s + "00000000" + "0000000000000000"
	} else if len(s) == 7 {
		s = s + "000000000" + "0000000000000000"
	} else if len(s) == 6 {
		s = s + "0000000000" + "0000000000000000"
	} else if len(s) == 5 {
		s = s + "00000000000" + "0000000000000000"
	} else if len(s) == 4 {
		s = s + "000000000000" + "0000000000000000"
	} else if len(s) == 3 {
		s = s + "0000000000000" + "0000000000000000"
	} else if len(s) == 2 {
		s = s + "00000000000000" + "0000000000000000"
	} else if len(s) == 1 {
		s = s + "000000000000000" + "0000000000000000"
	}

	if len(s) != 32 {
		panic("Invalid IPv6 address to parse")
	}
	return net.ParseIP(s[:4] + ":" + s[4:8] + ":" + s[8:12] + ":" + s[12:16] + ":" + s[16:20] + ":" + s[20:24] + ":" + s[24:28] + ":" + s[28:32])
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

	//fmt.Println(s, len(s))

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

var DEFIPModeValueToString = map[int64]string{
	0: "IPv4",
	1: "IPv6(64)",
	3: "IPv6(128)",
}

var IP = flag.String("ip", "10.71.20.115", "IP address of the remote device")
var Host = flag.String("hostname", "V8500", "Host name of the remote device")
var User = flag.String("username", "admin", "Username of the remote device")
var Password = flag.String("password", "", "Passwrod of the remote device")
var SFU = flag.String("sfu", "A", "SFU (A/B)")
var Table = flag.String("table", "all", "Address family to dump (v4/v6/all/arp/nd/alpm/alpm4/alpm6/defip)")

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

	if *Table != "v4" && *Table != "v6" && *Table != "all" && *Table != "arp" && *Table != "nd" && *Table != "alpm" && *Table != "alpm4" && *Table != "alpm6" && *Table != "defip" {
		fmt.Println("Valid table is: v4/v6/all/arp/nd/alpm/alpm4/alpm6/defip")
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

	if *Table == "v4" {
		DumpIPv4Entry(dev)
	} else if *Table == "v6" {
		DumpIPv6Entry(dev)
	} else if *Table == "all" {
		DumpIPv4Entry(dev)
		DumpIPv6Entry(dev)
	} else if *Table == "arp" {
		DumpIPv4HostEntry(dev)
	} else if *Table == "nd" {
		DumpIPv6HostEntry(dev)
	} else if *Table == "defip" {
		DumpL3DEFIPEntry(dev)
	} else if *Table == "alpm" {
		DumpL3DEFIPALPMEntry(dev)
	} else if *Table == "alpm4" {
		DumpL3IPv4DEFIPEntry(dev)
	} else if *Table == "alpm6" {
		DumpL3IPv664DEFIPEntry(dev)
	} else {
		fmt.Println("Invalid table: ", *Table, " to dump")
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

func DumpIPv4EntryByIndex(dev *rut.RUT, index int) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  fmt.Sprintf(" scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP_ALPM_IPV4 %d %d", index, index),
	})

	if err != nil {
		panic(err)
	}

	if !strings.Contains(res, "VALID=1") {
		fmt.Printf("Entry %d is invalid\n", index)
		return
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

type L3DEFIPHalfEntry struct {
	AF              int
	Index           int64
	VRF             int64
	Valid           bool
	SrcDiscard      bool
	DstDiscard      bool
	RPE             int64
	ReservedECMPPtr int64
	Priority        int64
	NexthopIndex    int64
	NH              Nexthop
	EG              ECMPGroup
	Mode            int64
	//Mask            int64
	//Key             int64
	IPAddrMask    net.IPMask
	IPAddrMaskLen int
	IPAddr        net.IP
	Hit           bool
	GlobalRoute   bool
	EntryType     int64
	ECMPPtr       int64
	ECMP          bool
	DefaultMiss   bool
	DefaultRoute  bool
	ALGHitIndex   int64
	ALGBktPtr     int64
}

func (re *L3DEFIPHalfEntry) String() string {
	var Out *color.Color
	if re.Hit {
		Out = color.New(color.FgGreen)
	} else {
		Out = color.New(color.FgRed)
	}
	if re.AF == IPV4 {
		if re.IPAddrMaskLen > 32 || re.NexthopIndex > threshold.MaxNexthopIndex {
			if re.IPAddrMaskLen > 32 {
				return Out.Sprintf("[%6d]: %39s/%-3d >> %20s", re.Index, re.IPAddr, re.IPAddrMaskLen, "is not a valid IPv4 Address")
			} else {
				return Out.Sprintf("[%6d]: %39s/%-3d >> has a invalid nexthop index: %d", re.Index, re.IPAddr, re.IPAddrMaskLen, re.NexthopIndex)
			}
		} else if !re.ECMP {
			return Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> NH[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", re.Index, re.IPAddr, re.IPAddrMaskLen, re.DstDiscard, re.SrcDiscard, re.NexthopIndex, re.NH.OIF.MAC, re.NH.DstMac, re.NH.OIF.Index, re.NH.OIF.Vid, re.NH.DstPort)
		} else {
			base := Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> is  ECMP Route, ECMP_PTR: %5d, BASE_PTR: %5d, MemberCount: %2d", re.Index, re.IPAddr, re.IPAddrMaskLen, re.DstDiscard, re.SrcDiscard, re.ECMPPtr, re.EG.ECMPBasePTR, re.EG.MemberCount)
			for i := 0; i < len(re.EG.Member); i++ {
				base += "\n"
				base += Out.Sprintf("%73s[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", "NH", re.EG.Member[i].Index, re.EG.Member[i].OIF.MAC, re.EG.Member[i].DstMac, re.EG.Member[i].OIF.Index, re.EG.Member[i].OIF.Vid, re.EG.Member[i].DstPort)
			}

			return base
		}
	} else if re.AF == IPV6 {
		if re.IPAddrMaskLen > 128 || re.NexthopIndex > threshold.MaxNexthopIndex {
			if re.IPAddrMaskLen > 128 {
				return Out.Sprintf("[%6d]: %39s/%-3d >> %20s", re.Index, re.IPAddr, re.IPAddrMaskLen, "is not a valid IPv6 Address")
			} else {
				return Out.Sprintf("[%6d]: %39s/%-3d >> has a invalid nexthop index: %d", re.Index, re.IPAddr, re.IPAddrMaskLen, re.NexthopIndex)
			}
		} else if !re.ECMP {
			return Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> NH[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", re.Index, re.IPAddr, re.IPAddrMaskLen, re.DstDiscard, re.SrcDiscard, re.NexthopIndex, re.NH.OIF.MAC, re.NH.DstMac, re.NH.OIF.Index, re.NH.OIF.Vid, re.NH.DstPort)
		} else {
			base := Out.Sprintf("[%6d]: %39s/%-3d (%5t|%5t) >> is  ECMP Route, ECMP_PTR: %5d, BASE_PTR: %5d, MemberCount: %2d", re.Index, re.IPAddr, re.IPAddrMaskLen, re.DstDiscard, re.SrcDiscard, re.ECMPPtr, re.EG.ECMPBasePTR, re.EG.MemberCount)
			for i := 0; i < len(re.EG.Member); i++ {
				base += "\n"
				base += Out.Sprintf("%73s[%5d]:{SMAC: %17s DMAC: %17s OIF: %4d VID: %4d DPORT: %3d}", "NH", re.EG.Member[i].Index, re.EG.Member[i].OIF.MAC, re.EG.Member[i].DstMac, re.EG.Member[i].OIF.Index, re.EG.Member[i].OIF.Vid, re.EG.Member[i].DstPort)
			}

			return base
		}
	}
	return fmt.Sprintf("Invalid route entry: %s/%d\n", re.IPAddr, re.IPAddrMaskLen)

	//	return ""

}

var HalfEntryFields = []string{
	"Index",
	"VRF_ID_",
	"VALID",
	"SRC_DISCARD",
	"DST_DISCARD",
	"RPE",
	"PRI",
	"RESERVED_ECMP_PTR",
	"NEXT_HOP_INDEX",
	",MODE",
	//"MASK",
	//"KEY",
	"IP_ADDR_MASK",
	"IP_ADDR",
	"HIT",
	"GLOBAL_ROUTE",
	"ENTRY_TYPE",
	",ECMP_PTR",
	"ECMP",
	"DEFAULT_MISS",
	"DEFAULT_ROUTE",
	"ALG_HIT_IDX",
	"ALG_BKT_PTR",
}

var ip1r = regexp.MustCompile("IP_ADDR1=(?P<f>[0x]?[[:alnum:]]+)")
var ip0r = regexp.MustCompile("IP_ADDR0=(?P<f>[0x]?[[:alnum:]]+)")

func ParseIPv6BestPrefixRouteFromDEFIP(entry string, index int) (net.IP, error) {
	var ip1 string
	var ip0 string
	if matches := ip1r.FindStringSubmatch(entry); len(matches) == 2 {
		ip1 = matches[1]
		if strings.HasPrefix(ip1, "0x") {
			ip1 = ip1[2:]
		}
	}

	if matches := ip0r.FindStringSubmatch(entry); len(matches) == 2 {
		ip0 = matches[1]
		if strings.HasPrefix(ip0, "0x") {
			ip0 = ip0[2:]
		}
	}

	ip := ip1 + ip0

	nip := FixIPv6AddressDEFIP(ip)

	return nip, nil
}

var ipmr0 = regexp.MustCompile("IP_ADDR_MASK0=(?P<f>[0x]?[[:alnum:]]+)")
var ipmr1 = regexp.MustCompile("IP_ADDR_MASK1=(?P<f>[0x]?[[:alnum:]]+)")

var HexToBin = map[rune]string{
	'0': "0000",
	'1': "0001",
	'2': "0010",
	'3': "0011",
	'4': "0100",
	'5': "0101",
	'6': "0110",
	'7': "0111",
	'8': "1000",
	'9': "1001",
	'a': "1010",
	'b': "1011",
	'c': "1100",
	'd': "1101",
	'e': "1110",
	'f': "1111",
}

func GetLeadingOnes(entry string) int {
	var res string
	for _, c := range entry {
		if c != 'f' && c != 'e' && c != 'c' && c != '8' && c != '0' {
			panic("Invalid mask")
		}
		res += HexToBin[c]
	}

	var count int
	for _, c := range res {
		if c == '1' {
			count++
		}
	}

	return count
}

func ParseIPv6BestPrefixRouteMaskLengthFromDEFIP(entry string, index int) (int, error) {
	var ipm1 string
	var ipm0 string
	if matches := ipmr1.FindStringSubmatch(entry); len(matches) == 2 {
		ipm1 = matches[1]
		if strings.HasPrefix(ipm1, "0x") {
			ipm1 = ipm1[2:]
		}
	}

	if matches := ipmr1.FindStringSubmatch(entry); len(matches) == 2 {
		ipm0 = matches[1]
		if strings.HasPrefix(ipm0, "0x") {
			ipm0 = ipm0[2:]
		}
	}

	ipm := ipm1 + ipm0

	return GetLeadingOnes(ipm), nil

}

func DumpL3DEFIPHalfEntry(entry string, index int) (*L3DEFIPHalfEntry, error) {
	fieldsReg := make(map[string]*regexp.Regexp, len(HalfEntryFields))
	for _, f := range HalfEntryFields {
		fieldsReg[f] = regexp.MustCompile(fmt.Sprintf("%s%d=(?P<f>[0x]?[[:alnum:]]+)", f, index))
	}
	/*
		var getHalfEntryVRF = regexp.MustCompile(fmt.Sprintf("VRF_ID_%d=[0x]*(?P<vrf>[0-9]+)", index))
		var getHalfEntryValid = regexp.MustCompile(fmt.Sprintf("VALID%d=(?P<valid>[0-9]+)", index))
		var getHalfEntrySrcDiscard = regexp.MustCompile(fmt.Sprintf("SRC_DISCARD%d=(?P<sd>[0-9]+)", index))
		var getHalfEntryDstDiscard = regexp.MustCompile(fmt.Sprintf("DST_DISCARD%d=(?P<dd>[0-9]+)", index))
		var getHalfEntryRPE = regexp.MustCompile(fmt.Sprintf("RPE%d=(?P<rpe>[[:alnum:]]+)", index))
		var getHalfEntryReservedECMPPtr = regexp.MustCompile(fmt.Sprintf("RESERVED_ECMP_PTR%d=[0x]*(?P<rep>[[:alnum:]]+)", index))
		var getHalfEntryPriority = regexp.MustCompile(fmt.Sprintf("PRI%d=[0x]*(?P<pri>[0-9]+)", index))
		var getHalfEntryNexthopIndex = regexp.MustCompile(fmt.Sprintf("NEXT_HOP_INDEX%d=[0x]?(?P<nhi>[[:alnum:]]+)", index))
		var getHalfEntryMode = regexp.MustCompile(fmt.Sprintf("MODE%d=[0x]?(?P<mode>([[:alnum:]]+)", index))
		var getHalfEntryMask = regexp.MustCompile(fmt.Sprintf("MASK%d=[0x]?(?P<mask>([[:alnum:]]+)", index))
		var getHalfEntryKey = regexp.MustCompile(fmt.Sprintf("KEY%d=[0x]?(?P<key>[[:alnum:]]+)", index))
		var getHalfEntryIPAddrMask = regexp.MustCompile(fmt.Sprintf("IP_ADDR_MASK%d=[0x]?(?P<mask>([[:alnum:]]+)", index))
		var getHalfEntryIPAddr = regexp.MustCompile(fmt.Sprintf("IP_ADDR%d=[0x]?(?P<ip>([[:alnum:]]+)", index))
		var getHalfEntryHit = regexp.MustCompile(fmt.Sprintf("HIT%d=(?P<hit>([[:alnum:]]+)", index))
		var getHalfEntryGlobalRoute = regexp.MustCompile(fmt.Sprintf("GLOBAL_ROUTE%d=(?P<gr>([[:alnum:]]+)", index))
		var getHalfEntryEntryType = regexp.MustCompile(fmt.Sprintf("ENTRY_TYPE%d=[0x]?(?P<et>([[:alnum:]]+)", index))
		var getHalfEntryECMPPtr = regexp.MustCompile(fmt.Sprintf("ECMP_PTR%d=[0x]?(?P<ept>([[:alnum:]]+)", index))
		var getHalfEntryECMP = regexp.MustCompile(fmt.Sprintf("ECMP%d=(?P<ecmp>([[:alnum:]]+)", index))
		var getHalfEntryDefaultMiss = regexp.MustCompile(fmt.Sprintf("DEFAULT_MISS%d=(?P<dfm>([[:alnum:]]+)", index))
		var getHalfEntryDefaultRoute = regexp.MustCompile(fmt.Sprintf("DEFAULT_ROUTE%d=(?P<dfr>([[:alnum:]]+)", index))
		var getHalfEntryALGHitIndex = regexp.MustCompile(fmt.Sprintf("ALG_HIT_IDX%d=[0x]?(?P<ahi>([[:alnum:]]+)", index))
		var getHalfEntryALGBktPtr = regexp.MustCompile(fmt.Sprintf("ALG_BKT_PTR%d=[0x]?(?P<abp>([[:alnum:]]+)", index))
	*/
	en := &L3DEFIPHalfEntry{}

	if res, err := match(entry, getEntryIndex); err != nil {
		panic(err)
	} else {
		index, err := strconv.ParseInt(res, 0, 64)
		if err != nil {
			panic(err)
		}
		en.Index = index
	}

	for k, v := range fieldsReg {
		if matches := v.FindStringSubmatch(entry); len(matches) == 2 {
			//fmt.Printf("%20s : %10s\n", k, matches[1])
			switch k {
			case "VRF_ID_":
				vrf, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.VRF = vrf
			case "VALID":
				valid, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				if valid == 1 {
					en.Valid = true
				} else {
					en.Valid = false
				}
			case "SRC_DISCARD":
				sd, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				if sd == 1 {
					en.SrcDiscard = true
				} else {
					en.SrcDiscard = false
				}
			case "DST_DISCARD":
				dd, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				if dd == 1 {
					en.DstDiscard = true
				} else {
					en.DstDiscard = false
				}
			case "RPE":
				rpe, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.RPE = rpe
			case "PRI":
				pri, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.Priority = pri
			case "RESERVED_ECMP_PTR":
				rep, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.ReservedECMPPtr = rep
				/*
					case "NEXT_HOP_INDEX":
						nhi, err := strconv.ParseInt(matches[1], 0, 64)
						if err != nil {
							panic(err)
						}
						en.NexthopIndex = nhi
						nh, err := ParseNexthopByIndex(en.NexthopIndex)
						if err != nil {
							panic(err)
						}

						en.NH = *nh
				*/
			case ",MODE":
				mode, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.Mode = mode

				var ipr *regexp.Regexp
				if index == 0 {
					ipr = regexp.MustCompile("IP_ADDR0=(?P<f>[0x]?[[:alnum:]]+)")

				} else {
					ipr = regexp.MustCompile("IP_ADDR1=(?P<f>[0x]?[[:alnum:]]+)")

				}
				if ms := ipr.FindStringSubmatch(entry); len(ms) == 2 {
					if en.Mode == 0 {
						en.AF = IPV4
						en.IPAddr = FixIPv4Address(ms[1])
					} else {
						en.AF = IPV6
						en.IPAddr, _ = ParseIPv6BestPrefixRouteFromDEFIP(entry, index)
					}
				}

				//ip mask
				var ipmr *regexp.Regexp
				if index == 0 {
					ipmr = regexp.MustCompile("IP_ADDR_MASK0=(?P<f>[0x]?[[:alnum:]]+)")

				} else {
					ipmr = regexp.MustCompile("IP_ADDR_MASK1=(?P<f>[0x]?[[:alnum:]]+)")

				}

				if ms := ipmr.FindStringSubmatch(entry); len(ms) == 2 {
					if en.Mode == 0 {
						en.IPAddrMask = FixIPv4NetMask(ms[1])
						en.IPAddrMaskLen, _ = en.IPAddrMask.Size()
					} else {
						en.IPAddrMaskLen, _ = ParseIPv6BestPrefixRouteMaskLengthFromDEFIP(entry, index)
					}
				}

				//ALG_BKT_PTR
				var abpr *regexp.Regexp
				if index == 0 {
					abpr = regexp.MustCompile("ALG_BKT_PTR0=(?P<f>[0x]?[[:alnum:]]+)")
				} else {
					if en.Mode == 1 || en.Mode == 3 { //IPv6 64 && IPv6 128
						abpr = regexp.MustCompile("ALG_BKT_PTR0=(?P<f>[0x]?[[:alnum:]]+)")
					} else {
						abpr = regexp.MustCompile("ALG_BKT_PTR1=(?P<f>[0x]?[[:alnum:]]+)")
					}
				}

				if ms := abpr.FindStringSubmatch(entry); len(ms) == 2 {
					abp, err := strconv.ParseInt(ms[1], 0, 64)
					if err != nil {
						panic(err)
					}
					en.ALGBktPtr = abp
				}

				//Nexthop.
				var nhpr *regexp.Regexp
				if index == 0 {
					nhpr = regexp.MustCompile("NEXT_HOP_INDEX0=(?P<f>[0x]?[[:alnum:]]+)")
				} else {
					if en.Mode == 1 || en.Mode == 3 {
						nhpr = regexp.MustCompile("NEXT_HOP_INDEX0=(?P<f>[0x]?[[:alnum:]]+)")
					} else {
						nhpr = regexp.MustCompile("NEXT_HOP_INDEX1=(?P<f>[0x]?[[:alnum:]]+)")
					}
				}

				if ms := nhpr.FindStringSubmatch(entry); len(ms) == 2 {
					nhi, err := strconv.ParseInt(ms[1], 0, 64)
					if err != nil {
						panic(err)
					}
					en.NexthopIndex = nhi
					nh, err := ParseNexthopByIndex(en.NexthopIndex)
					if err != nil {
						panic(err)
					}

					en.NH = *nh
				}

				//ecmp
				var ecmpr *regexp.Regexp
				if index == 0 {
					ecmpr = regexp.MustCompile("ECMP0=(?P<f>[0x]?[[:alnum:]]+)")
				} else {
					if en.Mode == 1 || en.Mode == 3 {
						ecmpr = regexp.MustCompile("ECMP0=(?P<f>[0x]?[[:alnum:]]+)")
					} else {
						ecmpr = regexp.MustCompile("ECMP1=(?P<f>[0x]?[[:alnum:]]+)")
					}
				}

				if ms := ecmpr.FindStringSubmatch(entry); len(ms) == 2 {
					ecmp, err := strconv.ParseInt(ms[1], 0, 64)
					if err != nil {
						panic(err)
					}
					if ecmp == 1 {
						en.ECMP = true
					} else {
						en.ECMP = false
					}
				}

				//ecmp_ptr
				var ecmpptrr *regexp.Regexp
				if index == 0 {
					ecmpptrr = regexp.MustCompile(",ECMP_PTR0=(?P<f>[0x]?[[:alnum:]]+)")
				} else {
					if en.Mode == 1 || en.Mode == 3 {
						ecmpptrr = regexp.MustCompile(",ECMP_PTR0=(?P<f>[0x]?[[:alnum:]]+)")
					} else {
						ecmpptrr = regexp.MustCompile(",ECMP_PTR1=(?P<f>[0x]?[[:alnum:]]+)")
					}
				}
				if ms := ecmpptrr.FindStringSubmatch(entry); len(ms) == 2 {
					ep, err := strconv.ParseInt(ms[1], 0, 64)
					if err != nil {
						panic(err)
					}
					en.ECMPPtr = ep
					en.EG, _ = DumpECMPGroupByIndex(en.ECMPPtr)
				}

				//GLOBAL route
				var grr *regexp.Regexp
				if index == 0 {
					grr = regexp.MustCompile("GLOBAL_ROUTE0=(?P<f>[0x]?[[:alnum:]]+)")
				} else {
					if en.Mode == 1 || en.Mode == 3 {
						grr = regexp.MustCompile("GLOBAL_ROUTE0=(?P<f>[0x]?[[:alnum:]]+)")
					} else {
						grr = regexp.MustCompile("GLOBAL_ROUTE1=(?P<f>[0x]?[[:alnum:]]+)")
					}
				}

				if ms := grr.FindStringSubmatch(entry); len(ms) == 2 {
					gr, err := strconv.ParseInt(ms[1], 0, 64)
					if err != nil {
						panic(err)
					}
					if gr == 1 {
						en.GlobalRoute = true
					} else {
						en.GlobalRoute = false
					}
				}

				//ecmp_ptr
				/*
					case "IP_ADDR_MASK":
						en.IPAddrMask = FixIPv4NetMask(matches[1])
						en.IPAddrMaskLen, _ = en.IPAddrMask.Size()
					case "IP_ADDR":
						fmt.Println("[", en.Mode, "]", entry)
						if en.Mode == 0 {
							en.AF = IPV4
							en.IPAddr = FixIPv4Address(matches[1])
						} else {
							en.AF = IPV6
							en.IPAddr, _ = ParseIPv6BestPrefixRouteFromDEFIP(entry, index)
						}
				*/
			case "HIT":
				hit, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				if hit == 1 {
					en.Hit = true
				} else {
					en.Hit = false
				}
				/*
					case "GLOBAL_ROUTE":
						gr, err := strconv.ParseInt(matches[1], 0, 64)
						if err != nil {
							panic(err)
						}
						if gr == 1 {
							en.GlobalRoute = true
						} else {
							en.GlobalRoute = false
						}
				*/
			case "ENTRY_TYPE":
				et, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.EntryType = et

				/*
					case ",ECMP_PTR":
						ep, err := strconv.ParseInt(matches[1], 0, 64)
						if err != nil {
							panic(err)
						}
						en.ECMPPtr = ep

						en.EG, _ = DumpECMPGroupByIndex(en.ECMPPtr)
					case "ECMP":
						ecmp, err := strconv.ParseInt(matches[1], 0, 64)
						if err != nil {
							panic(err)
						}
						if ecmp == 1 {
							en.ECMP = true
						} else {
							en.ECMP = false
						}
				*/
			case "DEFAULT_MISS":
				dm, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				if dm == 1 {
					en.DefaultMiss = true
				} else {
					en.DefaultMiss = false
				}
			case "DEFAULT_ROUTE":
				dr, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				if dr == 1 {
					en.DefaultRoute = true
				} else {
					en.DefaultRoute = false
				}
			case "ALG_HIT_IDX":
				ahi, err := strconv.ParseInt(matches[1], 0, 64)
				if err != nil {
					panic(err)
				}
				en.ALGHitIndex = ahi
				/*
					case "ALG_BKT_PTR":
						abp, err := strconv.ParseInt(matches[1], 0, 64)
						if err != nil {
							panic(err)
						}
						en.ALGBktPtr = abp
				*/
			}
		}
	}
	//fmt.Printf("%s", en)

	return en, nil
}

var ALPMIPv4 map[int64]string
var ALPMIPv664 map[int64]string
var ALPMIPv6128 map[int64]string

func DumpALPMIPv4DB(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP_ALPM_IPV4 0 393215  | grep VALID=1",
	})

	if err != nil {
		panic(err)
	}

	ALPMIPv4 = make(map[int64]string, 393215)

	for _, l := range strings.Split(res, "\n") {
		if res, err := match(l, getEntryIndex); err != nil {
			//fmt.Println("Cannot get entry index of: ", l)
			continue
		} else {
			index, err := strconv.ParseInt(res, 0, 64)
			if err != nil {
				panic(err)
			}
			ALPMIPv4[index] = l
		}
	}
}

func DumpALPMIPv664DB(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP_ALPM_IPV6_64 0 262143 | grep VALID=1",
	})

	if err != nil {
		panic(err)
	}

	ALPMIPv664 = make(map[int64]string, 262144)

	for _, l := range strings.Split(res, "\n") {
		if res, err := match(l, getEntryIndex); err != nil {
			//fmt.Println("Cannot get entry index of: ", l)
			continue
		} else {
			index, err := strconv.ParseInt(res, 0, 64)
			if err != nil {
				panic(err)
			}
			ALPMIPv664[index] = l
		}
	}
}

func DumpL3DEFIPEntry(dev *rut.RUT) {
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP 0 8191 | grep -E 'VALID(0|1)=1'",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	for _, l := range strings.Split(res, "\n") {
		if strings.Contains(l, "VALID0=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 0)
			//Currently just dump ipv4 entry for this half entry
			if en.Mode == 0 {
				fmt.Printf("%s\n", en)
			}
		}

		if strings.Contains(l, "VALID1=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 1)
			/* if en.Mode != 0 {
				//Currently just dump IPv4 entry.
				continue
			}
			*/
			if en.Mode == 0 {
				fmt.Printf("%s\n", en)
			}
		}
		if !strings.Contains(l, "VALID1=1") && !strings.Contains(l, "VALID0=1") {
			//fmt.Println("Invalid result: ", l)
		}

		/*
			r, _ := ParseRouteEntryString(l, IPV6)
			if r != nil {
				r.ParseNexthopInfo()
			}
			fmt.Println(r)
		*/
	}
}
func DumpL3DEFIPALPMEntry(dev *rut.RUT) {
	DumpALPMIPv4DB(dev)
	DumpALPMIPv664DB(dev)
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP 0 8191 | grep -E 'VALID(0|1)=1'",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		if strings.Contains(l, "VALID0=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 0)
			//Currently just dump ipv4 entry for this half entry
			if en.Mode == 0 {
				fmt.Printf("<<<Root[%d(0)(%6s)]: %s/%d (0x%-4x) Global: %5t -> VRF : %2d)>>>: \n", en.Index, DEFIPModeValueToString[en.Mode], en.IPAddr, en.IPAddrMaskLen, en.ALGBktPtr, en.GlobalRoute, en.VRF)
				fmt.Printf("%s\n", en)

				for bank := 0; bank < 4; bank++ {
					for eindex := 0; eindex < 6; eindex++ {
						alpmidx := eindex<<16 | int(en.ALGBktPtr<<2) | bank
						//DumpIPv4EntryByIndex(dev, alpmidx)
						if e, ok := ALPMIPv4[int64(alpmidx)]; !ok {
							//Do not display invlaid entry.
							//fmt.Printf("Entry %d does not exist\n", alpmidx)
						} else {
							if !strings.Contains(e, "VALID=1") {
								fmt.Printf("Entry %d is invalid\n", alpmidx)
								continue
							}
							r, _ := ParseRouteEntryString(e, IPV4)
							if r != nil {
								r.ParseNexthopInfo()
							}
							fmt.Println(r)
						}
					}
				}
			}
		}
		if strings.Contains(l, "VALID1=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 1)
			/* if en.Mode != 0 {
				//Currently just dump IPv4 entry.
				continue
			}
			*/
			fmt.Printf("<<<Root[%d(1)(%6s)]: %s/%d (0x%-4x) Global: %5t -> VRF : %2d)>>>: \n", en.Index, DEFIPModeValueToString[en.Mode], en.IPAddr, en.IPAddrMaskLen, en.ALGBktPtr, en.GlobalRoute, en.VRF)
			fmt.Printf("%s\n", en)
			if en.Mode == 0 {
				for bank := 0; bank < 4; bank++ {
					for eindex := 0; eindex < 6; eindex++ {
						alpmidx := eindex<<16 | int(en.ALGBktPtr<<2) | bank
						//DumpIPv4EntryByIndex(dev, alpmidx)
						if e, ok := ALPMIPv4[int64(alpmidx)]; !ok {
							//Do not display invlaid entry.
							//fmt.Printf("Entry %d does not exist\n", alpmidx)
						} else {
							if !strings.Contains(e, "VALID=1") {
								fmt.Printf("Entry %d is invalid\n", alpmidx)
								continue
							}
							r, _ := ParseRouteEntryString(e, IPV4)
							if r != nil {
								r.ParseNexthopInfo()
							}
							fmt.Println(r)
						}
					}
				}
			} else if en.Mode == 1 {
				for bank := 0; bank < 4; bank++ {
					for eindex := 0; eindex < 4; eindex++ {
						alpmidx := eindex<<16 | int(en.ALGBktPtr<<2) | bank
						//DumpIPv4EntryByIndex(dev, alpmidx)
						if e, ok := ALPMIPv664[int64(alpmidx)]; !ok {
							//Do not display invlaid entry.
							//fmt.Printf("Entry %d does not exist\n", alpmidx)
						} else {
							if !strings.Contains(e, "VALID=1") {
								fmt.Printf("Entry %d is invalid\n", alpmidx)
								continue
							}
							//fmt.Println(e)
							r, _ := ParseRouteEntryString(e, IPV6)
							if r != nil {
								r.ParseNexthopInfo()
							}
							fmt.Println(r)
						}
					}
				}

			}
		}

		if !strings.Contains(l, "VALID1=1") && !strings.Contains(l, "VALID0=1") {
			//fmt.Println("Invalid result: ", l)
		}

		/*
			r, _ := ParseRouteEntryString(l, IPV6)
			if r != nil {
				r.ParseNexthopInfo()
			}
			fmt.Println(r)
		*/
	}
}

func DumpL3IPv4DEFIPEntry(dev *rut.RUT) {
	DumpALPMIPv4DB(dev)
	DumpALPMIPv664DB(dev)
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP 0 8191 | grep -E 'VALID(0|1)=1'",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		if strings.Contains(l, "VALID0=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 0)
			//Currently just dump ipv4 entry for this half entry
			if en.Mode == 0 {
				fmt.Printf("<<<Root[%d(0)(%6s)]: %32s/%2d (0x%-4x) Global: %5t -> VRF : %2d)>>>: \n", en.Index, DEFIPModeValueToString[en.Mode], en.IPAddr, en.IPAddrMaskLen, en.ALGBktPtr, en.GlobalRoute, en.VRF)
				fmt.Printf("%s\n", en)

				for bank := 0; bank < 4; bank++ {
					for eindex := 0; eindex < 6; eindex++ {
						alpmidx := eindex<<16 | int(en.ALGBktPtr<<2) | bank
						//DumpIPv4EntryByIndex(dev, alpmidx)
						if e, ok := ALPMIPv4[int64(alpmidx)]; !ok {
							//Do not display invlaid entry.
							//fmt.Printf("Entry %d does not exist\n", alpmidx)
						} else {
							if !strings.Contains(e, "VALID=1") {
								fmt.Printf("Entry %d is invalid\n", alpmidx)
								continue
							}
							r, _ := ParseRouteEntryString(e, IPV4)
							if r != nil {
								r.ParseNexthopInfo()
							}
							fmt.Println(r)
						}
					}
				}
			}
		}
		if strings.Contains(l, "VALID1=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 1)
			/* if en.Mode != 0 {
				//Currently just dump IPv4 entry.
				continue
			}
			*/
			if en.Mode == 0 {
				fmt.Printf("<<<Root[%d(1)(%6s)]: %32s/%2d (0x%-4x) Global: %5t -> VRF : %2d)>>>: \n", en.Index, DEFIPModeValueToString[en.Mode], en.IPAddr, en.IPAddrMaskLen, en.ALGBktPtr, en.GlobalRoute, en.VRF)
				fmt.Printf("%s\n", en)

				for bank := 0; bank < 4; bank++ {
					for eindex := 0; eindex < 6; eindex++ {
						alpmidx := eindex<<16 | int(en.ALGBktPtr<<2) | bank
						//DumpIPv4EntryByIndex(dev, alpmidx)
						if e, ok := ALPMIPv4[int64(alpmidx)]; !ok {
							//Do not display invlaid entry.
							//fmt.Printf("Entry %d does not exist\n", alpmidx)
						} else {
							if !strings.Contains(e, "VALID=1") {
								fmt.Printf("Entry %d is invalid\n", alpmidx)
								continue
							}
							r, _ := ParseRouteEntryString(e, IPV4)
							if r != nil {
								r.ParseNexthopInfo()
							}
							fmt.Println(r)
						}
					}
				}
			}
		}

		if !strings.Contains(l, "VALID1=1") && !strings.Contains(l, "VALID0=1") {
			//fmt.Println("Invalid result: ", l)
		}

		/*
			r, _ := ParseRouteEntryString(l, IPV6)
			if r != nil {
				r.ParseNexthopInfo()
			}
			fmt.Println(r)
		*/
	}
}

func DumpL3IPv664DEFIPEntry(dev *rut.RUT) {
	DumpALPMIPv4DB(dev)
	DumpALPMIPv664DB(dev)
	res, err := dev.RunCommand(CTX, &command.Command{
		Mode: "shell",
		CMD:  " scontrol -f /proc/switch/ASIC/ctrl dump table 0 L3_DEFIP 0 8191 | grep -E 'VALID(0|1)=1'",
	})

	if err != nil {
		panic(err)
	}

	for _, l := range strings.Split(res, "\n") {
		if strings.Contains(l, "VALID1=1") {
			en, _ := DumpL3DEFIPHalfEntry(l, 1)
			/* if en.Mode != 0 {
				//Currently just dump IPv4 entry.
				continue
			}
			*/
			if en.Mode == 1 {
				fmt.Printf("<<<Root[%d(1)(%6s)]: %24s/%2d (0x%-4x) Global: %5t -> VRF : %2d)>>>: \n", en.Index, DEFIPModeValueToString[en.Mode], en.IPAddr, en.IPAddrMaskLen, en.ALGBktPtr, en.GlobalRoute, en.VRF)
				fmt.Printf("%s\n", en)

				for bank := 0; bank < 4; bank++ {

					for eindex := 0; eindex < 4; eindex++ {
						alpmidx := eindex<<16 | int(en.ALGBktPtr<<2) | bank
						//DumpIPv4EntryByIndex(dev, alpmidx)
						if e, ok := ALPMIPv664[int64(alpmidx)]; !ok {
							//Do not display invlaid entry.
							//fmt.Printf("Entry %d does not exist\n", alpmidx)
						} else {
							if !strings.Contains(e, "VALID=1") {
								fmt.Printf("Entry %d is invalid\n", alpmidx)
								continue
							}
							//fmt.Println(e)
							r, _ := ParseRouteEntryString(e, IPV6)
							if r != nil {
								r.ParseNexthopInfo()
							}
							fmt.Println(r)
						}
					}
				}

			}
		}

		if !strings.Contains(l, "VALID1=1") && !strings.Contains(l, "VALID0=1") {
			//fmt.Println("Invalid result: ", l)
		}

		/*
			r, _ := ParseRouteEntryString(l, IPV6)
			if r != nil {
				r.ParseNexthopInfo()
			}
			fmt.Println(r)
		*/
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

	return dev, nil
}
