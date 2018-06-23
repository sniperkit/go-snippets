// Copyright 2016-2018 Hyperchain Corp.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"os"
	"fmt"
)

type Account struct {
	Name string
	Age  int
}

var data map[string]Account = map[string]Account{
	"user1": {"zhangsan", 18},
	"user2": {"lisi", 20},
	"user3": {"wangwu", 31},
}

func streamEncoder() error {
	file, err := os.Create("json.dat")
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	for _, v := range data {
		encoder.Encode(v)
	}
	return nil
}

func streamDecoder() error {
	file, err := os.Open("json.dat")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var account Account
		decoder.Decode(&account)
		fmt.Println(account)
	}
	return nil
}

func main() {
	err := streamEncoder()
	if err != nil {
		fmt.Println(err)
		os.Exit(1);
	}
	fmt.Println("encode done")

	err = streamDecoder()
	if err != nil {
		fmt.Println(err)
		os.Exit(1);
	}
	fmt.Println("decode done")
}

