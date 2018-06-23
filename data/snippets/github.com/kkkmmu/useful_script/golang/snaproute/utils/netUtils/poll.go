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
	"syscall"
)

const MAXEPOLLEVENTS = 32

type EPoll struct {
	fd     int
	event  syscall.EpollEvent
	events [MAXEPOLLEVENTS]syscall.EpollEvent
}

func NewEPoll(fd int) (*EPoll, error) {
	var err error
	var eFd int
	if eFd, err = syscall.EpollCreate1(0); err != nil {
		return nil, err
	}

	e := EPoll{}
	e.fd = eFd
	e.event.Events = syscall.EPOLLOUT
	e.event.Fd = int32(fd)
	if err = syscall.EpollCtl(eFd, syscall.EPOLL_CTL_ADD, fd, &e.event); err != nil {
		e.Close()
		return nil, err
	}

	return &e, nil
}

func (e *EPoll) Close() error {
	return syscall.Close(e.fd)
}

func (e *EPoll) Wait(msec int) error {
	nevents, err := syscall.EpollWait(e.fd, e.events[:], msec)
	if err != nil {
		return err
	}
	if nevents <= 0 {
		return errors.New("i/o timeout")
	}
	return nil
}
