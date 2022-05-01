package transaction

import (
	"bytes"
	"gother/chapter7/internal/utils"
)

type TxInput struct {
	TxID   []byte
	OutIdx int
	PubKey []byte
	Sig    []byte
}

type TxOutput struct {
	Value      int
	HashPubKey []byte // 公钥Hash
}

// IsFromAddress 判断一个地址是否是输入地址
func (in *TxInput) IsFromAddress(pubKey []byte) bool {
	return bytes.Equal(in.PubKey, pubKey)
}

// IsToAddress 判断一个地址是否是输出地址
func (out *TxOutput) IsToAddress(pubKey []byte) bool {
	return bytes.Equal(out.HashPubKey, utils.PublicKeyHash(pubKey))
}
