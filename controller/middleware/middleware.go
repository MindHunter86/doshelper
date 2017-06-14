package middleware

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"
import "github.com/valyala/fasthttp"


type Middlewares interface {
	QueryIdentificatorAdd(fasthttp.RequestHandler) fasthttp.RequestHandler
}
type MiddlewareSubmodule struct {
	log *logrus.Logger
	middlewares Middlewares
}
func (self *MiddlewareSubmodule) Configure(ctx context.Context) (*MiddlewareSubmodule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.log = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.log.Debugln("Controller/MiddleWare submodule has been successfully initialized and configured!")

	return self,nil
}
func (self *MiddlewareSubmodule) createMiddlewares() error {
	if self.middlewares != nil { return util.Err_Controller_InvalidMiddle }

	self.middlewares = nil
	return nil
}
func (self *MiddlewareSubmodule) GetMiddlewares() Middlewares {
	return self.middlewares
}
