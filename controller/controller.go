package controller

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type ControllerModule struct {
	logout *logrus.Logger
}
func (self *ControllerModule) Configure(ctx context.Context) (*ControllerModule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.logout = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.logout.Debugln("Controller Module has been successfully initialized and configured!")

	return self,nil
}
func (self *ControllerModule) Destroy() error {
	self.logout.Debugln("Controller Module has been successfully destroyed!")
	return nil
}
