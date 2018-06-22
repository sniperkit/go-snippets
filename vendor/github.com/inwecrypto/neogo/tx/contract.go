package tx

import "github.com/inwecrypto/neogo/rpc"

// ContractTx contract transaction
type ContractTx Transaction

// NewContractTx create new contract transaction
func NewContractTx() *ContractTx {
	tx := &ContractTx{}

	tx.Type = ContractTransaction

	return tx
}

// Tx .
func (tx *ContractTx) Tx() *Transaction {
	return (*Transaction)(tx)
}

// CalcInputs .
func (tx *ContractTx) CalcInputs(outputs []*Vout, unspent []*rpc.UTXO) error {
	base := (*Transaction)(tx)

	vin, _, err := base.CalcInputs(outputs, unspent)

	if err != nil {
		return err
	}

	tx.Inputs = vin

	return nil
}
