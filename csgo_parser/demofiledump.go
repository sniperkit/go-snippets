package main

import (
	"fmt"
	"log"
)

const (
	NET_MAX_PAYLOAD         int32 = (262144 - 4)      // largest message that can be sent in bytes
	DEMO_RECORD_BUFFER_SIZE int32 = (2 * 1024 * 1024) // should fit string tables and server classes
	MAX_PLAYER_NAME_LENGTH  int32 = 128
	SIGNED_GUID_LEN         int32 = 32 // hashed CD key
	MAX_CUSTOM_FILES        int32 = 4
)

var (
	matchStartOccurred bool = false
	currentTick        int32
)

type StringTableData struct {
	Name       []byte
	MaxEntries int
}

// Player's information on server
type PlayerInfo struct {
	version         uint64
	xuid            uint64
	name            [MAX_PLAYER_NAME_LENGTH]byte
	userID          int
	guid            [SIGNED_GUID_LEN + 1]byte
	friendsID       uint32
	friendsName     [MAX_PLAYER_NAME_LENGTH]byte
	fakePlayer      bool
	isHLTV          bool
	customFiles     [MAX_CUSTOM_FILES]CRC32
	filesDownloaded uint8
}

type DemoFileDump struct {
	demoFile    DemoFile
	frameNumber int
}

func (d *DemoFileDump) Open(fileName string) bool {
	if d.demoFile.Open(fileName) == false {
		log.Fatal("Unable to open demo file")
		return false
	}

	return true
}

func (d *DemoFileDump) DoDump() {
	matchStartOccurred = false
	var demoFinished bool = false

	fmt.Println("--- Header debug ---")
	d.demoFile.debugHeader()
	fmt.Println("--------------------\n")

	for !demoFinished {
		fmt.Println("--- Read loop ---")
		// get data from command header
		var cmd, playerSlot uint8
		var tick int32 = 0
		d.demoFile.ReadCmdHeader(&cmd, &tick, &playerSlot)

		fmt.Println("Old tick: ", currentTick, " new tick: ", tick)

		currentTick = tick

		// command handling
		switch cmd {
		case dem_stop:
			demoFinished = true
			fmt.Println("Stopping demo at tick: ", tick)
		case dem_consolecmd:
			// read raw data
			d.demoFile.ReadRawData(nil, 0)
		case dem_datatables:
			// read some data
			//var data []byte = make([]byte, DEMO_RECORD_BUFFER_SIZE)
			// parse data
		case dem_stringtables:
			// read raw data
			// dump string tables
		case dem_usercmd:
			// read user command
			var dummy int32
			d.demoFile.ReadUserCmd(nil, &dummy)
		case dem_signon, dem_packet, dem_synctick:
			// handle packet
		}

	}
}
