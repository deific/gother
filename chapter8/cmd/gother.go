package main

import (
	"gother/chapter8/internal/cli"
	"gother/chapter8/internal/console"
	"os"
)

func main() {
	startCli()
}

func startCli() {
	defer os.Exit(0)
	cli := cli.CommandLine{}
	//cli.Run()

	console := console.Console{Cli: &cli}
	console.Start()
}
