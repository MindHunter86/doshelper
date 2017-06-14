package api

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"
import "github.com/buaazp/fasthttprouter"


type ApiSubmodule struct {
	log *logrus.Logger
	httprouter *fasthttprouter.Router
}
func (self *ApiSubmodule) Configure(ctx context.Context) (*ApiSubmodule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.log = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.log.Debugln("Controller/Api submodule has been successfully initialized and configured!")

	return self,nil
}
func (self *ApiSubmodule) GetRouter() *fasthttprouter.Router {
	return self.httprouter
}
func (self *ApiSubmodule) createRouter() error {
	if self.httprouter != nil { return util.Err_Controller_Api_InvalidRouter }

	self.httprouter = fasthttprouter.New()

	// api methods:
	self.httprouter.GET("/", nil)

	return nil
}
