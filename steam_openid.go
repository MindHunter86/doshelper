package main

import (
	"io/ioutil"
	"errors"
	"regexp"
	"net/http"
	"net/url"
	"strings"
)


var (
	appSteamID_Url string = "https://steamcommunity.com/openid/login"
	appSteamID_Mode string = "checkid_setup"
	appSteamID_NS string = "http://specs.openid.net/auth/2.0"
	appSteamID_Ident string = "http://specs.openid.net/auth/2.0/identifier_select"
	appSteamID_Validation *regexp.Regexp = regexp.MustCompile("^(http|https)://steamcommunity.com/openid/id/[0-9]{15,25}$")
	appSteamID_DightExtra *regexp.Regexp = regexp.MustCompile("\\D+")
)

type steamOpenID struct {
	net_proto string
	net_url string
	returnlink string
	data url.Values
}

func _steamOpenID( r *http.Request ) *steamOpenID {
	var self *steamOpenID = new(steamOpenID)

	if r.TLS == nil {
		self.net_proto = "http://"
	} else { self.net_proto = "https://" }
	self.net_url = self.net_proto + r.Header.Get("X-Forwarded-Host")

	uri := r.RequestURI
	if i := strings.Index( uri, "openid" ); i != -1 {
		uri = uri[0: i-1]
	}
	self.returnlink = self.net_url + uri

	switch r.Method {
	case "POST":
		self.data = r.Form
	case "GET":
		self.data = r.URL.Query()
	}
	return self
}

func (self *steamOpenID) AuthUrl() string {
	data := map[string]string {
		"openid.claimed_id": appSteamID_Ident,
		"openid.identity": appSteamID_Ident,
		"openid.mode": appSteamID_Mode,
		"openid.ns": appSteamID_NS,
		"openid.realm": self.net_url,
		"openid.return_to": self.returnlink,
	}

	var i uint8
	url := appSteamID_Url + "?"
	for k,v := range data {
		url += k + "=" + v
		if i != uint8( len(data)-1 ) { url += "&" }
		i++
	}
	return url
}

var (
	errSteamID_InvalidMode = errors.New("Mode must equal to \"id_res\".")
	errSteamID_InvalidRetLink = errors.New("The \"return_to url\" must match the url of current request.")
	err_SteamOpenID_InvalidNS = errors.New("Wrong NS in the Steam response.")
	err_SteamOpenID_NonValidate = errors.New("Unable validate Steam OpenID.")
	err_SteamOpenID_InvalidPattern = errors.New("Invalid Steam ID pattern.")
)

func (self *steamOpenID) ValidateAndGetId() ( string, error ) {
	if self.Mode() != "id_res" { return "",errSteamID_InvalidMode }
	if self.data.Get("openid.return_to") != self.returnlink { return "",errSteamID_InvalidRetLink }

	params := make(url.Values)
	params.Set("openid.assoc_handle", self.data.Get("openid.assoc_handle"))
	params.Set("openid.signed", self.data.Get("openid.signed"))
	params.Set("openid.sig", self.data.Get("openid.sig"))
	params.Set("openid.ns", self.data.Get("openid.ns"))

	for _,i := range strings.Split( self.data.Get("openid.signed"), "," ) {
		params.Set( "openid." + i, self.data.Get( "openid." + i ) )
	}
	params.Set("openid.mode", "check_authentication")

	rsp, e := http.PostForm( appSteamID_Url, params ); if e != nil { return "",e }
	defer rsp.Body.Close()
	rspContent, e := ioutil.ReadAll(rsp.Body); if e != nil { return "",e }

	response := strings.Split( string(rspContent), "\n" )
	if response[0] != "ns:" + appSteamID_NS { return "",err_SteamOpenID_InvalidNS }
	if strings.HasSuffix( response[1], "false" ) { return "",err_SteamOpenID_NonValidate }

	openidUrl := self.data.Get("openid.claimed_id")
	if !appSteamID_Validation.MatchString(openidUrl) { return "",err_SteamOpenID_InvalidPattern }

	return appSteamID_DightExtra.ReplaceAllString(openidUrl, ""), nil
}

func (self *steamOpenID) Mode() string {
	return self.data.Get("openid.mode")
}
