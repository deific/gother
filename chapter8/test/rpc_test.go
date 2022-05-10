package test

import (
	"fmt"
	rpc2 "gother/chapter8/internal/rpc"
	"log"
	"net/rpc"
	"testing"
)

func TestJsonRpcHttp(t *testing.T) {
	//连接远程rpc服务
	rpc, err := rpc.DialHTTP("tcp", "127.0.0.1:8333")
	if err != nil {
		log.Fatal(err)
	}
	var res rpc2.WalletInfoRes
	//调用远程方法
	//注意第三个参数是指针类型
	err2 := rpc.Call("JsonRpcService.WalletInfo", rpc2.WalletInfoReq{Address: "1Fvba7ojkjzV6hEMRcYBdiznqspDZHVnjg"}, &res)
	if err2 != nil {
		log.Fatal(err2)
	}
	fmt.Println(res)

}
