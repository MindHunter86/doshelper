package apicore

import "errors"
import "doshelpv2/users"
import "golang.org/x/net/context"

var (
	err_nil = errors.New("nil")
)

type ApiCore struct {
	users *users.Users
	handlers *handler
}
func InitModule( ctx context.Context ) ( *ApiCore, error ) {
	//	app := ctx.Value(doshelpv2.CTX_APP).(doshelpv2.AppCtx)
	return nil,nil
}
func (self *ApiCore) DeInitModule() error {
	return nil
}
func (self *ApiCore) GetUsers() ( *users.Users ) {
	return self.users
}
func (self *ApiCore) GetHandlers() ( *handler ) {
	return self.handlers
}
