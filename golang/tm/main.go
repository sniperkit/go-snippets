package main

import (
	"fmt"
	"group"
	"server"
)

var gGroup *group.Group

func main() {
	server.Start()
}

func init() {
	gGroup = group.New("TSL")
	st, _ := gGroup.AddTeam("SWITCH")
	lt, _ := gGroup.AddTeam("L3")
	liwei, _ := gGroup.AddMember("liwei")
	st.Join(liwei)
	pxy, _ := gGroup.AddMember("pengxiaoyan")
	st.Join(pxy)
	zfk, _ := gGroup.AddMember("zhangfengkun")
	st.Join(zfk)
	ly, _ := gGroup.AddMember("liyue")
	st.Join(ly)
	lhw, _ := gGroup.AddMember("liuhuawei")
	st.Join(lhw)
	lth, _ := gGroup.AddMember("liutonghai")
	lt.Join(lth)
	zyc, _ := gGroup.AddMember("zhouyuchao")
	lt.Join(zyc)
	yj, _ := gGroup.AddMember("yaojuan")
	lt.Join(yj)
	v24, _ := gGroup.AddProduct("V2224G")
	v2t, _ := v24.AddTask("Vlan Create Failed with id 4094")
	v27, _ := gGroup.AddProduct("V2708M")
	v7t, _ := v27.AddTask("Cannot get port status of pon link")
	v28, _ := gGroup.AddProduct("V2808M")
	v8t, _ := v28.AddTask("Cannot telnet to v2808m")
	v85, _ := gGroup.AddProduct("V8500")
	v5t, _ := v85.AddTask("OSPF neighbor adjacency cannot estiblish")
	v81, _ := gGroup.AddProduct("V8106")
	v1t, _ := v81.AddTask("BGP Extented ASN support requirement")

	liwei.Do(v2t)
	pxy.Do(v7t)
	zfk.Do(v8t)
	yj.Do(v5t)
	zyc.Do(v1t)

	fmt.Printf("%+v\n", gGroup)
}
