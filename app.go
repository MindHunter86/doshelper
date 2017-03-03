package main

import "sync"

var app *App
type App struct {
	sync.WaitGroup
	Users *UserBuf
	Socket *SockListener
	mgo *mgClient
}

func NewApp() ( *App, error ) {
	sl, e := NewSockListener( appNetProto, appNetPath ); if e != nil {
		return nil,e
	}
	return &App{
		Users: NewUserBuf(),
		Socket: sl,
	}, nil
}
func ( a *App ) Destroy() {
	a.Socket.Close()
	if a.mgo != nil { a.mgo.connDestroy() }
}
func ( a *App ) ThreadHTTPD() {
	a.Add(1)
	l := NewLogger( LPFX_HTTPD )
	l.PutOK("HTTPD goroutine has been inited!")

	r := NewHTTPRouter(l)
	r.Router.HandleFunc( "/", r.WebRoot )

	l.PutInf("Starting HTTP serving ...")
	for i := uint8(0); i < uint8(4); i++ {
		if e := a.Socket.HTTPServe( r.Router ); e != nil {
			l.PutWrn("Pre-fail state! HTTPServe error: " + e.Error())
			l.PutInf("Trying to restart HTTPServing ...")
			continue;
		}
		l.PutOK("HTTP serving has been stopped!")
		break;
	}

	l.PutOK("HTTPD goroutine has been destroyed!")
	a.Done()
}
