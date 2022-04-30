package blockchain

import (
	"fmt"
	"gother/chapter6/internal/utils"
)

// RunMine 挖矿
func (bc *Blockchain) RunMine() {
	transactionPool := CreateTransactionPool()
	candidateBlock := NewBlock(1, bc.LastHash, transactionPool.PubTx)
	if candidateBlock.ValidatePow() {
		bc.AddBlock(candidateBlock)
		err := ClearTransactionPool()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
