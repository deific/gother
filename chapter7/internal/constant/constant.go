package constant

const (
	BaseAddress = "Hi,gother"
	// Difficulty 难度常数
	Difficulty = 12
	ROOT       = "/Users/steven/workspace/open-git/gother/chapter7/"
	// 创始交易值
	InitCoin            = 1000
	GenesisPreHash      = "gother!"
	TransactionPoolFile = ROOT + "./tmp/transaction_pool.data"
	BCPath              = ROOT + "./tmp/blocks"
	BCFile              = ROOT + "./tmp/blocks/MANIFEST"
	UTXOPATH            = ROOT + "./tmp/utxo/"
	UTXOFile            = ROOT + "./tmp/utxo/MANIFEST"
	CheckSumLength      = 4
	NetworkVersion      = byte(0x00)
	Wallets             = ROOT + "./tmp/wallets/"
	WalletsRefList      = ROOT + "./tmp/ref_list/"
)
