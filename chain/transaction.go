package chain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/boryoku-tekina/makiko/utils"
)

var minerReward int = 100

// Transaction : represent a transaction
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// String function
func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("──Transaction %x:", tx.ID))

	for i, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("\t Input:\t %d", i))
		lines = append(lines, fmt.Sprintf("\t TXID:\t %x", input.ID))
		lines = append(lines, fmt.Sprintf("\t Out:\t %d", input.Out))
		lines = append(lines, fmt.Sprintf("\t Signature:\t %x", input.Signature))
		lines = append(lines, fmt.Sprintf("\t PubKey:\t %x", input.PubKey))
	}

	for i, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("\t Output:\t %d", i))
		lines = append(lines, fmt.Sprintf("\t Value:\t %d", output.Value))
		lines = append(lines, fmt.Sprintf("\t Script:\t %x", output.PubKeyHash))

	}

	return strings.Join(lines, "\n")

}

// SetID : setting id
// the tx is id is the hash of the encoded tx
// (encoded tx = byte representation of tx)
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	utils.HandleErr(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// Serialize : return a []byte representation of a transaction
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	utils.HandleErr(err)

	return encoded.Bytes()
}

// Hash : Hash the transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// IsCoinBase FUnction
// Determing what type of tx we have
func (tx *Transaction) IsCoinBase() bool {
	// coin base only have ONE input
	// return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
	return len(tx.Inputs) == 0
}

// CoinBaseTx : transaction of a coin
func CoinBaseTx(amount int, to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("New generated Coin to %s", to)
	}
	// txin := TxInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTxOutput(amount, to)

	// tx := Transaction{nil, []TxInput{txin}, []TxOutput{*txout}}
	tx := Transaction{nil, nil, []TxOutput{*txout}}

	tx.SetID()

	return &tx
}

// TrimmedCopy return a copy of the transaction
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, in := range tx.Inputs {
		inputs = append(inputs, TxInput{in.ID, in.Out, nil, nil})
	}

	for _, out := range tx.Outputs {
		outputs = append(outputs, TxOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Sign : Sign a transaction
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTxs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	for _, in := range tx.Inputs {
		if prevTxs[hex.EncodeToString(in.ID)].ID == nil {
			log.Panic("ERROR : Previous transaction is not correct")
		}
	}
	txCopy := tx.TrimmedCopy()

	for inID, in := range txCopy.Inputs {
		prevTx := prevTxs[hex.EncodeToString(in.ID)]
		txCopy.Inputs[inID].Signature = nil
		txCopy.Inputs[inID].PubKey = prevTx.Outputs[in.Out].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Inputs[inID].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		utils.HandleErr(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[inID].Signature = signature
	}
}
