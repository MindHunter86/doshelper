package util

import "reflect"
import "github.com/sirupsen/logrus"

type Application struct {
	Logout *logrus.Logger

	Config *AppConfig

	PTR_system System
	PTR_controller Controller
	PTR_model Model
}

type AppConfig struct {
	// p2p settings:
	P2PInfoHash string
	P2PDhtLstnPort uint16
	P2PDhtRouters string
	P2PDebug bool
	P2PDebugPort uint16

	// log settings:
	LogColorized bool
	LogTimestamps bool
	LogFormat string
}
// Merge default configuration with input configuration:
func (self *AppConfig) Configure(inpConfig *AppConfig) (*AppConfig, error) {
	if self == nil { return nil,Err_Glob_InvalidSelf }
	if inpConfig == nil { inpConfig = new(AppConfig) }

	// default values:
	self.P2PInfoHash = "deca7a89a1dbdc4b213de1c0d5351e92582f31fb"
	self.P2PDhtRouters = "router.utorrent.com:6881,router.magnets.im:6881,router.bittorrent.com:6881,dht.transmissionbt.com:6881,dht.aelitis.com:6881,router.bitcomet.com:6881"
	self.P2PDhtLstnPort = uint16(9000)
	self.P2PDebug = true
	self.P2PDebugPort = uint16(9001)

	self.LogColorized = true
	self.LogTimestamps = true
	self.LogFormat = "Mon, 02 Jan 2006 15:04:05 -0700"

	// merge config structs (mrgS# - merge struct #):
	mrgS1 := reflect.ValueOf(self).Elem()
	mrgS2 := reflect.ValueOf(inpConfig).Elem()

	s1fnum := mrgS1.NumField()
	for i := 0 ; i < s1fnum; i++ {
		// fldS2 - field of second struct (input struct)
		if fldS2 := mrgS2.Field(i); fldS2.IsNil() { mrgS1.Field(i).Set(reflect.ValueOf(fldS2)) }
	}

	return self,nil
}

type System interface {
	Destroy() error
}
type Controller interface {
	Destroy() error
}
type Model interface {
	Destroy() error
}
