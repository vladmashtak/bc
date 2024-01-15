package blockchain

import (
	"bbu/internal/block"
	"bbu/internal/iterator"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"log"
	"os"
)

const (
	dbPath     = "tmp/blocks"
	dbFile     = "tmp/blocks/MANIFEST"
	difficulty = 12
)

var (
	lastHashKey = []byte("last_hash")
	genesisData = block.MerkleRoot("First Transaction from Genesis")
)

type BlockChain struct {
	lastHash []byte
	db       *badger.DB
}

func NewBlockChain() (*BlockChain, error) {
	chain := &BlockChain{}

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		return nil, errors.Join(errors.New("NewBlockChain: open db"), err)
	}
	chain.db = db

	var lastHash []byte
	if _, err = os.Stat(dbFile); os.IsNotExist(err) { // check if file database exists
		lastHash, err = chain.createBlockChain()
	} else {
		lastHash, err = chain.continueBlockChain()
	}
	if err != nil {
		return nil, errors.Join(errors.New("NewBlockChain: "), err)
	}
	chain.lastHash = lastHash

	return chain, err
}

func (chain *BlockChain) createBlockChain() ([]byte, error) {
	var lastHash []byte
	err := chain.db.Update(func(txn *badger.Txn) error {
		fmt.Println("Genesis create")
		genesis := block.GenesisBlock(genesisData, difficulty)

		blockBytes, err := block.Serialize(genesis)
		if err != nil {
			return errors.Join(errors.New("createBlockChain: serialize genesis block"), err)
		}

		if err := txn.Set(genesis.Hash(), blockBytes); err != nil {
			return errors.Join(errors.New("createBlockChain: set initial block"), err)
		}

		if err := txn.Set(lastHashKey, genesis.Hash()); err != nil {
			return errors.Join(errors.New("createBlockChain: set initial last hash"), err)
		}

		lastHash = genesis.Hash()

		return nil
	})

	return lastHash, err
}

func (chain *BlockChain) continueBlockChain() ([]byte, error) {
	var lastHash []byte
	err := chain.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(lastHashKey)
		if err != nil {
			return errors.Join(errors.New("continueBlockChain: get last hash"), err)
		}

		lastHash, err = item.ValueCopy(nil)
		if err != nil {
			return errors.Join(errors.New("continueBlockChain: item copy value"), err)
		}

		return nil
	})

	return lastHash, err
}

func (chain *BlockChain) Close() error {
	if chain.db == nil {
		return errors.New("database not initialized")
	}

	return chain.db.Close()
}

func (chain *BlockChain) AddBlock(data string) {
	err := chain.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(lastHashKey)
		if err != nil {
			return errors.Join(errors.New("AddBlock: get last hash key"), err)
		}

		lastHash, err := item.ValueCopy(nil)
		if err != nil {
			return errors.Join(errors.New("AddBlock: value copy"), err)
		}

		newBlock := block.NewBlock(block.MerkleRoot(data), lastHash, difficulty)
		blockValue, err := block.Serialize(newBlock)
		if err != nil {
			return errors.Join(errors.New("AddBlock: serialize new block"), err)
		}

		err = txn.Set(newBlock.Hash(), blockValue)
		if err != nil {
			return errors.Join(errors.New("AddBlock: set new block"), err)
		}

		err = txn.Set(lastHashKey, newBlock.Hash())
		if err != nil {
			return errors.Join(errors.New("AddBlock: set new block"), err)
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}

}

func (chain *BlockChain) Iterator() *iterator.Iterator {
	return iterator.NewIterator(chain.db, chain.lastHash)
}
