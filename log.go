package main


import (
	"os"
	"log"
	"sync"
	"errors"
)


const (
	LPFX_CORE = iota
	LPFX_HTTPD
	LPFX_WEBROOT
	LPFX_NOTFOUND
	LPFX_USERMANAGE
)
const (
	LLEV_OK = iota
	LLEV_NOT
	LLEV_INF
	LLEV_WRN
	LLEV_ERR
)

type Logger struct {
	*log.Logger
	sync.RWMutex
	Prefix uint8
}


func NewLogger( prefix uint8 ) *Logger {
	return &Logger{
		Logger: log.New( os.Stdout, "", log.Ldate | log.Ltime | log.Lmicroseconds ),
		Prefix: prefix,
	}
}
func ( l *Logger ) SetPrefix( p uint8 ) {
	l.Lock(); l.Prefix = p; l.Unlock()
}
func ( l *Logger ) GetPrefix() string {
	l.RLock()
	defer l.RUnlock()

	switch l.Prefix {
	case LPFX_CORE:
		return "\x1b[36;1m[CORE]:\x1b[0m"
	case LPFX_HTTPD:
		return "\x1b[36;1m[HTTPD]:\x1b[0m"
	case LPFX_WEBROOT:
		return "\x1b[36;1m[HTTPD-WEBROOT]:\x1b[0m"
	case LPFX_USERMANAGE:
		return "\x1b[36;1m[HTTPD-MIDDLEUSER]:\x1b[0m"
	default:
		return ""
	}
}
func ( l *Logger ) PutOK( txt string ) {
	l.Logger.Println( l.GetPrefix(), "\x1b[32;1m✔\x1b[0m", txt )
}
func ( l *Logger ) PutNon( txt string ) {
	l.Logger.Println( l.GetPrefix(), "\x1b[31;1m✖\x1b[0m", txt )
}
func ( l *Logger ) PutInf( txt string ) {
	l.Logger.Println( l.GetPrefix(), "\x1b[34;1m🛈\x1b[0m", txt )
}
func ( l *Logger ) PutWrn( txt string ) error {
	l.Logger.Println( l.GetPrefix(), "\x1b[33;1m❢\x1b[0;33;22m", txt, "\x1b[0m" )
	return errors.New(txt)
}
func ( l *Logger ) PutErr( txt string ) error {
	l.Logger.Println( l.GetPrefix(), "\x1b[31;22m❢\x1b[0;31;1m", txt, "\x1b[0m" )
	return errors.New(txt)
}
