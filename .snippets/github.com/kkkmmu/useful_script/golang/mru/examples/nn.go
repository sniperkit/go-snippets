package main

import (
	"feature"
	"fmt"
	"group"
	"mcase"
	"newcache"
	"subgroup"
)

func main() {
	db := newcache.New("8499")
	db.AddCase(&mcase.Case{
		Group:    "L2",
		SubGroup: "Bridge",
		Feature:  "VLAN",
		Name:     "VLAN create",
	})

	db.AddCase(&mcase.Case{
		Group:    "L2",
		SubGroup: "Bridge",
		Feature:  "VLAN",
		Name:     "VLAN Delete",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv4",
		Feature:  "Static Route",
		Name:     "Static Route AddCase",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv4",
		Feature:  "Static Route",
		Name:     "Static Route Delete",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv4",
		Feature:  "Static Route",
		Name:     "Static Route Instance",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv4",
		Feature:  "L3 Interface",
		Name:     "AddCase Interface AddCaseress",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv4",
		Feature:  "L3 Interface",
		Name:     "Delete Interface AddCaseress",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "L3 Interface",
		Name:     "Delete Interface AddCaseress",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "L3 Interface",
		Name:     "AddCase Interface AddCaseress",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "Enable OSPF",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "Disable OSPF",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "OSPF Interface Parameters",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "OSPF Stub Area",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "BGP",
		Name:     "BGP Basic",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "BGP",
		Name:     "ibgp",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "BGP",
		Name:     "ebgp",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "BGP",
		Name:     "bgp summary",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "BGP",
		Name:     "bgp redistribution",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "OSPF NSSA Area",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "OSPF Redistribution",
	})

	db.AddCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "OSPF Summary",
	})

	db.AddCase(&mcase.Case{
		Group:    "System",
		SubGroup: "Connection",
		Feature:  "Telnet",
		Name:     "Telnet server",
	})

	db.AddCase(&mcase.Case{
		Group:    "System",
		SubGroup: "Connection",
		Feature:  "Telnet",
		Name:     "Telnet client",
	})

	db.AddCase(&mcase.Case{
		Group:    "System",
		SubGroup: "Connection",
		Feature:  "ssh",
		Name:     "ssh client",
	})

	db.AddCase(&mcase.Case{
		Group:    "System",
		SubGroup: "Connection",
		Feature:  "ssh",
		Name:     "ssh client",
	})

	db.AddCase(&mcase.Case{
		Group:    "System",
		SubGroup: "MGMT",
		Feature:  "SNMP",
		Name:     "SNMP server",
	})

	db.AddGroup(&group.Group{
		Name: "g1234",
	})

	db.AddGroup(&group.Group{
		Name: "g2",
	})

	db.AddGroup(&group.Group{
		Name: "g3",
	})

	db.AddSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg5",
	})

	db.AddSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg6",
	})

	db.AddSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg7",
	})

	db.AddSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg8",
	})

	db.AddFeature(&feature.Feature{
		Group:    "g1234",
		SubGroup: "sg5",
		Name:     "f111111",
	})

	db.AddFeature(&feature.Feature{
		Group:    "g1234",
		SubGroup: "sg5",
		Name:     "f211211",
	})

	db.AddFeature(&feature.Feature{
		Group:    "g1234",
		SubGroup: "sg5",
		Name:     "f211111",
	})

	db.AddFeature(&feature.Feature{
		Group:    "g1234",
		SubGroup: "sg5",
		Name:     "f213211",
	})

	fmt.Printf("DB:::::::::%q\n", db)
	fmt.Printf("CASE:::::::%q\n", db.GetAllCase())
	fmt.Printf("GROUP::::::%q\n", db.GetAllGroup())
	fmt.Printf("FEATURE::::%q\n", db.GetAllFeature())
	fmt.Printf("SUBGROUP:::%q\n", db.GetAllSubGroup())

	c, err := db.GetCase(&mcase.Case{
		Group:    "System",
		SubGroup: "MGMT",
		Feature:  "SNMP",
		Name:     "SNMP server",
	})
	fmt.Printf("GetCase:::%q %v\n", c, err)

	c, err = db.GetCase(&mcase.Case{
		Group:    "L3",
		SubGroup: "IPv6",
		Feature:  "OSPF",
		Name:     "OSPF Summary",
	})
	fmt.Printf("GetCase:::%q %v\n", c, err)

	g, err := db.GetGroup(&group.Group{
		Name: "g2",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	g, err = db.GetGroup(&group.Group{
		Name: "g3",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	db.DelGroup(&group.Group{
		Name: "g3",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	g, err = db.GetGroup(&group.Group{
		Name: "g3",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	g, err = db.GetGroup(&group.Group{
		Name: "L2",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	db.DelGroup(&group.Group{
		Name: "L2",
	})

	g, err = db.GetGroup(&group.Group{
		Name: "L2",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	g, err = db.GetGroup(&group.Group{
		Name: "L3",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)

	sg, err := db.GetSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg5",
	})
	fmt.Printf("GetSubGroup:::%q, %v\n", sg, err)

	db.DelSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg5",
	})
	fmt.Printf("GetSubGroup:::%q, %v\n", sg, err)

	sg, err = db.GetSubGroup(&subgroup.SubGroup{
		Group: "g1234",
		Name:  "sg5",
	})
	fmt.Printf("GetSubGroup:::%q, %v\n", sg, err)

	sg, err = db.GetSubGroup(&subgroup.SubGroup{
		Group: "System",
		Name:  "Connection",
	})
	fmt.Printf("GetSubGroup:::%q, %v\n", sg, err)

	f, err := db.GetFeature(&feature.Feature{
		Group:    "g1234",
		SubGroup: "sg5",
		Name:     "f213211",
	})
	fmt.Printf("GetFeature:::%q %v\n", f, err)

	f, err = db.GetFeature(&feature.Feature{
		Group:    "L3",
		SubGroup: "IPv4",
		Name:     "Static Route",
	})
	fmt.Printf("GetFeature:::%q %v\n", f, err)

	f, err = db.GetFeature(&feature.Feature{
		Group:    "L3",
		SubGroup: "IPv6",
		Name:     "OSPF",
	})
	fmt.Printf("GetFeature:::%q %v\n", f, err)
	db.DelFeature(&feature.Feature{
		Group:    "L3",
		SubGroup: "IPv6",
		Name:     "OSPF",
	})

	f, err = db.GetFeature(&feature.Feature{
		Group:    "L3",
		SubGroup: "IPv6",
		Name:     "OSPF",
	})
	fmt.Printf("GetFeature:::%q %v\n", f, err)

	g, err = db.GetGroup(&group.Group{
		Name: "L3",
	})
	fmt.Printf("GetGroup:::%q, %v\n", g, err)
}
