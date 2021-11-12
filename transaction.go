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
