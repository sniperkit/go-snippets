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
package lockStack

import (
	_ "fmt"
	"runtime"
	"sync"
	"utils/logging"
)

type MyLock struct {
	sync.RWMutex
	Logger *logging.Writer
}

// MyCaller returns the caller of the function that called it :)
func MyCaller(id int) string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(id, fpcs)
	if n == 0 {
		return "n/a" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}

	// return its name
	return fun.Name()
}

func (my *MyLock) Unlock() {
	for idx := 1; idx < 6; idx++ {
		my.Logger.Debug("Releasing Write Lock caller:", MyCaller(idx))
	}
	my.RWMutex.Unlock()
}

func (my *MyLock) Lock() {
	for idx := 1; idx < 6; idx++ {
		my.Logger.Debug("Acquiring Write Lock caller:", MyCaller(idx))
	}
	my.RWMutex.Lock()
}

func (my *MyLock) RUnlock() {
	for idx := 1; idx < 6; idx++ {
		my.Logger.Debug("Releasing Reader Lock caller:", MyCaller(idx))
	}
	my.RWMutex.RUnlock()
}

func (my *MyLock) RLock() {
	for idx := 1; idx < 6; idx++ {
		my.Logger.Debug("Acquiring Reader Lock caller:", MyCaller(idx))
	}
	my.RWMutex.RLock()
}
