package main

import "sync"
import "net/http"

//	Use CORS from here???
//	import "github.com/gorilla/handlers"

var app *App
type App struct {
	sync.WaitGroup
	Socket *SockListener
}

func NewApp() ( *App, error ) {
	sl, e := NewSockListener( appNetProto, appNetPath ); if e != nil {
		return nil,e
	}
	return &App{
		Socket: sl,
	}, nil
}
func ( a *App ) Destroy() {
	a.Socket.Close()
}
func ( a *App ) ThreadHTTPD() {
	a.Add(1)
	l := NewLogger( LPFX_HTTPD )
// 	logFile, _ := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	l.PutOK("HTTPD goroutine has been inited!")

	hr := newHttpRouter()
	webRootPage := http.HandlerFunc(hr.webRoot)
	webNotFoundPage := http.HandlerFunc(hr.webNotFound)

	hr.Handle( "/", hr.middleUserManage(webRootPage) )
	hr.NotFoundHandler = webNotFoundPage


// 	r := NewHTTPRouter(l)
// 	finalHandler := http.HandlerFunc( r.WebRoot )
// 	r.Router.Handle( "/", handlers.LoggingHandler( logFile, finalHandler ) )
// 	r.Router.NotFoundHandler = finalHandler

	l.PutInf("Starting HTTP serving ...")
	for i := uint8(0); i < uint8(4); i++ {
		if e := a.Socket.HTTPServe( hr.Router ); e != nil {
			l.PutWrn("Pre-fail state! HTTPServe error: " + e.Error())
			l.PutInf("Trying to restart HTTPServing ...")
			continue;
		}
		l.PutOK("HTTP serving has been stopped!")
		break;
	}

	l.PutOK("HTTPD goroutine has been destroyed!")
	a.Done()
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("2345678"))
}
