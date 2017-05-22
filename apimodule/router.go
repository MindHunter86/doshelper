package apimodule

import "github.com/valyala/fasthttp"
import "github.com/buaazp/fasthttprouter"

type apiRouter struct {
	*fasthttprouter.Router
}
func (self *apiRouter) configure() ( *apiRouter, error ) {
	if self.Router != nil { return nil,err_rtrAlreadyDefined }

	self.Router = fasthttprouter.New()
	self.GET("/", self.rt_index)
	return self,nil
}
func (self *apiRouter) rt_index(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Hello world! Index router."))
}
