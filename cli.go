package main

import (
	"fmt"
	"log"
	"os"
)

// 用来接收命令行参数并且控制区块链操作

const Usage = `
    addBlock --data DATA        "add data to blockchain"
    printChain                  "print all blockchain data"
`

type CLI struct {
	bc *BlockChain
}

// 接收参数按情况执行
func (cli *CLI) Run() {
	// 1. 得到命令
	args := os.Args
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}
	// 2. 分析命令
	// 3. 执行相应动作
	cmd := args[1]
	switch cmd {
	case "addBlock":
		// 添加区块
		if len(args) == 4 && args[2] == "--data" {
			data := args[3]
			cli.AddBlock(data)
		} else {
			log.Println("missing params")
			fmt.Printf(Usage)
		}
	case "printChain":
		// 打印区块
		cli.PrintBlockChain()
	default:
		fmt.Printf(Usage)
	}
}