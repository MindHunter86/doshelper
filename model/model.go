package model

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type ModelModule struct {
	logout *logrus.Logger
}
func (self *ModelModule) Configure(ctx context.Context) (*ModelModule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.logout = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.logout.Debugln("Model Module has been successfully initialized and configured!")

	return self,nil
}
func (self *ModelModule) Destroy() error {
	self.logout.Debugln("Model Module has been successfully destroyed!")
	return nil
}
