package blockchain

import (
	"crypto/sha256"
	"gother/internal/constant"
	"time"
)

type Block struct {
	Timestamp int64 `json:"timestamp"`
	Height    int64
	Hash      []byte `json:"hash"`
	PreHash   []byte `json:"preHash"`
	Target    []byte // 当前区块目标难度值
	Nonce     int64  // 随机数,用于证明该区块是合法区块,该随机数可以获得当前区块的hash值满足目标值
	Data      []byte `json:"data"`
}

func (b *Block) SetHash() {
	data := b.GetDataWithNonce(b.Nonce)
	// 加密后返回是固定长度的数组
	hash := sha256.Sum256(data)
	// 通过切片操作转换为切片赋值给Block
	b.Hash = hash[:]
}

func NewBlock(height int64, preHash, data []byte) *Block {
	block := &Block{
		Height:    height,
		Timestamp: time.Now().Unix(),
		PreHash:   preHash,
		Data:      data,
	}
	block.Target = block.GetTarget(constant.Difficulty)
	block.Nonce = block.FindNonce()
	block.SetHash()
	return block
}

// GenesisBlock 创世区块，每个区块都有上一个区块的hash,对于第一个创世区块,上一个区块hash为空
func GenesisBlock() *Block {
	genesisWords := "Hello, gother!"
	return NewBlock(0, []byte{}, []byte(genesisWords))
}
