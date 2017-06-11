package main

import (
	"os"
	"sync"

	"net"
	"net/http"
)
import "github.com/gorilla/mux"

type sockListener struct {
	sync.RWMutex
	net.Listener
	alive bool
}

func newSockListener( proto, path string ) ( *sockListener, error ) {
	s, e := net.Listen( proto, path ); if e != nil {
		return nil,e
	}
	return &sockListener{
		Listener: s,
		alive: true,
	}, nil
}
func ( sl *sockListener ) serveHTTP( r *mux.Router ) error {
	sl.RLock()
	if sl.alive == false { return ERR_SOCK_DEAD }
	sl.RUnlock()

	if e := http.Serve( sl.Listener, r ); e != nil {
		switch sl.alive {
		case true:
			return e
		default:
			return nil
		}
	}
	return nil
}
func ( sl *sockListener ) file() ( *os.File, error ) {
	return sl.Listener.(*net.UnixListener).File()
}
func ( sl *sockListener ) stop() error {
	sl.Lock()
	defer sl.Unlock()

	if sl.alive == false { return ERR_SOCK_DEAD }
	sl.alive = false
	return sl.Close()
}
