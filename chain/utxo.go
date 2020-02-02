package chain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/boryoku-tekina/makiko/utils"

	"github.com/boltdb/bolt"
)

var (
	utxoPrefix   = []byte("utxo-")
	prefixLenght = len(utxoPrefix)
)

// UTXOSet struct of
type UTXOSet struct {
}

// DeleteByPrefix : delete a keys/value by the given prefix
func DeleteByPrefix(prefix []byte) {

	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)

	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		db := OpenDatabase(fmt.Sprintf("%x", actualBlock.Hash))
		defer db.Close()

		if err := db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("block"))
			if err := bucket.ForEach(func(k, v []byte) error {
				if bytes.HasPrefix(k, prefix) {
					if err := bucket.Delete(k); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				return err
			}
			return nil
		}); err != nil {
			log.Panic("Error updating Db : ", err.Error())
		}

		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}

}

// Reindex function
func Reindex() {
	DeleteByPrefix(utxoPrefix)
	UTXO := FindUTXO()

	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)
	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		db := OpenDatabase(fmt.Sprintf("%x", actualBlock.Hash))
		defer db.Close()
		err := db.Update(func(tx *bolt.Tx) error {
			for txID, outs := range UTXO {
				key, err := hex.DecodeString(txID)
				if err != nil {
					return err
				}
				key = append(utxoPrefix, key...)
				bucket := tx.Bucket([]byte("block"))
				err = bucket.Put(key, outs.Serialize())
				utils.HandleErr(err)
			}
			return nil
		})
		utils.HandleErr(err)

		actualBlock = GetBlockByHash(actualBlock.PrevHash)

	}

}

// Update function
// update block
func Update(block *Block) {
	db := OpenDatabase(fmt.Sprintf("%x", block.Hash))
	defer db.Close()

	err := db.Update(func(txn *bolt.Tx) error {
		for _, tx := range block.Transactions {
			if !tx.IsCoinBase() {
				for _, in := range tx.Inputs {
					updateOuts := TxOutputs{}
					inID := append(utxoPrefix, in.ID...)
					item := txn.Bucket([]byte("block"))
					v := item.Get(inID)

					outs := DeserializeOutputs(v)

					for outIDx, out := range outs.Outputs {
						if outIDx != in.Out {
							updateOuts.Outputs = append(updateOuts.Outputs, out)
						}
					}
					if len(updateOuts.Outputs) == 0 {
						if err := item.Delete(inID); err != nil {
							log.Panic(err.Error())
						}
					} else {
						if err := item.Put(inID, updateOuts.Serialize()); err != nil {
							log.Panic(err.Error())
						}
					}
				}
			}
			newOutputs := TxOutputs{}
			for _, out := range tx.Outputs {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			txID := append(utxoPrefix, tx.ID...)
			item := txn.Bucket([]byte("block"))
			if err := item.Put(txID, newOutputs.Serialize()); err != nil {
				log.Panic(err)
			}

		}
		return nil
	})

	utils.HandleErr(err)
}

// CountTransactions : count unspent transactions inside of a block
func CountTransactions() int {
	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)
	counter := 0
	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		db := OpenDatabase(fmt.Sprintf("%x", actualBlock.Hash))
		defer db.Close()

		err := db.View(func(txn *bolt.Tx) error {
			it := txn.Bucket([]byte("block"))
			it.ForEach(func(k, v []byte) error {
				if bytes.HasPrefix(k, utxoPrefix) {
					counter++
				}
				return nil
			})
			return nil
		})
		utils.HandleErr(err)

		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}

	return counter
}

// FindUnspentTransactions : find unspent transactions
// return array of transactions that has not been spent yet
func FindUnspentTransactions(pubKeyHash []byte) []TxOutput {
	var UTXOs []TxOutput
	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)
	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}
		db := OpenDatabase(fmt.Sprintf("%x", actualBlock.Hash))
		defer db.Close()

		err := db.View(func(txn *bolt.Tx) error {
			it := txn.Bucket([]byte("block"))
			it.ForEach(func(k, v []byte) error {
				if bytes.HasPrefix(k, utxoPrefix) {
					v := it.Get(k)
					outs := DeserializeOutputs(v)

					for _, out := range outs.Outputs {
						if out.IsLockedWithKey(pubKeyHash) {
							UTXOs = append(UTXOs, out)
						}
					}

				}
				return nil
			})
			return nil
		})
		utils.HandleErr(err)

		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}

	return UTXOs
}

// FindSpendableOutputs : find spendable outputs
func FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0
	var lh []byte
	GetLastBlockHash(&lh)
	actualBlock := GetBlockByHash(lh)
	for {
		// break if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}

		db := OpenDatabase(fmt.Sprintf("%x", actualBlock.Hash))

		db.View(func(txn *bolt.Tx) error {
			it := txn.Bucket([]byte("block"))

			it.ForEach(func(k, v []byte) error {
				if bytes.HasPrefix(k, utxoPrefix) {
					v := it.Get(k)
					k = bytes.TrimPrefix(k, utxoPrefix)
					txID := hex.EncodeToString(k)
					outs := DeserializeOutputs(v)

					for outIdx, out := range outs.Outputs {
						if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
							accumulated += out.Value
							unspentOuts[txID] = append(unspentOuts[txID], outIdx)
						}
					}
				}
				return nil
			})
			return nil
		})

		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}

	return accumulated, unspentOuts
}
