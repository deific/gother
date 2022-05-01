package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"gother/chapter6/internal/constant"
	"log"
	"math/big"
	"os"
)

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

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// FileExists 判断文件是否存在
func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

func PublicKeyHash(publicKey []byte) []byte {
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

func CheckSum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:constant.CheckSumLength]
}

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)
	return []byte(encode)
}
func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{constant.NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}
func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constant.CheckSumLength]
	return pubKeyHash
}

func Sign(msg []byte, privKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, msg)
	Handle(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func Verify(msg []byte, pubKey []byte, signature []byte) bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(signature)
	r.SetBytes(signature[:(sigLen / 2)])
	s.SetBytes(signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(pubKey)
	x.SetBytes(pubKey[:(keyLen / 2)])
	y.SetBytes(pubKey[(keyLen / 2):])

	rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
	return ecdsa.Verify(&rawPubKey, msg, &r, &s)
}
