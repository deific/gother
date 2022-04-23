package main

import (
	"encoding/binary"
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
		fmt.Printf("blockchain Hash:%d\n", binary.BigEndian.Uint64(block.Hash)) // 转为十六进制输出
		fmt.Printf("blockchain PreHash:%d\n", binary.BigEndian.Uint64(block.PreHash))
		fmt.Printf("blockchain target:%d\n", binary.BigEndian.Uint64(block.Target))
		fmt.Printf("blockchain nonce:%d\n", block.Nonce)
		fmt.Printf("blockchain Data:%s\n", block.Data)
		fmt.Println("============================================")
	}
}
