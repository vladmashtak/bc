package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"math"
	"math/big"
	"strconv"
	"time"
)

// MerkleRoot - tmp type for merkle root
// it will be the hash of all transactions in the block
// in the next revision
type MerkleRoot []byte

type Block struct {
	hash       []byte
	merkleRoot MerkleRoot
	prevHash   []byte
	nonce      int64
	difficulty int64
	timestamp  int64
}

func NewBlock(root MerkleRoot, prevHash []byte, difficulty int64) *Block {
	block := &Block{
		hash:       []byte{},
		merkleRoot: root,
		prevHash:   prevHash,
		difficulty: difficulty,
		timestamp:  time.Now().UnixNano(),
	}

	block.runProofOfWork()

	return block
}

func (block *Block) String() string {
	return fmt.Sprintf(
		"Hash: %x\nMerkleRoot: %s\nPrev hash: %x\nNonce: %d\nTimestamp: %d\nPoW: %s\n",
		block.hash,
		string(block.merkleRoot),
		block.prevHash,
		block.nonce,
		block.timestamp,
		strconv.FormatBool(block.Validate()),
	)
}

func (block *Block) Hash() []byte {
	return block.hash
}

func (block *Block) PrevHash() []byte {
	return block.prevHash
}

func (block *Block) Validate() bool {
	var intHash big.Int

	target := big.NewInt(1)
	target.Lsh(target, uint(256-block.difficulty))

	record := block.getBytes(block.nonce)
	hash := sha256.Sum256(record)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(target) == -1
}

// Get array of bytes from the block
// create a counter (nonce) which starts at 0
// create a hash of the data plus the counter
// check the hash to see if it meets to set of requirements
func (block *Block) runProofOfWork() {
	var intHash big.Int
	var hash [32]byte
	var nonce int64

	// target is our threshold that needs to be overcome by doing PoW
	target := big.NewInt(1)
	target.Lsh(target, uint(256-block.difficulty))

	for nonce < math.MaxInt64 {
		record := block.getBytes(nonce)
		hash = sha256.Sum256(record)
		intHash.SetBytes(hash[:])

		cmp := intHash.Cmp(target)
		if cmp == -1 {
			break
		}

		nonce++
	}

	// save founded nonce and hash after PoW
	block.hash = hash[:]
	block.nonce = nonce
}

func (block *Block) getBytes(nonce int64) []byte {
	// concatenate all block fields into one byte array
	record := bytes.Join(
		[][]byte{
			block.prevHash,
			block.merkleRoot,
			block.toHex(nonce),
			block.toHex(block.difficulty),
			block.toHex(block.timestamp),
		},
		[]byte{},
	)

	return record
}

func (block *Block) toHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}

	return buff.Bytes()
}

func (block *Block) MarshalJSON() ([]byte, error) {
	var json = jsoniter.ConfigFastest
	return json.Marshal(struct {
		Hash       []byte     `json:"hash"`
		MerkleRoot MerkleRoot `json:"merkleRoot"`
		PrevHash   []byte     `json:"prevHash"`
		Nonce      int64      `json:"nonce"`
		Difficulty int64      `json:"difficulty"`
		Timestamp  int64      `json:"timestamp"`
	}{
		Hash:       block.hash,
		MerkleRoot: block.merkleRoot,
		PrevHash:   block.prevHash,
		Nonce:      block.nonce,
		Difficulty: block.difficulty,
		Timestamp:  block.timestamp,
	})
}

func (block *Block) UnmarshalJSON(data []byte) error {
	var json = jsoniter.ConfigFastest
	var blockJSON struct {
		Hash       []byte     `json:"hash"`
		MerkleRoot MerkleRoot `json:"merkleRoot"`
		PrevHash   []byte     `json:"prevHash"`
		Nonce      int64      `json:"nonce"`
		Difficulty int64      `json:"difficulty"`
		Timestamp  int64      `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &blockJSON); err != nil {
		return err
	}

	block.hash = blockJSON.Hash
	block.merkleRoot = blockJSON.MerkleRoot
	block.prevHash = blockJSON.PrevHash
	block.nonce = blockJSON.Nonce
	block.difficulty = blockJSON.Difficulty
	block.timestamp = blockJSON.Timestamp

	return nil
}

func Serialize(b *Block) ([]byte, error) {
	var json = jsoniter.ConfigFastest
	return json.Marshal(b)
}

func Deserialize(data []byte) (*Block, error) {
	var json = jsoniter.ConfigFastest
	b := &Block{}

	err := json.Unmarshal(data, b)

	return b, err
}

func GenesisBlock(root MerkleRoot, difficulty int64) *Block {
	return NewBlock(root, []byte{}, difficulty)
}
