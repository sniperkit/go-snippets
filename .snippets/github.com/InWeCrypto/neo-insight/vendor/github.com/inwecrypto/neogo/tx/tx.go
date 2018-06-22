package tx

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/inwecrypto/neogo/rpc"
)

// Asserts .
const (
	GasAssert = "0x602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
	NEOAssert = "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
)

// Err
var (
	ErrNoUTXO = errors.New("no enough utxo")
)

// Transaction types
const (
	MinerTransaction      byte = 0x00
	IssueTransaction      byte = 0x01
	ClaimTransaction      byte = 0x02
	EnrollmentTransaction byte = 0x20
	RegisterTransaction   byte = 0x40
	ContractTransaction   byte = 0x80
	PublishTransaction    byte = 0xd0
	InvocationTransaction byte = 0xd1
)

// Attr Usage
const (
	ContractHash   = byte(0x00)
	ECDH02         = byte(0x02)
	ECDH03         = byte(0x03)
	Script         = byte(0x20)
	Vote           = byte(0x30)
	CertURL        = byte(0x80)
	DescriptionURL = byte(0x81)
	Description    = byte(90)
	Hash1          = byte(0xa1)
	Hash2          = byte(0xa2)
	Hash3          = byte(0xa3)
	Hash4          = byte(0xa4)
	Hash5          = byte(0xa5)
	Hash6          = byte(0xa6)
	Hash7          = byte(0xa7)
	Hash8          = byte(0xa8)
	Hash9          = byte(0xa9)
	Hash10         = byte(0xaa)
	Hash11         = byte(0xab)
	Hash12         = byte(0xac)
	Hash13         = byte(0xad)
	Hash14         = byte(0xae)
	Hash15         = byte(0xaf)
	Remark         = byte(0xf0)
	Remark1        = byte(0xf1)
	Remark2        = byte(0xf2)
	Remark3        = byte(0xf3)
	Remark4        = byte(0xf4)
	Remark5        = byte(0xf5)
	Remark6        = byte(0xf6)
	Remark7        = byte(0xf7)
	Remark8        = byte(0xf8)
	Remark9        = byte(0xf9)
	Remark10       = byte(0xfa)
	Remark11       = byte(0xfb)
	Remark12       = byte(0xfc)
	Remark13       = byte(0xfd)
	Remark14       = byte(0xfe)
	Remark15       = byte(0xff)
)

// Serializable .
type Serializable interface {
	Read(reader io.Reader) error
	Write(writer io.Writer) error
}

// ToJSON .
type ToJSON interface {
	JSON() string
}

// Transaction neo transaction object
type Transaction struct {
	Type       byte         // transaction type
	Version    byte         // tx version
	Extend     Serializable // special transaction data type
	Attributes []*Attribute // transaction extra attributes
	Inputs     []*Vin       // transaction input utxos
	Outputs    []*Vout      // transaction change utxos
	Scripts    []*Scripts   // tx scripts
	SignData   []byte       // sign source data
	SignResult []byte       // sign result
	TxID       string       // transaction id
	RawData    []byte       // raw transaction
}

func (tx *Transaction) String() string {
	var buff bytes.Buffer

	buff.WriteString("{")

	switch tx.Type {
	case MinerTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "MinerTransaction"))
		}
	case IssueTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "IssueTransaction"))
		}
	case ClaimTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "ClaimTransaction"))
		}
	case EnrollmentTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "EnrollmentTransaction"))
		}
	case RegisterTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "RegisterTransaction"))
		}
	case ContractTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "ContractTransaction"))
		}
	case PublishTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "PublishTransaction"))
		}
	case InvocationTransaction:
		{
			buff.WriteString(fmt.Sprintf("\"type\":\"%s\"", "InvocationTransaction"))
		}
	}

	buff.WriteString(fmt.Sprintf(",\"version\": %d", tx.Version))

	if tx.Extend != nil {

		if tojson, ok := tx.Extend.(ToJSON); ok {
			buff.WriteString(fmt.Sprintf(",\"extend\": %s", tojson.JSON()))
		} else {
			jsondata, _ := json.Marshal(tx.Extend)

			buff.WriteString(fmt.Sprintf(",\"extend\": %s", jsondata))
		}
	}

	if tx.Attributes != nil {
		var attrjson []string

		for _, attr := range tx.Attributes {
			attrjson = append(attrjson, attr.JSON())
		}

		buff.WriteString(fmt.Sprintf(",\"attributes\": [%s]", strings.Join(attrjson, ",")))
	}

	if tx.Inputs != nil {
		jsondata, _ := json.Marshal(tx.Inputs)

		buff.WriteString(fmt.Sprintf(",\"inputs\": %s", jsondata))
	}

	if tx.Outputs != nil {
		jsondata, _ := json.Marshal(tx.Outputs)

		buff.WriteString(fmt.Sprintf(",\"outputs\": %s", jsondata))
	}

	if tx.Scripts != nil {
		var attrjson []string

		for _, attr := range tx.Scripts {
			attrjson = append(attrjson, attr.JSON())
		}

		buff.WriteString(fmt.Sprintf(",\"scripts\": [%s]", strings.Join(attrjson, ",")))
	}

	if tx.SignData != nil {
		jsondata := hex.EncodeToString(tx.SignData)

		buff.WriteString(fmt.Sprintf(",\"signdata\": \"%s\"", jsondata))
	}

	if tx.SignResult != nil {
		jsondata := hex.EncodeToString(tx.SignResult)

		buff.WriteString(fmt.Sprintf(",\"sign\": \"%s\"", jsondata))
	}

	if tx.RawData != nil {
		jsondata := hex.EncodeToString(tx.RawData)

		buff.WriteString(fmt.Sprintf(",\"raw\": \"%s\"", jsondata))
	}

	buff.WriteString(fmt.Sprintf(",\"txid\": \"%s\"", tx.TxID))

	buff.WriteString("}")

	var any interface{}

	json.Unmarshal(buff.Bytes(), &any)

	data, _ := json.MarshalIndent(any, "", "\t")

	return string(data)
}

type utxoSorter []*rpc.UTXO

func (s utxoSorter) Len() int      { return len(s) }
func (s utxoSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s utxoSorter) Less(i, j int) bool {

	ival, _ := s[i].Value()
	jval, _ := s[j].Value()

	return ival < jval
}

func calcTxInput(amount float64, asset string, unspent []*rpc.UTXO) ([]*rpc.UTXO, float64, error) {
	sort.Sort(utxoSorter(unspent))

	selected := make([]*rpc.UTXO, 0)
	vinvalue := float64(0)

	for _, utxo := range unspent {

		if utxo.Vout.Asset != asset {
			continue
		}

		var err error
		selected = append(selected, utxo)

		val, err := utxo.Value()

		if err != nil {
			return nil, 0, err
		}

		vinvalue += val

		if vinvalue >= amount {
			return selected, vinvalue, nil
		}
	}

	return selected, vinvalue, nil
}

func filter(unspent []*rpc.UTXO, spent []*rpc.UTXO) []*rpc.UTXO {
	result := make([]*rpc.UTXO, 0)

	for _, utxo := range unspent {

		for _, target := range spent {
			if target == utxo {
				goto Skip
			}
		}

		result = append(result, utxo)
	Skip:
	}

	return result
}

// CalcInputs calculate tx Inputs
func (tx *Transaction) CalcInputs(outputs []*Vout, unspent []*rpc.UTXO) ([]*Vin, []*rpc.UTXO, error) {

	inputs := make([]*Vin, 0)

	tx.Outputs = append(tx.Outputs, outputs...)

	for _, vout := range outputs {
		amount := vout.Value.Float64()

		selected, selectedAmount, err := calcTxInput(amount, vout.Asset, unspent)

		if err != nil {
			return nil, nil, err
		}

		if selectedAmount < amount {
			return nil, nil, ErrNoUTXO
		}

		for _, utxo := range selected {
			inputs = append(inputs, &Vin{
				Tx: utxo.TransactionID,
				N:  uint16(utxo.Vout.N),
			})
		}

		if selectedAmount > amount {
			tx.Outputs = append(tx.Outputs, &Vout{
				Asset:   vout.Asset,
				Value:   MakeFixed8(selectedAmount - amount),
				Address: selected[0].Vout.Address,
			})
		}

		unspent = filter(unspent, selected)
	}

	return inputs, unspent, nil
}

func (tx *Transaction) Write(writer io.Writer) error {

	if err := tx.writeSignData(writer); err != nil {
		return err
	}

	length := Varint(len(tx.Scripts))

	if err := length.Write(writer); err != nil {
		return err
	}

	for _, script := range tx.Scripts {
		if err := script.Write(writer); err != nil {
			return err
		}
	}

	return nil
}

func (tx *Transaction) Read(reader io.Reader) error {

	header := make([]byte, 2)

	_, err := reader.Read(header)

	if err != nil {
		return err
	}

	tx.Type = header[0]
	tx.Version = header[1]

	if tx.Extend != nil {
		if err := tx.Extend.Read(reader); err != nil {
			return err
		}
	}

	var length Varint

	if err := length.Read(reader); err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		attr := &Attribute{}

		if err := attr.Read(reader); err != nil {
			return err
		}

		tx.Attributes = append(tx.Attributes, attr)
	}

	if err := length.Read(reader); err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		vin := &Vin{}

		if err := vin.Read(reader); err != nil {
			return err
		}

		tx.Inputs = append(tx.Inputs, vin)
	}

	if err := length.Read(reader); err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		vout := &Vout{}

		if err := vout.Read(reader); err != nil {
			return err
		}

		tx.Outputs = append(tx.Outputs, vout)
	}

	if err := length.Read(reader); err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		scripts := &Scripts{}

		if err := scripts.Read(reader); err != nil {
			return err
		}

		tx.Scripts = append(tx.Scripts, scripts)
	}

	return nil
}

// Attribute .
type Attribute struct {
	Usage byte
	Data  []byte
}

// JSON .
func (attr *Attribute) JSON() string {
	return fmt.Sprintf("{ \"usage\":%d,\"data\": \"%s\" }", attr.Usage, hex.EncodeToString(attr.Data))
}

func (attr *Attribute) Read(reader io.Reader) error {

	usage := make([]byte, 1)

	_, err := reader.Read(usage)

	if err != nil {
		return err
	}

	attr.Usage = usage[0]

	var body []byte

	if attr.Usage <= ECDH03 || attr.Usage == Vote || (attr.Usage <= Hash15 && attr.Usage >= Hash1) {
		body = make([]byte, 32)
	} else {
		var length byte
		_, err = reader.Read([]byte{length})

		if err != nil {
			return err
		}

		body = make([]byte, length)
	}

	_, err = reader.Read(body)

	if err != nil {
		return err
	}

	attr.Data = body

	return nil
}

// Write .
func (attr *Attribute) Write(writer io.Writer) error {

	_, err := writer.Write([]byte{attr.Usage})

	if err != nil {
		return err
	}

	if !(attr.Usage <= ECDH03 || attr.Usage == Vote || (attr.Usage <= Hash15 && attr.Usage >= Hash1)) {
		_, err := writer.Write([]byte{byte(len(attr.Data))})

		if err != nil {
			return err
		}
	}

	_, err = writer.Write(attr.Data)

	if err != nil {
		return err
	}

	return nil
}

// Vin .
type Vin struct {
	Tx string `json:"tx"`
	N  uint16 `json:"n"`
}

func (vin *Vin) Read(reader io.Reader) error {
	txid := make([]byte, 32)

	_, err := reader.Read(txid)

	if err != nil {
		return err
	}

	txid = reverseBytes(txid)

	vin.Tx = fmt.Sprintf("0x%s", hex.EncodeToString(txid))

	data := make([]byte, 2)

	_, err = reader.Read(data)

	if err != nil {
		return err
	}

	vin.N = binary.LittleEndian.Uint16(data)

	return nil
}

func (vin *Vin) Write(writer io.Writer) error {

	data, err := hex.DecodeString(strings.TrimPrefix(vin.Tx, "0x"))

	if err != nil {
		return err
	}

	_, err = writer.Write(reverseBytes(data))

	if err != nil {
		return err
	}

	data = make([]byte, 2)

	binary.LittleEndian.PutUint16(data, vin.N)

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// Vout .
type Vout struct {
	Asset   string `json:"asset"`
	Value   Fixed8 `json:"value"`
	Address string `json:"address"`
}

func (vout *Vout) Read(reader io.Reader) error {
	assertID := make([]byte, 32)

	_, err := reader.Read(assertID)

	if err != nil {
		return err
	}

	assertID = reverseBytes(assertID)

	vout.Asset = fmt.Sprintf("0x%s", hex.EncodeToString(assertID))

	data := make([]byte, 8)

	_, err = reader.Read(data)

	if err != nil {
		return err
	}

	vout.Value = Fixed8(binary.LittleEndian.Uint64(data))

	data = make([]byte, 20)

	_, err = reader.Read(data)

	if err != nil {
		return err
	}

	vout.Address = encodeAddress(data)

	return nil
}

func (vout *Vout) Write(writer io.Writer) error {

	data, err := hex.DecodeString(strings.TrimPrefix(vout.Asset, "0x"))

	if err != nil {
		println("write err :", vout.Asset, vout.Value)
		return err
	}

	_, err = writer.Write(reverseBytes(data))

	if err != nil {
		return err
	}

	value := uint64(vout.Value)

	data = make([]byte, 8)

	binary.LittleEndian.PutUint64(data, value)

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	data, err = decodeAddress(vout.Address)

	if err != nil {
		return err
	}

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// Scripts .
type Scripts struct {
	StackScript  []byte
	RedeemScript []byte
}

// JSON .
func (scripts *Scripts) JSON() string {
	return fmt.Sprintf(
		`{"stackscript":"%s","redeemscript":"%s"}`,
		hex.EncodeToString(scripts.StackScript),
		hex.EncodeToString(scripts.RedeemScript))
}

func (scripts *Scripts) Read(reader io.Reader) error {

	var length Varint

	err := length.Read(reader)

	if err != nil {
		return err
	}

	buff := make([]byte, int(length))

	_, err = reader.Read(buff)

	if err != nil {
		return err
	}

	scripts.StackScript = buff

	err = length.Read(reader)

	if err != nil {
		return err
	}

	buff = make([]byte, int(length))

	_, err = reader.Read(buff)

	scripts.RedeemScript = buff

	return err
}

// WriteBytes .
func (scripts *Scripts) Write(writer io.Writer) error {

	length := Varint(len(scripts.StackScript))

	if err := length.Write(writer); err != nil {
		return err
	}

	_, err := writer.Write(scripts.StackScript)

	if err != nil {
		return err
	}

	length = Varint(len(scripts.RedeemScript))

	if err := length.Write(writer); err != nil {
		return err
	}

	_, err = writer.Write(scripts.RedeemScript)

	return err
}

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func decodeAddress(address string) ([]byte, error) {

	result, _, err := base58.CheckDecode(address)

	if err != nil {
		return nil, err
	}

	return result[0:20], nil
}

// DecodeAddress .
func DecodeAddress(address string) ([]byte, error) {
	return decodeAddress(address)
}

func encodeAddress(address []byte) string {
	return base58.CheckEncode(address, 0x17)
}

// EncodeAddress .
func EncodeAddress(address []byte) string {
	return encodeAddress(address)
}
