package chain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boryoku-tekina/makiko/utils"
	"github.com/boryoku-tekina/makiko/wallet"
)

// InitChain initialize the blockchain
// create the database and the first transaction ── Genesis
// store the last hash key/value to the Genesis db
// One block ──> One Database = one file ; so we can download one DB file instead of
// downloading all the database from the beginning and file by file
func InitChain() {
	if !DBExists() {
		fmt.Println("[INFO] : no chain yet, creating genesis block")
		CreateGenesisBlock()
	}
	fmt.Println("[INFO] : it means that there is already a chain database")
}

// AddBlock : Add New Block to the chain
// mine the block with the pending txs
func AddBlock(transactions []*Transaction) {
	var b Block
	b.Transactions = transactions
	b.Mine()
}

// PrintChain : print the chain
func PrintChain() {
	// iterating through all block
	// beginning from the last
	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)
	for {
		// if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			actualBlock.PrintBlockInfo()
			break
		}
		actualBlock.PrintBlockInfo()
		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}
}

// ValidChain return true if chain is valid
// if all block is connected
func ValidChain() bool {
	// var actualBlock Block
	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)

	for {
		// if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			validation := actualBlock.ValidateBlock()
			if validation == false {
				return false
			}
			break
		}
		validation := actualBlock.ValidateBlock()
		if validation == false {
			return false
		}
		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}
	return true
}

// DBExists function to check if database file exist
func DBExists() bool {
	if _, err := os.Stat("./DB/LastBlockHash.bc"); os.IsNotExist(err) {
		return false
	}

	return true
}

// SignTransaction function to sign a tx
func SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTxs := make(map[string]Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := FindTransaction(in.ID)
		utils.HandleErr(err)
		prevTxs[hex.EncodeToString(prevTX.ID)] = prevTX
	}
	tx.Sign(privKey, prevTxs)
}

// FindTransaction : find transaction by given id
func FindTransaction(ID []byte) (Transaction, error) {
	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)

	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		for _, tx := range actualBlock.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}
		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}

	return Transaction{}, errors.New("Transaction does not exist")
}

// NewTransaction : create new transaction from an address to another adress
func NewTransaction(from, to string, amount int) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	wallets, err := wallet.CreateWallets()
	utils.HandleErr(err)
	w := wallets.GetWallets(from)
	pubKeyHash := wallet.PublicKeyHash(w.PublicKey)

	acc, validOutputs := FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("Error: not enough funds")
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.HandleErr(err)

		for _, out := range outs {
			input := TxInput{txID, out, nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, *NewTxOutput(amount, to))

	if acc > amount {
		// changes
		outputs = append(outputs, *NewTxOutput(acc-amount, from))
	}

	tx := Transaction{nil, inputs, outputs}
	tx.ID = tx.Hash()
	SignTransaction(&tx, w.PrivateKey)

	return &tx

}

// FindUTXO find unspent transaction outputs
func FindUTXO() map[string]TxOutputs {
	UTXO := make(map[string]TxOutputs)
	spentTXOs := make(map[string][]int)

	var lh []byte
	GetLastBlockHash(&lh)
	block := GetBlockByHash(lh)

	for {
		// break if we are on the genesis block
		if bytes.Equal(block.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				// take the outputs
				// and put it in the UTXO[] map
				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}
			if tx.IsCoinBase() == false {
				for _, in := range tx.Inputs {
					inTxID := hex.EncodeToString(in.ID)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
				}
			}
		}
		block = GetBlockByHash(block.PrevHash)

	}
	return UTXO
}
