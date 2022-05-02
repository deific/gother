package blockchain

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"gother/chapter7/internal/constant"
	"gother/chapter7/internal/utils"
)

type UTXOSet struct {
	chain *Blockchain
	db    *badger.DB
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

func InitUTXOSet(chain *Blockchain) *UTXOSet {
	utxoSet := &UTXOSet{chain: chain}

	needReindex := !utils.FileExists(constant.GetNetworkFile(constant.UTXOFile))

	opts := badger.DefaultOptions(constant.GetNetworkPath(constant.UTXOPATH))
	opts.Logger = nil
	db, err := badger.Open(opts)
	utils.Handle(err)
	utxoSet.db = db

	if needReindex {
		fmt.Println("utxo set not found,reindex utxo set now....")
		count := utxoSet.Reindex()
		fmt.Printf("reindex utxo set succeed, utxo address size:%d \n", count)
	}

	utxoSet.loadUtxoSet()
	return utxoSet
}

func (u UTXOSet) FindSpendableOutputs(pubkeyHash []byte, amount int) (int, []UTXOItem) {
	addr := string(utils.PubHash2Address(pubkeyHash))
	uxto := u.GetUtxos(addr)
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

func (u *UTXOSet) Reindex() int {
	uxtoSet := u.FindAllUTXOs()
	u.UTXOS = uxtoSet.UTXOS
	return u.save()
}

func (u *UTXOSet) save() int {
	err := u.db.DropPrefix([]byte("utxo:"))
	utils.Handle(err)
	err = u.db.Update(func(txn *badger.Txn) error {
		for _, utxo := range u.UTXOS {
			var content bytes.Buffer
			encoder := gob.NewEncoder(&content)
			err := encoder.Encode(utxo)
			utils.Handle(err)
			err = txn.Set([]byte("utxo:"+utxo.Address), content.Bytes())
			utils.Handle(err)
		}
		return nil
	})
	utils.Handle(err)
	return len(u.UTXOS)
}

func (u *UTXOSet) loadUtxoSet() {
	err := u.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prifix := []byte("utxo:")

		for it.Seek(prifix); it.ValidForPrefix(prifix); it.Next() {
			var utxo UTXO
			item := it.Item()
			err := item.Value(func(val []byte) error {
				decoder := gob.NewDecoder(bytes.NewReader(val))
				err := decoder.Decode(&utxo)
				utils.Handle(err)
				return nil
			})
			utils.Handle(err)
			utxo.Address = string(item.Key())[5:]
			u.UTXOS = append(u.UTXOS, &utxo)
		}
		return nil
	})
	utils.Handle(err)
}

func (u *UTXOSet) Update(block *Block) {
	for _, tx := range block.Transactions {
		for _, in := range tx.Inputs {
			u.spentByTx(in.TxID, in.PubKey, in.OutIdx)
		}
		for idx, out := range tx.Outputs {
			u.unspentByTx(out.HashPubKey, tx.ID, idx, out.Value)
		}
	}
	u.save()
}

func (u *UTXOSet) spentByTx(txId []byte, pubKey []byte, outIdx int) {
	addr := utils.PubHash2Address(utils.PublicKeyHash(pubKey))
	var targetUtxo *UTXO
	for _, utxo := range u.UTXOS {
		if string(addr) == utxo.Address {
			targetUtxo = utxo
			break
		}
	}

	var leftUtxoItems []UTXOItem
	if targetUtxo != nil {
		for _, item := range targetUtxo.UxtoItems {
			if item.Txid == hex.EncodeToString(txId) && item.OutIdx == outIdx {
				continue
			}
			leftUtxoItems = append(leftUtxoItems, item)
		}
	} else {
		leftUtxoItems = targetUtxo.UxtoItems[:]
	}
	targetUtxo.UxtoItems = leftUtxoItems
}

func (u *UTXOSet) unspentByTx(hashPubKey, txId []byte, outIdx int, amount int) {
	addr := utils.PubHash2Address(hashPubKey)
	var targetUtxo *UTXO
	for _, utxo := range u.UTXOS {
		if string(addr) == utxo.Address {
			targetUtxo = utxo
			break
		}
	}

	if targetUtxo != nil {
		targetUtxo.UxtoItems = append(targetUtxo.UxtoItems, UTXOItem{Txid: hex.EncodeToString(txId), OutIdx: outIdx, Value: amount})
	} else {
		targetUtxo = &UTXO{Address: string(addr), UxtoItems: append([]UTXOItem{}, UTXOItem{Txid: hex.EncodeToString(txId), OutIdx: outIdx, Value: amount})}
		u.UTXOS = append(u.UTXOS, targetUtxo)
	}
}

func (u *UTXOSet) GetUtxos(address string) *UTXO {
	for _, utxo := range u.UTXOS {
		if utxo.Address == address {
			return utxo
		}
	}
	return nil
}

func (u *UTXOSet) FindAllUTXOs() *UTXOSet {
	unspentOuts := &UTXOSet{}
	spentTxs := make(map[string][]int)

	// 循环查找整个区块链，从后向前查找
	iter := u.chain.Iterator()

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
				utxo := unspentOuts.GetUtxos(addr)
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

func (u *UTXOSet) PrintUtxo(address string) {
	var utxos []*UTXO
	if address != "" {
		utxos = append(utxos, u.GetUtxos(address))
	} else {
		utxos = u.UTXOS
	}

	for _, utxo := range utxos {
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("address：%s utxoSize:%d \n", utxo.Address, len(utxo.UxtoItems))
		for _, item := range utxo.UxtoItems {
			fmt.Println("---------------------------------------------------")
			fmt.Printf("txId:%s \n", item.Txid)
			fmt.Printf("outIndex:%d\n", item.OutIdx)
			fmt.Printf("amount:%d\n", item.Value)
		}
	}
}
