package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// 用来接收命令行参数并且控制区块链操作

const Usage = `
    printChain                      "print all blockchain data"
    getBalance --address ADDRESS    "get address balance"
    send FROM TO AMOUNT MINER DATA  "send coin to one, the Miner write data"
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
	case "printChain":
		// 打印区块
		cli.PrintBlockChain()
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			addr := args[3]
			cli.GetBalance(addr)
		} else {
			log.Println("missing params")
			fmt.Printf(Usage)
		}
	case "send":
		if len(args) != 7 {
			log.Println("missing params")
			fmt.Printf(Usage)
			return
		}
		// send FROM TO AMOUNT MINER DATA
		from := args[2]
		to := args[3]
		amount, _ := strconv.ParseFloat(args[4], 64)
		miner := args[5]
		data := args[6]
		cli.Send(from, to, amount, miner, data)
	default:
		fmt.Printf(Usage)
	}
}
