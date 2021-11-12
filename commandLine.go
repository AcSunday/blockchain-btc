package main

import (
	"fmt"
	"log"
	"time"
)

func (cli *CLI) AddBlock(data string) {
	//cli.bc.AddBlock(data)
	log.Println("add block to blockchain finished")
}

func (cli *CLI) PrintBlockChain() {
	bc := cli.bc
	iterator := bc.NewIterator()
	for {
		// 返回区块，游标左移
		block := iterator.Next()
		date := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("===== 当前区块高度 %d =====\n", 0)
		fmt.Printf("前区块hash值: %x\n", block.PrevHash)
		fmt.Printf("当前终端版本: %d\n", block.Version)
		fmt.Printf("当前块产生时间: %s\n", date)
		fmt.Printf("当前块难度: %d\n", block.Difficulty)
		fmt.Printf("当前块随机数: %d\n", block.Nonce)
		fmt.Printf("当前区块hash值: %x\n", block.Hash)
		fmt.Printf("区块数据: %s\n", block.Transactions[0].TxInputs[0].Sig)

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
