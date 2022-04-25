package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"gother/internal/constant"
	"gother/internal/transaction"
	"time"
)

type Block struct {
	Timestamp    int64 `json:"timestamp"`
	Height       int64
	Hash         []byte                     `json:"hash"`
	PreHash      []byte                     `json:"preHash"`
	Target       []byte                     // 当前区块目标难度值
	Nonce        int64                      // 随机数,用于证明该区块是合法区块,该随机数可以获得当前区块的hash值满足目标值
	Transactions []*transaction.Transaction `json:"transactions"`
}

// GetTransactionSummary 获取区块中的交易信息的序列化信息
// Transaction的ID是一个交易的hash,所有交易的Hash则可以基于每个交易的Hash构建
func (b *Block) GetTransactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}

	summary := bytes.Join(txIDs, []byte{})
	return summary
}

// SetHash 设置区块的hash
func (b *Block) SetHash() {
	data := b.GetDataWithNonce(b.Nonce)
	// 加密后返回是固定长度的数组
	hash := sha256.Sum256(data)
	// 通过切片操作转换为切片赋值给Block
	b.Hash = hash[:]
}

func NewBlock(height int64, preHash []byte, transactions []*transaction.Transaction) *Block {
	block := &Block{
		Height:       height,
		Timestamp:    time.Now().Unix(),
		PreHash:      preHash,
		Transactions: transactions,
	}
	block.Target = block.GetTarget(constant.Difficulty)
	block.Nonce = block.FindNonce()
	block.SetHash()
	return block
}

// GenesisBlock 创世区块，每个区块都有上一个区块的hash,对于第一个创世区块,上一个区块hash为空
func GenesisBlock() *Block {
	var preHash = make([]byte, 8)
	binary.BigEndian.PutUint64(preHash, 0)
	return NewBlock(0, preHash, []*transaction.Transaction{transaction.BaseTransaction([]byte(constant.BaseAddress))})
}
