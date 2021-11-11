package main

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	var xiaoMing = Person{
		Name: "小明",
		Age:  18,
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(&xiaoMing)
	if err != nil {
		log.Panic("encode failed")
	}
	log.Printf("编码后的小明: %v\n", buffer.Bytes())

	// 解码
	decoder := gob.NewDecoder(bytes.NewReader(buffer.Bytes()))
	var daMing Person
	err = decoder.Decode(&daMing)
	if err != nil {
		log.Panic("decode failed")
	}
	log.Printf("解码后的小明: %v\n", daMing)
}
