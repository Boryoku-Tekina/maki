package chain

import (
	"bytes"
	"fmt"

	"github.com/boryoku-tekina/makiko/utils"
)

var (
	utxoPrefix   = []byte("utxo-")
	prefixLenght = len(utxoPrefix)
)

// GetUTXOOf function
// return all UTXO for the address
func GetUTXOOf(address string) TxOutputs {
	UTXOs := TxOutputs{}
	UTXOs.Outputs = nil

	// wallets, err := wallet.CreateWallets()
	// utils.HandleErr(err)
	// w := wallets.GetWallets(address)

	PubKeyHash := utils.Base58Decode([]byte(address))
	PubKeyHash = PubKeyHash[1 : len(PubKeyHash)-4]

	lh := GetLastBlockHash()
	actualBlock := GetBlockByHash(lh)
Parcour:
	for {
		// if we are on the genesis block
		if bytes.Equal(actualBlock.PrevHash, bytes.Repeat([]byte{0}, 32)) {
			break
		}
		if actualBlock.Transactions == nil {
			fmt.Println("there is no tx in this block")
		}

		for _, actualTx := range actualBlock.Transactions {
			if actualTx.IsCoinBase() {
				for _, out := range actualTx.Outputs {
					if out.IsLockedWithKey(PubKeyHash) {
						fmt.Println("get coin base for ", address, " appending it...")
						UTXOs.Outputs = append(UTXOs.Outputs, out)
						fmt.Println("actual UTXOs.OUtputs : ", UTXOs.Outputs)
					}
				}
			} else {
				for _, out := range actualTx.Outputs {
					if out.IsLockedWithKey(PubKeyHash) {
						isChange := false
						// for _, in := range actualTx.Inputs {
						// 	if bytes.Equal(in.PubKey, w.PublicKey) {
						// 		isChange = true
						// 	}
						// }
						// if the last outputs is locked with key pubkeyhash
						// it is a change
						if actualTx.Outputs[(len(actualTx.Outputs) - 1)].IsLockedWithKey(PubKeyHash) {
							isChange = true
						}
						if isChange == true {
							fmt.Println("get a CHANGE for ", address, "stoping...")
							UTXOs.Outputs = append(UTXOs.Outputs, out)
							fmt.Println("actual UTXOs.OUtputs : ", UTXOs.Outputs)
							break Parcour
						} else {
							fmt.Println("get a TRANSACTION for ", address, "appending and continue...")
							UTXOs.Outputs = append(UTXOs.Outputs, out)
							fmt.Println("actual UTXOs.OUtputs : ", UTXOs.Outputs)
						}
					}
				}
			}

		}

		actualBlock = GetBlockByHash(actualBlock.PrevHash)
	}
	return UTXOs
}

// if address contains in Inputs and OUtputs then reset the UTXO to the actual TxO
// if address contains only in Outpus then append actual UTXO with the actual TxO
