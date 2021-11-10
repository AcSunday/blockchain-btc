package main

import (
	"bytes"
	"crypto/sha256"
	"log"
	"math/big"
)

// 1. 定义pow结构
type ProofOfWork struct {
	block *Block
	// 一个非常大的数，它有很丰富的方法: 比较、赋值
	target *big.Int
}

// 2. 提供创建POW的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := &ProofOfWork{
		block:  block,
		target: big.NewInt(1),
	}

	// 我们指定的难度值，现在是一个string类型，需要转换
	targetStr := "0000100000000000000000000000000000000000000000000000000000000000"
	//targetStr := "4f0f23fd637149e0db995b17c8ca26d1209119480f5d466f239d9676262753f8"
	tmpInt, _ := new(big.Int).SetString(targetStr, 16)
	pow.target = tmpInt
	return pow
}

// 3. 提供不断计算hash的函数
func (pow *ProofOfWork) Run() (hash []byte, nonce uint64) {
	// 拼装数据(区块数据，还有不断变化的随机数)
	// 做hash运算
	// 与pow中的target进行比较

	var calcHash [32]byte
	b := pow.block

	for {
		tmp := [][]byte{
			Uint64ToByte(b.Version),
			b.PrevHash,
			b.MerkelRoot,
			Uint64ToByte(b.TimeStamp),
			Uint64ToByte(b.Difficulty),
			Uint64ToByte(nonce),
			b.Data,
		}
		blockInfo := bytes.Join(tmp, []byte{})

		calcHash = sha256.Sum256(blockInfo)
		tmpInt := new(big.Int).SetBytes(calcHash[:])
		if tmpInt.Cmp(pow.target) == -1 {
			log.Printf("miner found block, hash: %x, nonce: %d", calcHash, nonce)
			hash = calcHash[:]
			break
		}
		nonce++
	}

	return
}

// 4. 提供一个校验函数
