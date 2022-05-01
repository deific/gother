package blockchain

import (
	"bytes"
	"encoding/gob"
	"gother/chapter7/internal/constant"
	"gother/chapter7/internal/transaction"
	"gother/chapter7/internal/utils"
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

	err = ioutil.WriteFile(constant.TransactionPoolFile, content.Bytes(), 0644)
	utils.Handle(err)
}

func (tp *TransactionPool) LoadFile() error {
	if !utils.FileExists(constant.TransactionPoolFile) {
		return nil
	}

	var transactionPool TransactionPool

	fileContent, err := ioutil.ReadFile(constant.TransactionPoolFile)
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
	err := os.Remove(constant.TransactionPoolFile)
	return err
}
