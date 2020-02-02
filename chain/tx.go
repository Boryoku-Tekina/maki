package chain

import (
	"bytes"
	"encoding/gob"

	"github.com/boryoku-tekina/makiko/wallet"

	"github.com/boryoku-tekina/makiko/utils"
)

////////////////////////////////////////////////////////////
//──────────────────*TxInput Section*─────────────────────//
////////////////////////////////////////////////////////////

// TxInput : represent the input in transaction
type TxInput struct {
	ID        []byte
	Out       int
	Signature []byte
	PubKey    []byte
}

// UsesKey :
// return true if the public key hash in the wallet correspond to the given PubKeyHash
func (in *TxInput) UsesKey(PubKeyHash []byte) bool {
	lockingHash := wallet.PublicKeyHash(in.PubKey)

	return bytes.Compare(lockingHash, PubKeyHash) == 0
}

/*//////////////////////////////////////////////////////////
//─────────────────*TxOutput Section*─────────────────────//
//////////////////////////////////////////////////////////*/

// TxOutput : represent the output in transaction
type TxOutput struct {
	Value      int
	PubKeyHash []byte
}

// NewTxOutput :
// create New Transaction Output
func NewTxOutput(value int, address string) *TxOutput {
	txo := &TxOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}

// Lock : locking an output to a given address
// setting the public key hash of the output to the given address
func (out *TxOutput) Lock(address []byte) {
	PubKeyHash := utils.Base58Decode(address)
	PubKeyHash = PubKeyHash[1 : len(PubKeyHash)-4]
	out.PubKeyHash = PubKeyHash
}

// IsLockedWithKey function
func (out *TxOutput) IsLockedWithKey(PubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, PubKeyHash) == 0
}

////////////////////////////////////////////////////////////
//─────────────────*TxOutputs Section*────────────────────//
////////////////////////////////////////////////////////////

// TxOutputs : array of TxOutput
// we can use this to identify our transaction outputs
// and then sort it by unspent outpus
type TxOutputs struct {
	Outputs []TxOutput
}

// Serialize TxOutputs
// take the actual struct and encode it in bytes
func (outs TxOutputs) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(outs)
	utils.HandleErr(err)
	return buffer.Bytes()
}

// DeserializeOutputs : Deserialize Outputs
// take and encoded to byte struct
// and decode it to original Txouputs struct
func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs
	decode := gob.NewDecoder(bytes.NewReader(data))
	err := decode.Decode(&outputs)
	utils.HandleErr(err)
	return outputs
}
