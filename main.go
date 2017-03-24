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
	appLogPath string = "./doshelpv2.log"
	appLogBuf int = 128
)



///		WARNING!!!!
// APP DOES NOT SUPPORT URL REWRITE!!!!
///		WARNING!!!!

// Change mux? https://godoc.org/github.com/husobee/vestigo#CustomNotFoundHandlerFunc
//	NO: remove mux! We have only one route!!


// mysql relations http://stackoverflow.com/questions/260441/how-to-create-relationships-in-mysql


// siege performance testing



// PROBLEM!!!!!! LOG MUST BE DEFINED BEFORE NEWAPP() FUNCTION
func main() {
	app, ok := newApp(); if !ok {
		os.Exit(1)
	}
	application = app
	app.stdout_logger.wr( LLEV_OK, "CORE log system has been inited!")
	app.stdout_logger.wr( LLEV_OK, "Application started!")

	defer func() {
		app.stdout_logger.wr( LLEV_INF, "Reached app destroy function")
		app.Destroy()
		app.stdout_logger.wr( LLEV_INF, "Reached app wair function")
		app.Wait()
		app.stdout_logger.wr( LLEV_OK, "App has been destroyed!")
		app.file_logger.stop()
		app.file_logger.Wait()
	}()

	var sgn = make( chan os.Signal )
	signal.Notify( sgn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT )
	app.stdout_logger.wr( LLEV_OK, "Kernel signal catcher has been initialized!")

	go app.ThreadHTTPD()

	for {
		select {
		case <-sgn:
			app.stdout_logger.wr( LLEV_WRN, "Catched QUIT signal from kernel! Stopping prg...")
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
func ( hr *httpRouter ) middleUserManage( next http.Handler ) http.Handler {
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request ) {
		_, cooks, e := newClient2(&r.Header); if e != nil {
			hr.lgUserManage.wr( LLEV_WRN, e.Error() )
			http.Error( w, ERR_MDL_USERFAIL, http.StatusInternalServerError )
			return
		}
		for _, ck := range cooks {
			hr.lgUserManage.wr( LLEV_DBG, ck.String() )
			http.SetCookie(w,ck)
		}
		hr.lgUserManage.wr( LLEV_DBG, "BUF DEBUG:" )
		for key, val := range application.clients.clients {
			hr.lgUserManage.wr( LLEV_DBG, string(key) + " " + val.uuid )
			hr.lgUserManage.wr( LLEV_DBG, string(key) + " " + val.sec_link )
			hr.lgUserManage.wr( LLEV_DBG, string(key) + " " + val.addr )
			hr.lgUserManage.wr( LLEV_DBG, string(key) + " " + val.origin  )
		}
		next.ServeHTTP(w,r)
	})
}
//func ( hr *httpRouter ) middleUserManage( next http.Handler ) http.Handler {
//	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request ) {
//
//		cl, uid_c, e := newClient(&r.Header); if e != nil {
//			hr.lgUserManage.wr( LLEV_WRN, e.Error())
//			http.Error( w, ERR_MDL_USERFAIL, http.StatusInternalServerError )
//			return
//		} else if uid_c != nil { http.SetCookie(  w, uid_c ) }
//
//		var proto,host string
//		proto = r.Header.Get("X-Forwarded-Proto")
//		host = r.Header.Get("X-Forwarded-Host")
//
//		hwk_c, hwk, e := cl.getHwKey( r.Header.Get("X-Client-HWID"), proto, host ); if e != nil {
//			hr.lgUserManage.wr( LLEV_WRN, "CL_" + cl.uuid + ": genHW error: " + e.Error() )
//			http.Error( w, ERR_MDL_HASHFAIL, http.StatusInternalServerError )
//			return
//		} else if hwk_c != nil { http.SetCookie( w, hwk_c ) }
//
//		sl_c, e := cl.generateSecLink( r.Header.Get("X-SecureLink-Secret"), proto, host ); if e != nil {
//			hr.lgUserManage.wr( LLEV_WRN, "CL_" + cl.uuid + ": genSL error: " + e.Error() )
//			http.Error( w, ERR_MDL_HASHFAIL, http.StatusInternalServerError )
//			return
//		} else { http.SetCookie( w, sl_c ) }
//
//		hr.lgUserManage.wr( LLEV_INF, hwk)
//
//		next.ServeHTTP(w,r)
//	})
//}
func ( hr *httpRouter ) webRoot( w http.ResponseWriter, r *http.Request ) {
	hr.lgRoot.wr( LLEV_INF, "WebRoot")

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
