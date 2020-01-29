/*
*
* copyright - Bōryoku-tekina soro - 2020
* makiko -  malagasy cryptocurrency
 */

package main

import (
	"fmt"

	"github.com/Boryoku-tekina-soro/makiko/chain"
)

func main() {
	var block1 chain.Block

	block1.Data = []byte("First blodddck Data")

	var work chain.ProofOfWork

	work.Block = &block1

	// block1.Nonce, block1.Hash = work.Work()
	// work.Validate()

	// block1.RegisterToDB()
	// block1.SetAsLastBlock()

	var lastHash []byte

	chain.GetLastBlockHash(&lastHash)

	fmt.Printf("last hash : %x\n", lastHash)

}
