package system

import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type SystemModule struct {
	logout *logrus.Logger
}
func (self *SystemModule) Configure(ctx context.Context) (*SystemModule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.logout = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.logout.Debugln("System Module has been successfully initialized and configured!")

	return self,nil
}
func (self *SystemModule) Destroy() error {
	self.logout.Debugln("System Module has been successfully destroyed!")
	return nil
}
