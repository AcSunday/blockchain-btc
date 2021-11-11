package main

func main() {
	bc := NewBlockChain()
	bc.AddBlock("sunday给alex转1个BTC")
	bc.AddBlock("alex给egon转0.1个BTC")
	bc.AddBlock("momo给egon转0.1个BTC")

	/*
	for i, block := range bc.blocks {
		fmt.Printf("===== 当前区块高度 %d =====\n", i)
		fmt.Printf("前区块hash值: %x\n", block.PrevHash)
		fmt.Printf("当前区块hash值: %x\n", block.Hash)
		fmt.Printf("区块数据: %s\n", block.Data)
	}
	 */
}
