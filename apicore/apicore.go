package apicore

import "os"
import "log"
import "errors"
import "doshelpv2/appctx"
import dlog "doshelpv2/log"
import "doshelpv2/users"
import "golang.org/x/net/context"
import "github.com/valyala/fasthttp"

var (
	err_nil = errors.New("nil")
)

type ApiHandlers interface {
	Login(*fasthttp.RequestCtx)
}
type ApiCore struct {
	slogger *dlog.Logger
	users *users.Users
	Handlers ApiHandlers
}
func InitModule( ctx context.Context ) ( *ApiCore, error ) {
	var core *ApiCore = new(ApiCore)

	// Log testing:
	core.slogger = &dlog.Logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		Ch_message: ctx.Value(appctx.CTX_LOG_FILE).(*dlog.FileLogger).Mess_queue,
		Prefix: dlog.LPFX_MODCORE,
	}
	core.slogger.W( dlog.LLEV_DBG, "TEST TEST TEST FROM MODULE" )

	//Handlers initialization:
	core.Handlers = new(handler)

	return core,nil
}
func (self *ApiCore) DeInitModule() error {
	self.slogger.W( dlog.LLEV_OK, "YEAH! Module Exit!" )
	return nil
}
func (self *ApiCore) GetUsers() ( *users.Users ) {
	return self.users
}
