package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("sunday给alex转1个BTC")
	bc.AddBlock("alex给egon转0.1个BTC")
	bc.AddBlock("momo给egon转0.1个BTC")

	// 调用迭代器，返回每一个区块数据
	iterator := bc.NewIterator()
	for {
		// 返回区块，游标左移
		block := iterator.Next()
		fmt.Printf("===== 当前区块高度 %d =====\n", 0)
		fmt.Printf("前区块hash值: %x\n", block.PrevHash)
		fmt.Printf("当前区块hash值: %x\n", block.Hash)
		fmt.Printf("区块数据: %s\n", block.Data)

		if len(block.PrevHash) == 0 {
			break
		}
	}

	//for i, block := range bc.blocks {
	//	fmt.Printf("===== 当前区块高度 %d =====\n", i)
	//	fmt.Printf("前区块hash值: %x\n", block.PrevHash)
	//	fmt.Printf("当前区块hash值: %x\n", block.Hash)
	//	fmt.Printf("区块数据: %s\n", block.Data)
	//}
}
