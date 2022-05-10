package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func StartServer() {
	// 注册服务
	rpc.Register(JsonRpcService{})
	addr := "127.0.0.1:1234"
	// 启动 server
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("rpc started listen by: " + addr)

	for {
		// 不断连接服务
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error, %v", err)
			continue
		}
		// 使用 Goroutine：ServeConn runs the JSON-RPC server on a single connection.
		go jsonrpc.ServeConn(conn)
	}
}
