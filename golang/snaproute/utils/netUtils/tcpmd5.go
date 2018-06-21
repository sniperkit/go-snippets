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

// tcpmd5.go
package netUtils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"unsafe"
)

type tcpMD5Sig struct {
	ssFamily uint16
	ss       [126]byte
	pad1     uint16
	keylen   uint16
	pad2     uint32
	key      [80]byte
}

func NewTCPMD5Sig() *tcpMD5Sig {
	return &tcpMD5Sig{}
}

func getMD5Sig(ipAddr string, key string) (*tcpMD5Sig, error) {
	var ip []byte
	var family uint16
	idx := 0
	addr := net.ParseIP(ipAddr)
	if addr == nil {
		return nil, errors.New(fmt.Sprintf("%s is not a valid IP", ipAddr))
	}

	sig := NewTCPMD5Sig()
	if ip = addr.To4(); ip != nil {
		//fmt.Println("getMD5Sig - ip", ip, "is a v4 address")
		family = syscall.AF_INET
		idx = 2
	} else if ip = addr.To16(); ip != nil {
		//fmt.Println("getMD5Sig - ip", ip, "is a v6 address")
		family = syscall.AF_INET6
		idx = 6
	}

	if ip == nil {
		return nil, errors.New(fmt.Sprintf("%s is not a valid 4-byte or 16-byte IP", ipAddr))
	}

	sig.ssFamily = family
	copy(sig.ss[idx:], ip)

	sig.keylen = uint16(len(key))
	copy(sig.key[0:], []byte(key))

	return sig, nil
}

func SetSockoptTCPMD5(socket int, ipAddr, key string) error {
	//fmt.Println("SetSockoptTCPMD5 - start, socket", socket, "ip address", ipAddr, "md5 key", key)
	sig, err := getMD5Sig(ipAddr, key)
	if err != nil {
		fmt.Println("SetSockoptTCPMD5 - getMD5Sig failed with error", err)
		return err
	}
	_, _, errNo := syscall.Syscall6(syscall.SYS_SETSOCKOPT, uintptr(socket), uintptr(syscall.IPPROTO_TCP),
		uintptr(14), uintptr(unsafe.Pointer(sig)), unsafe.Sizeof(*sig), 0)

	//fmt.Println("SetSockoptTCPMD5 - syscall.Syscall6 returned", errNo)
	if errNo != 0 {
		err = errNo
	}

	return err
}

func getListenerFile(l *net.TCPListener) (*os.File, error) {
	//fmt.Println("getListenerFd - start")
	file, err := l.File()
	if err != nil {
		fmt.Println("getListenerFd - failed to get File for TCP listener with error", err)
		goto fail
	}

	if listener, err := net.FileListener(file); err == nil {
		defer listener.Close()
	} else {
		//fmt.Println("getListenerFd - failed to get FileListener for File", file, "with error", err)
		goto closefile
	}

	return file, nil

closefile:
	//fmt.Println("getListenerFd - file close")
	file.Close()
fail:
	fmt.Println("getListenerFd - return err")
	return nil, err
}

func SetTCPListenerMD5(l *net.TCPListener, ipAddr, key string) error {
	//fmt.Println("SetTCPListenerMD5 - start ip address:", ipAddr, "key:", key)
	file, err := getListenerFile(l)
	if err != nil {
		fmt.Println("SetTCPListenerMD5 - failed to get listener FD for ip address", ipAddr)
		return err
	}
	defer file.Close()

	return SetSockoptTCPMD5(int(file.Fd()), ipAddr, key)
}
