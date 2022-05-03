package blockchain

import (
	"bytes"
	"encoding/hex"
	"fmt"
	transaction2 "gother/chapter8/internal/transaction"
	"gother/chapter8/internal/utils"
	"log"
)

// RunMine 挖矿
func (bc *Blockchain) RunMine() *Block {
	transactionPool := CreateTransactionPool()
	if !bc.VerifyTransactions(transactionPool.PubTx) {
		log.Println("falls in transaction verification")
		err := ClearTransactionPool()
		utils.Handle(err)
		return nil
	}

	candidateBlock := NewBlock(1, bc.LastHash, transactionPool.PubTx)
	if candidateBlock.ValidatePow() {
		bc.AddBlock(candidateBlock)
		err := ClearTransactionPool()
		utils.Handle(err)
		return candidateBlock
	} else {
		fmt.Println("Block has invalid nonce.")
		return nil
	}
}

// VerifyTransactions 验证交易
func (bc *Blockchain) VerifyTransactions(txs []*transaction2.Transaction) bool {
	if len(txs) == 0 {
		return false
	}

	spentOutputs := make(map[string]int)
	for _, tx := range txs {
		pubKey := tx.Inputs[0].PubKey
		unspentOutputs := bc.FindUnspentTx(pubKey)
		inputAmount := 0
		outputAmount := 0

		for _, input := range tx.Inputs {
			// 如果交易输入已经是花费过，验证不通过
			if outIdx, ok := spentOutputs[hex.EncodeToString(input.TxID)]; ok && outIdx == input.OutIdx {
				return false
			}
			ok, amount := isInputRight(unspentOutputs, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TxID)] = input.OutIdx
		}

		for _, output := range tx.Outputs {
			outputAmount += output.Value
		}

		if inputAmount != outputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}
	}
	return true
}

func isInputRight(txs []transaction2.Transaction, in transaction2.TxInput) (bool, int) {
	for _, tx := range txs {
		if bytes.Equal(tx.ID, in.TxID) {
			return true, tx.Outputs[in.OutIdx].Value
		}
	}
	return false, 0
}
