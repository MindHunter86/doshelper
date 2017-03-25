package main

import (
	"os"
	"sync"
	"errors"

	"net"
	"net/http"
)
import "github.com/gorilla/mux"


type SockListener struct {
	sync.RWMutex
	net.Listener
	alive bool
}

func newSockListener( proto, path string ) ( *SockListener, error ) {
	s, e := net.Listen( proto, path ); if e != nil {
		return nil,e
	}
	return &SockListener{
		Listener: s,
		alive: true,
	}, nil
}
func ( sl *SockListener ) HTTPServe( r *mux.Router ) error {
	sl.RLock()
	if sl.alive == false { return errors.New("Listener is dead!") }
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
func ( sl *SockListener ) File() ( *os.File, error ) {
	return sl.Listener.(*net.UnixListener).File()
}
func ( sl *SockListener ) Close() error {
	sl.Lock()
	defer sl.Unlock()

	if sl.alive == false { return errors.New("Listener is dead!") }
	sl.alive = false
	return sl.Listener.Close()
}
