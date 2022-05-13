package rpc

import (
	"fmt"
	"net/http"
	"net/rpc"
)

func StartServer() {
	// 注册服务
	rpc.Register(&JsonRpcService{})
	// go自带的rpc，如果使用http协议，只能由go
	rpc.HandleHTTP()

	addr := ":8333"
	// 启动 server
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("rpc started listen by: " + addr)
}
