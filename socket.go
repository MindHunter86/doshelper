package main

import (
	"os"
	"errors"

	"net"
	"net/http"
)
import "github.com/gorilla/mux"


type SockListener struct {
	net.Listener
	alive bool
}

func NewSockListener( proto, path string ) ( *SockListener, error ) {
	s, e := net.Listen( proto, path ); if e != nil {
		return nil,e
	}
	return &SockListener{
		Listener: s,
		alive: true,
	}, nil
}
func ( sl *SockListener ) HTTPServe( r *mux.Router ) error {
// DATA RACE 
	if sl.alive == false { return errors.New("Listener is dead!") }

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
	if sl.alive == false { return errors.New("Listener is dead!") }
// DATA RACE
	sl.alive = false
	return sl.Listener.Close()
}
