package main

import (
	"log"
	"os"
)

func main() {
	//len1 := len(os.Args)
	for i, cmd := range os.Args {
		log.Printf("arg[%d] val:%v\n", i, cmd)
	}
}
