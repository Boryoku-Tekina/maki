package tests

import (
	"fmt"

	"github.com/boryoku-tekina/makiko/chain"
	"github.com/boryoku-tekina/makiko/utils"
	"github.com/boryoku-tekina/makiko/wallet"
)

// CreateWallet test the creation of wallet
func CreateWallet() {
	w, _ := wallet.CreateWallets()
	w.LoadFile()
	fmt.Println("creating wallet 3 times")
	w.AddWallet()
	w.AddWallet()
	w.AddWallet()
	fmt.Println("saving wallet file")
	w.SaveFile()
	all := w.GetAllAddresses()

	for _, addr := range all {
		fmt.Println(addr)
	}
}

// Transaction : test transactions functions
func Transaction() {
	chain.InitChain()
	cbtx := chain.CoinBaseTx(100, "1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", "")
	chain.AddBlock([]*chain.Transaction{cbtx})
	tx := chain.NewTransaction("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", "1C2xh3EhizWUXMagg83UpAcU4fk7UnUdYc", 50)
	chain.AddBlock([]*chain.Transaction{tx})
}

// GetFunds : get fund in address
func GetFunds(address string) {
	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	// UTXOs := chain.FindUTXO(pubKeyHash)

	// for _, out := range UTXOs {
	// 	balance += out.Value
	// }
	UTXOS := chain.FindUnspentTransactions(pubKeyHash)

	for _, tx := range UTXOS {
		balance += tx.Outputs[len(tx.Outputs)-1].Value
	}
	fmt.Println("balance of ", address, " = ", balance)
	// chain.PrintChain()
}
