package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"gother/chapter6/internal/constant"
	transaction2 "gother/chapter6/internal/transaction"
	"gother/chapter6/internal/utils"
	"runtime"
)

type Blockchain struct {
	//Blocks   []*Block
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain(address []byte) *Blockchain {
	var lastHash []byte

	if utils.FileExists(constant.BCFile) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constant.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("Genesis Created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)

		err = txn.Set([]byte("lh"), genesis.Hash)
		utils.Handle(err)

		err = txn.Set([]byte("ogprevhash"), genesis.PreHash)
		utils.Handle(err)

		lastHash = genesis.Hash
		return err
	})
	utils.Handle(err)

	blockchain := Blockchain{lastHash, db}
	return &blockchain
}

// LoadBlockChain 加载区块链
func LoadBlockChain() *Blockchain {
	if !utils.FileExists(constant.BCFile) {
		fmt.Println("No blockchain found,please create one first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constant.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return nil
	})
	utils.Handle(err)

	chain := Blockchain{lastHash, db}
	return &chain
}

func (bc *Blockchain) AddBlock(newBlock *Block) {
	var lastHash []byte

	// 查询数据库中最后一个区块hash
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	// 判断引用关系
	if !bytes.Equal(newBlock.PreHash, lastHash) {
		fmt.Println("This block is out of age")
		runtime.Goexit()
	}

	// 保存新区块
	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)
		bc.LastHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	it := BlockchainIterator{CurrentHash: bc.LastHash, Database: bc.Database}
	return &it
}

func (it *BlockchainIterator) Next() *Block {
	var block *Block

	err := it.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(it.CurrentHash)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerialize(val)
			return nil
		})
		utils.Handle(err)
		return nil
	})
	utils.Handle(err)

	it.CurrentHash = block.PreHash
	return block
}

// HasNext 是否还有下一个区块
func (it *BlockchainIterator) HasNext() bool {
	return !bytes.Equal(it.CurrentHash, []byte(constant.GenesisPreHash))
}

// Mine 模拟挖矿
func (bc *Blockchain) Mine(txs []*transaction2.Transaction) {
}

// CreateTransaction 创建交易
func (bc *Blockchain) CreateTransaction(fromPubKey, toHashPubKey []byte, amount int, privKey ecdsa.PrivateKey) (*transaction2.Transaction, bool) {
	var inputs []transaction2.TxInput
	var outputs []transaction2.TxOutput

	// 根据输入地址找出该地址的未花费输出
	balance, validOutputs := bc.FindSpendableOutputs(fromPubKey, amount)
	if balance < amount {
		fmt.Println("Not enough coins!")
		return &transaction2.Transaction{}, false
	}

	// 根据未花费输出，构建新交易的输入信息
	for txId, outIdx := range validOutputs {
		txID, err := hex.DecodeString(txId)
		utils.Handle(err)

		input := transaction2.TxInput{TxID: txID, OutIdx: outIdx, PubKey: fromPubKey}
		inputs = append(inputs, input)
	}

	// 构建输出信息
	outputs = append(outputs, transaction2.TxOutput{HashPubKey: toHashPubKey, Value: amount})
	// 如果可花费金额超出了amount,需要找零
	if balance > amount {
		outputs = append(outputs, transaction2.TxOutput{HashPubKey: utils.PublicKeyHash(fromPubKey), Value: balance - amount})
	}

	tx := transaction2.Transaction{Inputs: inputs, Outputs: outputs}
	tx.SetID()
	tx.Sign(privKey)
	return &tx, true
}

// FindSpendableOutputs 查找目标地址的可花费amount的未花费输出
func (bc *Blockchain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unSpentOuts := make(map[string]int)
	unSpentTxs := bc.FindUnspentTx(address)
	balance := 0

Work:
	for _, tx := range unSpentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.IsToAddress(address) {
				balance += out.Value
				unSpentOuts[txID] = outIdx
				if balance >= amount {
					break Work
				}
			}
		}
	}

	return balance, unSpentOuts
}

// FindUTXOs 查找指定地址的所有未花费交易输出
func (bc *Blockchain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTx(address)
	balance := 0

Work:
	for _, tx := range unspentTxs {
		txId := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.IsToAddress(address) {
				balance += out.Value
				unspentOuts[txId] = outIdx
				continue Work // one transaction can only have one output referred to address
			}
		}
	}

	return balance, unspentOuts
}

// FindUnspentTx 查找所有包含address未花费交易
func (bc *Blockchain) FindUnspentTx(address []byte) []transaction2.Transaction {
	var unSpentTx []transaction2.Transaction
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
				if out.IsToAddress(address) {
					unSpentTx = append(unSpentTx, *tx)
				}
			}
			// 构建已花费map缓存
			if !tx.IsBase() {
				// 每个交易的输入就是之前交易的输出，只要某个交易的输出被另一个交易的输入引用了，则认为已被花费
				for _, in := range tx.Inputs {
					// 如果某个交易的输入地址是被查找地址，则说明其对应
					if in.IsFromAddress(address) {
						inTxId := hex.EncodeToString(in.TxID)
						spentTxs[inTxId] = append(spentTxs[inTxId], in.OutIdx)
					}
				}
			}
		}

		if !iter.HasNext() {
			break
		}
	}

	return unSpentTx
}
