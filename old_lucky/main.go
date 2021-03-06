package main

import (
	"time"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"errors"
)
import dlog "doshelpv2/log"
import "github.com/gorilla/mux"

const (
	ERR_SQL_CONNFAIL = "Connection status is not good!"
	ERR_SQL_NOUSERS = "Could not find any user with "
	ERR_MDL_USERFAIL = "Sorry, but we could not identify you! =(\nTry again later"
	ERR_MDL_HASHFAIL = "Sorry, but we could not handle your query data!\nTry again later"
	ERR_USER_UUIDEMP = "User's UUID is empty! Logical error!"
	ERR_USER_HWIDEMP = "User's HWID is empty! Logical error!"
)
var (
	ERR_USER_NOUUID = errors.New("User's UUID is empty! Logical error!")
	ERR_MAIN_NOPARAM = errors.New("Received empty params! Function ferror!")
	ERR_DDOS_REJECTED = errors.New("Too many requests from unique client! User has been banned!")
	ERR_DDOS_BANNED = errors.New("Client is banned! Try again later.")
	ERR_SOCK_DEAD = errors.New("Current listener is dead!")
)


// Application config:
//	appNetProto - application network protocol (UNIX tested only);
//	appNetPath - application UNIX socket full path (Recomended /var/run for production);
//	appLogPath - application path for storing log file;
//	appLogBuf - application max log query (Imp: if log query is overloaded - message would be lost);
//	appDosBanTime - application ban time for users;
//	appDosReqTime - application required interval between requests;
const (
	appNetProto string = "unix"
	appNetPath string = "./doshelpv2.sock"
	appLogPath string = "./doshelpv2.log"
	appLogBuf int = 128
	appDosBanTime time.Duration = 60 * time.Second
	appDosReqTime time.Duration = 3 * time.Second
)


///		WARNING!!!!
// APP DOES NOT SUPPORT URL REWRITE!!!!
///		WARNING!!!!

// Change mux? https://godoc.org/github.com/husobee/vestigo#CustomNotFoundHandlerFunc
//	NO: remove mux! We have only one route!!
// mysql relations http://stackoverflow.com/questions/260441/how-to-create-relationships-in-mysql
// siege performance testing

// DONE:
// + removing using local app vars (Now, we have application variable for it);
// + remove data races ( 2xDataRaces in Socket Alive & client buf destroy - 25-Mar-17);
// + removing global defining in structs;
// + added steam openid support (long task);

// TODO:
// ? check & remove data races (regularly);
// ? remove log defines (Now we have helper function for it);
// - adding sql support for writing buffer values;
// - rewriting log prefix (It's hard to read current logfile);
// - adding writing hw key in "banned" messages (for future logfile grep);
// - replace current http methods in Interface;
// - remove mux import (We have only one route);
// - adding P2P supporting for buffer synchronization between running apps (can be not only tox);
// - adding full app restarting with saving sockets for application deploy;
// - adding minimal log level reporting in application config;
// - added application customization over flags or configuration;
// - replace net/http with fasthttp;
// - replace default http router with fasthttp router;
// - adding rpc connection support between server and client (golang and js);
// - rename all strcut method's variables

func main() {
	newApp()
	defer func() { application.destroy(); application.Wait() }()

	var sgn = make( chan os.Signal )
	signal.Notify( sgn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT )
	application.slogger.W( dlog.LLEV_OK, "Kernel signal catcher has been initialized!")

	go application.threadHTTPD()
	go application.rpcServe()
	go application.apiServe()
	application.slogger.W(dlog.LLEV_OK, "Application started!")

	for {
		select {
		case <-sgn:
			application.slogger.W( dlog.LLEV_ERR, "Catched QUIT signal from kernel! Stopping prg...")
			return
		}
	}
}


type httpRouter struct {
	*mux.Router
	lgRoot *dlog.Logger
	lgNotfound *dlog.Logger
	lgUserManage *dlog.Logger
}
func ( hr *httpRouter ) middleUserManage( next http.Handler ) http.Handler {
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request ) {
		_, cooks, e := newClient(&r.Header); if e != nil {
			hr.lgUserManage.W( dlog.LLEV_WRN, e.Error() )
			http.Error( w, ERR_MDL_USERFAIL, http.StatusInternalServerError )
			return
		}
		for _, ck := range cooks {
			hr.lgUserManage.W( dlog.LLEV_DBG, ck.String() )
			http.SetCookie(w,ck) // - Disable only for testing golucky project
		}
		next.ServeHTTP(w,r)
	})
}
func (self *httpRouter) webSteamOpenid( w http.ResponseWriter, r *http.Request ) {
	oid := _steamOpenID(r)

	switch oid.Mode() {
	case "":
		http.Redirect( w, r, oid.AuthUrl(), 301 )
	case "cancel":
		w.Write([]byte("Authorization cancelled!"))
	default:
		if steamid, e := oid.ValidateAndGetId(); e != nil {
			http.Error( w, e.Error(), http.StatusInternalServerError )
		} else { w.Write([]byte(steamid)) }
	}
}
func ( hr *httpRouter ) webRoot( w http.ResponseWriter, r *http.Request ) {}
func ( hr *httpRouter ) webNotFound( w http.ResponseWriter, r *http.Request ) {
	w.Write( []byte("Sorry, but this page has been killed!") )
}
