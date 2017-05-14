package main

import (
	"net"
	"net/rpc"

	"sync"
	"errors"
)


var (
	err_RpcService_DefinedListener = errors.New("Listener has been already defined!")
	err_RpcService_NonAliveListener = errors.New("Listener is not alive!")
	err_RpcService_UnexpectedClose = errors.New("Rpc Service has been closed!")
)

type rpcService struct {
	listener *net.Listener
	server *rpc.Server
	sync.RWMutex
	alive bool
}
func _rpcService() ( *rpcService, error ) {
	var rpc *rpcService = new(rpcService)

	if rpc.newListener(); e != nil { return nil,e }
	return rpc,nil
}
func (self *rpcService) newListener() error {
	if self.listener != nil { return err_RpcService_DefinedListener }

	var e error
	if self.listener, e := net.Listen( "tcp", ":8081" ); e == nil {
		self.Lock()
		self.alive = true
		self.Unlock()
	} else { return e }
}
func (self *rpcService) killListener() error {
	self.Lock()
	defer self.Unlock()

	if self.alive == true {
		self.alive = false
		return self.listener.Close()
	} else { return err_RpcService_NonAliveListener }
}

func (self *rpcService) rpvServe() error {

	if self.listener != nil {
		self.server.Accept(self.listener)
	} else { return err_RpcService_NonAliveListener }

	if self.Lock(); self.alive == true {
		return err_RpcService_UnexpectedClose
	}
}
