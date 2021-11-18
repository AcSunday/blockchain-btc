package main

import (
	"fmt"
	"log"
	"time"
)

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	iterator := bc.NewIterator()
	for {
		// 返回区块，游标左移
		block := iterator.Next()
		date := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("===== 当前区块高度 %d =====\n", 0)
		fmt.Printf("终端版本: %d\n", block.Version)
		fmt.Printf("前区块hash值: %x\n", block.PrevHash)
		fmt.Printf("梅克尔根hash值: %x\n", block.MerkelRoot)
		fmt.Printf("块产生时间: %s\n", date)
		fmt.Printf("块难度: %d\n", block.Difficulty)
		fmt.Printf("随机数: %d\n", block.Nonce)
		fmt.Printf("当前区块hash值: %x\n", block.Hash)
		fmt.Printf("当前区块数据: %s\n", block.Transactions[0].TxInputs[0].Sig)

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) GetBalance(addr string) {
	utxos := cli.bc.FindUTXOs(addr)

	amount := 0.0
	for _, utxo := range utxos {
		amount += utxo.Amount
	}
	log.Printf("%s balance: %f\n", addr, amount)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	// 1. 创建挖矿交易
	coinbase := NewCoinBaseTx(miner, data)
	// 2. 创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if coinbase == nil || tx == nil {
		return
	}
	// 3. 添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
}

func (cli *CLI) NewWallet() {
	//wallet := NewWallet()
	//address := wallet.NewAddress()
	//fmt.Printf("private key: %x\n", *wallet.Private)
	//fmt.Printf("public key: %x\n", wallet.PubKey)
	//fmt.Printf("address: %s\n", address)
	wallets := NewWallets()
	for address, _ := range wallets.WalletsMap {
		fmt.Printf("your new address: %s\n", address)
	}
}
