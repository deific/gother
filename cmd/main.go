package main

import (
	"fmt"
	"gother/internal/block"
	"time"
)

func main() {

	blockChain := block.CreateBlockChain()
	time.Sleep(time.Second)
	blockChain.AddBlock("Block 2")
	time.Sleep(time.Second)
	blockChain.AddBlock("Block 3")

	for _, block := range blockChain.Blocks {
		fmt.Printf("block is:%+v\n", block)
	}
}
