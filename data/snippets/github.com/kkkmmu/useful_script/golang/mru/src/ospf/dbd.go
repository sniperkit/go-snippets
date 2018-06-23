package ospf

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type DBD struct {
	Header           Header
	InterfaceMTU     uint16
	Options          uint8
	MasterSlave      uint8
	DDSequenceNumber uint32
	LSAHeader        []*LSAHeader
}

func UnMarshalDBD(b []byte, length int) (*DBD, error) {
	if length < 8 {
		return nil, errors.New("DBD packet is too short")
	}

	d := &DBD{}
	d.InterfaceMTU = binary.BigEndian.Uint16(b[0:2])
	d.Options = b[2]
	d.MasterSlave = b[3]
	d.DDSequenceNumber = binary.BigEndian.Uint32(b[4:8])

	if length-8 > 0 { //LSA Header Exist
		d.LSAHeader = make([]*LSAHeader, 0, (length-8)/20) //LSA Header Count
		for i := 0; i < (length - 8); i++ {
			if lsh, err := UnMarshalLSAHeader(b[8+i*20 : 8+(i+1)*20]); err != nil {
				panic(err)
			} else {
				d.LSAHeader = append(d.LSAHeader, lsh)
			}
		}
	}

	return d, nil
}

func (d DBD) Marshal() ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint16(b[0:2], uint16(d.InterfaceMTU))
	b[2] = d.Options
	b[3] = d.MasterSlave
	binary.BigEndian.PutUint32(b[4:8], uint32(d.DDSequenceNumber))
	return b, nil
}

func (d DBD) IsOSPF()  bool {
	return true
}

var MasterSlaveValueToString = map[uint8]string{
	1: "Master",
	2: "More",
	3: "More, Master",
	4: "Init",
	5: "Init, Master",
	6: "Init, More",
	7: "Init, More, Master",
}

func (d *DBD) String() string {
	var s string
	s += fmt.Sprintf("OSPF Database Description Packet: ")
	s += fmt.Sprintf("	 InterfaceMTU           : %d\n", d.InterfaceMTU)
	s += fmt.Sprintf("	 Options:               : 0x%x\n", d.Options)
	s += fmt.Sprintf("	 Master/Slave           : 0x%x(%s)\n", d.MasterSlave, MasterSlaveValueToString[d.MasterSlave])
	s += fmt.Sprintf("	 DDSequenceNumber       : 0x%x\n", d.DDSequenceNumber)
	s += fmt.Sprintf("	 LSAHeader              : \n")
	for _, lsh := range d.LSAHeader {
		s += fmt.Sprintf("                  %s\n", lsh)
	}

	return s
}


func (d DBD) UnMarshal(data []byte) (interface{}, error) {
	return nil, nil
}
