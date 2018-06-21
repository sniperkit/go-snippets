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

package netUtils

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func resolveAddress(network, address string) (net.Addr, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, address)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, address)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func favoriteAddrFamily(net string, laddr, raddr net.Addr, mode string) (family int, ipv6only bool) {
	switch net[len(net)-1] {
	case '4':
		return syscall.AF_INET, false
	case '6':
		return syscall.AF_INET6, true
	}

	if mode == "listen" && (laddr == nil || isWildcard(laddr)) {
		/*
			if supportsIPv4map {
				return syscall.AF_INET6, false
			}
		*/
		if laddr == nil {
			return syscall.AF_INET, false
		}
		return getFamily(laddr), false
	}

	if (laddr == nil || getFamily(laddr) == syscall.AF_INET) &&
		(raddr == nil || getFamily(raddr) == syscall.AF_INET) {
		return syscall.AF_INET, false
	}
	return syscall.AF_INET6, false
}

func Socket(family, sotype, proto int) (int, error) {
	s, err := syscall.Socket(family, sotype|syscall.SOCK_NONBLOCK|syscall.SOCK_CLOEXEC, proto)
	// On Linux the SOCK_NONBLOCK and SOCK_CLOEXEC flags were
	// introduced in 2.6.27 kernel and on FreeBSD both flags were
	// introduced in 10 kernel. If we get an EINVAL error on Linux
	// or EPROTONOSUPPORT error on FreeBSD, fall back to using
	// socket without them.
	switch err {
	case nil:
		return s, nil
	default:
		return -1, os.NewSyscallError("socket", err)
	case syscall.EPROTONOSUPPORT, syscall.EINVAL:
	}

	// See ../syscall/exec_unix.go for description of ForkLock.
	syscall.ForkLock.RLock()
	s, err = syscall.Socket(family, sotype, proto)
	if err == nil {
		syscall.CloseOnExec(s)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return -1, os.NewSyscallError("socket", err)
	}
	if err = syscall.SetNonblock(s, true); err != nil {
		syscall.Close(s)
		return -1, os.NewSyscallError("setnonblock", err)
	}
	return s, nil
}

func CloseSocket(socket int) {
	syscall.Close(socket)
}

func ConnectSocket(network, remote, local string) (int, error) {
	//fmt.Println("ConnectSocket: network=", network, "remote =", remote, "local =", local)
	var localAddr net.Addr
	var socketType, proto int
	netAddr, err := resolveAddress(network, remote)
	if err != nil {
		fmt.Println("ConnectSocket: resolveAddress for remote failed with error", err)
		return -1, err
	}

	if local != "" {
		localAddr, err = resolveAddress(network, local)
		if err != nil {
			fmt.Println("ConnectSocket: resolveAddress for local failed with error", err)
			return -1, err
		}
	}

	family, ipv6only := favoriteAddrFamily(network, localAddr, netAddr, "dial")
	proto = 0

	switch netAddr := netAddr.(type) {
	case *net.TCPAddr:
		localAddr, _ = localAddr.(*net.TCPAddr)
		socketType = syscall.SOCK_STREAM
	case *net.UDPAddr:
		localAddr, _ = localAddr.(*net.UDPAddr)
		socketType = syscall.SOCK_DGRAM
	default:
		//fmt.Println("ConnectSocket: remote is not TCPAddr or UDPAddr")
		return -1, &net.OpError{Op: "dial", Net: network, Source: localAddr, Addr: netAddr,
			Err: &net.AddrError{Err: "unexpected address type", Addr: remote}}
	}

	socket, err := Socket(family, socketType, proto)
	if err != nil {
		fmt.Println("ConnectSocket: Socket call failed")
		return -1, err
	}
	if err = SetDefaultConnectSockopts(socket); err != nil {
		fmt.Println("ConnectSocket: SetDefaultConnectSockopts failed")
		CloseSocket(socket)
		return -1, err
	}
	if err = SetSockoptIPv6Only(socket, family, socketType, ipv6only); err != nil {
		fmt.Println("ConnectSocket: SetSockoptIPv6Only failed")
		CloseSocket(socket)
		return -1, err
	}

	return socket, err
}

func SetSockoptIPv6Only(s, family, sotype int, ipv6only bool) (err error) {
	if family == syscall.AF_INET6 && sotype != syscall.SOCK_RAW {
		// Allow both IP versions even if the OS default
		// is otherwise.  Note that some operating systems
		// never admit this option.
		err = syscall.SetsockoptInt(s, syscall.IPPROTO_IPV6, syscall.IPV6_V6ONLY, boolint(ipv6only))
	}

	return err
}

func setDefaultSockopts(s int) error {
	// Allow broadcast.
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1))
}

func SetDefaultConnectSockopts(s int) error {
	return setDefaultSockopts(s)
}

func SetDefaultListenerSockopts(s int) error {
	// Allow reuse of recently-used addresses.
	err := setDefaultSockopts(s)
	if err != nil {
		return err
	}
	return os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1))
}

func Connect(socket int, network, remote, local string, timeout time.Duration) error {
	var lAddr net.Addr
	var rsa, lsa syscall.Sockaddr
	var deadline time.Time

	if timeout != 0 {
		deadline = time.Now().Add(timeout)
	}

	rAddr, err := resolveAddress(network, remote)
	if err != nil {
		return err
	}

	if local != "" {
		lAddr, err = resolveAddress(network, local)
		if err != nil {
			return err
		}
	}

	family, _ := favoriteAddrFamily(network, lAddr, rAddr, "dial")

	if lAddr != nil {
		lsa, err = sockaddr(lAddr, family)
		if err != nil {
			return err
		}

		if err := syscall.Bind(socket, lsa); err != nil {
			return os.NewSyscallError("bind", err)
		}
	}

	rsa, err = sockaddr(rAddr, family)
	if err != nil {
		return err
	}

	// Do not need to call fd.writeLock here,
	// because fd is not yet accessible to user,
	// so no concurrent operations are possible.
	switch err := syscall.Connect(socket, rsa); err {
	case syscall.EINPROGRESS, syscall.EALREADY, syscall.EINTR:
	case nil, syscall.EISCONN:
		if !deadline.IsZero() && deadline.Before(time.Now()) {
			return errors.New("i/o timeout")
		}
		return nil
	default:
		return os.NewSyscallError("connect", err)
	}

	ePoll, err := NewEPoll(socket)
	if err != nil {
		return err
	}
	defer ePoll.Close()

	for {
		waitTimeout := -1
		if !deadline.IsZero() {
			if deadline.Before(time.Now()) {
				return errors.New("i/o timeout")
			}
			duration := deadline.Sub(time.Now())
			waitTimeout = int(duration.Nanoseconds()) / 1000000
		}
		// Performing multiple connect system calls on a
		// non-blocking socket under Unix variants does not
		// necessarily result in earlier errors being
		// returned. Instead, once runtime-integrated network
		// poller tells us that the socket is ready, get the
		// SO_ERROR socket option to see if the connection
		// succeeded or failed. See issue 7474 for further
		// details.
		if err := ePoll.Wait(waitTimeout); err != nil {
			return err
		}
		nerr, err := syscall.GetsockoptInt(socket, syscall.SOL_SOCKET, syscall.SO_ERROR)
		if err != nil {
			return os.NewSyscallError("getsockopt", err)
		}

		switch err := syscall.Errno(nerr); err {
		case syscall.EINPROGRESS, syscall.EALREADY, syscall.EINTR:
		case syscall.Errno(0), syscall.EISCONN:
			return nil
		default:
			return os.NewSyscallError("getsockopt", err)
		}
	}
}

func ConvertFdToConn(socket int) (net.Conn, error) {
	file := os.NewFile(uintptr(socket), "")
	conn, err := net.FileConn(file)
	if err != nil {
		err1 := file.Close()
		if err1 != nil {
			return nil, errors.New(fmt.Sprintf("Failed to create net.Conn and close intermediate file, error:%s", err))
		}
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to close intermediate file with error:%s", err))
	}
	return conn, nil
}
