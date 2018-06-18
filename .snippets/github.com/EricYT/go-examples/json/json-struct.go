package main

import (
	"encoding/json"
	"fmt"
)

type WritePlan struct {
	ObjectId string   `json:"object_id"`
	Actions  []Action `json:"actions"`
}

type Action struct {
	ActionType string `json:"action_type"`
	Data       Data   `json:"data"`
}

type Data struct {
	// if field is empty that will be omited
	Disks          string `json:"disks,omitempty"`
	ExtentId       string `json:"extent_id,omitempty"`
	ExtentGroupId  string `json:"extent_group_id,omitempty"`
	OffsetAtExtent int64  `json:"offset_at_extent,omitempty"`

	// this will not omit the empty fields
	Crc []byte `json:"crc"`
}

func main() {
	fmt.Println("json go")
	data := `{"object_id":"o1", "actions":[{"action_type":"gopher", "data":{"disks":"d1", "extent_group_id":"eg1"}}]}`

	var action WritePlan
	err := json.Unmarshal([]byte(data), &action)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json Unmarshal is:%+v\n", action)

	output, err := json.Marshal(&action)
	if err != nil {
		panic(err)
	}
	fmt.Printf("json marshal is: %s\n", string(output))
}
