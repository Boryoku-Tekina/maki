package chain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"github.com/Boryoku-tekina-soro/makiko/utils"
	"github.com/boltdb/bolt"
)

const (
	dBPath = "DB/"
)

// Block represent one block structure
type Block struct {
	Hash      []byte
	Data      []byte
	PrevHash  []byte
	Timestamp time.Time
	Nonce     int
}

// OpenDatabase : opening a database
func OpenDatabase(s string) *bolt.DB {
	opts := bolt.DefaultOptions

	Path := dBPath + s + ".bc"

	db, err := bolt.Open(Path, 0777, opts)
	utils.HandleErr(err)
	return db
}

// MineBlock : create a new block and mine it
func MineBlock(data, prevHash []byte) *Block {
	// block := &Block{[]byte{}, txs, prevHash, 0}
	var block *Block
	block.Data = data
	block.PrevHash = prevHash

	pow := NewWork(block)
	nonce, hash := pow.Work()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// ValidateBlock : validate a passed block
func (b *Block) ValidateBlock() bool {
	pow := NewWork(b)
	return pow.Validate()
}

// RegisterToDB open new boltDB database and create the block value one key in it
func (b *Block) RegisterToDB() {

	db := OpenDatabase(fmt.Sprintf("%x", b.Hash))
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(b.Hash)
		utils.HandleErr(err)
		utils.HandleErr(err)
		bucket.Put([]byte("Hash"), b.Hash)
		bucket.Put([]byte("Data"), b.Data)
		bucket.Put([]byte("PrevHash"), b.PrevHash)
		bucket.Put([]byte("TimeStamp"), []byte(b.Timestamp.String()))
		bucket.Put([]byte("nonce"), []byte(strconv.Itoa(b.Nonce)))
		return err
	})
	utils.HandleErr(err)
	fmt.Println("block successfully registered to database")
}

// SetAsLastBlock set a block as the last block in the database
func (b *Block) SetAsLastBlock() {

	db := OpenDatabase("LastBlockHash")
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("LastBlockHash"))
		utils.HandleErr(err)
		bucket.Put([]byte("LastBlockHash"), b.Hash)
		return err
	})
	utils.HandleErr(err)

	fmt.Println("block successfully set as last block")
}

// GetLastBlockHash get the hash of the last block in the chain
// putting it in d
func GetLastBlockHash(d *[]byte) {
	var result []byte

	var b Block

	db := OpenDatabase("LastBlockHash")
	defer db.Close()

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("LastBlockHash"))
		result = bucket.Get([]byte("LastBlockHash"))
		b.Hash = result

		var buffer bytes.Buffer
		encoder := gob.NewEncoder(&buffer)
		errs := encoder.Encode(result)
		utils.HandleErr(errs)
		res := buffer.Bytes()
		// fmt.Printf("Inspecting res : %x\n", res[0:4])
		*d = res[4:]

		return nil
	})
	utils.HandleErr(err)
	fmt.Println("got last block hash")
}
