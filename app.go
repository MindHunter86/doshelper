package main

import "os"
import "sync"
import "net/http"
import "net/http/pprof"
import "log"
import "doshelpv2/appctx"
import dlog "doshelpv2/log"
import "doshelpv2/apimodule"
import "doshelpv2/apicore"

import "github.com/gorilla/mux"
import "golang.org/x/net/context"
//	Use CORS from here???
//	import "github.com/gorilla/handlers"


var application *app
type app struct {
	sync.WaitGroup
	clients *activeClients
	socket *sockListener
	rpc *rpcService
	api *apimodule.ApiModule
	core *apicore.ApiCore
	flogger *dlog.FileLogger
	slogger *dlog.Logger
}

func newApp() {
	var e error
	application = &app{ clients: &activeClients{} }

	if e = application.newFileLogger( appLogPath, appLogBuf ); e != nil { log.Fatalln("Application INIT problem:", e); return }
	application.slogger = application.newLogger(dlog.LPFX_CORE)
	application.flogger.Start()
	application.clients.init()

	// Activate module context buffer:
	var ctx context.Context
	ctx = context.WithValue( context.Background(), appctx.CTX_LOG_STD, application.slogger )
	ctx = context.WithValue( ctx, appctx.CTX_LOG_FILE, application.flogger )

	// Modules initialization:
	if application.socket, e = newSockListener( appNetProto, appNetPath ); e != nil { log.Fatalln("Application INIT problem:", e); return }
	if application.rpc, e = _rpcService(); e != nil { log.Fatalln("Application INIT problem:", e); return }
	if application.api, e = apimodule.InitModule(); e != nil { log.Fatalln("Application INIT problem:", e); return }
	if application.core, e = apicore.InitModule(ctx); e != nil { log.Fatalln("Application INIT problem:", e); return }
}
func ( a *app ) destroy() {
	// PRE PROD KILL ZONE
	a.rpc.killListener() // close all rpc connections
	a.core.DeInitModule() // remove all handlers and data
	a.api.DeInitModule() // close all apimodule connections

	// Test KILL ZONE
	a.socket.stop() // close all sockets, break http listen
	a.clients.destroy() // clean clients buffer, writing all data in SQL (in future)
	a.Wait() // Goroutines "Workers" waiting

	a.flogger.Stop() // stop file logger goroutine and wait it's closing
	a.flogger.Wait() // Log buffer waiting
}

func (self *app) apiServe() {
	self.Add(1)
	l := self.newLogger(dlog.LPFX_MODAPI)
	l.W( dlog.LLEV_OK, "APImodule goroutine has been inited!" )

	l.W( dlog.LLEV_INF, "Starting API serving ..." )
	for i := uint8(0); i < uint8(4); i++ {
		if e := self.api.Serve(); e != nil {
			l.W( dlog.LLEV_WRN, "Pre-fail state! apiServe error:" + e.Error() )
			l.W( dlog.LLEV_INF, "Trying to restart apiServing ..." )
		}
		l.W( dlog.LLEV_OK, "APImodule goroutine has been stopped!" )
		break
	}

	l.W( dlog.LLEV_OK, "APImodule goroutine has been destroyed!" )
	self.Done()
}
func ( self *app ) rpcServe() {
	self.Add(1)

	l := self.newLogger(dlog.LPFX_MODRPC)
	l.W( dlog.LLEV_OK, "RPC goroutine has been inited!" )

	l.W( dlog.LLEV_INF, "Starting RPC serving ..." )
	for i := uint8(0); i < uint8(4); i++ {
		if e := self.rpc.serve(); e != nil {
			l.W( dlog.LLEV_WRN, "Pre-fail state! rpcServe error:" + e.Error() )
			l.W( dlog.LLEV_INF, "Trying to restart rpcServing ..." )
			continue
		}
		l.W( dlog.LLEV_OK, "RPC serving has been stopped!" )
		break
	}

	l.W( dlog.LLEV_OK, "RPC goroutine has been destroyed!" )
	self.Done()
}
func ( a *app ) threadHTTPD() {
	a.Add(1)
	l := a.newLogger(dlog.LPFX_HTTPD)
	l.W( dlog.LLEV_OK, "HTTPD goroutine has been inited!")

	hr := a.newHttpRouter()
	webRootPage := http.HandlerFunc(hr.webRoot)
	webNotFoundPage := http.HandlerFunc(hr.webNotFound)

	hr.Handle( "/", hr.middleUserManage(webRootPage) )
	hr.NotFoundHandler = webNotFoundPage

	// GO TOOL PPROF DEBUGGING:
	hr.HandleFunc( "/debug/pprof/", pprof.Index )
	hr.HandleFunc( "/debug/pprof/cmdline", pprof.Cmdline )
	hr.HandleFunc( "/debug/pprof/profile", pprof.Profile )
	hr.HandleFunc( "/debug/pprof/symbol", pprof.Symbol )
	hr.HandleFunc( "/debug/pprof/trace", pprof.Trace )
	hr.HandleFunc( "/login", hr.webSteamOpenid )

	l.W( dlog.LLEV_INF, "Starting HTTP serving ...")
	for i := uint8(0); i < uint8(4); i++ {
		if e := a.socket.serveHTTP( hr.Router ); e != nil {
			l.W( dlog.LLEV_WRN, "Pre-fail state! HTTPServe error: " + e.Error())
			l.W( dlog.LLEV_INF, "Trying to restart HTTPServing ...")
			continue;
		}
		l.W( dlog.LLEV_OK, "HTTP serving has been stopped!")
		break;
	}

	l.W( dlog.LLEV_OK, "HTTPD goroutine has been destroyed!")
	a.Done()
}
// 2DELETE
//func final(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("2345678"))
//}
func ( a *app ) newLogger( prefix uint8 ) *dlog.Logger {
	return &dlog.Logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		Ch_message: a.flogger.Mess_queue,
		Prefix: prefix,
	}
}
func ( a *app ) newFileLogger( fpath string, logbuf int ) ( error ) {
	fd, e := os.OpenFile( fpath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600 )
	if e != nil { return e }

	a.flogger = &dlog.FileLogger{
		Logger: log.New( fd, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		Mess_queue: make( chan string, logbuf ),
		Stop_handle: make( chan bool ),
	}
	return nil
}
func ( a *app ) newHttpRouter() *httpRouter {
	return &httpRouter{
		Router: mux.NewRouter(),
		lgRoot: a.newLogger(dlog.LPFX_WEBROOT),
		lgNotfound: a.newLogger(dlog.LPFX_NOTFOUND),
		lgUserManage: a.newLogger(dlog.LPFX_USERMANAGE),
	}
}
