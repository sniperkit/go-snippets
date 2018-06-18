package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const explorer = "http://btc.blockr.io/api/v1/address/info/"

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

func FetchUrlByte(url string) []byte {

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal("Error while fetching url\n", err)
	}

	request.Header.Set("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) 
									Ubuntu Chromium/53.0.2785.143 Chrome/53.0.2785.143 Safari/537.36`)
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error while trying to get response\n", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatal("Error status not OK!\n", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal("Error reading body\n", err)
	}

	return body
}

func LoadJsonFromUrl(url string) Response {
	body := FetchUrlByte(url)
	res := Response{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal("Unmarchal failed !\n", err)
	}

	return res
}

func ReadFromFile(path string) []string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading the file\n", err)
	}
	return strings.Split(string(data), "\n")
}

func main() {
	if len(os.Args) > 2 || len(os.Args) < 2 {
		log.Fatal(`
		This current script accept only one argument.
		Usage: ./GoCheckBitcoinAddress addresses_path`)
	}

	data := ReadFromFile(os.Args[1])
	length := len(data) - 1
	fmt.Printf("\t\tWe have %d addresses to check their balances\n\n", length)
	fmt.Println("\tAddress\t\t\t \t\tBalance\t\t\tETA")
	fmt.Println("----------------------------------\t---------------------\t\t-------------")

	i, j := 0, 1

	for _, value := range data {
		if value == "" {
			continue
		} else {
			url := explorer + value
			res := LoadJsonFromUrl(url)
			fmt.Printf("\033[92m%s\033[0m\t\033[95m%.8f\tBTC\033[0m\t\tETA(%%): %.2f\n", res.Data.Address, res.Data.Balance, float64(j*100/length))
			if i == 5 {
				time.Sleep(1000 * time.Millisecond) // Wait 1s in order to escape Blockr.io's API restriction
				i = 0
			}
			i++
			j++
		}
	}
}
