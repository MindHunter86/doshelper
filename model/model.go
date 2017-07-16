package model

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type ModelModule struct {
	logger *logrus.Logger
}
func (self *ModelModule) Configure(ctx context.Context) (util.AppModule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.logger = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.logger.Debugln("Model Module has been successfully initialized and configured!")

	return self,nil
}
func (self *ModelModule) Load() error { return nil }
func (self *ModelModule) Unload() {
	self.logger.Debugln("Model Module has been successfully unloaded!")
}
