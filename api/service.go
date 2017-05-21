package apison

import "net"
import "sync"
import "errors"

import "github.com/valyala/fasthttp"
import "github.com/buaazp/fasthttprouter"

var (
	module_input_lsnproto = "tcp"
	module_input_lsnport = ":8082"
)
var (
	err_lsnAlreadyDefined = errors.New("Net.Listener has been already defined!")
	err_lsnNotAlive = errors.New("Module Listener is dead!")
)

type apiModule struct {
	lsn net.Listener
	rtr *apiRouter

	sync.RWMutex
	tcp_alive bool
}
func _apiModule() ( *apiModule, error ) {
	var module *apiModule = new(apiModule)

	module.rtr = _apiRouter()
	if e := module.newListener(); e != nil { return nil,e }
	return module,nil
}
func (self *apiModule) newListener() error {
	if self.lsn != nil { return err_lsnAlreadyDefined }

	var e error
	self.lsn, e = net.Listen( module_input_lsnproto, module_input_lsnport ); if e != nil {
		return e
	}

	self.Lock(); self.tcp_alive = true; self.Unlock()
	return nil
}
func (self *apiModule) killListener() error {
	self.Lock()
	defer self.Unlock()

	if self.tcp_alive == false { return err_lsnNotAlive }
	self.tcp_alive = false
	return self.lsn.Close()
}
func (self *apiModule) serve() error {
	var listener net.Listener
	var router func(ctx *fasthttp.RequestCtx)

	if self.RLock(); self.lsn != nil {
		listener, router = self.lsn, self.rtr.Handler
	} else { self.RUnlock(); return err_lsnNotAlive }
	self.RUnlock()

	e := fasthttp.Serve( listener, router )

	self.RLock(); defer self.RUnlock()
	if self.tcp_alive == false { return nil } //  it means, that exit called by application
	return e
}

type apiRouter struct {
	*fasthttprouter.Router
}
func _apiRouter() *apiRouter {
	var router *apiRouter = new(apiRouter)
	router.Router = fasthttprouter.New()
	router.GET("/", router.rt_index)
	return router
}
func (self *apiRouter) rt_index(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Hello world! Index router."))
}
