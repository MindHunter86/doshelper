package apicore

import "crypto/hmac"
import "crypto/sha1"

import dlog "doshelpv2/log"


type apiSigner struct {
	secret []byte
}
func (self *apiSigner) configure( secret []byte, logger *dlog.Logger ) ( *apiSigner, error ) {
	if self == nil { return nil,err_Signer_InvalidSigner }

	if len(secret) == 0 { return nil,err_Signer_InvalidInput }
	self.secret = secret

	logger.W( dlog.LLEV_DBG, "Hmac signer submodule has been inited!" )
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
