package apicore

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"


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
