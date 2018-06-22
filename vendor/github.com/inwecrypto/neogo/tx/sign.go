package tx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/apisit/rfc6979"
	"github.com/inwecrypto/neogo/script"
)

func publicKeyToBytes(pub *ecdsa.PublicKey) (b []byte) {
	/* See Certicom SEC1 2.3.3, pg. 10 */

	x := pub.X.Bytes()

	/* Pad X to 32-bytes */
	paddedx := append(bytes.Repeat([]byte{0x00}, 32-len(x)), x...)

	/* Add prefix 0x02 or 0x03 depending on ylsb */
	if pub.Y.Bit(0) == 0 {
		return append([]byte{0x02}, paddedx...)
	}

	return append([]byte{0x03}, paddedx...)
}

func rfc6979Sign(ecdsaPrivateKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {

	digest := sha256.Sum256(data)

	r, s, err := rfc6979.SignECDSA(ecdsaPrivateKey, digest[:], sha256.New)

	if err != nil {
		return nil, err
	}

	params := ecdsaPrivateKey.Curve.Params()
	curveOrderByteSize := params.P.BitLen() / 8
	rBytes, sBytes := r.Bytes(), s.Bytes()
	signature := make([]byte, curveOrderByteSize*2)
	copy(signature[curveOrderByteSize-len(rBytes):], rBytes)
	copy(signature[curveOrderByteSize*2-len(sBytes):], sBytes)

	return signature, nil
}

// Sign sign transaction
func (tx *Transaction) Sign(ecdsaPrivateKey *ecdsa.PrivateKey) ([]byte, string, error) {
	var buff bytes.Buffer

	if err := tx.writeSignData(&buff); err != nil {
		return nil, "", err
	}

	tx.SignData = make([]byte, len(buff.Bytes()))
	copy(tx.SignData, buff.Bytes())

	txid := sha256.Sum256(tx.SignData)
	txid = sha256.Sum256(txid[:])

	tx.TxID = hex.EncodeToString(reverseBytes(txid[:]))

	sign, err := rfc6979Sign(ecdsaPrivateKey, tx.SignData)

	if err != nil {
		return nil, "", err
	}

	tx.SignResult = sign

	var stackScriptBuffer bytes.Buffer

	signScript := script.New("test")

	signScript.EmitPushBytes(sign)

	signScript.Write(&stackScriptBuffer)

	stackScript := stackScriptBuffer.Bytes()

	address := publicKeyToBytes(&ecdsaPrivateKey.PublicKey)

	signScript.Reset()

	signScript.
		EmitPushBytes(address).
		Emit(script.CHECKSIG, nil)

	var redeemScriptBuffer bytes.Buffer

	signScript.Write(&redeemScriptBuffer)

	tx.Scripts = []*Scripts{
		&Scripts{
			StackScript:  stackScript,
			RedeemScript: redeemScriptBuffer.Bytes(),
		},
	}

	var rawTx bytes.Buffer

	if err := tx.Write(&rawTx); err != nil {
		return nil, "", err
	}

	return rawTx.Bytes(), tx.TxID, nil
}

func (tx *Transaction) writeSignData(writer io.Writer) error {
	_, err := writer.Write([]byte{tx.Type, tx.Version})

	if err != nil {
		return err
	}

	if tx.Extend != nil {
		if err := tx.Extend.Write(writer); err != nil {
			return err
		}
	}

	length := Varint(len(tx.Attributes))

	if err := length.Write(writer); err != nil {
		return err
	}

	for _, attr := range tx.Attributes {
		if err := attr.Write(writer); err != nil {
			return err
		}
	}

	length = Varint(len(tx.Inputs))

	if err := length.Write(writer); err != nil {
		return err
	}

	for _, input := range tx.Inputs {
		if err := input.Write(writer); err != nil {
			return err
		}
	}

	length = Varint(len(tx.Outputs))

	if err := length.Write(writer); err != nil {
		return err
	}

	for _, output := range tx.Outputs {
		if err := output.Write(writer); err != nil {
			return err
		}
	}

	return nil
}
