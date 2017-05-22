package main
import (
	"log"
	"sync"
)


const (
	LPFX_CORE = iota
	LPFX_HTTPD
	LPFX_WEBROOT
	LPFX_NOTFOUND
	LPFX_USERMANAGE
	LPFX_RPC
	LPFX_API
)
const (
	LLEV_DBG = iota
	LLEV_OK
	LLEV_NON
	LLEV_INF
	LLEV_WRN
	LLEV_ERR
)

type logger struct {
	*log.Logger
	ch_message chan string
	sync.RWMutex
	prefix uint8
}


func ( l *logger ) setPrefix( p uint8 ) {
	l.Lock(); l.prefix = p; l.Unlock()
}
func ( l *logger ) getPrefix( colo bool ) string {
	l.RLock()
	defer l.RUnlock()

	switch l.prefix {
	case LPFX_CORE:
		if colo { return "\x1b[36;1m[CORE]:\x1b[0m" } else { return "[CORE]: " }
	case LPFX_HTTPD:
		if colo { return "\x1b[36;1m[HTTPD]:\x1b[0m" } else { return "[HTTPD]: " }
	case LPFX_WEBROOT:
		if colo { return "\x1b[36;1m[HTTPD-WEBROOT]:\x1b[0m" } else { return "[HTTPD-WEBROOT]: " }
	case LPFX_USERMANAGE:
		if colo { return "\x1b[36;1m[HTTPD-MIDDLEUSER]:\x1b[0m" } else { return "[HTTPD-MIDDLEUSER]: " }
	case LPFX_RPC:
		if colo { return "\x1b[36;1m[RPC]:\x1b[0m" } else { return "[RPC]: " }
	case LPFX_API:
		if colo { return "\x1b[36;1m[API]:\x1b[0m" } else { return "[API]: " }
	default:
		return ""
	}
}
func ( l *logger ) w( lvl uint8, m string ) {
// log to file
	switch lvl {
	case LLEV_DBG:
		l.ch_message <- l.getPrefix(false) + "DBG: " + m
	case LLEV_INF:
		l.ch_message <- l.getPrefix(false) + "INF: " + m
	case LLEV_OK:
		l.ch_message <- l.getPrefix(false) + "SUC: " + m
	case LLEV_NON:
		l.ch_message <- l.getPrefix(false) + "FAI: " + m
	case LLEV_WRN:
		l.ch_message <- l.getPrefix(false) + "WRN: " + m
	case LLEV_ERR:
		l.ch_message <- l.getPrefix(false) + "ERR: " + m
	}

// log to stdout
	switch lvl {
	case LLEV_DBG:
		l.Println( l.getPrefix(true), m )
	case LLEV_OK:
		l.Println( l.getPrefix(true), "\x1b[32;1m✔\x1b[0m", m )
	case LLEV_NON:
		l.Println( l.getPrefix(true), "\x1b[31;1m✖\x1b[0m", m)
	case LLEV_INF:
		l.Println( l.getPrefix(true), "\x1b[34;1m🛈\x1b[0m", m)
	case LLEV_WRN:
		l.Println( l.getPrefix(true), "\x1b[33;1m❢\x1b[0;33;22m", m, "\x1b[0m")
	case LLEV_ERR:
		l.Println( l.getPrefix(true), "\x1b[31;22m❢\x1b[0;31;1m", m, "\x1b[0m")
	}
}


type fileLogger struct {
	sync.WaitGroup
	*log.Logger
	mess_queue chan string
	stop_handle chan bool
}

func ( fl *fileLogger ) start() {
	go func() {
		fl.Add(1)
		defer fl.Done()

		for {
			select{
			case m := <-fl.mess_queue:
				fl.Println(m)
			case <-fl.stop_handle:
				var buf_size uint8 = uint8( len(fl.mess_queue) )
				for ; buf_size != 0; buf_size-- {
					fl.Println(<-fl.mess_queue)
				}
				log.Println("Log worker has been stopped! Buf is empty")
				return
			}
		}
	}()
}
func ( fl *fileLogger ) stop() {
	close(fl.stop_handle)
}
