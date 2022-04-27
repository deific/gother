package main

import (
	"gother/chapter4/internal/cli"
	"os"
)

func main() {
	startCli()
}

func startCli() {
	defer os.Exit(0)
	cmd := cli.CommandLine{}
	cmd.Run()
}
