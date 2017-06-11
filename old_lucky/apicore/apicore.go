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
	err_glob_InvalidInputData = errors.New("Some input data is corrupted!")
	err_ReqStore_StoreAlreadyDefined = errors.New("Request Store has been already defined! You must free it before!")
	err_ReqStore_StoreIsUndefined = errors.New("Request Store is undefined!")
	err_ReqStore_ReqPutError = errors.New("Could not put request in store!")
	err_DB_DdAlreadyDefined = errors.New("sql.DB ahs been already defined!")
	err_DB_InvalidDb = errors.New("Invalid DB pointer in method!")
)

var (  // Must be in some "config struct" in module input context! (See InitModule method); (in main.go):
	module_super_secret = []byte("OKu2a1LXpmRoKrZS")
)

type ApiHandlers interface {
	Login(*fasthttp.RequestCtx)

	// apicore_v1
	CentrifugoConnection(*fasthttp.RequestCtx)
}
type ApiMiddlewares interface {
	QueryIdentificatorAdd(fasthttp.RequestHandler) fasthttp.RequestHandler
}
type ApiCore struct {
	slogger *dlog.Logger
	users *users.Users
	signer *apiSigner
	jsoner *apiJsoner
	database *apiDatabase
	// add request store && don;t forget to destroy memory store in DeInit function! 

	Handlers ApiHandlers
	Middlewares ApiMiddlewares

	// Must be in config struct
	sign_secret string
	sql_addr string
	sql_username, sql_password, sql_database string
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

	// !!! Must be in config struct: (in main.go)
	core.sign_secret = string(module_super_secret)
	core.sql_addr = "127.0.0.1:23306"
	core.sql_username = "golucky"
	core.sql_password = "wR3Lzp27jCh0KAMV"
	core.sql_database = "doshelpv2"

	// Validate input "config" ( In future, this must be as config struct in context; Not as global variable!!! ):
	if len(module_super_secret) == 0 { return nil,err_glob_InvalidInputData }
	if len(core.sign_secret) == 0 { return nil,err_glob_InvalidInputData }
	if len(core.sql_addr) == 0 { return nil,err_glob_InvalidInputData }
	if len(core.sql_username) == 0 || len(core.sql_password) == 0 { return nil,err_glob_InvalidInputData }
	if len(core.sql_database) == 0 { return nil,err_glob_InvalidInputData }

	// Submodules initialization:
	if core.signer, e = new(apiSigner).configure(apictx); e != nil { return nil,e }
	if core.jsoner, e = new(apiJsoner).configure(apictx); e != nil { return nil,e }
	if core.database, e = new(apiDatabase).configure(apictx); e != nil { return nil,e }

	// Exports initialization:
	if core.Middlewares, e = new(apiMiddleware).configure(apictx); e != nil { return nil,e }
	if core.Handlers, e = new(apiHandler).configure(apictx); e != nil { return nil,e } // Must be in end of list

	core.slogger.W( dlog.LLEV_DBG, "ApiCore module initialization has been complete!" )
	return core,nil
}
func (self *ApiCore) DeInitModule() error {
	var errs string
	if e := self.database.destroyConnection(); e != nil { errs += e.Error() + " " }

	if len(errs) != 0 { return errors.New("Some problems with de-initialization ApiCore module: " + errs) }
	self.slogger.W( dlog.LLEV_OK, "ApiCore module has been successfully destroyed!" )
	return nil
}
func (self *ApiCore) GetUsers() ( *users.Users ) {
	return self.users
}
