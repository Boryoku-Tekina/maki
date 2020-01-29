package chain

import (
	"github.com/boltdb/bolt"
)

// Chain represent a blockchain struct
type Chain struct {
	LastBlockHash []byte
	Database      bolt.DB
}

const (
	genesisPath = "tmp/Genesis.db"
)

// InitChain initialize the blockchain
// create the database and the first transaction ── Genesis
// store the last hash key/value to the Genesis db
// One block ──> One Database = one file ; so we can download one DB file instead of
// downloading all the database from the beginning
func InitChain() *Chain {

	return nil
}
