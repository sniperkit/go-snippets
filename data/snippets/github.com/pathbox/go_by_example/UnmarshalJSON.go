package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type IntBool bool

func (p *IntBool) UnmarshalJSON(b []byte) error {
	str := string(b)
	switch str {
	case "1":
		*p = true
	case "0":
		*p = false
	default:
		return fmt.Errorf("unexpected bool: %s", str)
	}
	return nil
}

func main() {
	var account struct {
		Gender IntBool `json:"gender"`
	}

	err := json.Unmarshal([]byte(`{"gender":1}`), &account)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(account.Gender == true)
}
