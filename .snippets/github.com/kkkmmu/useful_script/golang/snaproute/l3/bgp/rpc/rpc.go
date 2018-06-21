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

// server.go
package rpc

import (
	"asicdServices"
	"bfdd"
	"bgpd"
	"encoding/json"
	_ "fmt"
	"io/ioutil"
	"ndpd"
	"ribd"
	"strconv"
	"time"
	"utils/ipcutils"
	"utils/logging"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const ClientsFileName string = "clients.json"

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

func getClient(logger *logging.Writer, fileName string, process string) (*ClientJson, error) {
	var allClients []ClientJson

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Errf("Failed to open BGPd config file:%s, err:%s", fileName, err)
		return nil, err
	}

	json.Unmarshal(data, &allClients)
	for _, client := range allClients {
		if client.Name == process {
			return &client, nil
		}
	}

	logger.Errf("Did not find port for %s in config file:%s", process, fileName)
	return nil, nil
}

func StartServer(logger *logging.Writer, handler *BGPHandler, filePath string) {
	fileName := filePath + ClientsFileName
	clientJson, err := getClient(logger, fileName, "bgpd")
	if err != nil || clientJson == nil {
		logger.Err("Client info not found for bgpd")
		return
	}

	logger.Info("Client info for bgpd", clientJson)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	serverTransport, err := thrift.NewTServerSocket("localhost:" + strconv.Itoa(clientJson.Port))
	if err != nil {
		logger.Info("StartServer: NewTServerSocket failed with error:", err)
		return
	}
	processor := bgpd.NewBGPDServicesProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)
	logger.Info("Starting thrift server now...")
	err = server.Serve()
	if err != nil {
		logger.Info("Failed to start the listener, err:", err)
	}
	logger.Info("Starting thrift server failed")
	return
}

/*func createClientIPCHandles(logger *logging.Writer, port string) (thrift.TTransport, thrift.TProtocolFactory, error) {
	var clientTransport thrift.TTransport

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTBufferedTransportFactory(8192)
	clientTransport, err := thrift.NewTSocket("localhost:" + port)
	if err != nil {
		logger.Err("NewTSocket failed with error:", err)
		return nil, nil, err
	}

	clientTransport = transportFactory.GetTransport(clientTransport)
	err = clientTransport.Open()
	return clientTransport, protocolFactory, err
}*/

func connectToClient(logger *logging.Writer, clientTransport thrift.TTransport) error {
	return clientTransport.Open()
}

func StartAsicdClient(logger *logging.Writer, filePath string,
	asicdClient chan *asicdServices.ASICDServicesClient) {
	fileName := filePath + ClientsFileName
	clientJson, err := getClient(logger, fileName, "asicd")
	if err != nil || clientJson == nil {
		asicdClient <- nil
		return
	}

	clientTransport, protocolFactory, err := ipcutils.CreateIPCHandles("localhost:" +
		strconv.Itoa(clientJson.Port))
	if err != nil {
		logger.Infof("Failed to connect to ASICd, retrying until connection is successful")
		count := 0
		ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
		for _ = range ticker.C {
			clientTransport, protocolFactory, err =
				ipcutils.CreateIPCHandles("localhost:" +
					strconv.Itoa(clientJson.Port))
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				logger.Infof("Still can't connect to ASICd, retrying...")
			}
		}
	}

	client := asicdServices.NewASICDServicesClientFactory(clientTransport, protocolFactory)
	asicdClient <- client
}

func StartRibdClient(logger *logging.Writer, filePath string, ribdClient chan *ribd.RIBDServicesClient) {
	fileName := filePath + ClientsFileName
	clientJson, err := getClient(logger, fileName, "ribd")
	if err != nil || clientJson == nil {
		ribdClient <- nil
		return
	}

	clientTransport, protocolFactory, err := ipcutils.CreateIPCHandles("localhost:" + strconv.Itoa(clientJson.Port))
	if err != nil {
		logger.Infof("Failed to connect to RIBd, retrying until connection is successful")
		count := 0
		ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
		for _ = range ticker.C {
			clientTransport, protocolFactory, err = ipcutils.CreateIPCHandles("localhost:" + strconv.Itoa(clientJson.Port))
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				logger.Infof("Still can't connect to RIBd, retrying...")
			}
		}
	}

	client := ribd.NewRIBDServicesClientFactory(clientTransport, protocolFactory)
	ribdClient <- client
}

func StartNdpdClient(logger *logging.Writer, filePath string, ndpdClient chan *ndpd.NDPDServicesClient) {
	fileName := filePath + ClientsFileName
	clientJson, err := getClient(logger, fileName, "ndpd")
	if err != nil || clientJson == nil {
		ndpdClient <- nil
		return
	}

	clientTransport, protocolFactory, err := ipcutils.CreateIPCHandles("localhost:" + strconv.Itoa(clientJson.Port))
	if err != nil {
		logger.Infof("Failed to connect to NDPd, retrying until connection is successful")
		count := 0
		ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
		for _ = range ticker.C {
			clientTransport, protocolFactory, err = ipcutils.CreateIPCHandles("localhost:" + strconv.Itoa(clientJson.Port))
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				logger.Infof("Still can't connect to NDPd, retrying...")
			}
		}
	}

	client := ndpd.NewNDPDServicesClientFactory(clientTransport, protocolFactory)
	ndpdClient <- client
}

func StartBfddClient(logger *logging.Writer, filePath string, bfddClient chan *bfdd.BFDDServicesClient) {
	fileName := filePath + ClientsFileName
	clientJson, err := getClient(logger, fileName, "bfdd")
	if err != nil || clientJson == nil {
		bfddClient <- nil
		return
	}

	clientTransport, protocolFactory, err := ipcutils.CreateIPCHandles("localhost:" + strconv.Itoa(clientJson.Port))
	if err != nil {
		logger.Infof("Failed to connect to BFDd, retrying until connection is successful")
		count := 0
		ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
		for _ = range ticker.C {
			clientTransport, protocolFactory, err = ipcutils.CreateIPCHandles("localhost:" + strconv.Itoa(clientJson.Port))
			if err == nil {
				ticker.Stop()
				break
			}
			count++
			if (count % 10) == 0 {
				logger.Infof("Still can't connect to BFDd, retrying...")
			}
		}
	}

	client := bfdd.NewBFDDServicesClientFactory(clientTransport, protocolFactory)
	bfddClient <- client
}
