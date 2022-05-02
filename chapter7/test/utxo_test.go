package test

import (
	"gother/chapter7/internal/blockchain"
	"testing"
)

func TestUtxo(t *testing.T) {
	chain := blockchain.LoadBlockChain()
	utxoSet := blockchain.InitUTXOSet(chain)

	//utxoSet.PrintUtxo("")
	aa := utxoSet.FindAllUTXOs()
	aa.PrintUtxo("")
	//utxo := utxoSet.GetUtxos("123oiqYM6E2ryicTiHi8FzKamME5Fe74Z4")
	//
	//refList := wallet.LoadRefList()
	//refName, _ := refList.FindRefName(utxo.Address)
	//for _, item := range utxo.UxtoItems {
	//	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	//	fmt.Printf("address：%s rename:%s \n", utxo.Address, refName)
	//	fmt.Printf("txId:%s \n", item.Txid)
	//	fmt.Printf("outIndex:%d\n", item.OutIdx)
	//	fmt.Printf("amount:%d\n", item.Value)
	//}
}
