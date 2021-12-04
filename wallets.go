package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"github.com/btcsuite/btcutil/base58"
	"io/ioutil"
	"log"
	"os"
)

// 定义一个Wallets结构，它保存所有的wallet以及它的地址
type Wallets struct {
	//map[地址]钱包
	WalletsMap map[string]*Wallet
}

// 创建方法
func NewWallets() *Wallets {
	var ws Wallets
	ws.WalletsMap = make(map[string]*Wallet)
	ws.loadWallets()
	return &ws
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()
	ws.WalletsMap[address] = wallet

	ws.saveWallets()
	return address
}

// 保存方法，把新建的wallet添加进去
func (ws *Wallets) saveWallets() {
	var buffer bytes.Buffer
	gob.Register(elliptic.P256()) // 注册这个Curve
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile("wallet.dat", buffer.Bytes(), 0600)
	if err != nil {
		log.Panic(err)
	}
}

// 读取文件方法，把所有的wallet读出来
func (ws *Wallets) loadWallets() {
	// 文件不存在直接退出
	_, err := os.Stat("wallet.dat")
	if err != nil && os.IsNotExist(err) {
		return
	}

	content, err := ioutil.ReadFile("wallet.dat")
	if err != nil {
		log.Panic(err)
	}

	// 解码
	gob.Register(elliptic.P256()) // 注册这个Curve
	decoder := gob.NewDecoder(bytes.NewReader(content))
	err = decoder.Decode(ws)
	if err != nil {
		log.Panic(err)
	}
}

// 获取所有的address
func (ws *Wallets) GetAllAddress() []string {
	var ret []string
	for address := range ws.WalletsMap {
		ret = append(ret, address)
	}
	return ret
}

// 通过地址返回公钥的hash
func GetPubKeyFromAddress(addr string) []byte {
	// 1. 解码
	addrByte := base58.Decode(addr) // 25字节
	// 2. 截取出公钥hash：取出version（1字节），取出校验码（4字节）
	pubKeyHash := addrByte[1 : len(addrByte)-4]

	return pubKeyHash
}
