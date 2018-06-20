package base

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	dproto "mdts/protocols/discovery"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
)

var isDev bool = true

func init() {
	if os.Getenv("ETCD_ENV") == "release" {
		isDev = false
	}
}

func Hash(t string, un uint64) string {
	// 8 + len(t)
	b := make([]byte, 8+len(t))
	bp := copy(b, t)

	binary.BigEndian.PutUint64(b[bp:], un)

	hash := sha256.New()
	hash.Write(b)

	return hex.EncodeToString(hash.Sum(nil))
}

type Service struct {
	Name    string
	Type    string
	Hash    string
	Info    *dproto.BrokerInfo
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
}

func NewService(t string, hostport string, endpoints []string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	un := time.Now().UnixNano()
	hash := Hash(t, uint64(un))

	return &Service{
		Name: "broker/" + t + "_" + hash,
		Type: t,
		Hash: hash,
		Info: &dproto.BrokerInfo{
			HostPort:  hostport,
			Type:      t,
			Hash:      hash,
			StartTime: un,
		},
		stop:   make(chan error),
		client: cli,
	}, nil
}

func (s *Service) Start() error {
	ch, err := s.keepAlive()
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("server closed")
		case ka, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				if isDev {
					log.Printf("Recv reply from service: %s, ttl:%d", s.Name, ka.TTL)
				}
			}
		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {

	info := &s.Info

	key := s.Name
	value, _ := json.Marshal(info)

	// minimum lease TTL is 5-second
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		return nil, err
	}

	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID

	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {

	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("servide:%s stop\n", s.Name)
	return err
}
