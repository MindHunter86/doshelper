package main

import "os"
import "sync"
import "net/http"
import "log"

import "github.com/gorilla/mux"
//	Use CORS from here???
//	import "github.com/gorilla/handlers"

var app *App
type App struct {
	sync.WaitGroup
	clients *activeClients
	Socket *SockListener
	file_logger *fileLogger
}

	appLogPath string = "./app.log"
	appLogBuf int = 128
func NewApp() ( *App, error ) {
	sl, e := NewSockListener( appNetProto, appNetPath ); if e != nil {
		return nil,e
	}
	return &App{
		clients: &activeClients{},
		Socket: sl,
	}, nil
}
func ( a *App ) Destroy() {
	a.Socket.Close()
}
func ( a *App ) ThreadHTTPD() {
	a.Add(1)
	l := a.newLogger(LPFX_HTTPD)
	l.wr( LLEV_OK, "HTTPD goroutine has been inited!")

	hr := a.newHttpRouter()
	webRootPage := http.HandlerFunc(hr.webRoot)
	webNotFoundPage := http.HandlerFunc(hr.webNotFound)

	hr.Handle( "/", hr.middleUserManage(webRootPage) )
	hr.NotFoundHandler = webNotFoundPage


// 	logFile, _ := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
// 	r := NewHTTPRouter(l)
// 	finalHandler := http.HandlerFunc( r.WebRoot )
// 	r.Router.Handle( "/", handlers.LoggingHandler( logFile, finalHandler ) )
// 	r.Router.NotFoundHandler = finalHandler


	l.wr( LLEV_INF, "Starting HTTP serving ...")
	for i := uint8(0); i < uint8(4); i++ {
		if e := a.Socket.HTTPServe( hr.Router ); e != nil {
			l.wr( LLEV_WRN, "Pre-fail state! HTTPServe error: " + e.Error())
			l.wr( LLEV_INF, "Trying to restart HTTPServing ...")
			continue;
		}
		l.wr( LLEV_OK, "HTTP serving has been stopped!")
		break;
	}

	l.wr( LLEV_OK, "HTTPD goroutine has been destroyed!")
	a.Done()
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("2345678"))
}

func ( a *App ) newLogger( prefix uint8 ) *Logger {
	return &Logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		ch_message: a.file_logger.mess_queue,
		prefix: prefix,
	}
}

func ( a *App ) newFileLogger( fpath string, logbuf int ) ( *fileLogger, error ) {
	fd, e := os.OpenFile( fpath, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0600 )
	if e != nil { return nil,e }

	return &fileLogger{
		Logger: log.New( fd, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		mess_queue: make( chan string,  ),
		stop_handle: make( chan bool ),
	},nil
}


func ( a *App ) newHttpRouter() *httpRouter {
	return &httpRouter{
		Router: mux.NewRouter(),
		lgRoot: a.newLogger(LPFX_WEBROOT),
		lgNotfound: a.newLogger(LPFX_NOTFOUND),
		lgUserManage: a.newLogger(LPFX_USERMANAGE),
	}
}
