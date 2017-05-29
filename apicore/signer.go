package apicore

import "crypto/hmac"
import "crypto/sha1"

import "doshelpv2/log"
import "doshelpv2/appctx"

import "golang.org/x/net/context"


type apiSigner struct {
	*ApiCore
	secret []byte
}
func (self *apiSigner) configure(ctx context.Context) ( *apiSigner, error ) {
	if self == nil { return nil,err_glob_InvalidSelf }
	if ctx == nil { return nil,err_glob_InvalidContext }

	self.ApiCore = ctx.Value(appctx.CTX_MOD_APICORE).(*ApiCore)
	self.secret = []byte(self.sign_secret)

	self.slogger.W( log.LLEV_DBG, "Hmac submodule has been initialized and configured!")
	return self,nil
}
func (self *apiSigner) sign( message []byte ) []byte {
	mac := hmac.New( sha1.New, self.secret )
	mac.Write( message )
	return mac.Sum(nil)
}
func (self *apiSigner) checkSign( message, sign []byte ) bool {
	mac := hmac.New( sha1.New, self.secret )
	mac.Write( message )
	return hmac.Equal( mac.Sum(nil), sign )
}
