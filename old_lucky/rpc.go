package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"sync"
	"errors"
)


var (
	err_RpcService_DefinedListener = errors.New("Listener has been already defined!")
	err_RpcService_NonAliveListener = errors.New("Listener is not alive!")
)

type rpcService struct {
	listener net.Listener
	server *rpc.Server
	sync.RWMutex
	alive bool
}
func _rpcService() ( *rpcService, error ) {
	var service *rpcService = new(rpcService)

	service.server = rpc.NewServer()
	if e := service.newListener(); e != nil { return nil,e }
	return service,nil
}
func (self *rpcService) newListener() error {
	if self.listener != nil { return err_RpcService_DefinedListener }

	var e error
	if self.listener, e = net.Listen( "tcp", ":8081" ); e != nil {
		return e
	}

	self.Lock()
	self.alive = true
	self.Unlock()

	return nil
}
func (self *rpcService) killListener() error {
	self.Lock()
	defer self.Unlock()

	if self.alive == false { return err_RpcService_NonAliveListener }
	self.alive = false
	return self.listener.Close()
}

func (self *rpcService) serve() error {
	if self.listener == nil  { return err_RpcService_NonAliveListener }

	var e error
	for {
		conn, e := self.listener.Accept(); if e != nil { break }
		go self.server.ServeCodec( jsonrpc.NewServerCodec(conn) )
	}

	self.RLock()
	defer self.RUnlock()
	if self.alive == true { return e }

	return nil 
}
