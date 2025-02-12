package main

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
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

	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			log.Println("miner verify tx failed")
			return
		}
	}

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
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []*TxOutput {
	var utxos = make([]*TxOutput, 0, 4)
	transactions, spentOutputs := bc.FindUTXOTransactions(pubKeyHash)
	// 1. 遍历区块
	// 2. 遍历交易
	// 3. 遍历output，找到和自己相关的UTXO(在添加output之前，检查是否已经消耗过)
	// 4. 遍历input，找到自己花费过的UTXO(把自己消耗过的给标识出来)

	for _, tx := range transactions {
		for i, output := range tx.TxOutputs {
			// 在这里做一个过滤，将所有消耗过的output和当前即将要添加的output对比一下
			// 如果相同则跳过
			key := fmt.Sprintf("%x:%d", tx.TxID, i)
			if _, ok := spentOutputs[key]; ok { // 当前准备添加的output已经消耗了，不要加了
				continue
			}

			// 这个output和我们的目标地址相同，加到返回的UTXOs切片中
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				utxos = append(utxos, output)
			}
		}
	}

	return utxos
}

// 找到足够转账额的UTXO
//  @return map[string][]int 以map[TxID][]int{outputIndex1, outputIndex2 ...}形式返回
//  @return float64 返回需要的余额或者总余额
func (bc *BlockChain) FindNeedUTXOs(senderPubKeyHash []byte, amount float64) (map[string][]int, float64) {
	var utxos = make(map[string][]int)
	var totalAmount float64
	transactions, spentOutputs := bc.FindUTXOTransactions(senderPubKeyHash)

	for _, tx := range transactions {
		for i, output := range tx.TxOutputs {
			// 在这里做一个过滤，将所有消耗过的output和当前即将要添加的output对比一下
			key := fmt.Sprintf("%x:%d", tx.TxID, i)
			if _, ok := spentOutputs[key]; ok { // 当前准备添加的output已经消耗了，不要加了
				continue
			}

			// 这个output和我们的目标地址相同，加到返回的utxos map中
			if bytes.Equal(senderPubKeyHash, output.PubKeyHash) {
				utxos[string(tx.TxID)] = append(utxos[string(tx.TxID)], i)
				totalAmount += output.Amount
				if totalAmount >= amount { // 目前找到的utxo余额足够支付，直接return
					return utxos, totalAmount
				}
			}
		}
	}

	return utxos, totalAmount
}

func (bc *BlockChain) FindUTXOTransactions(pubKeyHash []byte) ([]*Transaction, map[string]struct{}) {
	var txs = make([]*Transaction, 0, 8)
	var spentOutputs = make(map[string]struct{})
	// 1. 遍历区块
	// 2. 遍历交易
	// 3. 遍历output，找到和自己相关的UTXO
	// 4. 遍历input，找到自己花费过的UTXO(把自己消耗过的给标识出来)

	it := bc.NewIterator()
	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			for _, output := range tx.TxOutputs {
				// 这个output和我们的目标地址相同，加到返回的txs切片中
				if bytes.Equal(pubKeyHash, output.PubKeyHash) {
					txs = append(txs, tx)
					break
				}
			}

			// 如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if tx.IsCoinBase() {
				continue
			}

			//map[交易id:索引下标]struct{}
			for _, input := range tx.TxInputs {
				// 判断一下当前这个input和目标地址是否一致，如果相同说明是消耗过的output 则加进来
				if bytes.Equal(HashPubKey(input.PubKey), pubKeyHash) {
					key := fmt.Sprintf("%x:%d", input.TxID, input.Index)
					spentOutputs[key] = struct{}{}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return txs, spentOutputs
}

// 根据id查找交易本身，需要遍历整个区块链
func (bc *BlockChain) FindTransactionByTxid(txID []byte) (*Transaction, error) {
	// 1. 遍历区块链
	// 2. 遍历交易
	// 3. 比较交易，找到了直接退出
	// 4. 如果没找到，返回空Transaction，同时返回错误状态

	it := bc.NewIterator()
	for {
		block := it.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.TxID, txID) {
				return tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return nil, errors.New("invalid txid, not found tx")
}

// 签名交易
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	prevTxs := make(map[string]*Transaction)
	// 找到所有引用的交易
	// 1. 根据inputs来找，有多少input就遍历多少次
	// 2. 找到目标交易
	// 3. 添加到prevTxs里面
	for _, input := range tx.TxInputs {
		// 根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTxid(input.TxID)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[string(input.TxID)] = tx
	}
	tx.Sign(privateKey, prevTxs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {

	if tx.IsCoinBase() {
		return true
	}

	prevTxs := make(map[string]*Transaction)
	// 找到所有引用的交易
	// 1. 根据inputs来找，有多少input就遍历多少次
	// 2. 找到目标交易
	// 3. 添加到prevTxs里面
	for _, input := range tx.TxInputs {
		// 根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTxid(input.TxID)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[string(input.TxID)] = tx
	}

	return tx.Verify(prevTxs)
}
