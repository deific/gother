package console

import (
	"bufio"
	"flag"
	"fmt"
	"gother/chapter8/internal/cli"
	"os"
	"strings"
)

type Console struct {
	Cli *cli.CommandLine
}

const prompt = "gother > "
const (
	cmd_balance = "balance"
)

func (c *Console) Start() {
	in := bufio.NewReader(os.Stdin) // 声明并初始化读取器
	for {
		fmt.Printf(prompt)
		input, _, err := in.ReadLine()
		if err != nil {
			continue
		}
		args := strings.Split(string(input), " ")
		c.parseAndRunCmd("blockchaininfo", map[string]string{}, args, func(args map[string]*string) {
			c.Cli.Info()
		})
	}
}

func (cl *Console) parseAndRunCmd(subCmdName string, argNames map[string]string, args []string, runCmd func(args map[string]*string)) {
	if subCmdName != args[0] {
		return
	}

	subCmd := flag.NewFlagSet(subCmdName, flag.ContinueOnError)
	var params = make(map[string]*string)
	for argName, argUsage := range argNames {
		param := subCmd.String(argName, "", argUsage)
		params[argName] = param
	}

	err := subCmd.Parse(args[1:])
	if err != nil {
		return
	}

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
