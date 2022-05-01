package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob"
	"gother/chapter6/internal/constant"
	"gother/chapter6/internal/utils"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte
	// 将交易信息序列化后，hash计算得出交易hash
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

// IsBase 判断一个交易是否是创始交易
// 根据创始交易的特殊性，创始交易的输入没有之前的交易，因此创始交易中将输入设置为空地址，index=-1,根据特殊性判断是否创世交易
func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}

// BaseTransaction 创世区块的交易，同时也是会记录因为打包成功，产生奖励给矿工的交易
func BaseTransaction(toAddress []byte) *Transaction {
	txIn := TxInput{[]byte{}, -1, []byte{}, nil}
	txOut := TxOutput{constant.InitCoin, toAddress}
	tx := Transaction{[]byte("This is the Base Transaction!"), []TxInput{txIn}, []TxOutput{txOut}}
	return &tx
}

func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutIdx, nil, nil})
	}

	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) PlainHash(index int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[index].PubKey = prevPubKey
	return txCopy.TxHash()
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}
	for idx, input := range tx.Inputs {
		// 分别对每一个输入进行签名
		plainHash := tx.PlainHash(idx, input.PubKey)
		signature := utils.Sign(plainHash, privKey)
		tx.Inputs[idx].Sig = signature
	}
}

func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainHash := tx.PlainHash(idx, input.PubKey)
		if !utils.Verify(plainHash, input.PubKey, input.Sig) {
			return false
		}
	}
	return true
}
