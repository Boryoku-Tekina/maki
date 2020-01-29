package chain

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"github.com/Boryoku-tekina-soro/makiko/utils"
)

// take the data from the block

// create a counter (nonce) which start at 0

// create a hash of the data plus the counter

// check the hash to see if it meet a set of requirements

// REQUIREMENTS:
// THe first few bytes of the hash must contain 0

// Difficulty of consensus
const Difficulty = 2

// ProofOfWork struct
type ProofOfWork struct {
	Block *Block
}

// NewWork return a new pow
func NewWork(b *Block) *ProofOfWork {
	pow := &ProofOfWork{b}

	return pow
}

// initData Function
func (pow *ProofOfWork) initData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			[]byte(pow.Block.Timestamp.String()),
			utils.ToHex(int64(nonce)),
			utils.ToHex(int64(Difficulty)),
		},
		[]byte{},
	)
	return data
}

// Work the PoW algorithm
func (pow *ProofOfWork) Work() (int, []byte) {
	var hash [32]byte
	nonce := 0
	for {
		data := pow.initData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash)

		if bytes.HasPrefix(hash[:], bytes.Repeat([]byte{0}, Difficulty)) {
			break
		}
		nonce++
	}
	fmt.Println()
	return nonce, hash[:]
}

// Validate the work
func (pow *ProofOfWork) Validate() bool {
	data := pow.initData(pow.Block.Nonce)

	hash := sha256.Sum256(data)

	if !bytes.HasPrefix(hash[:], bytes.Repeat([]byte{0}, Difficulty)) {
		fmt.Println("hash does not satisfy difficulty requirements")
		return false
	}
	if !bytes.Equal(hash[:], pow.Block.Hash) {
		fmt.Println("hash is not correct : maybe this nonce does not provide a valid hash")
		return false
	}
	fmt.Println("good news : work valid")
	return true
}
