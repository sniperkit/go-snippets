package rpc

import (
	"bytes"
	"encoding/hex"
	"strconv"
)

// Value .
type Value struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// Nep5Result .
type Nep5Result struct {
	State       string   `json:"state"`
	GasConsumed string   `json:"gas_consumed"`
	Script      string   `json:"script"`
	Stack       []*Value `json:"stack"`
}

// Asset .
type Asset struct {
	Asset string `json:"asset"`
	Value string `json:"value"`
}

// Balance .
type Balance struct {
	Balance   string `json:"balance"`
	Confirmed string `json:"confirmed"`
}

// AccountSate neo account state
type AccountSate struct {
	Version    int         `json:"version"`
	ScriptHash string      `json:"script_hash"`
	Frozen     bool        `json:"frozen"`
	Votes      interface{} `json:"votes"`
	Balances   []Asset     `json:"balances"`
}

// L10NString localization string
type L10NString struct {
	Lang string `json:"lang"`
	Name string `json:"name"`
}

// AssetState neo asset state
type AssetState struct {
	Version    int          `json:"version"`
	ID         string       `json:"id"`
	Type       string       `json:"type"`
	Name       []L10NString `json:"name"`
	Amount     string       `json:"amount"`
	Available  string       `json:"available"`
	Precision  float64      `json:"precision"`
	Owner      string       `json:"owner"`
	Admin      string       `json:"admin"`
	Issuer     string       `json:"issuer"`
	Expiration float64      `json:"expiration"`
	Frozen     bool         `json:"frozen"`
}

// Script .
type Script struct {
	Invocation   string `json:"Invocation"`
	Verification string `json:"Verification"`
}

// Block .
type Block struct {
	Confirmations     int64         `json:"Confirmations"`
	Hash              string        `json:"Hash"`
	Index             int64         `json:"Index"`
	Merkleroot        string        `json:"Merkleroot"`
	NextBlockHash     string        `json:"Nextblockhash"`
	NextConsensus     string        `json:"Nextconsensus"`
	Nonce             string        `json:"Nonce"`
	PreviousBlockHash string        `json:"Previousblockhash"`
	Size              int64         `json:"Size"`
	Time              int64         `json:"Time"`
	Version           int64         `json:"Version"`
	Script            Script        `json:"Script"`
	Transactions      []Transaction `json:"Tx"`
}

// Vout .
type Vout struct {
	Address string `json:"Address"`
	Asset   string `json:"Asset"`
	N       int    `json:"N"`
	Value   string `json:"Value"`
}

// Vin .
type Vin struct {
	TransactionID string `json:"Txid"`
	Vout          int    `json:"Vout"`
}

// Claim .
type Claim struct {
	TransactionID string `json:"txid"`
	Vout          int    `json:"vout"`
}

// Transaction .
type Transaction struct {
	ID         string        `json:"Txid"`
	Size       int64         `json:"Size"`
	Type       string        `json:"Type"`
	Version    int64         `json:"Version"`
	Attributes []interface{} `json:"Attributes"` //
	Script     interface{}   `json:"Script"`
	Gas        string        `json:"Gas"`
	Claims     []Claim       `json:"Claims"`
	Vin        []Vin         `json:"Vin"`
	Vout       []Vout        `json:"Vout"`
	SysFee     string        `json:"Sys_fee"`
	NetFee     string        `json:"Net_fee"`
	Scripts    []Script      `json:"Scripts"`
	Nonce      int64         `json:"Nonce"`
}

// UTXO .
type UTXO struct {
	TransactionID string `json:"txid"`
	Vout          Vout   `json:"vout"`
	CreateTime    string `json:"createTime"`
	SpentTime     string `json:"spentTime"`
	Block         int64  `json:"block"`
	SpentBlock    int64  `json:"spentBlock"`
	value         *float64
	Gas           string `json:"gas"`
}

// Value get utxo value
func (utxo *UTXO) Value() (float64, error) {

	if utxo.value == nil {

		val, err := strconv.ParseFloat(utxo.Vout.Value, 64)

		if err != nil {
			return 0, err
		}

		utxo.value = &val
	}

	return *utxo.value, nil
}

// TxHex get utxo txid hex value
func (utxo *UTXO) TxHex() ([]byte, error) {
	return hex.DecodeString(utxo.TransactionID)
}

type utxoByValue []*UTXO

func (a utxoByValue) Len() int {
	return len(a)
}
func (a utxoByValue) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a utxoByValue) Less(i, j int) bool {

	vali, _ := a[1].Value()
	valj, _ := a[1].Value()

	return vali < valj
}

// TxAttr .
type TxAttr struct {
	Usage byte
	Data  []byte
}

// Bytes get attr bytes
func (attr *TxAttr) Bytes() []byte {
	var buff bytes.Buffer

	buff.WriteByte(attr.Usage)
	buff.Write(attr.Data)

	return buff.Bytes()
}

// BlockFee block
type BlockFee struct {
	ID         int64
	SysFee     float64
	NetFee     float64
	CreateTime string
}

// Unclaimed .
type Unclaimed struct {
	Unavailable string
	Available   string
	Claims      []*UTXO
}
