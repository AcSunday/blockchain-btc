package main

import (
	"fmt"
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
func NewBlockChain(addr string) *BlockChain {
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
			genesisBlock := GenesisBlock(addr)
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
func GenesisBlock(addr string) *Block {
	coinbase := NewCoinBaseTx(addr, "BTC创世块，老牛逼了")
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// 6. 添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	// 获取最后一个区块的hash
	db := bc.db
	lastHash := bc.tail

	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockChainBucket))
		if bucket == nil {
			log.Panic("bucket must be not nil")
		}

		// a. 创建新的区块
		block := NewBlock(txs, lastHash)
		// b. 添加到区块链到DB中
		bucket.Put(block.Hash, block.Serialize())
		bucket.Put([]byte(lastHashKey), block.Hash)
		bc.tail = block.Hash

		return nil
	})
}

// 找到指定地址的所有的UTXO
func (bc *BlockChain) FindUTXOs(addr string) []*TxOutput {
	var utxos = make([]*TxOutput, 0, 4)
	var spentOutputs = make(map[string]struct{})
	// 1. 遍历区块
	// 2. 遍历交易
	// 3. 遍历output，找到和自己相关的UTXO(在添加output之前，检查是否已经消耗过)
	// 4. 遍历input，找到自己花费过的UTXO(把自己消耗过的给标识出来)

	it := bc.NewIterator()
	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			for i, output := range tx.TxOutputs {
				// 在这里做一个过滤，将所有消耗过的output和当前即将要添加的output对比一下
				// 如果相同则跳过
				key := fmt.Sprintf("%x:%d", tx.TxID, i)
				if _, ok := spentOutputs[key]; ok { // 当前准备添加的output已经消耗了，不要加了
					continue
				}

				// 这个output和我们的目标地址相同，加到返回的UTXOs切片中
				if output.PubKeyHash == addr {
					utxos = append(utxos, output)
				}
			}

			// 如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if tx.IsCoinBase() {
				continue
			}

			//map[交易id:索引下标]struct{}
			for _, input := range tx.TxInputs {
				// 判断一下当前这个input和目标地址是否一致，如果相同说明是消耗过的output 则加进来
				if input.Sig == addr {
					key := fmt.Sprintf("%x:%d", input.TxID, input.Index)
					spentOutputs[key] = struct{}{}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return utxos
}

func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]int, float64) {
	var utxos = make(map[string][]int)
	var totalAmount float64
	var spentOutputs = make(map[string]struct{})

	it := bc.NewIterator()
	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			for i, output := range tx.TxOutputs {
				// 在这里做一个过滤，将所有消耗过的output和当前即将要添加的output对比一下
				key := fmt.Sprintf("%x:%d", tx.TxID, i)
				if _, ok := spentOutputs[key]; ok { // 当前准备添加的output已经消耗了，不要加了
					continue
				}

				// 这个output和我们的目标地址相同，加到返回的UTXOs切片中
				if output.PubKeyHash == from {
					if utxos[string(tx.TxID)] == nil {
						utxos[string(tx.TxID)] = make([]int, 0, 4)
					}
					utxos[string(tx.TxID)] = append(utxos[string(tx.TxID)], i)
					totalAmount += output.Amount
				}
			}

			// 如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if tx.IsCoinBase() {
				continue
			}

			//map[交易id:索引下标]struct{}
			for _, input := range tx.TxInputs {
				// 判断一下当前这个input和目标地址是否一致，如果相同说明是消耗过的output 则加进来
				if input.Sig == from {
					key := fmt.Sprintf("%x:%d", input.TxID, input.Index)
					spentOutputs[key] = struct{}{}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return utxos, totalAmount
}
