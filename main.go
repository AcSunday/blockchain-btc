package main

func main() {
	bc := NewBlockChain("0xasAGagdGHfsdq42dw0092")
	cli := &CLI{bc: bc}
	cli.Run()
}
