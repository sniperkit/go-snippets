//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package server

import (
	"arpd"
	"github.com/garyburd/redigo/redis"
	"models/objects"
	"strconv"
	"utils/dbutils"
)

type arpDbEntry struct {
	IpAddr  string
	L3IfIdx int
}

func (server *ARPServer) initiateDB() error {
	var err error
	server.dbHdl = dbutils.NewDBUtil(server.logger)
	err = server.dbHdl.Connect()
	if err != nil {
		server.logger.Err("Failed to create the DB handle")
		return err
	}
	return nil
}

func (server *ARPServer) getArpGlobalConfig() {
	var dbObj objects.ArpGlobal

	objList, err := dbObj.GetAllObjFromDb(server.dbHdl)
	if err != nil {
		server.logger.Err("DB Query init failed during Arp Initialization")
		return
	}

	if objList == nil {
		server.logger.Debug("No Config object found in DB for Arp")
		return
	}
	obj := arpd.NewArpGlobal()
	dbObject := objList[0].(objects.ArpGlobal)
	objects.ConvertarpdArpGlobalObjToThrift(&dbObject, obj)
	server.logger.Info("Timeout : ", int(obj.Timeout))
	arpConf := ArpConf{
		RefTimeout: int(obj.Timeout),
	}
	server.processArpConf(arpConf)
}

func (server *ARPServer) updateArpCacheFromDB() {
	server.logger.Debug("Populate ARP Cache from DB entries")
	if server.dbHdl != nil {
		keyPattern := "ArpCacheEntry#*"
		keys, err := redis.Strings(redis.Values(server.dbHdl.Do("KEYS", keyPattern)))
		if err != nil {
			server.logger.Err("Failed to get all keys from DB")
			return
		}
		for idx := 0; idx < len(keys); idx++ {
			var obj arpDbEntry
			val, err := redis.Values(server.dbHdl.Do("HGETALL", keys[idx]))
			if err != nil {
				server.logger.Err("Failed to get ARP entry for key:", keys[idx])
				continue
			}
			err = redis.ScanStruct(val, &obj)
			if err != nil {
				server.logger.Err("Failed to get values corresponding to ARP entry key:", keys[idx])
				continue
			}
			server.logger.Debug("Data Retrived From DB IP:", obj.IpAddr, "L3IfIdx:", obj.L3IfIdx)
			server.logger.Debug("Adding arp cache entry for ", obj.IpAddr)
			ent := server.arpCache[obj.IpAddr]
			ent.MacAddr = "incomplete"
			ent.Counter = (server.minCnt + server.retryCnt + 1)
			ent.L3IfIdx = obj.L3IfIdx
			server.arpCache[obj.IpAddr] = ent
		}
	} else {
		server.logger.Err("DB handler is nil")
	}
	server.logger.Debug("Arp Cache after restoring: ", server.arpCache)
}

func (server *ARPServer) refreshArpDB() {
	if server.dbHdl != nil {
		keyPattern := "ArpCacheEntry#*"
		keys, err := redis.Strings(redis.Values(server.dbHdl.Do("KEYS", keyPattern)))
		if err != nil {
			server.logger.Err("Failed to get all keys from DB")
			return
		}
		for idx := 0; idx < len(keys); idx++ {
			_, err := server.dbHdl.Do("DEL", keys[idx])
			if err != nil {
				server.logger.Err("Failed to Delete ARP entry for key:", keys[idx])
				continue
			}
		}
	} else {
		server.logger.Err("DB handler is nil")
	}
}

func (server *ARPServer) deleteArpEntryInDB(ipAddr string) {
	if server.dbHdl != nil {
		keyPattern := "ArpCacheEntry#" + ipAddr + "*"
		keys, err := redis.Strings(redis.Values(server.dbHdl.Do("KEYS", keyPattern)))
		if err != nil {
			server.logger.Err("Failed to get all keys from DB")
			return
		}
		for idx := 0; idx < len(keys); idx++ {
			_, err := server.dbHdl.Do("DEL", keys[idx])
			if err != nil {
				server.logger.Err("Failed to Delete ARP entry for key:", keys[idx])
				continue
			}
		}
	} else {
		server.logger.Err("DB handler is nil")
	}
}

func (server *ARPServer) storeArpEntryInDB(ip string, l3IfIdx int) {
	if server.dbHdl != nil {
		key := "ArpCacheEntry#" + ip + "#" + strconv.Itoa(l3IfIdx)
		obj := arpDbEntry{
			IpAddr:  ip,
			L3IfIdx: l3IfIdx,
		}
		_, err := server.dbHdl.Do("HMSET", redis.Args{}.Add(key).AddFlat(&obj)...)
		if err != nil {
			server.logger.Err("Failed to add entry to db : ", ip, l3IfIdx, err)
			return
		}
		return
	} else {
		server.logger.Err("DB handler is nil")
	}
}
