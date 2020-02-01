package chain

import (
	"bytes"
	"fmt"
	"os"
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
