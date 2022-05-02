package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"github.com/dgraph-io/badger"
	"gother/chapter7/internal/constant"
	"gother/chapter7/internal/utils"
)

type UTXOSet struct {
	UTXOS []*UTXO
}

type UTXO struct {
	Address   string
	UxtoItems []UTXOItem
}

type UTXOItem struct {
	Txid   string
	OutIdx int
	Value  int
}

func (u *UTXOSet) getUtxo(address string) *UTXO {
	for _, utxo := range u.UTXOS {
		if utxo.Address == address {
			return utxo
		}
	}
	return nil
}

func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, []UTXOItem) {
	addr := string(utils.PubHash2Address(pubkeyHash))
	uxto := u.getUtxo(addr)
	var unspentItems []UTXOItem

	accumulated := 0
	for _, item := range uxto.UxtoItems {
		accumulated += item.Value
		unspentItems = append(unspentItems, item)
		if accumulated >= amount {
			break
		}
	}
	return accumulated, unspentItems
}

func (u *UTXOSet) Reindex() {
	chain := LoadBlockChain()
	defer chain.Database.Close()

	uxtoSet := chain.FindAllUTXOs()

	opts := badger.DefaultOptions(constant.UTXOFile)
	opts.Logger = nil

	db, err := badger.Open(opts)
	defer db.Close()
	utils.Handle(err)

	db.Update(func(txn *badger.Txn) error {
		for _, utxo := range uxtoSet.UTXOS {
			var content bytes.Buffer
			encoder := gob.NewEncoder(&content)
			encoder.Encode(utxo.UxtoItems)
			err := txn.Set([]byte(utxo.Address), content.Bytes())
			utils.Handle(err)
		}
		return nil
	})
	utils.Handle(err)
}

func (u *UTXOSet) GetUtxos(address string) []UTXOItem {
	opts := badger.DefaultOptions(constant.UTXOFile)
	opts.Logger = nil

	db, err := badger.Open(opts)
	defer db.Close()
	utils.Handle(err)
	var utxoItems []UTXOItem
	db.View(func(txn *badger.Txn) error {
		val, err := txn.Get([]byte(address))
		utils.Handle(err)

		val.Value(func(val []byte) error {
			decoder := gob.NewDecoder(bytes.NewReader(val))
			err = decoder.Decode(&utxoItems)
			utils.Handle(err)
			return nil
		})
		return nil
	})
	utils.Handle(err)
	return utxoItems
}

func (bc *Blockchain) FindAllUTXOs() *UTXOSet {
	unspentOuts := &UTXOSet{}
	spentTxs := make(map[string][]int)

	// 循环查找整个区块链，从后向前查找
	iter := bc.Iterator()

	for {
		block := iter.Next()
		// 查找每个区块上的交易
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
			// 循环查找交易中的输出有没有作为其他交易的输入被花费掉
		IterOutputs:
			for outIndex, out := range tx.Outputs {
				// 如果当前交易是被花费过的交易，则判断该输出是否已花费交易中使用
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIndex {
							continue IterOutputs
						}
					}
				}

				addr := string(utils.PubHash2Address(out.HashPubKey)) //hex.EncodeToString(out.HashPubKey)
				utxo := unspentOuts.getUtxo(addr)
				uxtoItem := UTXOItem{Txid: txID, OutIdx: outIndex, Value: out.Value}
				if utxo == nil {
					utxo = &UTXO{Address: addr, UxtoItems: append([]UTXOItem{}, uxtoItem)}
					unspentOuts.UTXOS = append(unspentOuts.UTXOS, utxo)
				} else {
					utxo.UxtoItems = append(utxo.UxtoItems, uxtoItem)
				}
			}

			// 构建已花费map缓存
			if !tx.IsBase() {
				// 每个交易的输入就是之前交易的输出，只要某个交易的输出被另一个交易的输入引用了，则认为已被花费
				for _, in := range tx.Inputs {
					// 如果某个交易的输入地址是被查找地址，则说明其对应
					inTxId := hex.EncodeToString(in.TxID)
					spentTxs[inTxId] = append(spentTxs[inTxId], in.OutIdx)
				}
			}
		}

		if !iter.HasNext() {
			break
		}
	}

	return unspentOuts
}
