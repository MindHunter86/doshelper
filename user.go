package main

import (
	"bytes"
	"errors"
	"strings"

	"net/http"
	"crypto/md5"
	"encoding/base64"
)
import gouuid "github.com/satori/go.uuid"

type User struct {
	req *http.Request
	Uuid, secure_hash string
}
func NewUser( r *http.Request ) *User {
	u := &User{ req: r }
	return u
}
func ( u *User ) ParseOrCreateUUID() *http.Cookie {
	uuid_c, _ := u.req.Cookie("uuid")

	if len( uuid_c.String() ) <= 0 {
		u.Uuid = gouuid.NewV4().String()

		var https bool = false
		switch u.req.URL.Scheme {
		case "https":
			https = true
		}

		return &http.Cookie{
			Name: "uuid",
			Value: u.Uuid,
			Path: "/",
			Domain: u.req.URL.Host,
			Secure: https,
			HttpOnly: true,
		}
	}
	u.Uuid = uuid_c.Value
	return nil
}
func ( u *User ) GetSecureHash() string {
	sh, e := u.req.Cookie("sl"); if e != nil { return "" }
	return sh.Value
}
func ( u *User ) GenSecureHash() ( *http.Cookie, error ) {
// buf - $remote_addr:$cookie_uuid:$user_agent:$secret
	if len(u.Uuid) <= 0 { return nil,errors.New("User's UUID is empty! Logical error!") }
	if len( u.GetSecureHash() ) != 0 { return nil,errors.New("SL cookie was already defined!") }

	var buf bytes.Buffer
	buf.WriteString( u.req.Header.Get("X-Forwarded-For") + u.Uuid + u.req.UserAgent() )
	buf.WriteString( u.req.Header.Get("X-SecureLink-Secret") )

	t1 := md5.Sum( buf.Bytes() )
	t2 := base64.StdEncoding.EncodeToString( t1[:] )
	t3 := strings.Replace( strings.Replace( t2, "+", "-", -1 ), "/", "_", -1 )
	t4 := strings.Replace( t3, "=", "", -1 )

	var https bool = false
	switch u.req.URL.Scheme {
	case "https":
		https = true
	}

	return &http.Cookie{
		Name: "sl",
		Value: t4,
		Path: "/",
		Domain: u.req.URL.Host,
		Secure: https,
		HttpOnly: true,
	}, nil
}
