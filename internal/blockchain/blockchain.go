package blockchain

type Blockchain struct {
	Blocks []*Block
}

// AddBlock 向区块链上追加区块
func (c *Blockchain) AddBlock(data string) {
	height := len(c.Blocks)
	newBlock := NewBlock(int64(height), c.Blocks[height-1].Hash, []byte(data))
	c.Blocks = append(c.Blocks, newBlock)
}

// CreateBlockChain 创建blockchain
func CreateBlockChain() *Blockchain {
	blockChain := &Blockchain{}
	blockChain.Blocks = append(blockChain.Blocks, GenesisBlock())
	return blockChain
}
