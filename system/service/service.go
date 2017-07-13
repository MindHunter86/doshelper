package service

import "sync/atomic"

import "golucky/system/util"
import "golucky/system/service/p2p"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


const (
	StatusReady = uint32(iota)
	StatusRunning
	StatusStopping
	StatusFailed
)

type BaseService struct {
	util.Service
	status uint32
}
func (self *BaseService) Status() uint32 {
	return atomic.LoadUint32(&self.status)
}
func (self *BaseService) SetStatus(status uint32) {
	atomic.StoreUint32(&self.status, status)
}
func (self *BaseService) Stop() {
	self.SetStatus(StatusStopping)
}
func (self *BaseService) IsNeedStop() bool {
	return self.Status() == StatusStopping
}


type ServiceSubmodule struct {
	log *logrus.Logger
	services map[uint8]*BaseService
}
func (self *ServiceSubmodule) Configure(ctx context.Context) (*ServiceSubmodule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.log = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.log.Debugln("Service submodule has been successfully initialized and configured!")

	var e error
	self.services = make(map[uint8]*BaseService)

	// Service pre-load list:
	// Ex1:
	//   self.PreloadService(new(p2p.P2PService).Configure(ctx))
	// Ex2:
	//   for {
	//     if e = self.PreloadService(new(p2p.P2PService).Configure(ctx)); e != nil { break }
	//     if e = self.PreloadService(new(p2p.P3PService).Configure(ctx)); e != nil { break }
	//     if e = self.PreloadService(new(p2p.P4PService).Configure(ctx)); e != nil { break }
	//     if e = self.PreloadService(new(p2p.P5PService).Configure(ctx)); e != nil { break }
	//     break
	//   }
	//   if e != nil { log.Println("SHIT!", e) }
	for { // "module error catcher":
		if e = self.PreloadService(new(p2p.P2PService).Configure(ctx)); e != nil { break }
		break
	}
	if e != nil { self.log.WithField("error", e).Errorln(util.Err_Service_ConfigureError) }

	return self,nil
}
func (self *ServiceSubmodule) PreloadService(service_ident uint8, service_ptr util.Service, service_error error) error {
	self.services[service_ident] = &BaseService{
		Service: service_ptr,
		status: StatusReady,
	}
	if service_error != nil {
		self.services[service_ident].status = StatusFailed;
		return service_error
	}
	return nil
}
