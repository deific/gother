package script

import (
	"bytes"
	"errors"
	"fmt"
)

// OP_1 through OP_16
const (
	OP_1 = 81 + iota
	OP_2 //82
	OP_3 //83
	OP_4 //..
	OP_5
	OP_6
	OP_7
	OP_8
	OP_9
	OP_10
	OP_11
	OP_12
	OP_13
	OP_14 //..
	OP_15 //95
	OP_16 //96
)

// OP codes other than OP_1 through OP_16, used in P2SH Multisig transanctions.
const (
	OP_0             = 0
	OP_PUSHDATA1     = 76
	OP_PUSHDATA2     = 77
	OP_DUP           = 118
	OP_EQUAL         = 135
	OP_EQUALVERIFY   = 136
	OP_HASH160       = 169
	OP_CHECKSIG      = 172
	OP_CHECKMULTISIG = 174
)

// NewMofNReedeemScript 创建 m/n的多签锁定脚本
func NewMofNReedeemScript(m int, n int, publicKeys [][]byte) ([]byte, error) {
	//Check we have valid numbers for M and N
	if n < 1 || n > 7 {
		return nil, errors.New("n must be between 1 and 7 (inclusive) for valid, standard P2SH multisig transaction as per Bitcoin protocol")
	}
	if m < 1 || m > n {
		return nil, errors.New("m must be between 1 and N (inclusive)")
	}

	//Check we have N public keys as necessary.
	if len(publicKeys) != n {
		return nil, errors.New(fmt.Sprintf("Need exactly %d public keys to create P2SH address for %d-of-%d multisig transaction. Only %d keys provided.", n, m, n, len(publicKeys)))
	}

	//Get OP Code for m and n.
	//81 is OP_1, 82 is OP_2 etc.
	//80 is not a valid OP_Code, so we floor at 81
	mOPCode := OP_1 + (m - 1)
	nOPCode := OP_1 + (n - 1)
	// //Multisig redeemScript format:
	//<OP_m> <A pubkey> <B pubkey> <C pubkey>... <OP_n> OP_CHECKMULTISIG
	var redeemScript bytes.Buffer
	redeemScript.WriteByte(byte(mOPCode))
	for _, pubKey := range publicKeys {
		redeemScript.WriteByte(byte(len(pubKey)))
		redeemScript.Write(pubKey)
	}
	redeemScript.WriteByte(byte(nOPCode))
	redeemScript.WriteByte(byte(OP_CHECKMULTISIG))
	return redeemScript.Bytes(), nil
}

// NewP2SHScriptPubKey creates a scriptPubKey for a P2SH transaction given the redeemScript hash
func NewP2SHScriptPubKey(redeemScriptHash []byte) ([]byte, error) {
	//P2SH scriptSig format:
	//<OP_HASH160> <Hash160(redeemScript)> <OP_EQUAL>
	var scriptPubKey bytes.Buffer
	scriptPubKey.WriteByte(byte(OP_HASH160))
	scriptPubKey.WriteByte(byte(len(redeemScriptHash)))
	scriptPubKey.Write(redeemScriptHash)
	scriptPubKey.WriteByte(byte(OP_EQUAL))
	return scriptPubKey.Bytes(), nil
}
