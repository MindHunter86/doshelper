package main

import "os"
import "sync"
import "net/http"
import "net/http/pprof"
import "log"

import "github.com/gorilla/mux"
//	Use CORS from here???
//	import "github.com/gorilla/handlers"

var application *app
type app struct {
	sync.WaitGroup
	clients *activeClients
	socket *sockListener
	rpc *rpcService
	flogger *fileLogger
	slogger *logger
}

func newApp() {
	var e error
	application = &app{ clients: &activeClients{} }

	if e = application.newFileLogger( appLogPath, appLogBuf ); e != nil { log.Fatalln("Application INIT problem:", e); return }
	if application.socket, e = newSockListener( appNetProto, appNetPath ); e != nil { log.Fatalln("Application INIT problem:", e); return }
	if application.rpc, e = _rpcService(); e != nil { log.Fatalln("Application INIT problem:", e); return }

	application.slogger = application.newLogger(LPFX_CORE)
	application.flogger.start()
	application.clients.init()
}
func ( a *app ) destroy() {
	a.rpc.killListener() // close all rpc connections
	a.socket.stop() // close all sockets, break http listen
	a.clients.destroy() // clean clients buffer, writing all data in SQL (in future)
	a.Wait() // Goroutines "Workers" waiting

	a.flogger.stop() // stop file logger goroutine and wait it's closing
	a.flogger.Wait() // Log buffer waiting
}

func ( self *app ) rpcServe() {
	self.Add(1)

	l := self.newLogger(LPFX_RPC)
	l.w( LLEV_OK, "RPC goroutine has been inited!" )

	l.w( LLEV_INF, "Starting RPC serving ..." )
	for i := uint8(0); i < uint8(4); i++ {
		if e := self.rpc.serve(); e != nil {
			l.w( LLEV_WRN, "Pre-fail state! rpcServe error:" + e.Error() )
			l.w( LLEV_INF, "Trying to restart rpcServing ..." )
			continue
		}
		l.w( LLEV_OK, "RPC serving has been stopped!" )
		break
	}

	l.w( LLEV_OK, "RPC goroutine has been destroyed!" )
	self.Done()
}
func ( a *app ) threadHTTPD() {
	a.Add(1)
	l := a.newLogger(LPFX_HTTPD)
	l.w( LLEV_OK, "HTTPD goroutine has been inited!")

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

	l.w( LLEV_INF, "Starting HTTP serving ...")
	for i := uint8(0); i < uint8(4); i++ {
		if e := a.socket.serveHTTP( hr.Router ); e != nil {
			l.w( LLEV_WRN, "Pre-fail state! HTTPServe error: " + e.Error())
			l.w( LLEV_INF, "Trying to restart HTTPServing ...")
			continue;
		}
		l.w( LLEV_OK, "HTTP serving has been stopped!")
		break;
	}

	l.w( LLEV_OK, "HTTPD goroutine has been destroyed!")
	a.Done()
}
// 2DELETE
//func final(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("2345678"))
//}
func ( a *app ) newLogger( prefix uint8 ) *logger {
	return &logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		ch_message: a.flogger.mess_queue,
		prefix: prefix,
	}
}
func ( a *app ) newFileLogger( fpath string, logbuf int ) ( error ) {
	fd, e := os.OpenFile( fpath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600 )
	if e != nil { return e }

	a.flogger = &fileLogger{
		Logger: log.New( fd, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		mess_queue: make( chan string, logbuf ),
		stop_handle: make( chan bool ),
	}
	return nil
}
func ( a *app ) newHttpRouter() *httpRouter {
	return &httpRouter{
		Router: mux.NewRouter(),
		lgRoot: a.newLogger(LPFX_WEBROOT),
		lgNotfound: a.newLogger(LPFX_NOTFOUND),
		lgUserManage: a.newLogger(LPFX_USERMANAGE),
	}
}
