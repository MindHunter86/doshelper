package system

import "golucky/system/util"
import "golucky/system/service"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type SystemModule struct {
	logger *logrus.Logger
	done_pipe <-chan struct{}

	ptr_sub_service *service.ServiceSubmodule
}
func (self *SystemModule) Configure(ctx context.Context) (util.AppModule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.done_pipe = ctx.Done()
	self.logger = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)

	var e error
	if self.ptr_sub_service, e = new(service.ServiceSubmodule).Configure(ctx); e != nil { return nil,e }

	self.logger.Debugln("System Module has been successfully initialized and configured!")
	return self,nil
}
func (self *SystemModule) Load() error {
	return self.ptr_sub_service.Run()
}
func (self *SystemModule) Unload() {
	self.logger.Debugln("System Module has been successfully unloaded!")
}



//type ControllerModule struct {
//	logout *logrus.Logger
//
//	ptr_sub_api *api.ApiSubmodule
//	ptr_sub_middle *middleware.MiddlewareSubmodule
//}
//func (self *ControllerModule) Configure(ctx context.Context) (*ControllerModule, error) {
//	if self == nil { return nil,util.Err_Glob_InvalidSelf }
//	if ctx == nil { return nil,util.Err_Glob_InvalidContext }
//
//	var e error
//
//	// logger initialization:
//	self.logout = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
//
//	// submodules initialization:
//	if self.ptr_sub_middle, e = new(middleware.MiddlewareSubmodule).Configure(ctx); e == nil {
//		ctx = context.WithValue(ctx, util.CTX_CNTR_MIDDLE, self.ptr_sub_middle)
//	} else { return nil,e }
//
//	if self.ptr_sub_api, e = new(api.ApiSubmodule).Configure(ctx); e != nil { return nil,e }
//
//	self.logout.Debugln("Controller Module has been successfully initialized and configured!")
//	return self,nil
//}
//func (self *ControllerModule) Destroy() error {
//	self.logout.Debugln("Controller Module has been successfully destroyed!")
//	return nil
//}
