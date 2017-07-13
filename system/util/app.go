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

	// service submodule settings:
	ServiceMaxErrors uint8
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
		if fldS2 := mrgS2.Field(i); !self.isEmptyValue(fldS2) { mrgS1.Field(i).Set(reflect.Value(fldS2)) }
	}

	return self,nil
}
func (self *AppConfig) isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
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
