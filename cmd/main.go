package main

import (
	"fmt"
	"gother/internal/blockchain"
	"time"
)

func main() {

	blockChain := blockchain.CreateBlockChain()
	time.Sleep(time.Second)
	blockChain.AddBlock("Block 2")
	time.Sleep(time.Second)
	blockChain.AddBlock("Block 3")

	for _, block := range blockChain.Blocks {
		fmt.Printf("blockchain Height:%d\n", block.Height)
		fmt.Printf("blockchain Timestamp:%d\n", block.Timestamp)
		fmt.Printf("blockchain Hash:%x\n", block.Hash) // 转为十六进制输出
		fmt.Printf("blockchain PreHash:%x\n", block.PreHash)
		fmt.Printf("blockchain Data:%s\n", block.Data)
		fmt.Println("============================================")
	}
}
