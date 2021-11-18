package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"
)

// 演示如何使用ecdsa生成公钥和私钥
// 签名和校验

func main() {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	// 生成公钥
	pubKey := privateKey.PublicKey

	data := "hello world"
	hash := sha256.Sum256([]byte(data))

	// 签名
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Panic(err)
	}

	log.Printf("pubkey %v\n", pubKey)
	log.Printf("r len: %d %v\n", len(r.Bytes()), r.Bytes())
	log.Printf("s len: %d %v\n", len(s.Bytes()), s.Bytes())

	// 把r, s进行序列化传输
	signature := append(r.Bytes(), s.Bytes()...)
	// ....

	// 1. 定义两个辅助的big.Int
	r1 := new(big.Int)
	s1 := new(big.Int)
	// 2. 拆分我们signature，前半部分给r，后半部分给s
	idx := len(signature) / 2
	r1.SetBytes(signature[:idx])
	s1.SetBytes(signature[idx:])

	// 校验需要三个东西: 数据、签名、公钥
	//verify := ecdsa.Verify(&pubKey, hash[:], r, s)
	verify := ecdsa.Verify(&pubKey, hash[:], r1, s1)
	log.Println(verify)
}
