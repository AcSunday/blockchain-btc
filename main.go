package main

func main() {
	bc := NewBlockChain("1HDUPBxoxYwqDY2r78XJiE4TX98ra4hWvm")
	cli := &CLI{bc: bc}
	cli.Run()
}
