package constant

var Network string = "test"

const (
	BaseAddress = "Hi,gother"
	// Difficulty 难度常数
	Difficulty = 12
	BaseDir    = "./tmp/"
	// 创始交易值
	InitCoin            = 1000
	GenesisPreHash      = "gother!"
	TransactionPoolFile = "/transaction_pool.data"
	BCPath              = "/blocks"
	BCFile              = "/blocks/MANIFEST"
	UTXOPATH            = "/utxo/"
	UTXOFile            = "/utxo/MANIFEST"
	CheckSumLength      = 4
	P2PKH_PRIFIX        = byte(0x00)
	P2SH_PRIFIX         = byte(0x05)
	Wallets             = "/wallets/"
	WalletsRefList      = "/ref_list/"
)

func GetNetworkFile(fileName string) string {
	return BaseDir + Network + fileName
}

func GetNetworkPath(path string) string {
	return BaseDir + Network + path
}
