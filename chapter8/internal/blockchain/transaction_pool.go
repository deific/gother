package blockchain

import (
	"bytes"
	"encoding/gob"
	"gother/chapter8/internal/constant"
	"gother/chapter8/internal/transaction"
	"gother/chapter8/internal/utils"
	"io/ioutil"
	"os"
)

type TransactionPool struct {
	PubTx []*transaction.Transaction
}

func (tp *TransactionPool) AddTransaction(tx *transaction.Transaction) {
	tp.PubTx = append(tp.PubTx, tx)
}

func (tp *TransactionPool) SaveFile() {
	var content bytes.Buffer
	encoder := gob.NewEncoder(&content)

	err := encoder.Encode(tp)
	utils.Handle(err)

	err = ioutil.WriteFile(constant.GetNetworkFile(constant.TransactionPoolFile), content.Bytes(), 0644)
	utils.Handle(err)
}

func (tp *TransactionPool) LoadFile() error {
	if !utils.FileExists(constant.GetNetworkFile(constant.TransactionPoolFile)) {
		return nil
	}

	var transactionPool TransactionPool

	fileContent, err := ioutil.ReadFile(constant.GetNetworkFile(constant.TransactionPoolFile))
	utils.Handle(err)

	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&transactionPool)

	if err != nil {
		return err
	}
	tp.PubTx = transactionPool.PubTx
	return nil
}

func CreateTransactionPool() *TransactionPool {
	transactionPool := TransactionPool{}
	err := transactionPool.LoadFile()
	utils.Handle(err)
	return &transactionPool
}

func ClearTransactionPool() error {
	err := os.Remove(constant.GetNetworkFile(constant.TransactionPoolFile))
	return err
}
