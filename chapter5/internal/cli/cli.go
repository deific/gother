package cli

import (
	"flag"
	"fmt"
	"gother/chapter5/internal/blockchain"
	"gother/chapter5/internal/utils"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
}

func (cl *CommandLine) Run() {
	cl.checkArgs()

	cl.parseAndRunCmd("create", map[string]string{"address": "The address refer to the owner of blockchain"}, func(args map[string]*string) {
		cl.create(*args["address"])
	})

	cl.parseAndRunCmd("balance", map[string]string{"address": "Who need to get balance amount"}, func(args map[string]*string) {
		cl.balance(*args["address"])
	})

	cl.parseAndRunCmd("info", map[string]string{}, func(args map[string]*string) {
		cl.info()
	})

	cl.parseAndRunCmd("send", map[string]string{"from": "Source address", "to": "Destination address", "amount": "Amount to send"}, func(args map[string]*string) {
		amount, err := strconv.Atoi(*args["amount"])
		utils.Handle(err)
		cl.send(*args["from"], *args["to"], amount)
	})

	cl.parseAndRunCmd("mine", map[string]string{}, func(args map[string]*string) {
		cl.mine()
	})
}

func (cl *CommandLine) parseAndRunCmd(subCmdName string, args map[string]string, runCmd func(args map[string]*string)) {
	if subCmdName != os.Args[1] {
		return
	}

	subCmd := flag.NewFlagSet(subCmdName, flag.ExitOnError)
	var params = make(map[string]*string)
	for argName, argUsage := range args {
		param := subCmd.String(argName, "", argUsage)
		params[argName] = param
	}

	err := subCmd.Parse(os.Args[2:])
	utils.Handle(err)

	if subCmd.Parsed() {
		fmt.Printf("run cmd: %s%v \n", subCmdName, printArgs(params))
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
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("All you need is to first create a blockchain and declare the owner.")
	fmt.Println("And then you can make transactions.")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("create -address ADDRESS                   ----> Creates a blockchain with the owner you input")
	fmt.Println("balance -address ADDRESS                            ----> Back the balance of the address you input")
	fmt.Println("info                                      ----> Prints the blocks in the chain")
	fmt.Println("send -from FROADDRESS -to TOADDRESS -amount AMOUNT  ----> Make a transaction and put it into candidate block")
	fmt.Println("mine                                                ----> Mine and add a block to the chain")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
}

func (cl *CommandLine) create(address string) {
	newChain := blockchain.InitBlockChain([]byte(address))
	newChain.Database.Close()
	fmt.Println("Finished creating blockchain, and the owner is: ", address)
}

func (cl *CommandLine) balance(address string) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	balance, _ := chain.FindUTXOs([]byte(address))

	fmt.Printf("Address: %s, Balance: %d \n", address, balance)
}

func (cl *CommandLine) send(fromAddress string, toAddress string, amount int) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	tx, ok := chain.CreateTransaction([]byte(fromAddress), []byte(toAddress), amount)
	if !ok {
		fmt.Println("Failed to create transaction")
		return
	}

	txPool := blockchain.CreateTransactionPool()
	txPool.AddTransaction(tx)
	txPool.SaveFile()
	fmt.Println("Success")
}

func (cl *CommandLine) mine() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	chain.RunMine()
	fmt.Println("Finished mine")
}

func (cl *CommandLine) info() {
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
