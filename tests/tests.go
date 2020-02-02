package tests

import (
	"fmt"

	"github.com/boryoku-tekina/makiko/utils"

	"github.com/boryoku-tekina/makiko/chain"
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

// GetBalanceOf : get fund in address
func GetBalanceOf(address string) int {
	balance := 0
	pubKeyHash := utils.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	balance = chain.GetAmountOf(address)
	fmt.Printf("balance of %s == %d\n", address, balance)
	return balance
}

// Donate : give coin base transaction to an address
func Donate(address string, amount int) {
	T1 := chain.CoinBaseTx(amount, address, "Donation")
	var Txs []*chain.Transaction
	Txs = append(Txs, T1)
	chain.AddBlock(Txs)
	fmt.Printf("\n\n[INFO] : Donation of %d for %s DONE!\n\n", amount, address)
}

// Send : send amount of coin from an address to other
func Send(from, to string, amount int) {
	Tx := chain.NewTransaction(from, to, 70)
	var txs []*chain.Transaction
	txs = append(txs, Tx)
	chain.AddBlock(txs)
	fmt.Printf("[INFO]: Sending %d coins from %s to %s DONE!\n", amount, from, to)
}
