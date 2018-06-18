/* Expected output:

{success {198aMn6ZYAczwrE5NvNTUMyJ5qkfy4g3Hi false 8000.00236957 0 8000.00236957 33 {2009-02-22T10:50:53Z 0f0fbcc18fd0d090ad3402574df8404cec1176bc000f9aa0dc19f8d832ff94db 0 400 432048} {2016-08-11T16:33:24Z c653ac1e7a66117e097dd16cfd122b24b0f92060d096d8754a1e550d4c64f520 0 0.0002 12549} true} 200 }
>------------------------------<
Status: success
Code: 200
Message:
Address: 198aMn6ZYAczwrE5NvNTUMyJ5qkfy4g3Hi
Is known: false
Balance: 8000.002370
Balance Multi Signature: 0.000000
Total Recieved: 8000.002370
NB Transactions: 33.000000
First Transaction Tx: 0f0fbcc18fd0d090ad3402574df8404cec1176bc000f9aa0dc19f8d832ff94db
Last Transaction Tx: c653ac1e7a66117e097dd16cfd122b24b0f92060d096d8754a1e550d4c64f520
*/

/*
NB: The default concrete Go types are:
	- bool for JSON booleans
	- float64 for JSON numbers
	- string for JSON strings
	- nil for JSON null
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

var my_json string = `{
		"status":"success",
		"data":{
			"address":"198aMn6ZYAczwrE5NvNTUMyJ5qkfy4g3Hi",
			"is_unknown":false,
			"balance":8000.00236957,
			"balance_multisig":0,
			"totalreceived":8000.00236957,
			"nb_txs":33,
			"first_tx":
				{
					"time_utc":"2009-02-22T10:50:53Z",
					"tx":"0f0fbcc18fd0d090ad3402574df8404cec1176bc000f9aa0dc19f8d832ff94db",
					"block_nb":"5219",
					"value":400,
					"confirmations":432048
				},
			"last_tx":
				{
					"time_utc":"2016-08-11T16:33:24Z",
					"tx":"c653ac1e7a66117e097dd16cfd122b24b0f92060d096d8754a1e550d4c64f520",
					"block_nb":"424718",
					"value":0.0002,
					"confirmations":12549
				},
			"is_valid":true
		},
		"code":200,
		"message":""
}`

type Response struct {
	Status  string       `json:"status"`
	Data    ResponseData `json:"data"`
	Code    float64      `json:"code"`
	Message string       `json:"message"`
}

type ResponseData struct {
	Address         string           `json:"address"`
	IsKnown         bool             `json:"is_unknown"`
	Balance         float64          `json:"balance"`
	BalanceMultiSig float64          `json:"balance_multisig"`
	TotalRecieved   float64          `json:"totalreceived"`
	NbTxs           float64          `json:"nb_txs"`
	FirstTxs        ResponseFirstTxs `json:"first_tx"`
	LastTxs         ResponseLastTxs  `json:"last_tx"`
	IsValid         bool             `json:"is_valid"`
}

type ResponseFirstTxs struct {
	Time          string  `json:"time_utc"`
	Tx            string  `json:"tx"`
	BlockNb       string  `json:"block_nb"`
	Value         float64 `json:"value"`
	Confirmations int64   `json:"confirmations"`
}

type ResponseLastTxs struct {
	Time          string  `json:"time_utc"`
	Tx            string  `json:"tx"`
	BlockNb       string  `json:block_nb"`
	Value         float64 `json:value"`
	Confirmations int64   `json:confirmations"`
}

func main() {
	res := Response{}
	err := json.Unmarshal([]byte(my_json), &res)
	if err != nil {
		log.Fatal("Unmarchal failed!\n", err)
	}
	fmt.Println(res)
	fmt.Println(">------------------------------<")
	fmt.Printf("Status: %s\n", res.Status)
	fmt.Printf("Code: %.0f\n", res.Code)
	fmt.Printf("Message: %s\n", res.Message)
	fmt.Printf("Address: %s\n", res.Data.Address)
	fmt.Printf("Is known: %v\n", res.Data.IsKnown)
	fmt.Printf("Balance: %f\n", res.Data.Balance)
	fmt.Printf("Balance Multi Signature: %f\n", res.Data.BalanceMultiSig)
	fmt.Printf("Total Recieved: %f\n", res.Data.TotalRecieved)
	fmt.Printf("NB Transactions: %f\n", res.Data.NbTxs)
	fmt.Printf("First Transaction Tx: %s\n", res.Data.FirstTxs.Tx)
	fmt.Printf("Last Transaction Tx: %s\n", res.Data.LastTxs.Tx)

}
