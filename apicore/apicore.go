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
	err_glob_InvalidSelf = errors.New("Invalid self struct in configure method! Is self initialized?")
	err_glob_InvalidContext = errors.New("Invalid input context in configure method!")
	err_Signer_InvalidSigner = errors.New("Self Signer is invalid!")
	err_Signer_InvalidInput = errors.New("Method input is invalid!")
)

var (
	module_super_secret = []byte("OKu2a1LXpmRoKrZS")
)

type ApiHandlers interface {
	Login(*fasthttp.RequestCtx)
	HmacTest(*fasthttp.RequestCtx)
}
type ApiCore struct {
	slogger *dlog.Logger
	users *users.Users
	signer *apiSigner
	Handlers ApiHandlers
}
func InitModule( ctx context.Context ) ( *ApiCore, error ) {
	var e error
	var core *ApiCore = new(ApiCore)
	var apictx context.Context = context.WithValue( context.Background(), appctx.CTX_MOD_APICORE, core )

	// Log testing:
	core.slogger = &dlog.Logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		Ch_message: ctx.Value(appctx.CTX_LOG_FILE).(*dlog.FileLogger).Mess_queue,
		Prefix: dlog.LPFX_MODCORE,
	}
	core.slogger.W( dlog.LLEV_DBG, "TEST TEST TEST FROM MODULE" )

	// Submodules initialization:
	core.signer, e = new(apiSigner).configure( module_super_secret, core.slogger ) // slogger ONLY for debugging! (XXX)
	if e != nil { return nil,e }

	//Handlers initialization:
	core.Handlers, e = new(apiHandler).configure( apictx ); if e != nil { return nil,e }

	return core,nil
}
func (self *ApiCore) DeInitModule() error {
	self.slogger.W( dlog.LLEV_OK, "YEAH! Module Exit!" )
	return nil
}
func (self *ApiCore) GetUsers() ( *users.Users ) {
	return self.users
}
