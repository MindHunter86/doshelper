package apicore

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"
import "github.com/valyala/fasthttp"


type apiMiddleware struct {
	*ApiCore
}
func (self *apiMiddleware) configure(ctx context.Context) (*apiMiddleware, error) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*ApiCore)
	self.slogger.W( log.LLEV_DBG, "Middleware submodule has been initialized and configured!" )

	return self,nil
}
func (self *apiMiddleware) QueryIdentificatorAdd(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		self.slogger.W( log.LLEV_DBG, "Middleware ..." )
		next(ctx); return
	})
}
