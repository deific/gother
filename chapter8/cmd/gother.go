package main

import (
	"flag"
	"gother/chapter8/internal/cli"
	"gother/chapter8/internal/console"
	"os"
)

func main() {
	startCli()
}

func startCli() {
	defer os.Exit(0)

	cli := cli.New()
	defer cli.Close()

	c := flag.Bool("console", false, "console")
	flag.Parse()
	if *c {
		console := console.Console{}
		console.Start(cli)
	} else {
		cli.Run(os.Args[1:])
	}
}
