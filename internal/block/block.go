package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

type Block struct {
	Timestamp int64  `json:"timestamp"`
	Hash      []byte `json:"hash"`
	PreHash   []byte `json:"preHash"`
	Data      []byte `json:"data"`
}

type BlockChain struct {
	Blocks []*Block
}

func (b *Block) SetHash() {
	block := bytes.Join([][]byte{ToHexForInt(b.Timestamp), b.PreHash, b.Data}, []byte{})
	// 加密后返回是固定长度的数组
	hash := sha256.Sum256(block)
	// 通过切片操作转换为切片赋值给Block
	b.Hash = hash[:]
}

func NewBlock(preHash, data []byte) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		PreHash:   preHash,
		Data:      data,
	}
	block.SetHash()
	return block
}

// GenesisBlock 创世区块，每个区块都有上一个区块的hash,对于第一个创世区块,上一个区块hash为空
func GenesisBlock() *Block {
	genesisWords := "Hello, gother!"
	return NewBlock([]byte{}, []byte(genesisWords))
}

// AddBlock 向区块链上追加区块
func (c *BlockChain) AddBlock(data string) {
	newBlock := NewBlock(c.Blocks[len(c.Blocks)-1].Hash, []byte(data))
	c.Blocks = append(c.Blocks, newBlock)
}

// CreateBlockChain 创建blockchain
func CreateBlockChain() *BlockChain {
	blockChain := &BlockChain{}
	blockChain.Blocks = append(blockChain.Blocks, GenesisBlock())
	return blockChain
}

func ToHexForInt(input int64) []byte {
	// 创建一个字节buffer,未初始化零值
	buff := new(bytes.Buffer)
	// 按大端顺序，写入数据
	err := binary.Write(buff, binary.BigEndian, input)
	if err != nil {
		log.Panicln(err)
	}
	// 返回转换后的字节数组
	return buff.Bytes()
}
