package apimodule

import "net"
import "sync"

import "github.com/valyala/fasthttp"

type apiService struct {
	lsn net.Listener

	sync.RWMutex
	tcp_alive bool
}
func (self *apiService) configure() ( *apiService, error ) {
	self.Lock(); self.tcp_alive = false; self.Unlock()

	if e := self.newListener(); e != nil { return nil,e }
	return self,nil
}
func (self *apiService) kill() error {
	self.Lock()
	defer self.Unlock()

	if self.tcp_alive == false { return err_lsnNotAlive }
	self.tcp_alive = false
	return self.lsn.Close()
}
func (self *apiService) serve() error {
	var listener net.Listener
	var router func(ctx *fasthttp.RequestCtx)

	if self.RLock(); self.lsn != nil {
		listener, router = self.lsn, apimodule.router.Handler
	} else { self.RUnlock(); return err_lsnNotAlive }
	self.RUnlock()

	e := fasthttp.Serve( listener, router )

	self.RLock(); defer self.RUnlock()
	if self.tcp_alive == false { return nil } //  it means, that exit called by application
	return e
}
func (self *apiService) newListener() error {
	if self.lsn != nil { return err_lsnAlreadyDefined }

	var e error
	self.lsn, e = net.Listen( module_input_lsnproto, module_input_lsnport ); if e != nil {
		return e
	}

	self.Lock(); self.tcp_alive = true; self.Unlock()
	return nil
}
