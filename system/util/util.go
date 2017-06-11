package util


import "github.com/sirupsen/logrus"
//import log "github.com/sirupsen/logrus"
import "golang.org/x/net/context"

type utilSubModule struct {
	logout *logrus.Logger
}
func (self *utilSubModule) configure(ctx context.Context) (*utilSubModule, error) {
	if self == nil { return nil,Err_Glob_InvalidSelf }
	if ctx == nil { return nil,Err_Glob_InvalidContext }
	return self,nil
}
