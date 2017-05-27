package apicore

import "errors"
import "bytes"
import "regexp"

import "github.com/valyala/fasthttp"

var (
	err_SteamOid_NullSelf = errors.New("steamOid isn't defined!")
	err_SteamOid_InvalidMode = errors.New("Mode must equal to \"id_res\".")
	err_SteamOid_InvalidReturn = errors.New("The \"return_to url\" must match the url of current request.")
	err_SteamOid_InvalidNS = errors.New("Wrong NS in the Steam response.")
	err_SteamOid_NonValidate = errors.New("Unable validate Steam OpenID.")
	err_SteamOid_InvalidPattern = errors.New("Invalid Steam ID pattern.")
	err_SteamOid_SteamNon200 = errors.New("Steam is unavailable now! Please, try again later.")
)

var (
	steamoid_url []byte = []byte("https://steamcommunity.com/openid/login")
	steamoid_mode []byte = []byte("checkid_setup")
	steamoid_ns []byte = []byte("http://specs.openid.net/auth/2.0")
	steamoid_identity []byte = []byte("http://specs.openid.net/auth/2.0/identifier_select")
	steamoid_params []byte = []byte("assoc_handle,signed,sig,ns")
	steamoid_valide *regexp.Regexp = regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	steamoid_dightex *regexp.Regexp = regexp.MustCompile("\\D+")
)

type steamOid struct {
	realm, return_to []byte
	data *fasthttp.Args
}
func (self *steamOid) configure( uri *fasthttp.URI, header fasthttp.RequestHeader ) ( *steamOid, error ) {
	if self == nil { return nil,err_SteamOid_NullSelf }

	self.realm = append( uri.Scheme(), []byte("://") ...)
	self.realm = append( self.realm, header.Peek("X-Forwarded-Host") ...)

	self.return_to = uri.RequestURI()
	if i := bytes.Index( self.return_to, []byte("openid") ); i != -1 { self.return_to = self.return_to[0:i-1] }
	self.return_to = append( self.realm, self.return_to ...)

	self.data = uri.QueryArgs()
	return self,nil
}
func (self *steamOid) authUrl() []byte {
	oid_data := map[string][]byte {
		"openid.claimed_id": steamoid_identity,
		"openid.identity": steamoid_identity,
		"openid.mode": steamoid_mode,
		"openid.ns": steamoid_ns,
		"openid.realm": self.realm,
		"openid.return_to": self.return_to,
	}

	var i uint8
	var query []byte
	for k,v := range oid_data {
		query = append( query,  append( []byte( k + "=" ), v ...) ...)
		if i != uint8( len(oid_data) - 1 ) {
			query = append( query, []byte("&") ...)
		}
		i++
	}
	query = append( []byte("?"), query ...)
	return append( steamoid_url, query ...)
}
func (self *steamOid) validate() ( []byte, error ) {
	if ! bytes.Equal( self.data.Peek("openid.mode"), []byte("id_res") ) { return nil, err_SteamOid_InvalidMode }
	if ! bytes.Equal( self.data.Peek("openid.return_to"), self.return_to ) { return nil, err_SteamOid_InvalidReturn }

	postform := fasthttp.AcquireArgs()
	for _,v := range bytes.Split( steamoid_params, []byte(",") ) {
		postform.SetBytesV( "openid."+string(v), self.data.Peek("openid."+string(v)) )
	}
	for _,v := range bytes.Split( self.data.Peek("openid.signed"), []byte(",") ) {
		postform.SetBytesV( "openid."+string(v), self.data.Peek("openid."+string(v)) )
	}
	postform.Set("openid.mode", "check_authentication")

	rspCode, rspBody, e := fasthttp.Post( nil, string(steamoid_url), postform )
	if e != nil { return nil,e }
	if rspCode != 200 { return nil,err_SteamOid_SteamNon200 }

	response := bytes.Split( rspBody, []byte("\n") )
	if ! bytes.Equal( response[0], append( []byte("ns:"), steamoid_ns ...) ) { return nil,err_SteamOid_InvalidNS }
	if bytes.HasSuffix( response[1], []byte("false") ) { return nil,err_SteamOid_NonValidate }

	if ! steamoid_valide.Match( self.data.Peek("openid.claimed_id") ) { return nil,err_SteamOid_InvalidPattern }
	return steamoid_dightex.ReplaceAll( self.data.Peek("openid.claimed_id"), []byte("") ), nil
}
