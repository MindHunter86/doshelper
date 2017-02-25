package main

import (
	"sync"
	"errors"

	"net"
	"net/http"

	"os"
	"os/signal"
	"syscall"
)
import "github.com/gorilla/mux"



type SockListener struct {
	net.Listener
	alive bool
}

func NewSockListener( proto, path string ) ( *SockListener, error ) {
	s, e := net.Listen( proto, path ); if e != nil {
		return nil,e
	}
	return &SockListener{
		Listener: s,
		alive: true,
	}, nil
}

func ( sl *SockListener ) HTTPServe( r *mux.Router ) error {
	if sl.alive == false { return errors.New("Listener is dead!") }

	if e := http.Serve( sl.Listener, r ); e != nil {
		switch sl.alive {
		case true:
			return e
		default:
			return nil
		}
	}
	return nil
}
func ( sl *SockListener ) File() ( *os.File, error ) {
	return sl.Listener.(*net.UnixListener).File()
}
func ( sl *SockListener ) Close() error {
	if sl.alive == false { return errors.New("Listener is dead!") }
	sl.alive = false
	return sl.Listener.Close()
}

type App struct {
	sync.WaitGroup
	Socket *SockListener
}

func NewApp( netproto, netpath string ) ( *App, error ) {
	s, e := NewSockListener( netproto, netpath ); if e != nil {
		return nil,e
	}
	return &App{
		Socket: s,
	}, nil
}
func ( a *App ) Destroy() {
	a.Socket.Close()
}

func main() {
	const (
		appNetProto string = "unix"
		appNetPath string = "./dostest.sock"
	)

	l := NewLogger( LPFX_CORE )
	l.PutOK("Core log system has been inited!")

//	l.PutInf("Creating new app ...")
	app, e := NewApp( appNetProto, appNetPath ); if e != nil {
		l.PutNon("Could not create App!")
		l.PutErr(e.Error()); return
	} else { l.PutOK("App has been created!") }
	defer func() {
		app.Wait()
		app.Destroy()
		l.PutOK("App has been destroyed!")
	}()

//	l.PutInf("Kernel signal catcher initialisation ...")
	var sgn = make( chan os.Signal )
	signal.Notify( sgn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT )
	l.PutOK("Kernel signal catcher has been initialised!")


//	l.PutInf("Starting HTTP goroutine ...")
	go func( a *App ) {
		a.Add(1)
		l := NewLogger( LPFX_HTTPD )
		l.PutOK("HTTPD goroutine has been inited!")

		r := NewHTTPRouter(l)
		r.Router.HandleFunc( "/", r.WebRoot )

		l.PutInf("Starting HTTP serving ...")
		for i := uint8(0); i < uint8(4); i++ {
			if e := a.Socket.HTTPServe( r.Router ); e != nil {
				l.PutWrn("Pre-fail state! HTTPServe error: " + e.Error())
				l.PutInf("Trying to restart HTTPServing ...")
				continue;
			}
			l.PutOK("HTTP serving has been stoped!")
			break;
		}

		l.PutOK("HTTPD goroutine has been destroyed!")
		a.Done()
	}( app )

	for {
		select {
		case <-sgn:
			l.PutWrn("Catched QUIT signal from kernel! Stopping prg...")
			app.Socket.Close()
			return
		}
	}
}

type HTTPRouter struct {
	Router *mux.Router
	wroot_l *Logger
}
func NewHTTPRouter( l *Logger ) *HTTPRouter {
	return &HTTPRouter{
		Router: mux.NewRouter(),
		wroot_l: NewLogger( LPFX_WEBROOT ),
	}
}
func ( hr *HTTPRouter ) WebRoot( w http.ResponseWriter, r *http.Request ) {
	hr.wroot_l.PutInf("New connection!")
	return
}
