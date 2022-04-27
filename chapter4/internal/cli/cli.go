package cli

import (
	"flag"
	"fmt"
	"gother/chapter4/internal/blockchain"
	"gother/chapter4/internal/utils"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
}

func (cl *CommandLine) Run() {
	cl.checkArgs()

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	infoCmd := flag.NewFlagSet("info", flag.ExitOnError)

	createCmdAddress := createCmd.String("address", "", "The address refer to the owner of blockchain")
	balanceAddress := balanceCmd.String("address", "", "Who need to get balance amount")

	switch os.Args[1] {
	case "create":
		err := createCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "balance":
		err := balanceCmd.Parse(os.Args[2:])
		utils.Handle(err)
	case "info":
		err := infoCmd.Parse(os.Args[2:])
		utils.Handle(err)
	}

	if createCmd.Parsed() {
		if *createCmdAddress == "" {
			createCmd.Usage()
			runtime.Goexit()
		}
		cl.create(*createCmdAddress)
	}

	if balanceCmd.Parsed() {
		if *balanceAddress == "" {
			balanceCmd.Usage()
			runtime.Goexit()
		}
		cl.balance(*balanceAddress)
	}

	if infoCmd.Parsed() {
		cl.info()
	}
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
