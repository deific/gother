package blockchain

import (
	"bytes"
	"crypto/sha256"
	"gother/chapter8/internal/utils"
	"math"
	"math/big"
)

func (b *Block) GetTarget(diff int) []byte {
	// 目标难度值为1，256位字节
	target := big.NewInt(1)
	// 左移 256 - 难度值，通过调节难度值改变最终目标，左移位数越多目标值空间越大，越容易找到
	// 如果左移256位，即256个字节全为1，最大值，所有找到的hash几乎都满足
	// 如果左移256 - diff，即前diff个字节值为0，其余都为1，找到的hash数值小于该值时就符合要求
	// 通过调整diff值的大小，就能够调整满足要求的hash的空间范围，diff值越大，空间范围越小则难度越大，diff值越小空间范围大难度就小
	target.Lsh(target, uint(256-diff))
	return target.Bytes()
}

func (b *Block) GetDataWithNonce(nonce int64) []byte {
	data := bytes.Join([][]byte{
		utils.ToHexForInt(b.Timestamp),
		b.PreHash,
		utils.ToHexForInt(nonce),
		b.Target,
		b.GetTransactionSummary(),
		b.MTree.RootNode.Data,
	}, []byte{})
	return data
}

func (b *Block) FindNonce() int64 {
	var intHash, intTarget big.Int
	var hash [32]byte
	var nonce int64

	nonce = 0
	intTarget.SetBytes(b.Target)

	// 从0开始寻找随机数
	for nonce < math.MaxInt64 {
		data := b.GetDataWithNonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])

		// 如果区块当前的hash小于目标难度hash值，则认为找到了合适的nonce
		if intHash.Cmp(&intTarget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce
}

// ValidatePow 验证区块的工作量证明是否有效
func (b *Block) ValidatePow() bool {
	var intHash, intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)

	data := b.GetDataWithNonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(&intTarget) == -1
}
