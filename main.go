/*
*
* copyright - Bōryoku-tekina soro - 2020
* makiko -  malagasy cryptocurrency
 */

package main

import (
	"github.com/boryoku-tekina/makiko/chain"
)

func main() {

	// step 1
	chain.InitChain()
	var block1 chain.Block
	block1.Data = []byte("Block 1 Data")
	block1.Mine()
	var block2 chain.Block
	block2.Data = []byte("Block 2 Data")
	block2.Mine()

	// step 2
	var lastHash []byte
	chain.GetLastBlockHash(&lastHash)

	block := chain.GetBlockByHash(lastHash)

	block.PrintBlockInfo()

}
