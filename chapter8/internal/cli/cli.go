package cli

import (
	"flag"
	"fmt"
	"gother/chapter8/internal/blockchain"
	"gother/chapter8/internal/constant"
	"gother/chapter8/internal/utils"
	"gother/chapter8/internal/wallet"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
}

func (cl *CommandLine) Run() {
	cl.checkArgs()

	cl.parseAndRunCmd("createwallet", map[string]string{"refname": "The refname of the wallet, and this is optimal"}, func(args map[string]*string) {
		cl.createWallet(*args["refname"])
	})
	cl.parseAndRunCmd("walletinfo", map[string]string{"refname": "The refname of the wallet", "address": "The address of the wallet"}, func(args map[string]*string) {
		if *args["refname"] != "" {
			cl.walletInfoByRefName(*args["refname"])
		} else {
			cl.walletInfo(*args["address"])
		}
	})
	cl.parseAndRunCmd("walletslist", map[string]string{}, func(args map[string]*string) {
		cl.walletsList()
	})

	cl.parseAndRunCmd("createblockchain", map[string]string{
		"refname": "The refname refer to the owner of blockchain",
		"address": "The address refer to the owner of blockchain"}, func(args map[string]*string) {
		if *args["refname"] != "" {
			cl.createByRefName(*args["refname"])
		} else {
			cl.create(*args["address"])
		}
	})

	cl.parseAndRunCmd("balance", map[string]string{
		"refname": "Who need to get balance amount",
		"address": "Who need to get balance amount"}, func(args map[string]*string) {
		if *args["refname"] != "" {
			cl.balanceByRefName(*args["refname"])
		} else {
			cl.Balance(*args["address"])
		}
	})
	cl.parseAndRunCmd("balance2", map[string]string{
		"refname": "Who need to get balance amount",
		"address": "Who need to get balance amount"}, func(args map[string]*string) {
		if *args["refname"] != "" {
			cl.balanceByRefName2(*args["refname"])
		} else {
			cl.balance2(*args["address"])
		}
	})

	cl.parseAndRunCmd("blockchaininfo", map[string]string{}, func(args map[string]*string) {
		cl.Info()
	})

	cl.parseAndRunCmd("send", map[string]string{"from": "Source address", "to": "Destination address", "amount": "Amount to send"}, func(args map[string]*string) {
		amount, err := strconv.Atoi(*args["amount"])
		utils.Handle(err)
		cl.send(*args["from"], *args["to"], amount)
	})
	cl.parseAndRunCmd("sendbyrefname", map[string]string{"from": "Source address", "to": "Destination address", "amount": "Amount to send"}, func(args map[string]*string) {
		amount, err := strconv.Atoi(*args["amount"])
		utils.Handle(err)
		cl.send(cl.getAddressByRefName(*args["from"]), cl.getAddressByRefName(*args["to"]), amount)
	})

	cl.parseAndRunCmd("mine", map[string]string{}, func(args map[string]*string) {
		cl.mine()
	})

	cl.parseAndRunCmd("getutxos", map[string]string{}, func(args map[string]*string) {
		cl.GetUtxos()
	})
}

func (cl *CommandLine) parseAndRunCmd(subCmdName string, args map[string]string, runCmd func(args map[string]*string)) {
	if subCmdName != os.Args[1] {
		return
	}

	subCmd := flag.NewFlagSet(subCmdName, flag.ContinueOnError)
	var params = make(map[string]*string)
	for argName, argUsage := range args {
		param := subCmd.String(argName, "", argUsage)
		params[argName] = param
	}

	testNetwork := subCmd.String("test", "", "")

	err := subCmd.Parse(os.Args[2:])
	if err != nil {
		return
	}

	if subCmd.Parsed() {
		fmt.Printf("run cmd: %s%v \n", subCmdName, printArgs(params))
		// 初始化网络
		cl.initNetWork(testNetwork)
		runCmd(params)
		return
	}
	fmt.Printf("run cmd %s error\n", subCmdName)
}

func printArgs(args map[string]*string) string {
	var argValues string
	for key, value := range args {
		argValues = argValues + " -" + key + " " + *value
	}
	return argValues
}

func (cl *CommandLine) checkArgs() {
	if len(os.Args) < 2 {
		cl.printUsage()
		runtime.Goexit()
	}
}

func (cl *CommandLine) printUsage() {
	fmt.Println("Welcome to Leo Cao's tiny blockchain system, usage is as follows:")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a wallet.")
	fmt.Println("And then you can use the wallet address to create a blockchain and declare the owner.")
	fmt.Println("Make transactions to expand the blockchain.")
	fmt.Println("In addition, don't forget to run mine function after transatcions are collected.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("createwallet -refname REFNAME                       ----> Creates and save a wallet. The refname is optional.")
	fmt.Println("walletinfo -refname NAME -address P2PKHAddress           ----> Print the information of a wallet. At least one of the refname and address is required.")
	fmt.Println("walletsupdate                                       ----> Registrate and update all the wallets (especially when you have added an existed .wlt file).")
	fmt.Println("walletslist                                         ----> List all the wallets found (make sure you have run walletsupdate first).")
	fmt.Println("createblockchain -refname NAME -address ADDRESS     ----> Creates a blockchain with the owner you input (address or refname).")
	fmt.Println("balance -refname NAME -address ADDRESS              ----> Back the balance of a wallet using the address (or refname) you input.")
	fmt.Println("blockchaininfo                                      ----> Prints the blocks in the chain.")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block.")
	fmt.Println("sendbyrefname -from NAME1 -to NAME2 -amount AMOUNT  ----> Make a transaction and put it into candidate block using refname.")
	fmt.Println("mine                                                ----> Mine and add a block to the chain.")
	fmt.Println("---------------------------------------------------------------------------------------------------------------------------------------------------------")
}

func (cl *CommandLine) initNetWork(network *string) {
	if *network == "" {
		fmt.Printf("not set network,use default network :%s \n", constant.Network)
	} else {
		fmt.Printf("use  network :%s \n", *network)
		constant.Network = *network
	}

	if !utils.FileExists(constant.GetNetworkPath(constant.BCPath)) {
		os.Mkdir(constant.GetNetworkPath(constant.BCPath), 0644)
	}
	if !utils.FileExists(constant.GetNetworkPath(constant.UTXOPATH)) {
		os.Mkdir(constant.GetNetworkPath(constant.UTXOPATH), 0644)
	}
	if !utils.FileExists(constant.GetNetworkPath(constant.WalletsRefList)) {
		os.Mkdir(constant.GetNetworkPath(constant.WalletsRefList), 0644)
	}
	if !utils.FileExists(constant.GetNetworkPath(constant.Wallets)) {
		os.Mkdir(constant.GetNetworkPath(constant.Wallets), 0644)
	}
}

func (cl *CommandLine) getAddressByRefName(refName string) string {
	refList := wallet.LoadRefList()
	address, err := refList.FindAddress(refName)
	utils.Handle(err)
	return address
}
func (cl *CommandLine) createWallet(refName string) *wallet.Wallet {
	newWallet := wallet.NewWallet()
	newWallet.SaveWallet()
	refList := wallet.LoadRefList()
	refList.BindRef(string(newWallet.P2PKHAddress()), refName)
	refList.Save()
	fmt.Printf("Succeed create wallet:%s %s\n", refName, string(newWallet.P2PKHAddress()))
	return newWallet
}

func (cl *CommandLine) walletInfoByRefName(refName string) {
	cl.walletInfo(cl.getAddressByRefName(refName))
}
func (cli *CommandLine) walletInfo(address string) {
	wlt := wallet.LoadWallet(address)
	refList := wallet.LoadRefList()
	fmt.Printf("Wallet address:%x\n", wlt.P2PKHAddress())
	fmt.Printf("Public Key:%x\n", wlt.PublicKey)
	fmt.Printf("Reference Name:%s\n", (*refList)[address])
}

func (cli *CommandLine) walletsUpdate() {
	refList := wallet.LoadRefList()
	refList.Update()
	refList.Save()
	fmt.Println("Succeed in updating wallets.")
}

func (cli *CommandLine) walletsList() {
	refList := wallet.LoadRefList()
	for address, _ := range *refList {
		wlt := wallet.LoadWallet(address)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Wallet address:%s\n", address)
		fmt.Printf("Public Key:%x\n", wlt.PublicKey)
		fmt.Printf("Reference Name:%s\n", (*refList)[address])
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
	}
}

func (cl *CommandLine) create(address string) {
	newChain := blockchain.CreateBlockChain(utils.Address2PubHash([]byte(address)))
	newChain.Database.Close()

	fmt.Println("Finished creating blockchain, and the owner is: ", address)
}
func (cl *CommandLine) createByRefName(refName string) {
	cl.create(cl.getAddressByRefName(refName))
}

func (cl *CommandLine) Balance(address string) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	utxo := chain.UtxoSet.GetUtxos(address)

	fmt.Printf("P2PKHAddress: %s, Balance: %d \n", address, utxo.GetBalance())
}
func (cl *CommandLine) balanceByRefName(refName string) {
	cl.Balance(cl.getAddressByRefName(refName))
}

func (cl *CommandLine) balance2(address string) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	wlt := wallet.LoadWallet(address)
	balance := chain.GetBalance(utils.PubHash2Address(utils.PublicKeyHash(wlt.PublicKey)))

	fmt.Printf("P2PKHAddress: %s, Balance: %d \n", address, balance)
}
func (cl *CommandLine) balanceByRefName2(refName string) {
	cl.balance2(cl.getAddressByRefName(refName))
}

func (cl *CommandLine) send(fromAddress string, toAddress string, amount int) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	fromWallet := wallet.LoadWallet(fromAddress)
	tx, ok := chain.CreateTransaction(fromWallet.PublicKey, utils.Address2PubHash([]byte(toAddress)), amount, fromWallet.PrivateKey)
	if !ok {
		fmt.Println("Failed to create transaction")
		return
	}

	txPool := blockchain.CreateTransactionPool()
	txPool.AddTransaction(tx)
	txPool.SaveFile()
	fmt.Println("Success")
}
func (cl *CommandLine) sendByRefName(fromRefName string, toRefName string, amount int) {
	cl.send(cl.getAddressByRefName(fromRefName), cl.getAddressByRefName(toRefName), amount)
}

func (cl *CommandLine) mine() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	newBlock := chain.RunMine()

	chain.UtxoSet.Update(newBlock)
	fmt.Println("Finished mine")
}

func (cl *CommandLine) Info() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	iter := chain.Iterator()
	for {
		block := iter.Next()
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Printf("Timestamp:%d\n", block.Timestamp)
		fmt.Printf("Previous hash:%x\n", block.PreHash)
		fmt.Printf("Transactions:%v\n", block.Transactions)
		fmt.Printf("hash:%x\n", block.Hash)
		fmt.Printf("Pow: %s\n", strconv.FormatBool(block.ValidatePow()))
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println()
		if !iter.HasNext() {
			break
		}
	}
}

func (cl *CommandLine) GetUtxos() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	refList := wallet.LoadRefList()
	for _, utxo := range chain.UtxoSet.UTXOS {
		refName, _ := refList.FindRefName(utxo.Address)
		for _, item := range utxo.UxtoItems {
			fmt.Println("--------------------------------------------------------------------------------------------------------------")
			fmt.Printf("address：%s rename:%s \n", utxo.Address, refName)
			fmt.Printf("txId:%s \n", item.Txid)
			fmt.Printf("outIndex:%d\n", item.OutIdx)
			fmt.Printf("amount:%d\n", item.Value)
		}
	}
	fmt.Println()
	fmt.Println("Finished getUtxos")
}
