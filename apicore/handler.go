package apicore

import "github.com/valyala/fasthttp"

type handler struct {
}
func (self *handler) Login( ctx *fasthttp.RequestCtx ) {
	var steamoid *steamOid
	var e error

	steamoid, e = new(steamOid).configure( ctx.Request.URI(), ctx.Request.Header )
	if e != nil { ctx.Error( e.Error(), fasthttp.StatusInternalServerError ) }

	switch string(steamoid.data.Peek("openid.mode")) {
	case "":
//		ctx.Redirect( string(steamoid.authUrl()), 302 )
		ctx.Write(steamoid.authUrl())
	case "cancel":
		ctx.Write([]byte("Authorization cancelled!"))
	default:
		steamid, e := steamoid.validate(); if e != nil {
			ctx.Error( e.Error(), fasthttp.StatusInternalServerError )
		}
		ctx.Write([]byte(steamid))

		var cookie *fasthttp.Cookie = fasthttp.AcquireCookie()
		cookie.SetKey("claimed_id")
		cookie.SetPath("/")
		cookie.SetValueBytes(steamid)
		cookie.SetDomain(".gotest.mh00.info")
		cookie.SetSecure(false)
		cookie.SetHTTPOnly(true)

		ctx.Response.Header.SetCookie(cookie)
		ctx.Redirect("http://golucky.gotest.mh00.info/", 307)
	}
}
