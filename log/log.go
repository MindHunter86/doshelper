package log

import "log"
import "sync"

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

type Logger struct {
	*log.Logger
	Ch_message chan string
	sync.RWMutex
	Prefix uint8
}

func ( l *Logger ) SetPrefix( p uint8 ) {
	l.Lock(); l.Prefix = p; l.Unlock()
}
func ( l *Logger ) GetPrefix( colo bool ) string {
	l.RLock()
	defer l.RUnlock()

	switch l.Prefix {
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
func ( l *Logger ) W( lvl uint8, m string ) {
// log to file
	switch lvl {
	case LLEV_DBG:
		l.Ch_message <- l.GetPrefix(false) + "DBG: " + m
	case LLEV_INF:
		l.Ch_message <- l.GetPrefix(false) + "INF: " + m
	case LLEV_OK:
		l.Ch_message <- l.GetPrefix(false) + "SUC: " + m
	case LLEV_NON:
		l.Ch_message <- l.GetPrefix(false) + "FAI: " + m
	case LLEV_WRN:
		l.Ch_message <- l.GetPrefix(false) + "WRN: " + m
	case LLEV_ERR:
		l.Ch_message <- l.GetPrefix(false) + "ERR: " + m
	}

// log to stdout
	switch lvl {
	case LLEV_DBG:
		l.Println( l.GetPrefix(true), m )
	case LLEV_OK:
		l.Println( l.GetPrefix(true), "\x1b[32;1m✔\x1b[0m", m )
	case LLEV_NON:
		l.Println( l.GetPrefix(true), "\x1b[31;1m✖\x1b[0m", m)
	case LLEV_INF:
		l.Println( l.GetPrefix(true), "\x1b[34;1m🛈\x1b[0m", m)
	case LLEV_WRN:
		l.Println( l.GetPrefix(true), "\x1b[33;1m❢\x1b[0;33;22m", m, "\x1b[0m")
	case LLEV_ERR:
		l.Println( l.GetPrefix(true), "\x1b[31;22m❢\x1b[0;31;1m", m, "\x1b[0m")
	}
}


type FileLogger struct {
	sync.WaitGroup
	*log.Logger
	Mess_queue chan string
	Stop_handle chan bool
}

func ( fl *FileLogger ) Start() {
	go func() {
		fl.Add(1)
		defer fl.Done()

		for {
			select{
			case m := <-fl.Mess_queue:
				fl.Println(m)
			case <-fl.Stop_handle:
				var buf_size uint8 = uint8( len(fl.Mess_queue) )
				for ; buf_size != 0; buf_size-- {
					fl.Println(<-fl.Mess_queue)
				}
				log.Println("Log worker has been stopped! Buf is empty")
				return
			}
		}
	}()
}
func ( fl *FileLogger ) Stop() {
	close(fl.Stop_handle)
}
