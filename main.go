package main

import (
	"sync"

	"database/sql"

	"net/http"

	"os"
	"os/signal"
	"syscall"
)
import "github.com/gorilla/mux"
import _ "github.com/go-sql-driver/mysql"

//const (
//	ERR_
//)







///		WARNING!!!!
// APP DOES NOT SUPPORT URL REWRITE!!!!
///		WARNING!!!!

// Change mux? https://godoc.org/github.com/husobee/vestigo#CustomNotFoundHandlerFunc




var app *App

type App struct {
	sync.WaitGroup
	Users *UserBuf
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
func ( a *App ) ThreadHTTPD() {
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
		l.PutOK("HTTP serving has been stopped!")
		break;
	}

	l.PutOK("HTTPD goroutine has been destroyed!")
	a.Done()
}


func main() {
	const (
		appNetProto string = "unix"
		appNetPath string = "/run/doshelpv2.sock"
		appSqlHost string = "lambda.mh00.net:23306"
		appSqlUser string = "doshelper"
		appSqlPass string = "pyIF236NBfZLXUu1"
		appSqlDb string = "doshelpv2"
	)

	l := NewLogger( LPFX_CORE )
	l.PutOK("Core log system has been inited!")

	app, e := NewApp( appNetProto, appNetPath ); if e != nil {
		l.PutNon("Could not create App!")
		l.PutErr(e.Error()); return
	} else { l.PutOK("App has been created!") }
	defer func() {
		app.Wait()
		app.Destroy()
		l.PutOK("App has been destroyed!")
	}()

	var sgn = make( chan os.Signal )
	signal.Notify( sgn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT )
	l.PutOK("Kernel signal catcher has been initialised!")

	go app.ThreadHTTPD()

	for {
		select {
		case <-sgn:
			l.PutWrn("Catched QUIT signal from kernel! Stopping prg...")
			app.Socket.Close()
			return
		}
	}
}


type sqlClient struct {
	dbconn *sql.DB
}
// host port user pass db string
func newSqlClient( h,u,p,d string ) ( *sqlClient, error ) {
	db, e := sql.Open( "mysql", u + ":" + p + "@" + h + "/" + d ); if e != nil { return nil,e }
	if e := db.Ping(); e != nil { return nil,e }
	return &sqlClient{
		dbconn: db,
	}, nil
}
func ( sc *sqlClient ) connCheck() error {
	return sc.dbconn.Ping()
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

	u, ok := getOrCreateUser( r, app.Users )
	switch ok {
	case false:
		u_c := u.ParseOrCreateUUID(); if u_c == nil {
			hr.wroot_l.PutNon("Could not generate UUID for user!!! Something error!");
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.SetCookie( w, u_c )
		hr.wroot_l.PutInf( "New user: " + u.Uuid )
	case true:
	// MAKE SOME SECURE CHECKS ( VALIDATE ALL COOKIES!!!! )
		hr.wroot_l.PutInf( "User: " + u.Uuid )
	}

// 	if hwid_c, e := u.GenHWID(); e != nil {
// 		ht.wroot_l.PutNon( "Colud not set User's HWID cookie for " + u.Uuid )
// 		w.Write( []byte("Sorry, but you're a bot =(") )
// 		w.WriteHeader(http.StatusTeapot)
// 		return
// 	}

	c, e := u.GenSecureHash(); if e != nil {
	// Working with SESSION cookie LIFETIME ????
		hr.wroot_l.PutNon( "Could not set User's SL cookie for " + u.Uuid )
		w.Write( []byte("Sorry, but you are a bot =(") )	// If SL cookie is set - you are bot. NoNOK!
		w.WriteHeader(http.StatusTeapot)
		return
	}

	http.SetCookie( w, c )
	http.Redirect( w, r, r.Header.Get("Origin"), 301 )
}

