package t2s

import (
	"conf"
	"context"
	"errors"
	"log"
	"mdts/brokerSDK/base"
	pb "mdts/brokerSDK/t2s/broker"
	bmsg "mdts/protocols/brokermsg"
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

var confMap = map[string]string{
	"endpoints": "",
}

var (
	errConflictOfID = errors.New("Error Conflict Of Version")
	endpoints       = []string{"150.95.157.181:2385", "150.95.157.181:2381", "150.95.157.181:2383"}
)

func init() {
	conf.InitConfMapFromEnv(confMap)

	if confMap["endpoints"] != "" {
		endpoints = strings.Split(confMap["endpoints"], ",")
		for i, v := range endpoints {
			endpoints[i] = strings.TrimSpace(v)
		}
	}
}

// Server provide a method to run a grpc server
type Server struct {
	t    base.Transformer
	once sync.Once
	s    *base.Service
}

// normally, registerAddr == hostport , but if broker is run in docker, registerAddr is the public ip and port

func NewServer(T string, registerAddr string, t base.Transformer) *Server {
	service, err := base.NewService(T, registerAddr, endpoints)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	server := &Server{
		t: t,
		s: service,
	}

	return server
}

// TransforDataToService implement service T2SBroker.TransforDataToService
func (svr *Server) TransforDataToService(ctx context.Context, p *bmsg.ParamToService) (*bmsg.ResultToService, error) {
	if svr.t.ID() != p.GetVersion() {
		return nil, errConflictOfID
	}

	tfr, err := svr.t.TransTo(p.GetAPICODE(), p.GetData())
	if err != nil {
		return nil, err
	}

	return &bmsg.ResultToService{
		State:  bmsg.EnumState_SUCCESS,
		Method: tfr.Method,
		Head:   tfr.Head,
		Body:   tfr.Body,
		URL:    tfr.URL,
	}, nil
}

// TransforDataFromService implement service T2SBroker.TransforDataFromService
func (svr *Server) TransforDataFromService(ctx context.Context, p *bmsg.ParamFromService) (*bmsg.ResultFromService, error) {
	if svr.t.ID() != p.GetVersion() {
		return nil, errConflictOfID
	}

	tfr, err := svr.t.TransFrom(p.GetAPICODE(), p.GetData())
	if err != nil {
		return nil, err
	}

	return &bmsg.ResultFromService{
		State: bmsg.EnumState_SUCCESS,
		Head:  tfr.Head,
		Body:  tfr.Body,
	}, nil
}

// Run run server
func (svr *Server) Run(hostport string) (err error) {
	svr.once.Do(func() {
		err = svr.run(hostport)
	})

	return err
}

func (svr *Server) run(hostport string) error {
	go func() {
		log.Println("Discovery Service Start.")
		svr.s.Start()
		log.Println("Discovery Service Stop.")
	}()

	lis, err := net.Listen("tcp", hostport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterT2SBrokerServer(grpcServer, svr)
	return grpcServer.Serve(lis)
}
