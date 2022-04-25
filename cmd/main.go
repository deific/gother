package main

import (
	"encoding/binary"
	"fmt"
	"gother/internal/blockchain"
	"gother/internal/constant"
	"gother/internal/transaction"
)

func main() {
	startChain()
}

func startChain() {

	// 1.创始区块和创始交易
	chain := blockchain.CreateBlockChain()
	// balance=1000
	balance, _ := chain.FindUTXOs([]byte(constant.BaseAddress))
	fmt.Println("Balance of gother:", balance)

	// 2.创建一个交易，从BaseAddress 转出 100
	txPool := make([]*transaction.Transaction, 0)
	tmpTx, ok := chain.CreateTransaction([]byte(constant.BaseAddress), []byte("A"), 100)
	if ok {
		txPool = append(txPool, tmpTx)
	}
	chain.Mine(txPool)
	// balance=1000 - 100 = 900
	balance, _ = chain.FindUTXOs([]byte(constant.BaseAddress))
	fmt.Println("Balance of gother:", balance)

	// 创建一个交易，从BaseAddress 转出 100
	txPool = make([]*transaction.Transaction, 0)
	tmpTx, ok = chain.CreateTransaction([]byte(constant.BaseAddress), []byte("B"), 200)
	if ok {
		txPool = append(txPool, tmpTx)
	}
	// 打包区块
	chain.Mine(txPool)
	// balance=900 - 200 = 700
	balance, _ = chain.FindUTXOs([]byte(constant.BaseAddress))
	fmt.Println("Balance of gother:", balance)

	printChain(chain)
}

func printChain(blockChain *blockchain.Blockchain) {
	for _, block := range blockChain.Blocks {
		fmt.Printf("blockchain Height:%d\n", block.Height)
		fmt.Printf("blockchain Timestamp:%d\n", block.Timestamp)
		fmt.Printf("blockchain Hash:%d\n", binary.BigEndian.Uint64(block.Hash)) // 转为十六进制输出
		fmt.Printf("blockchain PreHash:%d\n", binary.BigEndian.Uint64(block.PreHash))
		fmt.Printf("blockchain target:%d\n", binary.BigEndian.Uint64(block.Target))
		fmt.Printf("blockchain nonce:%d\n", block.Nonce)
		fmt.Printf("Proof of Work validation::%s\n", block.ValidatePow())
		fmt.Println("============================================")
	}
}
