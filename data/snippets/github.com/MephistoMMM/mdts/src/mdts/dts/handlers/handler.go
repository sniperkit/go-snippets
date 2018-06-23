package handlers

import (
	"context"
	"fmt"
	"log"
	s2t "mdts/brokerSDK/s2t/broker"
	t2s "mdts/brokerSDK/t2s/broker"
	"mdts/dts/discovery"
	"mdts/dts/request"
	bmsg "mdts/protocols/brokermsg"
	pts "mdts/protocols/dtsproto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// ConnectBroker ...
func ConnectBroker(Type string) (*grpc.ClientConn, error) {
	info, ok := discovery.EtcdMaster.GetRandomBrokerInfo(Type)
	if !ok {
		return nil, fmt.Errorf("Non Broker %s", Type)
	}

	conn, err := grpc.Dial(info.HostPort, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// TransforDataToThird ...
func TransforDataToThird(TID string, APICODE string, Data []byte) (*bmsg.ResultToThird, error) {
	conn, err := ConnectBroker(TID)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := s2t.NewS2TBrokerClient(conn)
	result, err := client.TransforDataToThird(context.Background(), &bmsg.ParamToThird{
		TID:     TID,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// TransforDataFromThird ...
func TransforDataFromThird(TID string, APICODE string, Data []byte) (*bmsg.ResultFromThird, error) {
	conn, err := ConnectBroker(TID)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := s2t.NewS2TBrokerClient(conn)
	result, err := client.TransforDataFromThird(context.Background(), &bmsg.ParamFromThird{
		TID:     TID,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// TransforDataToService ...
func TransforDataToService(Version string, APICODE string, Data []byte) (*bmsg.ResultToService, error) {
	conn, err := ConnectBroker(Version)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := t2s.NewT2SBrokerClient(conn)
	result, err := client.TransforDataToService(context.Background(), &bmsg.ParamToService{
		Version: Version,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// TransforDataFromService ...
func TransforDataFromService(Version string, APICODE string, Data []byte) (*bmsg.ResultFromService, error) {
	conn, err := ConnectBroker(Version)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := t2s.NewT2SBrokerClient(conn)
	result, err := client.TransforDataFromService(context.Background(), &bmsg.ParamFromService{
		Version: Version,
		APICODE: APICODE,
		Data:    Data,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// HandleS2T ...
func HandleS2T(c *gin.Context) {
	head, body, err := pts.ParseS2T(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of TID From Etcd
	log.Printf("Request To TID: %s.", head.TID)

	result, err := TransforDataToThird(head.TID, head.APICODE, body)
	if err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}
	if result.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(1) Error!",
		})
		return
	}

	_, byt, err := request.PostBytes(result.GetURL(), result.GetHead()["Content-Type"], result.GetBody())
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	result2, err := TransforDataFromThird(head.TID, head.APICODE, byt)
	if err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}
	if result2.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(2) Error!",
		})
		return
	}

	c.Data(200, "application/json", result2.GetBody())
}

// HandleT2S ...
func HandleT2S(c *gin.Context) {
	head, body, err := pts.ParseT2S(c)
	if err != nil {
		errStr := fmt.Sprintf("Invalid Request From Third: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	// Get Info Of Version From Etcd
	log.Printf("Request To Version: %s.", head.Version)

	result, err := TransforDataToService(head.Version, head.APICODE, body)
	if err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}
	if result.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(1) Error!",
		})
		return
	}

	_, byt, err := request.PostBytes(result.GetURL(), result.GetHead()["Content-Type"], result.GetBody())
	if err != nil {
		errStr := fmt.Sprintf("Invalid Response From Service: %v.", err)
		log.Println(errStr)
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: errStr,
		})
		return
	}

	result2, err := TransforDataFromService(head.Version, head.APICODE, byt)
	if err != nil {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: err.Error(),
		})
		return
	}
	if result2.GetState() == bmsg.EnumState_FAILED {
		c.JSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Broker(2) Error!",
		})
		return
	}

	c.Data(200, "application/json", result2.GetBody())
}
