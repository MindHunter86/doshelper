package apimodule

import "doshelpv2/appctx"
import "doshelpv2/apicore"
import "golang.org/x/net/context"
import "github.com/valyala/fasthttp"
import "github.com/buaazp/fasthttprouter"

type apiRouter struct {
	*fasthttprouter.Router
}
func (self *apiRouter) configure( ctx context.Context ) ( *apiRouter, error ) {
	if self.Router != nil { return nil,err_rtrAlreadyDefined }

	var core *apicore.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*apicore.ApiCore)
	if core == nil { return nil,err_Init_InvalidCtxPointer }

	self.Router = fasthttprouter.New()
	self.GET("/", self.rt_index)
	self.GET("/login", core.Handlers.Login)

	// Api v1 methods:
	self.GET("/api/v1", nil)
	self.GET("/api/v1/centrifugo/connection", core.Handlers.CentrifugoConnection)

	return self,nil
}
func (self *apiRouter) rt_index(ctx *fasthttp.RequestCtx) {
	ctx.Write([]byte("Hello world! Index router."))
}
