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
	err_Init_InvalidInputData = errors.New("Some data from module input is corrupted!")
)

var (  // Must be in some "config struct" in module input context! (See InitModule method)
	module_super_secret = []byte("OKu2a1LXpmRoKrZS")
)

type ApiHandlers interface {
	Login(*fasthttp.RequestCtx)

	// apicore_v1
	CentrifugoConnection(*fasthttp.RequestCtx)
}
type ApiCore struct {
	slogger *dlog.Logger
	users *users.Users
	signer *apiSigner
	jsoner *apiJsoner

	Handlers ApiHandlers

	sign_secret []byte
}
func InitModule( ctx context.Context ) ( *ApiCore, error ) {
	var e error
	var core *ApiCore = new(ApiCore)
	var apictx context.Context = context.WithValue( context.Background(), appctx.CTX_MOD_APICORE, core )

	// Log initialization:
	core.slogger = &dlog.Logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		Ch_message: ctx.Value(appctx.CTX_LOG_FILE).(*dlog.FileLogger).Mess_queue,
		Prefix: dlog.LPFX_MODCORE,
	}
	core.slogger.W( dlog.LLEV_DBG, "ApiCore module initialization has been complete!" )

	// Validate input "config" ( In future, this must be as config struct in context; Not as global variable!!! ):
	if len(module_super_secret) == 0 { return nil,err_Init_InvalidInputData }

	// Submodules initialization:
	if core.signer, e = new(apiSigner).configure(apictx); e != nil { return nil,e }
	if core.jsoner, e = new(apiJsoner).configure(apictx); e != nil { return nil,e }
	if core.Handlers, e = new(apiHandler).configure(apictx); e != nil { return nil,e } // Must be in end of list

	return core,nil
}
func (self *ApiCore) DeInitModule() error {
	self.slogger.W( dlog.LLEV_OK, "YEAH! Module Exit!" )
	return nil
}
func (self *ApiCore) GetUsers() ( *users.Users ) {
	return self.users
}
