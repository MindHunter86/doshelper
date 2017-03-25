package main

import (
	"time"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"errors"
)
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

// TODO:
// - remove log defines (Now we have helper function for it);
// - adding sql support for writing buffer values;
// - rewriting log prefix (It's hard to read current logfile);
// - adding writing hw key in "banned" messages (for future logfile grep);
// - replace current http methods in Interface;
// - remove mux import (We have only one route);
// - adding P2P supporting for buffer synchronization between running apps;


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
			app.stdout_logger.wr( LLEV_ERR, "Catched QUIT signal from kernel! Stopping prg...")
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
		_, cooks, e := newClient(&r.Header); if e != nil {
			hr.lgUserManage.wr( LLEV_WRN, e.Error() )
			http.Error( w, ERR_MDL_USERFAIL, http.StatusInternalServerError )
			return
		}
		for _, ck := range cooks {
			hr.lgUserManage.wr( LLEV_DBG, ck.String() )
			http.SetCookie(w,ck)
		}
		next.ServeHTTP(w,r)
	})
}
func ( hr *httpRouter ) webRoot( w http.ResponseWriter, r *http.Request ) {
	hr.lgRoot.wr( LLEV_INF, "WebRoot")

}
func ( hr *httpRouter ) webNotFound( w http.ResponseWriter, r *http.Request ) {
	w.Write( []byte("Sorry, but page was killed!") )
}
