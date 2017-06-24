package service

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"

type ServiceSubmodule struct {
	log *logrus.Logger
}
func (self *ServiceSubmodule) Configure(ctx context.Context) (*ServiceSubmodule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.log = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.log.Debugln("Service submodule has been successfully initialized and configured!")

	return self,nil
}
