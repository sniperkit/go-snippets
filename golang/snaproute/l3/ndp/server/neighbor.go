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
package server

import (
	"l3/ndp/config"
	"l3/ndp/debug"
	"math/rand"
	"strings"
	"time"
)

type NDP_OPERATION byte

const (
	IGNORE NDP_OPERATION = 1
	CREATE NDP_OPERATION = 2
	DELETE NDP_OPERATION = 3
	UPDATE NDP_OPERATION = 4
)

const (
	MAX_UNICAST_SOLICIT            uint8 = 3
	MAX_MULTICAST_SOLICIT                = 3
	MAX_ANYCAST_DELAY_TIMER              = 1
	MAX_NEIGHBOR_ADVERTISEMENT           = 3
	DELAY_FIRST_PROBE_TIME               = 5 // this is in seconds
	MIN_RANDOM_FACTOR                    = 0.5
	MAX_RANDOM_FACTOR                    = 1.5
	RECOMPUTE_BASE_REACHABLE_TIMER       = 1 // this is in hour
)

const (
	_ = iota
	INCOMPLETE
	REACHABLE
	STALE
	DELAY
	PROBE
)

type NeighborInfo struct {
	BaseReachableTimer   float32
	RetransTimerConfig   uint32
	ReachableTimeConfig  uint32
	RecomputeBaseTimer   *time.Timer
	ReachableTimer       *time.Timer
	RetransTimer         *time.Timer
	DelayFirstProbeTimer *time.Timer
	InvalidationTimer    *time.Timer
	FastProbeTimer       *time.Timer
	ProbesSent           uint8
	FastProbesMultiplier uint8
	StopFastProbe        bool
	State                int
	LinkLayerAddress     string // this is our neighbor port mac address
	IpAddr               string
	ReturnCh             chan config.PacketData // NDP Server communicator
	IfIndex              int32                  // Physical Port where to send the packet on timer expiry
	counter              PktCounter
	pktRcvdTime          time.Time // last received packet time
}

/*
 *  Delete Cache completely
 */
func (c *NeighborInfo) DeInit() {
	// stopping all three timers in accending order
	debug.Logger.Debug("De-Init neighbor", c.IpAddr)
	c.StopReTransmitTimer()
	c.StopReachableTimer()
	c.StopFastProbeTimer()
	c.StopDelayProbeTimer()
	c.StopInvalidTimer()
	c.StopReComputeBaseTimer()
	c.counter.Rcvd = 0
	c.counter.Send = 0
}

/*
 *  Helper function to randomize BASE_REACHABLE_TIME
 */
func computeBase(reachableTime uint32) float32 {
	return float32(reachableTime) + ((rand.Float32() * MIN_RANDOM_FACTOR) + MIN_RANDOM_FACTOR)
}

/*
 *  Initialize cache with default values..
 */
func (c *NeighborInfo) InitCache(reachableTime, retransTime uint32, nbrKey string, pktCh chan config.PacketData, ifIndex int32) {
	c.ReachableTimeConfig = reachableTime
	c.RetransTimerConfig = retransTime
	c.BaseReachableTimer = computeBase(reachableTime)
	// Set the multiplier
	c.FastProbesMultiplier = 1
	c.State = INCOMPLETE
	ipMacStr := strings.Split(nbrKey, "_")
	c.IpAddr = ipMacStr[1]
	c.LinkLayerAddress = ipMacStr[0] // this is mac address
	c.ProbesSent = uint8(0)
	c.ReturnCh = pktCh
	c.IfIndex = ifIndex
	// Once initalized start reachable timer... And also start one hour timer for re-computing BaseReachableTimer
	// Reachable timer will handle Fast Probe Timer
	c.RchTimer()
	c.ReComputeBaseReachableTimer()
	c.counter.Rcvd = 0
	c.counter.Send = 0
	debug.Logger.Debug("Neighbor timers are ReachableTimeConfig:", c.ReachableTimeConfig, "RetransTimerConfig:", c.RetransTimerConfig,
		"BaseReachableTimer:", c.BaseReachableTimer)
}

func (nbr *NeighborInfo) populateNbrInfo(ifIndex int32, intfRef string) *config.NeighborConfig {
	nbrInfo := &config.NeighborConfig{}
	nbrInfo.IpAddr = nbr.IpAddr
	nbrInfo.MacAddr = nbr.LinkLayerAddress
	nbrInfo.IfIndex = ifIndex
	nbrInfo.Intf = intfRef
	return nbrInfo
}

func (nbr *NeighborInfo) updatePktRxStateInfo() {
	nbr.pktRcvdTime = time.Now()
	nbr.counter.Rcvd++
}
