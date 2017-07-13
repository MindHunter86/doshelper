package util

import "golang.org/x/net/context"

const (
	SERVICE_PTR_P2P = uint8(iota)
	SERVICE_PTR_HTTP
	SERVICE_PTR_RPC
)
var SERVICE_PTR = map[uint8]string {
	SERVICE_PTR_P2P: "P2P",
	SERVICE_PTR_HTTP: "HTTP",
	SERVICE_PTR_RPC: "RPC",
}

type Service interface {
	// Method for service configuration.
	// After successful configure, service status will change to "StatusReady"
	Configure(context.Context) (uint8, Service, error)

	// Run service. Avaliable only after configuring (service status must be "StatusReady")
	Start() error

	// Stop service. Avaliable only after starting (service status must be "StatusRunning")
	Stop() error

	// Destroy service. If service status is "ServiceRunning", it wiil be stopped before.
	Destroy() error
}
