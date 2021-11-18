package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

// 这里的钱包是一结构，每一个钱包保存了公钥，私钥对

type Wallet struct {
	Private *ecdsa.PrivateKey
	//PubKey  *ecdsa.PublicKey
	// 约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分(参考r,s传递)
	PubKey []byte
}

// 创建钱包
func NewWallet() *Wallet {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 生成公钥，拼接X和Y
	pubKeyOrig := privateKey.PublicKey
	pubKey := append(pubKeyOrig.X.Bytes(), pubKeyOrig.Y.Bytes()...)
	return &Wallet{
		Private: privateKey,
		PubKey:  pubKey,
	}
}

// 生成地址
func (w *Wallet) NewAddress() string {
	pubKey := w.PubKey

	rip160HashValue := HashPubKey(pubKey)

	// 拼接version
	version := byte(00)
	payload := append([]byte{version}, rip160HashValue...)

	// checksum
	checkCode := CheckSum(payload)

	// payload append完后是25字节的数据
	payload = append(payload, checkCode...)

	// go语言有一个库，叫做btcd，这个是go语言实现的BTC全节点源码
	address := base58.Encode(payload)

	return address
}

func HashPubKey(data []byte) []byte {
	hash := sha256.Sum256(data)

	// 理解为编码器
	rip160Hasher := ripemd160.New()
	_, err := rip160Hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	// 返回rip160的hash结果
	rip160HashValue := rip160Hasher.Sum(nil)
	return rip160HashValue
}

func CheckSum(data []byte) []byte {
	// 两次sha256
	hash1 := sha256.Sum256(data)
	hash2 := sha256.Sum256(hash1[:])
	checkCode := hash2[:4] // 前4字节校验码
	return checkCode
}
