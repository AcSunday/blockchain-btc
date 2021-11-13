package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

// 1. 定义交易结构
// 2. 提供创建交易方法
// 3. 创建挖矿交易
// 4. 根据交易调整程序

const Reward = 12.5

type Transaction struct {
	TxID      []byte      // 交易ID
	TxInputs  []*TxInput  // 交易输入数组
	TxOutputs []*TxOutput // 交易输出数组
}

type TxInput struct {
	TxID  []byte // 引用的交易ID
	Index int64  // 引用的output的索引值
	Sig   string // 解锁脚本，我们用地址来模拟
}

type TxOutput struct {
	Amount     float64 // 转账金额
	PubKeyHash string  // 锁定脚本，我们用地址模拟
}

func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TxID = hash[:]
}

// 实现一个函数，判断当前的交易是否为挖矿交易
func (tx *Transaction) IsCoinBase() bool {
	// 1. 交易的input只有一个
	// 2. 交易id为空
	// 3. 交易的index为-1
	if len(tx.TxInputs) == 1 {
		input := tx.TxInputs[0]
		if bytes.Equal(input.TxID, []byte{}) && input.Index == -1 {
			return true
		}
	}
	return false
}

func NewCoinBaseTx(addr string, data string) *Transaction {
	// 挖矿交易的特点:
	// 1. 只有一个input
	// 2. 无需引用交易id
	// 3. 无需引用output 的 index

	// 矿工由于挖矿时无需指定签名，所以Sig这个字段可以由矿工自由填写数据，一般是填写矿池的名字
	input := &TxInput{
		TxID:  []byte{},
		Index: -1,
		Sig:   data,
	}
	output := &TxOutput{
		Amount:     Reward,
		PubKeyHash: addr,
	}

	// 对于挖矿交易来说，只有一个input和一个output
	tx := &Transaction{
		TxID:      []byte{},
		TxInputs:  []*TxInput{input},
		TxOutputs: []*TxOutput{output},
	}
	tx.SetHash()
	return tx
}
