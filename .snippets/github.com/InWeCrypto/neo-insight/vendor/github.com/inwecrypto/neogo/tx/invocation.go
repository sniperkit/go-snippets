package tx

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/inwecrypto/neogo/rpc"
)

// InvocationTx .
type InvocationTx Transaction

type invocationTx struct {
	Script []byte `json:"script"`
	Gas    Fixed8 `json:"gas"`
}

// NewInvocationTx .
func NewInvocationTx(script []byte, gas float64) *InvocationTx {
	return &InvocationTx{
		Type:    InvocationTransaction,
		Version: 1,
		Extend: &invocationTx{
			Script: script,
			Gas:    MakeFixed8(gas),
		},
	}
}

// JSON .
func (tx *invocationTx) JSON() string {
	return fmt.Sprintf(`{ "script":"%s","gas":%d }`, hex.EncodeToString(tx.Script), tx.Gas)
}

// Tx .
func (tx *InvocationTx) Tx() *Transaction {
	return (*Transaction)(tx)
}

// CheckFromWitness .
func (tx *InvocationTx) CheckFromWitness(fromScriptHash []byte) {
	tx.Attributes = append(tx.Attributes, &Attribute{
		Usage: Script,
		Data:  fromScriptHash,
	})
}

// CalcInputs .
func (tx *InvocationTx) CalcInputs(outputs []*Vout, unspent []*rpc.UTXO) error {
	invocation := tx.Extend.(*invocationTx)

	base := (*Transaction)(tx)

	inputs, unselected, err := base.CalcInputs(outputs, unspent)

	if err != nil {
		return err
	}

	tx.Inputs = inputs

	unspent = unselected

	if invocation.Gas.Float64() < 1 {
		invocation.Gas = Fixed8(0) // zero gas
		for _, utxo := range unspent {
			if utxo.Vout.Asset == GasAssert {
				val, err := utxo.Value()
				if err != nil {
					return err
				}

				tx.Outputs = append(tx.Outputs, &Vout{
					Asset:   GasAssert,
					Value:   MakeFixed8(val),
					Address: utxo.Vout.Address,
				})

				tx.Inputs = append(tx.Inputs, &Vin{
					Tx: utxo.TransactionID,
					N:  uint16(utxo.Vout.N),
				})

				return nil
			}
		}
	}

	amount := invocation.Gas.Float64()

	selected, selectedAmount, err := calcTxInput(amount, GasAssert, unspent)

	if err != nil {
		return err
	}

	if selectedAmount < amount {
		return ErrNoUTXO
	}

	for _, utxo := range selected {
		tx.Inputs = append(tx.Inputs, &Vin{
			Tx: utxo.TransactionID,
			N:  uint16(utxo.Vout.N),
		})
	}

	if selectedAmount > amount {
		tx.Outputs = append(tx.Outputs, &Vout{
			Asset:   GasAssert,
			Value:   MakeFixed8(selectedAmount - amount),
			Address: selected[0].Vout.Address,
		})
	}

	return nil
}

func (tx *invocationTx) Write(writer io.Writer) error {

	length := Varint(len(tx.Script))

	if err := length.Write(writer); err != nil {
		return err
	}

	_, err := writer.Write(tx.Script)

	if err != nil {
		return err
	}

	return tx.Gas.Write(writer)
}

func (tx *invocationTx) Read(reader io.Reader) error {

	var length Varint

	if err := length.Read(reader); err != nil {
		return err
	}

	buff := make([]byte, int(length))

	_, err := reader.Read(buff)

	if err != nil {
		return err
	}

	tx.Script = buff

	return tx.Gas.Read(reader)
}

// ToInvocationAddress neo wallet address to invocation address
func ToInvocationAddress(address string) string {
	bytesOfAddress, _ := decodeAddress(address)

	return hex.EncodeToString(reverseBytes(bytesOfAddress))
}
