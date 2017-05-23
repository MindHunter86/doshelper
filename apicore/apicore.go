package apicore

import "errors"

var (
	err_nil = errors.New("nil")
)

type ApiCore struct {
}
func (self *ApiCore) InitModule() ( *ApiCore, error ) {
}
func (self *ApiCore) DeInitModule() error {
}
