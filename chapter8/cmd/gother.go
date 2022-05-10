package main

import (
	"flag"
	"gother/chapter8/internal/cli"
	"gother/chapter8/internal/console"
	"gother/chapter8/internal/rpc"
	"os"
	"time"
)

func main() {
	startCli()
}

func startCli() {
	defer os.Exit(0)

	cli := cli.New()
	defer cli.Close()

	c := flag.Bool("console", false, "-console or --console, if exist will start console")
	r := flag.Bool("rpc", false, "-rpc or --rpc, if exist will start rpc server")

	flag.Parse()

	if *r {
		go rpc.StartServer()
	}

	time.Sleep(time.Second)
	if *c {
		console := console.Console{}
		console.Start(cli)
	} else {
		cli.Run(os.Args[1:])
	}
}
