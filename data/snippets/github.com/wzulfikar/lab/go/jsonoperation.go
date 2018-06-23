package main

import (
	"encoding/json"
	"fmt"
)

type context struct {
	Ctx  string            `json:"ctx"`
	Data map[string]string `json:"data"`
}

func main() {
	ctx := &context{
		"merchant.checkout",
		map[string]string{
			"hello": "asa",
		},
	}

	b, err := json.Marshal(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("json of ctx:", string(b))

	ctxJson := []byte(`{"ctx":"merchant.checkout","data":{"hello":"from json"}}`)
	if err := json.Unmarshal(ctxJson, ctx); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("data of ctx", ctx)
}
