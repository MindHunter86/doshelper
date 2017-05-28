package apicore

import "doshelpv2/log"
import "doshelpv2/appctx"

import "github.com/valyala/fasthttp"
import "golang.org/x/net/context"


type apiHandler struct {
	*ApiCore
}
func (self *apiHandler) configure( ctx context.Context ) ( *apiHandler, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value( appctx.CTX_MOD_APICORE ).(*ApiCore)
	self.slogger.W( log.LLEV_DBG, "ApiHandler submodule has been initialized and configured!" )

	return self,nil
}
func (self *apiHandler) Login( ctx *fasthttp.RequestCtx ) {
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

		//
		// authid, e := steamoid.validate(); if e != nil { ctx.Error( e.Error(), nil ) } // XXX - nil
	}
}
func (self *apiHandler) CentrifugoConnection(ctx *fasthttp.RequestCtx) {
}
