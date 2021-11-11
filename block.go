package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"
)

// 0. 定义结构
type Block struct {
	// 版本号
	Version uint64
	// 1. 前区块hash
	PrevHash []byte
	// Merkel根（梅克尔根，这就是一个hash值，我们先不管，V4再介绍）
	MerkelRoot []byte
	// 时间戳
	TimeStamp uint64
	// 难度值
	Difficulty uint64
	// 随机数，也就是挖矿要找的数据
	Nonce uint64
	// 2. 当前区块hash，正常BTC区块中没有当前区块的hash，我们是为了方便做了简化
	Hash []byte
	// 3. 数据
	Data []byte
}

// 1. 补充区块字段
// 2. 更新计算hash函数
// 3. 优化代码

// 实现一个辅助函数，功能是将uint64转成[]byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panicln(err)
	}
	return buffer.Bytes()
}

// 2. 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkelRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
		Hash:       []byte{},
		Data:       []byte(data),
	}
	//block.SetHash()

	// 创建一个pow对象
	pow := NewProofOfWork(block)
	// 查找随机数，不停的进行hash运算
	hash, nonce := pow.Run()
	block.Hash = hash
	block.Nonce = nonce

	return block
}

func (b *Block) ToBytes() []byte {
	return []byte{}
}

/*
// 3. 生成hash
func (b *Block) SetHash() {
	// 1. 拼装数据
	//var blockInfo []byte
	//blockInfo = append(blockInfo, Uint64ToByte(b.Version)...)
	//blockInfo = append(blockInfo, b.PrevHash...)
	//blockInfo = append(blockInfo, b.MerkelRoot...)
	//blockInfo = append(blockInfo, Uint64ToByte(b.TimeStamp)...)
	//blockInfo = append(blockInfo, Uint64ToByte(b.Difficulty)...)
	//blockInfo = append(blockInfo, Uint64ToByte(b.Nonce)...)
	//blockInfo = append(blockInfo, b.Data...)

	tmp := [][]byte{
		Uint64ToByte(b.Version),
		b.PrevHash,
		b.MerkelRoot,
		Uint64ToByte(b.TimeStamp),
		Uint64ToByte(b.Difficulty),
		Uint64ToByte(b.Nonce),
		b.Data,
	}
	blockInfo := bytes.Join(tmp, []byte{})

	// 2. sha256
	hash := sha256.Sum256(blockInfo)
	b.Hash = hash[:]
}
*/
