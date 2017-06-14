package controller

import "golucky/system/util"
import "golucky/controller/api"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type ControllerModule struct {
	logout *logrus.Logger

	ptr_sub_api *api.ApiSubmodule
}
func (self *ControllerModule) Configure(ctx context.Context) (*ControllerModule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	// logger initialization:
	self.logout = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)

	// submodules initialization:
	var e error
	if self.ptr_sub_api, e = new(api.ApiSubmodule).Configure(ctx); e != nil { return nil,e }

	self.logout.Debugln("Controller Module has been successfully initialized and configured!")
	return self,nil
}
func (self *ControllerModule) Destroy() error {
	self.logout.Debugln("Controller Module has been successfully destroyed!")
	return nil
}
