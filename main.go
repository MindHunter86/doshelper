package main

import (
	"net/http"

	"os"
	"os/signal"
	"syscall"
)
import "github.com/gorilla/mux"

const (
	ERR_SQL_CONNFAIL = "Connection status is not good!"
	ERR_SQL_NOUSERS = "Could not find any user with "
	ERR_MDL_USERFAIL = "Sorry, but we could not identify you! =(\nTry again later"
	ERR_MDL_HASHFAIL = "Sorry, but we could not handle your query data!\nTry again later"
)

const (
	appNetProto string = "unix"
	appNetPath string = "./doshelpv2.sock"
)



///		WARNING!!!!
// APP DOES NOT SUPPORT URL REWRITE!!!!
///		WARNING!!!!

// Change mux? https://godoc.org/github.com/husobee/vestigo#CustomNotFoundHandlerFunc


// mysql relations http://stackoverflow.com/questions/260441/how-to-create-relationships-in-mysql


// siege performance

func main() {

	l := NewLogger( LPFX_CORE )
	l.PutOK("Core log system has been inited!")

	app, e := NewApp(); if e != nil {
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


type httpRouter struct {
	*mux.Router
	lgRoot *Logger
	lgNotfound *Logger
	lgUserManage *Logger
}
func newHttpRouter() *httpRouter {
	return &httpRouter{
		Router: mux.NewRouter(),
		lgRoot: NewLogger( LPFX_WEBROOT ),
		lgNotfound: NewLogger( LPFX_NOTFOUND ),
		lgUserManage: NewLogger( LPFX_USERMANAGE ),
	}
}
func ( hr *httpRouter ) middleUserManage( next http.Handler ) http.Handler {
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request ) {

		cl, uid_c, e := newClient(&r.Header); if e != nil {
			hr.lgUserManage.PutNon(e.Error())
			http.Error( w, ERR_MDL_USERFAIL, http.StatusPreconditionFailed )
			return
		} else if uid_c != nil { http.SetCookie(  w, uid_c ) }

		var proto,host string
		proto = r.Header.Get("X-Forwarded-Proto")
		host = r.Header.Get("X-Forwarded-Host")

		hwk_c, hwk, e := cl.getHwKey( r.Header.Get("X-Client-HWID"), proto, host ); if e != nil {
			hr.lgUserManage.PutNon( "CL_" + cl.uuid + ": genHW error: " + e.Error() )
			http.Error( w, ERR_MDL_HASHFAIL, http.StatusInternalServerError )
			return
		} else if hwk_c != nil { http.SetCookie( w, hwk_c ) }

		sl_c, e := cl.generateSecLink( r.Header.Get("X-SecureLink-Secret"), proto, host ); if e != nil {
			hr.lgUserManage.PutNon( "CL_" + cl.uuid + ": genSL error: " + e.Error() )
			http.Error( w, ERR_MDL_HASHFAIL, http.StatusInternalServerError )
			return
		} else { http.SetCookie( w, sl_c ) }

		hr.lgUserManage.PutInf(hwk)

		next.ServeHTTP(w,r)
	})
}
func ( hr *httpRouter ) webRoot( w http.ResponseWriter, r *http.Request ) {
	hr.lgRoot.PutInf("WebRoot")

}
func ( hr *httpRouter ) webNotFound( w http.ResponseWriter, r *http.Request ) {
	w.Write( []byte("Sorry, but page was killed!") )
}



// 	u, ok := getOrCreateUser(r)
// 	switch ok {
// 	case false:
// 		u_c := u.ParseOrCreateUUID(); if u_c == nil {
// 			hr.wroot_l.PutNon("Could not generate UUID for user!!! Something error!");
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		http.SetCookie( w, u_c )
// 		hr.wroot_l.PutInf( "New user: " + u.Uuid )
// 	case true:
// 	// MAKE SOME SECURE CHECKS ( VALIDATE ALL COOKIES!!!! )
// 	// WRITE USER IN CACHE!!!
// 		hr.wroot_l.PutInf( "User: " + u.Uuid )
// 	}
// 
// 	hwid_c, e := u.getOrCreateHWID(); if e != nil { hr.wroot_l.PutNon(e.Error()); w.WriteHeader(http.StatusTeapot); return }
// 	hr.wroot_l.PutInf( "User " + u.Uuid + " has HWID - " + hwid_c.Value )
// 	http.SetCookie( w, hwid_c )
// 
// // 	if hwid_c, e := u.GenHWID(); e != nil {
// // 		ht.wroot_l.PutNon( "Colud not set User's HWID cookie for " + u.Uuid )
// // 		w.Write( []byte("Sorry, but you're a bot =(") )
// // 		w.WriteHeader(http.StatusTeapot)
// // 		return
// // 	}
// 
// 	c, e := u.GenSecureHash(); if e != nil {
// 	// Working with SESSION cookie LIFETIME ????
// 		hr.wroot_l.PutNon( "Could not set User's SL cookie for " + u.Uuid )
// 		w.Write( []byte("Sorry, but you are a bot =(") )	// If SL cookie is set - you are bot. NoNOK!
// 		w.WriteHeader(http.StatusTeapot)
// 		return
// 	}
// 
// 	http.SetCookie( w, c )
// 	http.Redirect( w, r, r.Header.Get("Origin"), 301 )
