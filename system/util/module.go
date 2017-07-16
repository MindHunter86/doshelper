package util

import "golang.org/x/net/context"

const ( // application modules:
	APPMODULE_CONTROLLER = uint8(iota)
	APPMODULE_SYSTEM
	APPMODULE_MODEL
)
var AppModules = map[uint8]string { // module namings for logging:
	APPMODULE_CONTROLLER: "CONTROLLER",
	APPMODULE_SYSTEM: "SYSTEM",
	APPMODULE_MODEL: "MODEL",
}

type AppModule interface {
	Configure(context.Context) (AppModule, error)
	Load() error
	Unload()
}
