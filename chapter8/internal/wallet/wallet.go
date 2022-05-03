package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"gother/chapter8/internal/constant"
	"gother/chapter8/internal/script"
	"gother/chapter8/internal/utils"
	"io/ioutil"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	prvKey, pubKey := NewPairKey()
	return &Wallet{PrivateKey: prvKey, PublicKey: pubKey}
}

func (w *Wallet) P2PKHAddress() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

func (w *Wallet) P2SHAddress() []byte {
	reedeemScript, err := script.NewMofNReedeemScript(1, 1, [][]byte{w.PublicKey})
	utils.Handle(err)
	scriptHash := utils.Hash160(reedeemScript)
	return utils.ScriptHash2Address(scriptHash)
}

func (w *Wallet) MSignAddress(m int, n int, publicKeys []string) ([]byte, []byte) {
	publicKeyArray := make([][]byte, len(publicKeys))
	for i, pubKey := range publicKeys {
		publicKeyArray[i], _ = hex.DecodeString(pubKey)
	}

	reedeemScript, err := script.NewMofNReedeemScript(m, n, publicKeyArray)
	utils.Handle(err)
	scriptHash := utils.Hash160(reedeemScript)
	return utils.ScriptHash2Address(scriptHash), reedeemScript
}

func (w *Wallet) SaveWallet() {
	filename := constant.GetNetworkPath(constant.Wallets) + string(w.P2PKHAddress()) + ".wlt"
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(w)
	utils.Handle(err)

	err = ioutil.WriteFile(filename, content.Bytes(), 0644)
	utils.Handle(err)
}
func LoadWallet(address string) *Wallet {
	filename := constant.GetNetworkPath(constant.Wallets) + address + ".wlt"
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("no wallet with such address"))
	}
	var w Wallet
	gob.Register(elliptic.P256())
	fileContent, err := ioutil.ReadFile(filename)
	utils.Handle(err)

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&w)
	utils.Handle(err)
	return &w
}

// NewPairKey 创建椭圆曲线秘钥对的生成函数
func NewPairKey() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	utils.Handle(err)
	// 公钥是椭圆曲线的一个点，将2个点坐标拼接保存起来
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}
