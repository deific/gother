package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"gother/chapter8/internal/constant"
	transaction2 "gother/chapter8/internal/transaction"
	"gother/chapter8/internal/utils"
	"runtime"
)

type Blockchain struct {
	//Blocks   []*Block
	LastHash []byte
	network  string
	UtxoSet  *UTXOSet
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func CreateBlockChain(address []byte) *Blockchain {
	var lastHash []byte

	if utils.FileExists(constant.GetNetworkFile(constant.BCFile)) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constant.GetNetworkPath(constant.BCPath))
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

	blockchain := Blockchain{lastHash, constant.Network, nil, db}
	utxoSet := InitUTXOSet(&blockchain)
	blockchain.UtxoSet = utxoSet

	return &blockchain
}

// LoadBlockChain 加载区块链
func LoadBlockChain() *Blockchain {
	if !utils.FileExists(constant.GetNetworkFile(constant.BCFile)) {
		fmt.Println("No blockchain found,please create one first")
		return nil
	}

	var lastHash []byte

	opts := badger.DefaultOptions(constant.GetNetworkPath(constant.BCPath))
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

	chain := Blockchain{lastHash, constant.Network, nil, db}
	chain.UtxoSet = InitUTXOSet(&chain)
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
		return
		//runtime.Goexit()
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

// CreateTransaction 创建交易
func (bc *Blockchain) CreateTransaction(fromPubKey, toHashPubKey []byte, amount int, privKey ecdsa.PrivateKey) (*transaction2.Transaction, bool) {
	var inputs []transaction2.TxInput
	var outputs []transaction2.TxOutput

	// 根据输入地址找出该地址的未花费输出
	balance, validOutputs := bc.UtxoSet.FindSpendableOutputs(fromPubKey, amount)
	//balance, validOutputs := bc.FindSpendableOutputs(fromPubKey, amount)
	if balance < amount {
		fmt.Println("Not enough coins!")
		return &transaction2.Transaction{}, false
	}

	// 根据未花费输出，构建新交易的输入信息
	for _, item := range validOutputs {
		txID, err := hex.DecodeString(item.Txid)
		utils.Handle(err)

		input := transaction2.TxInput{TxID: txID, OutIdx: item.OutIdx, PubKey: fromPubKey}
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

func (bc *Blockchain) GetBalance(address []byte) int {
	balance := 0
	utxos := bc.UtxoSet.GetUtxos(string(address))
	for _, item := range utxos.UxtoItems {
		balance += item.Value
	}
	return balance
}
