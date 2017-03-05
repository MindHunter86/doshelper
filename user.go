package main

import (
	"sync"
	"bytes"
	"errors"
	"strings"

	"net/http"
	"crypto/md5"
	"encoding/base64"
)
import gouuid "github.com/satori/go.uuid"


const (
	ERR_USER_UUIDEMP = "User's UUID is empty! Logical error!"
	ERR_USER_HWIDEMP = "User's HWID is empty! Logical error!"
)
var (
	ERR_USER_NOUUID = errors.New("User's UUID is empty! Logical error!")
	ERR_MAIN_NOPARAM = errors.New("Received empty params! Function ferror!")
)

// User cacher ( Key/value "database" )
//
// Key: user's HWID
// Value: *User
// return &UserBuf{ clients: make(map[string]client)	}
type activeClients struct {
	sync.RWMutex
	clients map[string]client
}
func ( ac *activeClients ) get( hwk string ) ( *client, bool ) {
	ac.RLock()
	cl, ok := ac.clients[hwk]
	ac.RUnlock()
	return &cl, ok
}
func ( ac *activeClients ) validate( hwk string ) bool {
	ac.RLock()
	_, ok := ac.clients[hwk]
	ac.RUnlock()
	return ok
}
func ( ac *activeClients ) put( hwk string, cl client ) {
	ac.Lock()
	ac.clients[hwk] = cl
	ac.Unlock()
}


// type sqlClient struct {
// 	conn *sql.DB
// }
// func NewSQLClient( host,user,pass,database string ) ( *sql.DB, error ) {
// 	db, e := sql.Open( "msyql", user + ":" + p + "@" + host + "/" + database )
// 	if e != nil || db.Ping() != nil { return nil,e }
// 	return &sqlClient{ conn: db }
// }
// func ( sc *sqlClient ) connectionCheck() error { return sc.conn.Ping() }
// func ( sc *sqlClient ) connectionClose() {}

// func ( ub *UserBuf ) UserSave( u client ) {
// 	app.Add(0)
// 
// 	hwid := u.getHWID()
// 	if ub.userValidate(hwid) {
// 	//	TRUE - let's only put in DB
// 	} else {
// 	// FALSE - let's put in cache, then put in DB	
// 		ub.userPut( hwid, u )
// 	}
// 
// 	app.Done()
// }
// func ( ub *UserBuf ) UserLoad( hwid string ) {}


// Client in webHandler
// Get HWhash;
//	if no hash - createUser ( create user with received values + put in cache );
//	if hash is OK, search in cache.
//		if hash is't cached; get value from DB
//			if no record in DB - DROP user
//			else update cache
//	if hash is not OK
//		create user in cache

//???
//	defer in handler - update in db ???


func userCreate() {}
// Create with received values in cache then in db
// or get from DB & compare with received values. Making updates in cache in db
func userUpdate() {}
// Rewrite user in cache, update in db

/*
	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Forwarded-Host $host;
	proxy_set_header X-Forwarded-Proto $scheme;
	proxy_set_header X-Client-UUID $cookie_uuid;
	proxy_set_header X-Client-HWID $cookie_hwid;
	proxy_set_header X-Client-SecureLink $cookie_sl;
	proxy_set_header X-SecureLink-Secret $md5secret;
*/

// SECURE LINK HASH = 22 symb
// HWID HASH = 22 symb

type client struct {
	uuid, sec_link string
	addr, user_agent, origin, referer string
}
func newClient2( h *http.Header ) ( *client, []*http.Cookie, error ) {
	var cl *client
	hwk_h := h.Get("X-Client-HWID")

// Cache & DB validation:
	if len(hwk_h) != 0 {
		var ok bool = false
		if cl, ok = app.clients.get(hwk_h); !ok { cl = &client{}; hwk_h = "" }
	}

	var hwk_c *http.Cookie
	var cooks []*http.Cookie
	var host string = h.Get("X-Forwarded-Host")
	var proto string = h.Get("X-Forwarded-Proto")
	var mdsec string = h.Get("X-SecureLink-Secret")

	if app.clients.validate(hwk_h) == false {
		hwk_c, e := cl.generateHwKey( proto, host ); if e != nil { return nil,nil,e }
		cooks = append( cooks, hwk_c )
	}
	uid_c, e := cl.generateUid( proto, host ); if e != nil { return nil,nil,e }		// auto put in struct cl.uuid
	cooks = append( cooks, uid_c )
	scl_c, e := cl.generateSecLink( mdsec, proto, host ); if e != nil { return nil,nil,e } // auto put in struct cl.sec_link
	cooks = append( cooks, scl_c )

	if len( h.Get("X-Client-UUID") ) == 0 && len( h.Get("X-Client-SecureLink") ) == 0 {
		cl.addr = h.Get("X-Real-IP")
		cl.user_agent = h.Get("User-Agent")
		cl.origin = h.Get("Origin")
		cl.referer = h.Get("Referer")
	}

	app.clients.put( hwk_c.Value, *cl )
	return cl,cooks,nil
}
func newClient( h *http.Header ) ( *client, *http.Cookie, error ) {
	cl := &client{
		sec_link: h.Get("X-Client-SecureLink"),
		addr: h.Get("X-Real-IP"),
		user_agent: h.Get("User-Agent"),
		origin: h.Get("Origin"),
		referer: h.Get("Referer"),
	}

	cl.uuid = h.Get("X-Client-UUID")
	switch len(cl.uuid) {
	case 0:
		if uid_c, e := cl.generateUid( h.Get("X-Forwarded-Proto"), h.Get("X-Forwarded-Host") ); e != nil {
			return nil,nil,e
		} else { return cl, uid_c, nil }
	default:
		return cl, nil, nil
	}
}
func ( cl *client ) generateUid( scheme, host string ) ( *http.Cookie, error ) {
	if len(scheme) == 0 || len(host) == 0 { return nil,ERR_MAIN_NOPARAM }

	var https bool = false
	if scheme == "https" { https = true }
	cl.uuid = gouuid.NewV4().String()

	return &http.Cookie{
		Name: "uuid",
		Value: cl.uuid,
		Path: "/",
		Domain: host,
		Secure: https,
		HttpOnly: true,
	}, nil
}
func ( cl *client ) getHwKey( hwk, scheme, host string ) ( *http.Cookie, string, error ) {
	switch len(hwk) {
	case 0:
		hwk_c, e := cl.generateHwKey( scheme, host ); if e != nil { return nil,"",e }
		return hwk_c, hwk_c.Value, nil
	default:
		return nil, hwk, nil
	}
}
func ( cl *client ) generateHwKey( scheme, host string ) ( *http.Cookie, error ) {
	if len(cl.uuid) == 0 { return nil,ERR_USER_NOUUID }
	if len(scheme) == 0 || len(host) == 0 {
		return nil,ERR_MAIN_NOPARAM
	}

	var buf bytes.Buffer
	buf.WriteString( cl.addr + cl.user_agent )

	t1 := md5.Sum( buf.Bytes() )
	t2 := base64.StdEncoding.EncodeToString( t1[:] )
	t3 := strings.Replace( strings.Replace( t2, "+", "-", -1 ), "/", "_", -1 )
	t4 := strings.Replace( t3, "=", "", -1 )

	var https bool = false
	if scheme == "https" { https = true }

	return &http.Cookie{
		Name: "hwid",
		Value: t4,
		Path: "/",
		Domain: host,
		Secure: https,
		HttpOnly: true,
	}, nil
}
func ( cl *client ) generateSecLink( secret, scheme, host string ) ( *http.Cookie, error ) {
	if len(cl.uuid) == 0 { return nil,ERR_USER_NOUUID }
	if len(secret) == 0 || len(scheme) == 0 || len(host) == 0 {
		return nil,ERR_MAIN_NOPARAM
	}

	var buf bytes.Buffer
	buf.WriteString( cl.addr + cl.uuid + cl.user_agent + secret )

	t1 := md5.Sum( buf.Bytes() )
	t2 := base64.StdEncoding.EncodeToString( t1[:] )
	t3 := strings.Replace( strings.Replace( t2, "+", "-", -1 ), "/", "_", -1 )
	cl.sec_link = strings.Replace( t3, "=", "", -1 )

	var https bool = false
	if scheme == "https" { https = true }

	return &http.Cookie{
		Name: "sl",
		Value: cl.sec_link,
		Path: "/",
		Domain: host,
		Secure: https,
		HttpOnly: true,
	}, nil
}
