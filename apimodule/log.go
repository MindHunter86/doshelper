package apimodule

import "os"
import "log"
import dlog "doshelpv2/log"
import "doshelpv2/appctx"
import "golang.org/x/net/context"

type apiLogger struct {
	*dlog.Logger
}
func (self *apiLogger) configure( ctx context.Context ) ( *apiLogger, error ) {
	ctxLogFlQueue := ctx.Value(appctx.CTX_LOG_FILE).(*dlog.FileLogger).Mess_queue
	if ctxLogFlQueue == nil { return nil,err_Init_InvalidCtxPointer }

	self.Logger = new(dlog.Logger)
	self.Logger.Logger = log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds)
	self.Ch_message = ctxLogFlQueue
	self.Prefix = dlog.LPFX_MODAPI

	self.W( dlog.LLEV_DBG, "Module's log support has been started." )
	return self,nil
}
