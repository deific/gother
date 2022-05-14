package rpc

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"net/http"
	"net/rpc"
	"reflect"
	"sync"
)

type Server struct {
	services *serviceRegistry
	run      int32
	codecs   mapset.Set
}

type serviceRegistry struct {
	mu       sync.Mutex
	services map[string]service
}
type service struct {
	name    string
	srv     interface{}
	srvType string
}

func NewServer() *Server {
	server := &Server{
		codecs: mapset.NewSet(),
		run:    1,
	}

	rpcService := &RpcService{server: server}
	server.RegisterName("/rpc", rpcService)
	return server
}

func (s *Server) RegisterName(name string, service interface{}) {
	s.services.registerName(name, service)
}

func (r *serviceRegistry) registerName(name string, srv interface{}) error {
	if srv == nil {
		return fmt.Errorf("no service")
	}

	srvType := reflect.ValueOf(srv)
	if name == "" {
		return fmt.Errorf("no service name for type#{srvType.Type().String()}")
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.services == nil {
		r.services = make(map[string]service)
	}
	svc, ok := r.services[name]
	if !ok {
		svc = service{
			name:    name,
			srvType: srvType.String(),
			srv:     srv,
		}
		r.services[name] = svc
	}
	return nil
}

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
