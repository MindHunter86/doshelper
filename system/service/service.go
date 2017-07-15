package service

import "sync"
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
	service util.Service
	status uint32
	error_ch chan error
}
func (self *BaseService) Status() uint32 {
	return atomic.LoadUint32(&self.status)
}
func (self *BaseService) SetStatus(status uint32) {
	atomic.StoreUint32(&self.status, status)
}
func (self *BaseService) IsNeedStop() bool {
	return self.Status() == StatusStopping
}


type ServiceSubmodule struct {
	log *logrus.Logger
	wgroup sync.WaitGroup
	services map[uint8]*BaseService

	cnf_maxErrors uint8
}
func (self *ServiceSubmodule) Configure(ctx context.Context) (*ServiceSubmodule, error) {
	if self == nil { return nil,util.Err_Glob_InvalidSelf }
	if ctx == nil { return nil,util.Err_Glob_InvalidContext }

	self.log = ctx.Value(util.CTX_MAIN_LOGGER).(*logrus.Logger)
	self.cnf_maxErrors = ctx.Value(util.CTX_MAIN_CONFIG).(*util.AppConfig).ServiceMaxErrors
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
// run all services conigured in self.services:
func (self *ServiceSubmodule) Run(done <-chan struct{}) error {
	for id, bService := range self.services {
		switch bService.Status() {
		case StatusReady:
			go self.bootstrap(id)
			self.log.Debugln("Service " + util.SERVICE_PTR[id] + " has been successfully bootstraped!")
		case StatusFailed:
			self.log.Warnln("Service " + util.SERVICE_PTR[id] + " has FAILURE status. You must reset it, before starting again!")
		case StatusStopping:
			self.log.Warnln("Service " + util.SERVICE_PTR[id] + " has not stopped yet. You can not run it now!")
		case StatusRunning:
			self.log.Infoln("Service " + util.SERVICE_PTR[id] + " is running now!")
		}
	}

	// wating signal for main.go for closing all running services:
	self.log.Infoln("WAIT...")
	<-done
	self.log.Warnln("DEBUG! Cached MainGO signal")
	return nil
}
func (self *ServiceSubmodule) bootstrap(id uint8) {
	self.wgroup.Add(1)

	for i := uint8(0); i < self.cnf_maxErrors; i++ {
		// bootstrap initialization:
		self.services[id].error_ch = make(chan error)
		self.services[id].SetStatus(StatusRunning)

		// service bootstrap:
		go func(self *ServiceSubmodule, id uint8) {
			self.wgroup.Add(1)

			if e := self.services[id].service.Start(); e != nil { self.services[id].error_ch <-e }
			close(self.services[id].error_ch)

			self.wgroup.Done()
		}(self, id)

		// catch error or close() method:
		e := <-self.services[id].error_ch
		if e != nil {
			self.services[id].SetStatus(StatusFailed)
			self.log.WithField("error", e).Warnln("Service " + util.SERVICE_PTR[id] + " has been unexpectedly closed!")
			continue
		}

		// no errors, service has been closed normaly, exit...:
		self.services[id].SetStatus(StatusReady)
		break
	}

	self.wgroup.Done()
}
// stop service with "id":
func (self *ServiceSubmodule) Stop(id uint32) error {
	return nil
}
// stop all services; destroy submodule:
func (self *ServiceSubmodule) Destroy() error {
	self.log.Debugln("Waiting for services closing...")
	self.Wait()

	self.log.Infoln("All services has been stopped!")
	return nil
}
func (self *ServiceSubmodule) PreloadService(service_ident uint8, service_ptr util.Service, service_error error) error {
	self.services[service_ident] = &BaseService{
		service: service_ptr,
		status: StatusReady,
	}
	if service_error != nil {
		self.services[service_ident].SetStatus(StatusFailed)
		return service_error
	}
	return nil
}
