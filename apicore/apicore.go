package apicore

import "errors"
import "doshelpv2/users"

var (
	err_nil = errors.New("nil")
)

type ApiCore struct {
	users *users.Users
}
func (self *ApiCore) InitModule() ( *ApiCore, error ) {
	return nil,nil
}
func (self *ApiCore) DeInitModule() error {
	return nil
}
func (self *ApiCore) GetUsers() ( *users.Users ) {
	return self.users
}
