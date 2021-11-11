package main

import (
	"github.com/boltdb/bolt"
	"log"
)

type BlockChainIterator struct {
	db                 *bolt.DB
	currentHashPointer []byte
}

func (bc *BlockChain) NewIterator() *BlockChainIterator {
	return &BlockChainIterator{
		db:                 bc.db,
		currentHashPointer: bc.tail,
	}
}

// 迭代器是属于区块链的
// Next方法是属于迭代器的
func (it *BlockChainIterator) Next() *Block {
	// 1. 返回当前区块
	// 2. 指针前移
	var block *Block

	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockChainBucket))
		if bucket == nil {
			log.Panic("bucket must be not nil")
		}
		blockTmp := bucket.Get(it.currentHashPointer)
		// 解码动作
		block = Deserialize(blockTmp)
		// 游标hash左移
		it.currentHashPointer = block.PrevHash

		return nil
	})
	return block
}
