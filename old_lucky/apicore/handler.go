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
		authid,e := steamoid.validate(); if e != nil {
			// TODO: move error (e.Error()) in jsoner!!!
			ctx.Error(e.Error(), fasthttp.StatusInternalServerError)
		}

		var dataCookie *fasthttp.Cookie = fasthttp.AcquireCookie()
		dataCookie.SetKey("claimed_id")
		dataCookie.SetPath("/")
		dataCookie.SetDomain(".gotest.mh00.info")
		dataCookie.SetSecure(false)
		dataCookie.SetHTTPOnly(true)
		dataCookie.SetValueBytes(authid)

		ctx.Response.Header.SetCookie(dataCookie)
		ctx.Redirect("http://golucky.gotest.mh00.info/", 307)

		// XXX STOPPED HERE!; Create cookie, create JSON response; SIGN json response; Send cookie (2 cookie - data && sign )
		// TODO ; ENCRYPT COOKIE ????

	}
}
func (self *apiHandler) CentrifugoConnection(ctx *fasthttp.RequestCtx) {
}
