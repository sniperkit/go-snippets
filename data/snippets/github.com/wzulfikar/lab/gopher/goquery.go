package main

import "github.com/davecgh/go-spew/spew"

type ParamsAllTransactions struct {
	Address string `url:"address"`
	Hash    string `url:"hash,omitempty"`
	Id      string `url:"hash,omitempty"`
}

var TransactionsPath = struct {
	All         string
	Incoming    string
	Unconfirmed string
}{
	"/account/transfers/all",
	"/account/transfers/incoming",
	"/account/unconfirmedTransactions",
}

var Paths = struct {
	Transactions []apiPath
}{
	Transactions: []apiPath{
		{"incom", "asdfa"},
	},
}

func main() {

	spew.Dump(Paths.Transactions)
}
