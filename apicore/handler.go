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
	}
}
func (self *apiHandler) HmacTest( ctx *fasthttp.RequestCtx ) {
	var input1 []byte = []byte("TestString; Super Secret!")
	var input2 []byte = []byte("TestString; SuperSuperSecret!")

	var signed1 = self.signer.sign(input1)
	var signed2 = self.signer.sign(input2)

//	var logical_boolean1 = self.signer.checkSign( input1, signed1 )
//	var logical_boolean2 = self.signer.checkSign( input1, signed2 )
//	self.signer.secret = []byte("12341234123")
//	var logical_boolean3 = self.signer.checkSign( input1, signed1 )

	self.slogger.W( log.LLEV_DBG, "1-st mess: " + string(input1) )
	self.slogger.W( log.LLEV_DBG, "1-st sign: " + string(signed1) )
	self.slogger.W( log.LLEV_DBG, "2-nd mess: " + string(input2) )
	self.slogger.W( log.LLEV_DBG, "2-nd sign: " + string(signed2) )
//	self.slogger.W( log.LLEV_DBG, "Check1(OK): " + logical_boolean1 )
//	self.slogger.W( log.LLEV_DBG, "Check2(NON): " + logical_boolean2 )
//	self.slogger.W( log.LLEV_DBG, "Check3(???): " + logical_boolean3 )

	ctx.Write([]byte("Test completed! Check logs."))
}
