package blockchain

import (
	"encoding/hex"
	"gother/internal/transaction"
)

type Blockchain struct {
	Blocks []*Block
}

// AddBlock 向区块链上追加区块
func (c *Blockchain) AddBlock(transactions []*transaction.Transaction) {
	height := len(c.Blocks)
	newBlock := NewBlock(int64(height), c.Blocks[height-1].Hash, transactions)
	c.Blocks = append(c.Blocks, newBlock)
}

// CreateBlockChain 创建blockchain
func CreateBlockChain() *Blockchain {
	blockChain := &Blockchain{}
	blockChain.Blocks = append(blockChain.Blocks, GenesisBlock())
	return blockChain
}

// CreateTransaction 创建交易
func CreateTransaction(Value int64, fromAddress []byte, toAddress []byte) *transaction.Transaction {
	// 根据输入地址找出该地址的未花费输出

	// 根据Value选择合适的未花费输出

	return nil
}
func (bc *Blockchain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTx(address)
	balance := 0

Work:
	for _, tx := range unspentTxs {
		txId := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.IsToAddress(address) {
				balance += out.Value
				unspentOuts[txId] = outIdx
				continue Work
			}
		}
	}

	return balance, unspentOuts
}

func (bc *Blockchain) FindUnspentTx(address []byte) []transaction.Transaction {
	var unSpentTx []transaction.Transaction
	spentTxs := make(map[string][]int)

	// 循环查找整个区块链，从后向前查找
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]
		// 选好查找该区块上的交易输出
		for _, tx := range block.Transactions {
			if tx.IsBase() {
				continue
			}
			txID := hex.EncodeToString(tx.ID)
		IterOutputs:
			for outIndex, out := range tx.Outputs {
				// 如果已花费，则继续查找该花费交易中是否存在未花费的out
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIndex {
							continue IterOutputs
						}
					}
				}
				if out.IsToAddress(address) {
					unSpentTx = append(unSpentTx, *tx)
				}
			}
			// 构建已花费map缓存
			if !tx.IsBase() {
				for _, in := range tx.Inputs {
					if in.IsFromAddress(address) {
						inTxId := hex.EncodeToString(in.TxID)
						spentTxs[inTxId] = append(spentTxs[inTxId], in.OutIdx)
					}
				}
			}
		}
	}

	return unSpentTx
}
