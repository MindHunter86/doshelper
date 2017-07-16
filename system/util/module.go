package util

import "golang.org/x/net/context"

const ( // application modules:
	APPMODULE_MODEL = uint8(iota)
	APPMODULE_SYSTEM
	APPMODULE_CONTROLLER
)
var AppModules = map[uint8]string { // module namings for logging:
	APPMODULE_MODEL: "MODEL",
	APPMODULE_SYSTEM: "SYSTEM",
	APPMODULE_CONTROLLER: "CONTROLLER",
}

type AppModule interface {
	Configure(context.Context) (AppModule, error)
	Load() error
	Unload()
}
