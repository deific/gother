package transaction

import "bytes"

type TxInput struct {
	TxID        []byte
	OutIdx      int
	FromAddress []byte
}

type TxOutput struct {
	Value     int
	ToAddress []byte
}

// IsFromAddress 判断一个地址是否是输入地址
func (in *TxInput) IsFromAddress(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

// IsToAddress 判断一个地址是否是输出地址
func (out *TxOutput) IsToAddress(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
