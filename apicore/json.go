package apicore

import "time"
import "encoding/json"

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"

const (	// var group only for Error identificators:
	api_errno_UndefinedUserId = uint8(iota)
)
var errorMessages []string = []string{	// var for Error messages:
	"Неверный идентификатор пользователя!",	// api_error_UndefinedUserId
}

type message struct {
	id uint64
	err string
	errno uint8
	timestamp int64
}
type apiJsoner struct {
	*ApiCore
}
func (self *apiJsoner) configure(ctx context.Context) ( *apiJsoner, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*ApiCore)
	self.slogger.W( log.LLEV_DBG, "Jsoner submodule has been initialized and configured!" )
	return self,nil
}
func (self *apiJsoner) failureMessage( respid uint64, err uint8 ) ( []byte, error )  {
	return json.Marshal(&message{	// &message OR message ???
		id: respid,
		errno: err,
		err: errorMessages[err],
		timestamp: time.Now().Unix(),
	})
}
func (self *apiJsoner) generalMessage() {}
