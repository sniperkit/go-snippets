//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package ringBuffer

const (
	DefCapacity int = 10
)

type RingBuffer struct {
	wPtr   int
	rPtr   int
	buffer []interface{}
}

func (rB *RingBuffer) SetRingBufferCapacity(size int) {
	rB.verifyInit()
	rB.IncCapacity(size)
}

func (rB *RingBuffer) verifyInit() {
	if rB.buffer == nil {
		rB.buffer = make([]interface{}, DefCapacity)
		for i := range rB.buffer {
			rB.buffer[i] = nil
		}
		rB.wPtr, rB.rPtr = -1, 0
	}
}

func (rB *RingBuffer) IncCapacity(size int) {
	if size == len(rB.buffer) {
		return
	} else if size < len(rB.buffer) {
		rB.buffer = rB.buffer[0:size]
	}
	newbuffer := make([]interface{}, size-len(rB.buffer))
	for i := range newbuffer {
		newbuffer[i] = nil
	}
	rB.buffer = append(rB.buffer, newbuffer...)
}

func (rB *RingBuffer) GetRingBufferCapacity() int {
	return len(rB.buffer)
}

func (rB *RingBuffer) Set(ptr int, intf interface{}) {
	p := rB.Modulo(ptr)
	rB.buffer[p] = intf
}

func (rB *RingBuffer) Modulo(idx int) int {
	return idx % len(rB.buffer)
}

func (rB *RingBuffer) Get(ptr int) interface{} {
	p := rB.Modulo(ptr)
	return rB.buffer[p]
}

func (rB *RingBuffer) InsertIntoRingBuffer(intf interface{}) (int, interface{}) {
	rB.verifyInit()
	oldVal := rB.Get(rB.wPtr + 1)
	rB.Set(rB.wPtr+1, intf)
	old := rB.wPtr
	rB.wPtr = rB.Modulo(rB.wPtr + 1)
	if old != -1 && rB.wPtr == rB.rPtr {
		rB.rPtr = rB.Modulo(rB.rPtr + 1)
	}
	return rB.wPtr, oldVal
}

func (rB *RingBuffer) DeleteFromRingBuffer() interface{} {
	rB.verifyInit()
	if rB.wPtr == -1 {
		return nil
	}
	intf := rB.Get(rB.rPtr)
	if rB.wPtr == rB.rPtr {
		rB.wPtr = -1
		rB.rPtr = 0
	} else {
		rB.rPtr = rB.Modulo(rB.rPtr + 1)
	}
	return intf
}

func (rB *RingBuffer) PeekIntoRingBuffer() interface{} {
	rB.verifyInit()
	if rB.wPtr == -1 {
		return nil
	}
	intf := rB.Get(rB.rPtr)
	return intf
}

func (rB *RingBuffer) GetListOfEntriesFromRingBuffer() []interface{} {
	if rB.wPtr == -1 {
		return []interface{}{}
	}

	intfSlice := make([]interface{}, 0, rB.GetRingBufferCapacity())
	for idx := 0; idx < rB.GetRingBufferCapacity(); idx++ {
		ptr := rB.Modulo(rB.rPtr + idx)
		intfSlice = append(intfSlice, rB.buffer[ptr])
		if ptr == rB.wPtr {
			break
		}
	}
	return intfSlice
}

func (rB *RingBuffer) UpdateEntryInRingBuffer(intf interface{}, idx int) {
	ptr := rB.Modulo(idx)
	rB.buffer[ptr] = intf
}

func (rB *RingBuffer) GetEntryFromRingBuffer(idx int) interface{} {
	ptr := rB.Modulo(idx)
	return rB.buffer[ptr]
}

func (rB *RingBuffer) FlushRingBuffer() {
	rB.verifyInit()
	if rB.wPtr == -1 {
		return
	}

	rB.wPtr = -1
	rB.rPtr = 0
	return
}
