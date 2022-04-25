package blockchain

import (
	"encoding/hex"
	"fmt"
	"gother/internal/transaction"
	"gother/internal/utils"
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

// Mine 模拟挖矿
func (bc *Blockchain) Mine(txs []*transaction.Transaction) {
	bc.AddBlock(txs)
}

// CreateTransaction 创建交易
func (bc *Blockchain) CreateTransaction(fromAddress []byte, toAddress []byte, amount int) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput
	// 根据输入地址找出该地址的未花费输出
	balance, vaildOutputs := bc.FindSpendableOutputs(fromAddress, amount)
	if balance < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}

	// 根据未花费输出，构建新交易的输入信息
	for txId, outIdx := range vaildOutputs {
		txID, err := hex.DecodeString(txId)
		utils.Handle(err)

		input := transaction.TxInput{TxID: txID, OutIdx: outIdx, FromAddress: fromAddress}
		inputs = append(inputs, input)
	}

	// 构建输出信息
	outputs = append(outputs, transaction.TxOutput{ToAddress: toAddress, Value: amount})
	// 如果可花费金额超出了amount,需要找零
	if balance > amount {
		outputs = append(outputs, transaction.TxOutput{ToAddress: fromAddress, Value: balance - amount})
	}

	tx := transaction.Transaction{Inputs: inputs, Outputs: outputs}
	tx.SetID()
	return &tx, true
}

// FindSpendableOutputs 查找目标地址的可花费amount的未花费输出
func (bc *Blockchain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unSpentOuts := make(map[string]int)
	unSpentTxs := bc.FindUnspentTx(address)
	balance := 0

Work:
	for _, tx := range unSpentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.IsToAddress(address) {
				balance += out.Value
				unSpentOuts[txID] = outIdx
				if balance >= amount {
					break Work
				}
			}
		}
	}

	return balance, unSpentOuts
}

// FindUTXOs 查找指定地址的所有未花费交易输出
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
				continue Work // one transaction can only have one output referred to address
			}
		}
	}

	return balance, unspentOuts
}

// FindUnspentTx 查找所有包含address未花费交易
func (bc *Blockchain) FindUnspentTx(address []byte) []transaction.Transaction {
	var unSpentTx []transaction.Transaction
	spentTxs := make(map[string][]int)

	// 循环查找整个区块链，从后向前查找
	for idx := len(bc.Blocks) - 1; idx >= 0; idx-- {
		block := bc.Blocks[idx]
		// 查找每个区块上的交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			// 循环查找交易中的输出有没有作为其他交易的输入被花费掉
		IterOutputs:
			for outIndex, out := range tx.Outputs {
				// 如果当前交易是被花费过的交易，则判断该输出是否已花费交易中使用
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
				// 每个交易的输入就是之前交易的输出，只要某个交易的输出被另一个交易的输入引用了，则认为已被花费
				for _, in := range tx.Inputs {
					// 如果某个交易的输入地址是被查找地址，则说明其对应
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
