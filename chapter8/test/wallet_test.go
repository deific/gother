package test

import (
	"encoding/hex"
	"fmt"
	"gother/chapter8/internal/wallet"
	"strings"
	"testing"
)

func TestAddress(t *testing.T) {
	wlt := wallet.NewWallet()
	fmt.Printf("address : %s \n", wlt.P2PKHAddress())
	wlt2 := wallet.NewWallet()
	fmt.Printf("address : %s \n", wlt2.P2SHAddress())
	testM := 2
	testN := 3
	testPublicKeyStrs := "04a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd,046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187,0411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e83"
	//testAddress := "347N1Thc213QqfYCz3PZkjoJpNv5b14kBd"
	testPublicKeys := strings.Split(testPublicKeyStrs, ",")
	//testRedeemScriptHex := "524104a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd41046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187410411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e8353ae"
	//testRedeemScriptHex := "524104a882d414e478039cd5b52a92ffb13dd5e6bd4515497439dffd691a0f12af9575fa349b5694ed3155b136f09e63975a1700c9f4d4df849323dac06cf3bd6458cd41046ce31db9bdd543e72fe3039a1f1c047dab87037c36a669ff90e28da1848f640de68c2fe913d363a51154a0c62d7adea1b822d05035077418267b1a1379790187410411ffd36c70776538d079fbae117dc38effafb33304af83ce4894589747aee1ef992f63280567f52f5ba870678b4ab4ff6c8ea600bd217870a8b4f1f09f3a8e8353ae"

	mSignAddr, sh := wlt2.MSignAddress(testM, testN, testPublicKeys)
	fmt.Printf("address : %s %s\n", mSignAddr, hex.EncodeToString(sh))
}
