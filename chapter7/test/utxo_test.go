package test

import (
	"gother/chapter7/internal/blockchain"
	"testing"
)

func TestUtxo(t *testing.T) {
	utxoSet := blockchain.UTXOSet{}
	utxoSet.Reindex()

	utxoSet.GetUtxos("123oiqYM6E2ryicTiHi8FzKamME5Fe74Z4")
}
