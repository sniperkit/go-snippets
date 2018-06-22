package tx

import (
	"encoding/json"
	"io"

	"github.com/inwecrypto/neogo/rpc"
)

// ClaimTx .
type ClaimTx Transaction

type claimTx struct {
	Inputs []*Vin `json:"inputs"`
}

// NewClaimTx .
func NewClaimTx() *ClaimTx {
	tx := &ClaimTx{
		Type: ClaimTransaction,
	}

	return tx
}

// JSON .
func (tx *ClaimTx) JSON() string {
	data, _ := json.Marshal(tx.Inputs)

	return string(data)
}

// Tx get basic transaction object
func (tx *ClaimTx) Tx() *Transaction {
	return (*Transaction)(tx)
}

// Claim .
func (tx *ClaimTx) Claim(amount float64, to string, claims []*rpc.UTXO) error {

	var inputs []*Vin

	for _, utxo := range claims {
		if utxo.Vout.Asset != NEOAssert {
			continue
		}

		inputs = append(inputs, &Vin{
			Tx: utxo.TransactionID,
			N:  uint16(utxo.Vout.N),
		})
	}

	if len(inputs) == 0 {
		return ErrNoUTXO
	}

	tx.Extend = &claimTx{
		Inputs: inputs,
	}

	tx.Outputs = []*Vout{
		&Vout{
			Asset:   GasAssert,
			Value:   MakeFixed8(amount),
			Address: to,
		},
	}

	return nil
}

func (tx *claimTx) Write(writer io.Writer) error {

	length := Varint(len(tx.Inputs))

	if err := length.Write(writer); err != nil {
		return err
	}

	for _, vin := range tx.Inputs {
		if err := (*Vin)(vin).Write(writer); err != nil {
			return err
		}
	}

	return nil
}

func (tx *claimTx) Read(reader io.Reader) error {
	var length Varint

	if err := length.Read(reader); err != nil {
		return err
	}

	for i := 0; i < int(length); i++ {
		var vin Vin

		if err := vin.Read(reader); err != nil {
			return err
		}

		tx.Inputs = append(tx.Inputs, &vin)
	}

	return nil
}
