package apimodule

import "errors"

var (
	module_input_lsnproto = "tcp"
	module_input_lsnport = ":8082"
)
var (
	err_lsnAlreadyDefined = errors.New("Net.Listener has been already defined!")
	err_lsnNotAlive = errors.New("Module Listener is dead!")
	err_rtrAlreadyDefined = errors.New("Router has been already defined")
)

var apimodule *ApiModule

type ApiModule struct {
	service *apiService
	router *apiRouter
}
func InitModule() ( *ApiModule, error ) {
	var e error
	if apimodule != nil { apimodule = new(ApiModule) }

	if apimodule.router, e = new(apiRouter).configure(); e != nil { return nil,e }
	if apimodule.service, e = new(apiService).configure(); e != nil { return nil,e }
	return apimodule,nil
}
func (self *ApiModule) DeInitModule() error {
	return self.service.kill()
}
