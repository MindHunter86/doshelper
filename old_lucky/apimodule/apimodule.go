package apimodule

import "errors"
import "doshelpv2/appctx"
import "doshelpv2/apicore"
import "golang.org/x/net/context"

var (
	module_input_lsnproto = "tcp"
	module_input_lsnport = ":8082"
)
var (
	err_lsnAlreadyDefined = errors.New("Net.Listener has been already defined!")
	err_lsnNotAlive = errors.New("Module Listener is dead!")
	err_rtrAlreadyDefined = errors.New("Router has been already defined")
	err_Init_InvalidInput = errors.New("Input variables in init function are corrupted!")
	err_Init_InvalidCtxPointer = errors.New("Received pointer from Context is corrupted!")
)

var apimodule *ApiModule

type ApiModule struct {
	service *apiService
	router *apiRouter
	log *apiLogger

	core *apicore.ApiCore
}
func InitModule( ctx context.Context ) ( *ApiModule, error ) {
	var e error
	if apimodule == nil { apimodule = new(ApiModule) }
	if ctx == nil { return nil,err_Init_InvalidInput }

	apimodule.core = ctx.Value(appctx.CTX_MOD_APICORE).(*apicore.ApiCore)
	if apimodule.core == nil { return nil,err_Init_InvalidCtxPointer }

	if apimodule.log, e = new(apiLogger).configure(ctx); e != nil { return nil,e }
	if apimodule.router, e = new(apiRouter).configure(ctx); e != nil { return nil,e }
	if apimodule.service, e = new(apiService).configure(); e != nil { return nil,e }
	return apimodule,nil
}
func (self *ApiModule) Serve() error {
	return apimodule.service.serve()
}
func (self *ApiModule) DeInitModule() error {
	return self.service.kill()
}
