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

package keepalive

import (
	"encoding/json"
	"fmt"
	nanomsg "github.com/op/go-nanomsg"
	"infra/sysd/sysdCommonDefs"
	"io/ioutil"
	"strconv"
	"sysd"
	"time"
	"utils/ipcutils"
)

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type SysdClient struct {
	ipcutils.IPCClientBase
	ClientHdl *sysd.SYSDServicesClient
}

type KeepAlive struct {
	sysdClient SysdClient
	name       string
	status     int32
}

const (
	KA_ACTIVE   = 1 // Default status of a daemon
	KA_INTERVAL = 1 // KA message will be sent every 1 second
)

type DaemonStatusNotifier struct {
	DaemonStatusCh chan sysdCommonDefs.DaemonStatus
	subSocket      *nanomsg.SubSocket
	socketCh       chan []byte
}

func InitKeepAlive(name string, paramsDir string) {
	var clientsList []ClientJson

	paramsFile := paramsDir + "clients.json"
	bytes, err := ioutil.ReadFile(paramsFile)
	if err != nil {
		fmt.Println("Error in reading configuration file")
		return
	}

	err = json.Unmarshal(bytes, &clientsList)
	if err != nil {
		fmt.Println("Error in Unmarshalling Json")
		return
	}

	ka := new(KeepAlive)
	ka.name = name
	ka.status = KA_ACTIVE

	for _, client := range clientsList {
		if client.Name == "sysd" {
			fmt.Println("found sysd at port", client.Port)
			ka.sysdClient.Address = "localhost:" + strconv.Itoa(client.Port)
			ka.sysdClient.TTransport, ka.sysdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(ka.sysdClient.Address)
			if err != nil {
				fmt.Println(ka.name, " Failed to connect to sysd, retrying until connection is successful")
				count := 0
				ticker := time.NewTicker(time.Duration(1000) * time.Millisecond)
				for _ = range ticker.C {
					ka.sysdClient.TTransport, ka.sysdClient.PtrProtocolFactory, err = ipcutils.CreateIPCHandles(ka.sysdClient.Address)
					if err == nil {
						ticker.Stop()
						break
					}
					count++
					if (count % 100) == 0 {
						fmt.Println(ka.name, " Still can't connect to sysd, retrying...")
					}
				}
			}
		}
	}
	if ka.sysdClient.TTransport != nil && ka.sysdClient.PtrProtocolFactory != nil {
		ka.sysdClient.ClientHdl = sysd.NewSYSDServicesClientFactory(ka.sysdClient.TTransport, ka.sysdClient.PtrProtocolFactory)
		ka.sysdClient.IsConnected = true
		fmt.Println(ka.name, " connected to sysd")
	}
	if ka.sysdClient.ClientHdl != nil && ka.sysdClient.IsConnected {
		fmt.Println("Initialized KA for ", ka.name, " status ", ka.status)
		retryTimer := time.NewTicker(time.Second * KA_INTERVAL)
		for t := range retryTimer.C {
			_ = t
			ka.sysdClient.ClientHdl.PeriodicKeepAlive(ka.name)
		}
	} else {
		fmt.Println("Failed to start KA for ", ka.name)
	}
	return
}

func (statusNotifier *DaemonStatusNotifier) ProcessStatusNotifications(rxBuf []byte) error {
	var msg sysdCommonDefs.Notification
	err := json.Unmarshal(rxBuf, &msg)
	if err != nil {
		return err
	}
	if msg.Type == sysdCommonDefs.KA_DAEMON {
		var dStatus sysdCommonDefs.DaemonStatus
		err = json.Unmarshal(msg.Payload, &dStatus)
		if err == nil {
			statusNotifier.DaemonStatusCh <- dStatus
		} else {
			return err
		}
	}
	return nil
}

func (statusNotifier *DaemonStatusNotifier) ReceiveStatusNotifications() error {
	for {
		select {
		case rxBuf := <-statusNotifier.socketCh:
			if rxBuf != nil {
				statusNotifier.ProcessStatusNotifications(rxBuf)
			}
		}
	}
	return nil
}

func (statusNotifier *DaemonStatusNotifier) StartDaemonStatusListner() error {
	go statusNotifier.ReceiveStatusNotifications()
	for {
		rxBuf, err := statusNotifier.subSocket.Recv(0)
		if err == nil {
			statusNotifier.socketCh <- rxBuf
		}
	}
	return nil
}

func (statusNotifier *DaemonStatusNotifier) SetupDaemonStatusSub() error {
	var err error
	if statusNotifier.subSocket, err = nanomsg.NewSubSocket(); err == nil {
		if err = statusNotifier.subSocket.Subscribe(""); err == nil {
			if _, err = statusNotifier.subSocket.Connect(sysdCommonDefs.PUB_SOCKET_ADDR); err == nil {
				if err = statusNotifier.subSocket.SetRecvBuffer(1024 * 1024); err == nil {
					return nil
				}
			}
		}
	}
	return err
}

func InitDaemonStatusListener() *DaemonStatusNotifier {
	statusNotifier := new(DaemonStatusNotifier)
	statusNotifier.socketCh = make(chan []byte, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	statusNotifier.DaemonStatusCh = make(chan sysdCommonDefs.DaemonStatus, sysdCommonDefs.SYSD_TOTAL_KA_DAEMONS)
	err := statusNotifier.SetupDaemonStatusSub()
	if err == nil {
		return statusNotifier
	}
	return nil
}
