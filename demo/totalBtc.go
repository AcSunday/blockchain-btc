package main

import "fmt"

func main() {
	fmt.Println("hello")
	// 1. 每21万个区块减半
	// 2. 最初奖励50个BTC
	// 3. 用一个循环来判断，累加

	total := 0.0
	blockInterval := 21.0 // 单位 万
	currentReward := 50.0

	for currentReward > 0 {
		// 每一个区间内的总量
		amount1 := blockInterval * currentReward
		//currentReward /= 2
		currentReward *= 0.5 // 除效率低，我们使用等价的乘法
		total += amount1
	}

	fmt.Println("BTC总量:", total, "万")
}
