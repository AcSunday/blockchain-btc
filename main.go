package main

func main() {
	bc := NewBlockChain("sunday")
	cli := &CLI{bc: bc}
	cli.Run()
}
