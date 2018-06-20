package discovery

import (
	"color"
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"mdts/dts/conf"
	dproto "mdts/protocols/discovery"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var EtcdMaster *Master

func init() {
	master, err := NewMaster(conf.EndPoints, conf.EtcdPath)
	if err != nil {
		log.Fatalln(err)
	}

	EtcdMaster = master
}

// Master watch etcd to implement service discovery
// etcd path: brokers/<type>_<hash: hash(time)>
type Master struct {
	Path   string
	locker sync.RWMutex
	Nodes  map[string]*BrokerMap
	Client *clientv3.Client
}

// BrokerMap : Master.Nodes { "type": BrokerMap.Brokers{ "hash": BrokerInfo}}
type BrokerMap struct {
	Type    string
	Brokers *SortedMap
}

func NewBrokerMap(key string) *BrokerMap {
	return &BrokerMap{
		Type:    key,
		Brokers: NewSortedMap(),
	}
}

func (bm *BrokerMap) AddBroker(hash string, info *dproto.BrokerInfo) {
	bm.Brokers.AddKV(hash, info)
}

func (bm *BrokerMap) DeleteBroker(hash string) (info *dproto.BrokerInfo) {
	v := bm.Brokers.DeleteKV(hash)
	info, _ = v.(*dproto.BrokerInfo)
	return info
}

func (bm *BrokerMap) Len() int {
	return bm.Brokers.Len()
}

func NewMaster(endpoints []string, watchPath string) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})

	if err != nil {
		return nil, err
	}

	master := &Master{
		Path:   watchPath,
		Nodes:  make(map[string]*BrokerMap),
		Client: cli,
	}

	return master, err
}

func (m *Master) AddBroker(info *dproto.BrokerInfo) {
	m.locker.Lock()
	defer m.locker.Unlock()

	bm, ok := m.Nodes[info.Type]
	if !ok {
		bm = NewBrokerMap(info.Type)
		m.Nodes[info.Type] = bm
	}

	bm.AddBroker(info.Hash, info)
}

func (m *Master) DeleteBroker(Type string, hash string) {
	m.locker.Lock()
	defer m.locker.Unlock()

	bm, ok := m.Nodes[Type]
	if !ok {
		return
	}

	bm.DeleteBroker(hash)
	if bm.Len() == 0 {
		delete(m.Nodes, Type)
	}
}

// GetRandomBrokerInfo get random BrokerInfo from Nodes map, if not exist, return nil false
func (m *Master) GetRandomBrokerInfo(Type string) (*dproto.BrokerInfo, bool) {
	m.locker.RLock()
	defer m.locker.RUnlock()

	bm, ok := m.Nodes[Type]
	if !ok {
		return nil, ok
	}

	n := rand.Intn(bm.Len())
	info, _ := bm.Brokers.GetIV(n).(*dproto.BrokerInfo)
	return info, true
}

func parseBrokerInfo(value []byte) *dproto.BrokerInfo {
	info := &dproto.BrokerInfo{}
	err := json.Unmarshal(value, info)
	if err != nil {
		log.Println(err)
	}
	return info
}

func parseBrokerInfoFromKey(key string) (Type string, hash string) {
	slashPostion := strings.IndexByte(key, '/')
	key = key[slashPostion+1:]
	items := strings.Split(key, "_")

	return items[0], items[1]
}

func (m *Master) WatchNodes() {
	// TODO: Fetch Exist Nodes First
	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				log.Printf("[%s] %q : %q\n", color.Green("PUT"), ev.Kv.Key, ev.Kv.Value)
				info := parseBrokerInfo(ev.Kv.Value)
				m.AddBroker(info)
			case clientv3.EventTypeDelete:
				log.Printf("[%s] %q\n", color.Red("DELETE"), ev.Kv.Key)
				Type, hash := parseBrokerInfoFromKey(string(ev.Kv.Key))
				m.DeleteBroker(Type, hash)
			}
		}
	}
}

func (m *Master) Start() error {

	resp, err := m.Client.Get(context.Background(), m.Path, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, ev := range resp.Kvs {
		log.Printf("[%s] %s : %s\n", color.Cyan("FOUND"), ev.Key, ev.Value)
		info := parseBrokerInfo(ev.Value)
		m.AddBroker(info)
	}

	log.Println("Discovery Master Start.")
	m.WatchNodes()
	log.Println("Discovery Master End.")

	return nil
}
