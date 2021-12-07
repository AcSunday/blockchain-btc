package main

func main() {
	bc := NewBlockChain("1HhH22Ugs1yap3oaAdnnLiFbrEVj45pHwC")
	cli := &CLI{bc: bc}
	cli.Run()
}
