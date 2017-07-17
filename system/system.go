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
	self.ptr_sub_service.Destroy()
	self.logger.Debugln("System Module has been successfully unloaded!")
}
