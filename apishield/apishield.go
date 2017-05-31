package apishield

import "os"
import "errors"
import golog "log"

import "doshelpv2/log"
import "doshelpv2/appctx"
import "golang.org/x/net/context"


var (
	err_mod_InvalidSelf = errors.New("Invalid self struct in configure method! Is self initialized?")
	err_mod_InvalidContext = errors.New("Invalid input context in configure method!")
	err_mod_InvalidInputData = errors.New("Some input data is corrupted!")
)


type ApiShieldModule interface {
	Configure(context.Context) (*ApiShield, error)
	Destroy() error
}
type ApiShield struct { // main module struct
	log *log.Logger
}
func (self *ApiShield) Configure(ctx context.Context) (*ApiShield, error) {
	if self == nil { return nil,err_mod_InvalidSelf }
	if ctx == nil { return nil,err_mod_InvalidContext }

	// module log subsystem initialization:
	self.logConfigure(ctx.Value(appctx.CTX_LOG_FILE).(*log.FileLogger))

	self.log.W(log.LLEV_OK, "ApiShield module has been initialized and configured!")
	return self,nil
}
func (self *ApiShield) Destroy() error {
	self.log.W(log.LLEV_OK, "ApiShield module has been destroyed!")
	return nil
}

func (self *ApiShield) logConfigure(fl *log.FileLogger) {
	self.log = &log.Logger {
		Logger: golog.New(os.Stdout, "", golog.Ldate | golog.Ltime | golog.Lmicroseconds),
		Ch_message: fl.Mess_queue,
		Prefix: log.LPFX_MODSHIELD,
	}
}
