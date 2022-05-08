package console

import (
	"gother/chapter8/internal/cli"
)

type Console struct {
	cli *cli.CommandLine
}

const prompt = "gother > "
const (
	cmd_balance = "balance"
)

func (c *Console) Start() {
	//in := bufio.NewReader(os.Stdin) // 声明并初始化读取器

}
