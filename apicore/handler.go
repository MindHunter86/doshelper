package apicore

import "github.com/valyala/fasthttp"

type handler struct {
}
func (self *handler) login( ctx *fasthttp.RequestCtx ) {
	var steamoid *steamOid
	var e error

	steamoid, e = new(steamOid).configure( ctx.Request.URI(), ctx.Request.Header )
	if e != nil { ctx.Error( e.Error(), fasthttp.StatusInternalServerError ) }

	switch string(steamoid.data.Peek("openid.mode")) {
	case "":
		ctx.Redirect( string(steamoid.authUrl()), 302 )
	case "cancel":
		ctx.Write([]byte("Authorization cancelled!"))
	default:
		if steamid, e := steamoid.validate(); e != nil {
			ctx.Error( e.Error(), fasthttp.StatusInternalServerError )
		} else { ctx.Write([]byte(steamid)) }
	}
}
