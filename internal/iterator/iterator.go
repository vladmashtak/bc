package iterator

import (
	"bbu/internal/block"
	"errors"
	"github.com/dgraph-io/badger"
	"log"
)

type Iterator struct {
	currentHash []byte
	db          *badger.DB
}

func NewIterator(db *badger.DB, lastHash []byte) *Iterator {
	return &Iterator{
		currentHash: lastHash,
		db:          db,
	}
}

func (i *Iterator) Next() *block.Block {
	var currentBlock *block.Block
	err := i.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(i.currentHash)
		if err != nil {
			return errors.Join(errors.New("iterator Next: get item from view"), err)
		}

		value, err := item.ValueCopy(nil)
		if err != nil {
			return errors.Join(errors.New("iterator Next: copy value"), err)
		}

		currentBlock, err = block.Deserialize(value)
		if err != nil {
			return errors.Join(errors.New("iterator Next: deserialize value"), err)
		}

		return nil
	})

	if err != nil {
		log.Println(err)
		return currentBlock
	}

	if currentBlock != nil {
		i.currentHash = currentBlock.PrevHash()
	}

	return currentBlock
}
