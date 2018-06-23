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

package dbutils

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/nu7hatch/gouuid"
	"models/events"
	"models/objects"
	"reflect"
	"strings"
	"sync"
	"time"
	"utils/logging"
)

const (
	DB_CONNECT_TIME_INTERVAL   = 2
	DB_CONNECT_RETRY_LOG_COUNT = 100
)

type DBNotConnectedError struct {
	network string
	address string
}

func (e DBNotConnectedError) Error() string {
	return fmt.Sprintf("Not connected to DB at %s%s", e.network, e.address)
}

type DBUtil struct {
	redis.Conn
	logger  logging.LoggerIntf
	network string
	address string
	DbLock  sync.RWMutex
}

type DBIntf interface {
	Connect() error
	Disconnect()
	StoreObjectInDb(objects.ConfigObj) error
	StoreObjectDefaultInDb(objects.ConfigObj) error
	DeleteObjectFromDb(objects.ConfigObj) error
	GetObjectFromDb(objects.ConfigObj, string) (objects.ConfigObj, error)
	GetKey(objects.ConfigObj) string
	GetAllObjFromDb(objects.ConfigObj) ([]objects.ConfigObj, error)
	CompareObjectsAndDiff(objects.ConfigObj, map[string]bool, objects.ConfigObj) ([]bool, error)
	CompareObjectDefaultAndDiff(objects.ConfigObj, objects.ConfigObj) ([]bool, error)
	UpdateObjectInDb(objects.ConfigObj, objects.ConfigObj, []bool) error
	MergeDbAndConfigObj(objects.ConfigObj, objects.ConfigObj, []bool) (objects.ConfigObj, error)
	GetBulkObjFromDb(obj objects.ConfigObj, startIndex, count int64) (error, int64, int64, bool, []objects.ConfigObj)
	Publish(string, interface{}, interface{})
	StoreValInDb(interface{}, interface{}, interface{}) error
	DeleteValFromDb(interface{}) error
	GetAllKeys(interface{}) (interface{}, error)
	GetValFromDB(key interface{}, field interface{}) (val interface{}, err error)
	StoreEventObjectInDb(events.EventObj) error
	GetEventObjectFromDb(events.EventObj, string) (events.EventObj, error)
	GetAllEventObjFromDb(events.EventObj) ([]events.EventObj, error)
	MergeDbAndConfigObjForPatchUpdate(objects.ConfigObj, objects.ConfigObj, []objects.PatchOpInfo) (objects.ConfigObj, []bool, error)
	StoreUUIDToObjKeyMap(objKey string) (string, error)
	DeleteUUIDToObjKeyMap(uuid, objKey string) error
	GetUUIDFromObjKey(objKey string) (string, error)
	GetObjKeyFromUUID(uuid string) (string, error)
	MergeDbObjKeys(obj, dbObj objects.ConfigObj) (objects.ConfigObj, error)
}

func NewDBUtil(logger logging.LoggerIntf) *DBUtil {
	return &DBUtil{
		logger:  logger,
		network: "tcp",
		address: ":6379",
		DbLock:  sync.RWMutex{},
	}
}

func (db *DBUtil) Connect() error {
	retryCount := 0
	ticker := time.NewTicker(DB_CONNECT_TIME_INTERVAL * time.Second)
	for _ = range ticker.C {
		retryCount += 1
		dbHdl, err := redis.Dial(db.network, db.address)
		if err != nil {
			if retryCount%DB_CONNECT_RETRY_LOG_COUNT == 0 {
				if db.logger != nil {
					db.logger.Err(fmt.Sprintln("Failed to dial out to Redis server. Retrying connection. Num retries = ", retryCount))
				}
			}
		} else {
			// ping to ensure that that the server is up and running
			// this is the suggested way to determine that redis is 'ready'
			dbHdl.Send("PING")
			dbHdl.Flush()
			response, err := dbHdl.Receive()
			var pongReply interface{} = "PONG"
			db.logger.Info(fmt.Sprintln("Received Response From Redis Server %#v", response))
			if err == nil && reflect.DeepEqual(response, pongReply) {
				db.Conn = dbHdl
				break
			}
		}
	}
	return nil
}

func (db *DBUtil) Disconnect() {
	if db.Conn != nil {
		db.Close()
	}
}

func (db *DBUtil) StoreObjectInDb(obj objects.ConfigObj) error {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.StoreObjectInDb(db.Conn)
}

func (db *DBUtil) StoreObjectDefaultInDb(obj objects.ConfigObj) error {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.StoreObjectDefaultInDb(db.Conn)
}

func (db *DBUtil) DeleteObjectFromDb(obj objects.ConfigObj) error {
	if db.Conn == nil {
		return DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.DeleteObjectFromDb(db.Conn)
}
func (db *DBUtil) DeleteObjectWithKeyFromDb(key interface{}) error {
	if db.Conn == nil {
		return DBNotConnectedError{db.network, db.address}
	}
	list, err := redis.Strings(db.Do("KEYS", key))
	if err != nil {
		fmt.Println("Failed to get all object keys from db for key", key, " error:")
		return err
	}
	for _, k := range list {
		_, err = db.Do("DEL", k)
		if err != nil {
			fmt.Println("Failed to delete obj from DB for key", k, " error:", err)
			return err
		}
	}
	return nil
}
func (db *DBUtil) GetObjectFromDb(obj objects.ConfigObj, objKey string) (objects.ConfigObj, error) {
	if db.Conn == nil {
		return obj, DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.GetObjectFromDb(objKey, db.Conn)
}

func (db *DBUtil) GetKey(obj objects.ConfigObj) string {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.GetKey()
}

func (db *DBUtil) GetAllObjFromDb(obj objects.ConfigObj) ([]objects.ConfigObj, error) {
	if db.Conn == nil {
		return make([]objects.ConfigObj, 0), DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.GetAllObjFromDb(db.Conn)
}

func (db *DBUtil) CompareObjectsAndDiff(obj objects.ConfigObj, updateKeys map[string]bool, inObj objects.ConfigObj) (
	[]bool, error) {
	if db.Conn == nil {
		return make([]bool, 0), DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.CompareObjectsAndDiff(updateKeys, inObj)
}

func (db *DBUtil) CompareObjectDefaultAndDiff(obj objects.ConfigObj, inObj objects.ConfigObj) (
	[]bool, error) {
	if db.Conn == nil {
		return make([]bool, 0), DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.CompareObjectDefaultAndDiff(inObj)
}

func (db *DBUtil) UpdateObjectInDb(obj, inObj objects.ConfigObj, attrSet []bool) error {
	if db.Conn == nil {
		return DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.UpdateObjectInDb(inObj, attrSet, db.Conn)
}

func (db *DBUtil) MergeDbAndConfigObj(obj, dbObj objects.ConfigObj, attrSet []bool) (objects.ConfigObj, error) {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.MergeDbAndConfigObj(dbObj, attrSet)
}

func (db *DBUtil) MergeDbAndConfigObjForPatchUpdate(obj, dbObj objects.ConfigObj, patchInfo []objects.PatchOpInfo) (objects.ConfigObj, []bool, error) {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.MergeDbAndConfigObjForPatchUpdate(dbObj, patchInfo)
}

func (db *DBUtil) GetBulkObjFromDb(obj objects.ConfigObj, startIndex, count int64) (error, int64, int64, bool,
	[]objects.ConfigObj) {
	if db.Conn == nil {
		return DBNotConnectedError{db.network, db.address}, 0, 0, false, make([]objects.ConfigObj, 0)
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.GetBulkObjFromDb(startIndex, count, db.Conn)
}

func (db *DBUtil) Publish(op string, channel interface{}, msg interface{}) {
	if db.Conn != nil {
		defer db.DbLock.Unlock()
		db.DbLock.Lock()
		db.Do(op, channel, msg)
	}
}

func (db *DBUtil) StoreEventObjectInDb(obj events.EventObj) error {
	if db.Conn == nil {
		return DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.StoreObjectInDb(db.Conn)
}

func (db *DBUtil) GetEventObjectFromDb(obj events.EventObj, objKey string) (events.EventObj, error) {
	if db.Conn == nil {
		return obj, DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.GetObjectFromDb(objKey, db.Conn)
}

func (db *DBUtil) GetAllEventObjFromDb(obj events.EventObj) ([]events.EventObj, error) {
	if db.Conn == nil {
		return make([]events.EventObj, 0), DBNotConnectedError{db.network, db.address}
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.GetAllObjFromDb(db.Conn)
}

func (db *DBUtil) StoreValInDb(key interface{}, val interface{}, field interface{}) error {
	if db.Conn != nil {
		defer db.DbLock.Unlock()
		db.DbLock.Lock()
		_, err := db.Do("HMSET", key, field, val)
		return err
	}
	err := errors.New("DB Connection handler is nil")
	return err
}

func (db *DBUtil) DeleteValFromDb(key interface{}) error {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	_, err := db.Do("DEL", key)
	if err != nil {
		db.logger.Err(fmt.Sprintln("Failed to delete entry with key ", key, "entry in db ", err))
		return err
	}
	return nil
}

func (db *DBUtil) GetAllKeys(pattern interface{}) (val interface{}, err error) {
	if db.Conn != nil {
		defer db.DbLock.Unlock()
		db.DbLock.Lock()
		val, err = db.Do("KEYS", pattern)
		return val, err
	}
	err = errors.New("DB Connection handler is nil")
	return val, err
}

func (db *DBUtil) GetValFromDB(key interface{}, field interface{}) (val interface{}, err error) {
	if db.Conn != nil {
		defer db.DbLock.Unlock()
		db.DbLock.Lock()
		val, err := db.Do("HGET", key, field)
		return val, err
	}
	err = errors.New("DB Connection handler is nil")
	return val, err
}

func (db *DBUtil) StoreUUIDToObjKeyMap(objKey string) (string, error) {
	UUId, err := uuid.NewV4()
	if err != nil {
		db.logger.Err(fmt.Sprintln("Failed to get UUID ", err))
		return "", err
	}
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	_, err = db.Do("SET", UUId.String(), objKey)
	if err != nil {
		db.logger.Err(fmt.Sprintln("Failed to insert uuid to objkey entry in db ", err))
		return "", err
	}
	objKeyWithUUIDPrefix := "UUID" + objKey
	_, err = db.Do("SET", objKeyWithUUIDPrefix, UUId.String())
	if err != nil {
		db.logger.Err(fmt.Sprintln("Failed to insert objkey to uuid entry in db ", err))
		return "", err
	}
	return UUId.String(), nil
}

func (db *DBUtil) DeleteUUIDToObjKeyMap(uuid, objKey string) error {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	_, err := db.Do("DEL", uuid)
	if err != nil {
		db.logger.Err(fmt.Sprintln("Failed to delete uuid to objkey entry in db ", err))
		return err
	}
	objKeyWithUUIDPrefix := "UUID" + objKey
	_, err = db.Do("DEL", objKeyWithUUIDPrefix)
	if err != nil {
		db.logger.Err(fmt.Sprintln("Failed to delete objkey to uuid entry in db ", err))
		return err
	}
	return nil
}

func (db *DBUtil) GetUUIDFromObjKey(objKey string) (string, error) {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	objKeyWithUUIDPrefix := "UUID" + objKey
	uuid, err := redis.String(db.Do("GET", objKeyWithUUIDPrefix))
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func (db *DBUtil) GetObjKeyFromUUID(uuid string) (string, error) {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	objKey, err := redis.String(db.Do("GET", uuid))
	if err != nil {
		return "", err
	}
	objKey = strings.TrimRight(objKey, "UUID")
	return objKey, nil
}

func (db *DBUtil) MergeDbObjKeys(obj, dbObj objects.ConfigObj) (objects.ConfigObj, error) {
	defer db.DbLock.Unlock()
	db.DbLock.Lock()
	return obj.MergeDbObjKeys(dbObj)
}
