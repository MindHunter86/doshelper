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
)
const (
	LLEV_DBG = iota
	LLEV_OK
	LLEV_NON
	LLEV_INF
	LLEV_WRN
	LLEV_ERR
)

type Logger struct {
	*log.Logger
	ch_message chan string
	sync.RWMutex
	prefix uint8
}


func ( l *Logger ) SetPrefix( p uint8 ) {
	l.Lock(); l.prefix = p; l.Unlock()
}
func ( l *Logger ) getPrefix( colo bool ) string {
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
	default:
		return ""
	}
}
func ( l *Logger ) wr( lvl uint8, m string ) {
// log to file
	l.ch_message <- l.getPrefix(false) + m

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
	*log.Logger
	sync.WaitGroup
	mess_queue chan string
	stop_handle chan bool
}

func ( fl *fileLogger ) start() {
	go func() {
		fl.Add(1)
		for {
			select{
			case m := <- fl.mess_queue:
				fl.Println(m)
			case <- fl.stop_handle:
				fl.Println("Log worker has been stopped!")
				return
			}
		}
		fl.Done()
	}()
}
func ( fl *fileLogger ) stop() {
	fl.Println("Trying to close File Logger...")
	close(fl.stop_handle)
}
