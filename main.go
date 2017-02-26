package main

import (
	"sync"
	"errors"

	"bytes"
	"strings"
	"crypto/md5"
	"encoding/base64"

	"net/http"

	"os"
	"os/signal"
	"syscall"
)
import "github.com/gorilla/mux"
import gouuid "github.com/satori/go.uuid"


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
			l.PutOK("HTTP serving has been stopped!")
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

	u := NewUser(r)
	if u_c := u.ParseOrCreateUUID(); u_c != nil {
		hr.wroot_l.PutInf( "New user: " + u.Uuid )
		http.SetCookie( w, u_c )
	} else { hr.wroot_l.PutInf( "User: " + u.Uuid ) }

	c, e := u.GenSecureHash(); if e != nil {
		hr.wroot_l.PutNon( "Could not set User's SL cookie for " + u.Uuid )
		w.Write( []byte("Sorry, but you are bot =(") )	// If SL cookie is set - you are bot. NoNOK!
		w.WriteHeader(http.StatusTeapot)
		return
	}

	http.SetCookie( w, c )
// RETURN TO ORIGIN!
	w.Write( []byte("OK") )
}


type User struct {
	req *http.Request
	Uuid, secure_hash string
}
func NewUser( r *http.Request ) *User {
	u := &User{ req: r }
	return u
}
func ( u *User ) ParseOrCreateUUID() *http.Cookie {
	uuid_c, _ := u.req.Cookie("uuid")

	if len( uuid_c.String() ) <= 0 {
		u.Uuid = gouuid.NewV4().String()

		var https bool = false
		switch u.req.URL.Scheme {
		case "https":
			https = true
		}

		return &http.Cookie{
			Name: "uuid",
			Value: u.Uuid,
			Path: "/",
			Domain: u.req.URL.Host,
			Secure: https,
			HttpOnly: true,
		}
	}
	u.Uuid = uuid_c.Value
	return nil
}
func ( u *User ) GetSecureHash() string {
	sh, e := u.req.Cookie("sl"); if e != nil { return "" }
	return sh.Value
}
func ( u *User ) GenSecureHash() ( *http.Cookie, error ) {
// buf - $remote_addr:$cookie_uuid:$user_agent:$secret
	if len(u.Uuid) <= 0 { return nil,errors.New("User's UUID is empty! Logical error!") }
	if len( u.GetSecureHash() ) != 0 { return nil,errors.New("SL cookie was already defined!") }

	var buf bytes.Buffer
	buf.WriteString( u.req.Header.Get("X-Forwarded-For") + u.Uuid + u.req.UserAgent() )
	buf.WriteString( u.req.Header.Get("X-SecureLink-Secret") )

	t1 := md5.Sum( buf.Bytes() )
	t2 := base64.StdEncoding.EncodeToString( t1[:] )
	t3 := strings.Replace( strings.Replace( t2, "+", "-", -1 ), "/", "_", -1 )
	t4 := strings.Replace( t3, "=", "", -1 )

	var https bool = false
	switch u.req.URL.Scheme {
	case "https":
		https = true
	}

	return &http.Cookie{
		Name: "sl",
		Value: t4,
		Path: "/",
		Domain: u.req.URL.Host,
		Secure: https,
		HttpOnly: true,
	}, nil
}
