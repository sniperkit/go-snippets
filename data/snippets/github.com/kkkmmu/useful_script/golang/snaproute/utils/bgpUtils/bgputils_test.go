//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package bgpUtils

import (
	"fmt"
	"testing"
)

type ASPathData struct {
	AsPath   string
	MatchStr string
	Expected bool
}

var ASPathInfo []ASPathData

func TestEncodeExtCommunity(t *testing.T) {
	fmt.Println("****TestEncodeExtCommunity()****")
	var extComm ExtCommunity
	var err error
	var val uint64
	extComm = ExtCommunity{"Route-Target", "200:10"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "65535:4294967295"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "200.1.1.22:10"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "70000:300"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "65535.210:300"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "75535.210:300"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)

}
func TestGetCommunityValue(t *testing.T) {
	fmt.Println("****TestGetCommunityValue()****")
	comm := ""
	comm = "200:10"
	commVal, _ := GetCommunityValue(comm)
	fmt.Println("Comm bytes for inp:", comm, " is - ", commVal)
	comm = "65535:4294967295"
	commVal, _ = GetCommunityValue(comm)
	fmt.Println("Comm bytes for inp:", comm, " is - ", commVal)
}
func TestInitASPathInfo(t *testing.T) {
	fmt.Println("****TestInitASPathInfo()****")
	ASPathInfo = make([]ASPathData, 0)
	ASPathInfo = append(ASPathInfo, ASPathData{"^.{0}$", "400", false})
	ASPathInfo = append(ASPathInfo, ASPathData{"^.{0}$", "", true})
	ASPathInfo = append(ASPathInfo, ASPathData{"^4 [0-9]*$", "4 200", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*400.*", "400", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*400.*", "200 400", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*400.*", "200 4000", false})
	ASPathInfo = append(ASPathInfo, ASPathData{".*400.*", "200 4000,400", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*400.*", "200 3400", false})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200 300 400.*", "200 400", false})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200 300 400.*", "1 300 200 400", false})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200 300 400.*", "200 3000 300 400", false})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200 300 400.*", "3 200 300 400", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200 300 400.*", "3 200 300 4000", false})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200 300 400.*", "2 4 200 300 400 5 6", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200.* 300 400.*", "200 2 300 400", true})
	ASPathInfo = append(ASPathInfo, ASPathData{".*200.* 300 400.*", "200 2 300 200 400", false})
}
func TestMatchAsPath(t *testing.T) {
	fmt.Println("***TestMatchAsPath()********")
	for _, v := range ASPathInfo {
		inAsPathRegex, _ := GetAsPathRegex(v.AsPath)
		fmt.Println("aspath regex for inp ", v, " is ", inAsPathRegex)
		if MatchASPath(inAsPathRegex, v.MatchStr) != v.Expected {
			fmt.Println("aspath regex match for inp:", v.MatchStr, ":", MatchASPath(inAsPathRegex, v.MatchStr), " not the same as expected:", v.Expected)
		}
	}
}
