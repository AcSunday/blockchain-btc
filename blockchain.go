package main

import (
	"github.com/boltdb/bolt"
	"log"
)

const (
	blockChainDB     = "blockChain.db"
	blockChainBucket = "blockBucket"
	lastHashKey      = "LastHashKey"
)

// 4. 引入区块链
type BlockChain struct {
	// 定一个区块链切片
	//blocks []*Block
	db   *bolt.DB
	tail []byte // 存储最后一个区块的hash
}

// 5. 定义一个区块链
func NewBlockChain() *BlockChain {
	// 最后一个区块的hash，从DB读出来的
	var lastHash []byte

	// 1. 打开数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)
	if err != nil {
		log.Panic("open db failed")
	}
	//defer db.Close()

	// 2. 找到抽屉bucket(如果没有就创建)
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockChainBucket))
		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockChainBucket))
			if err != nil {
				log.Panic("create bucket failed")
			}

			// 创建一个创世块，并作为第一个区块添加到区块链中
			genesisBlock := GenesisBlock()
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte(lastHashKey), genesisBlock.Hash)
		}

		lastHash = bucket.Get([]byte(lastHashKey))

		return nil
	})

	return &BlockChain{
		db:   db,
		tail: lastHash,
	}
}

// 定义一个创世块
func GenesisBlock() *Block {
	return NewBlock("BTC创世块，老牛逼了", []byte{})
}

// 6. 添加区块
func (bc *BlockChain) AddBlock(data string) {
	// 获取最后一个区块的hash
	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockChainBucket))
		if bucket == nil {
			log.Panic("bucket must be not nil")
		}

		// a. 创建新的区块
		block := NewBlock(data, lastHash)
		// b. 添加到区块链到DB中
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte(lastHashKey), block.Hash)
		bc.tail = block.Hash

		return nil
	})
}
