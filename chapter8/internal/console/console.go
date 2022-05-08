package console

import (
	"bufio"
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
		c.Cli.Run(args)
	}
}
