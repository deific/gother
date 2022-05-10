package rpc

import (
	"fmt"
	"net/http"
	"net/rpc"
)

func StartServer() {
	// 注册服务
	rpc.Register(JsonRpcService{})
	rpc.HandleHTTP()

	addr := ":8333"
	// 启动 server
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("rpc started listen by: " + addr)
	//
	//for {
	//	// 不断连接服务
	//	conn, err := listener.Accept()
	//	if err != nil {
	//		log.Printf("accept error, %v", err)
	//		continue
	//	}
	//	fmt.Println("new client accept")
	//	// 使用 Goroutine：ServeConn runs the JSON-RPC server on a single connection.
	//	go jsonrpc.ServeConn(conn)
	//}
}
