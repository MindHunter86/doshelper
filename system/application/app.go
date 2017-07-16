package application

import "syscall"
import "os"
import "os/signal"
import "sync"

import "golucky/model"
import "golucky/system"
//import "golucky/controller"
import "golucky/system/util"

import "golang.org/x/net/context"
import "github.com/sirupsen/logrus"


type baseModule struct {
	util.AppModule
	error_ch chan error
}


type Application struct {
	logger *logrus.Logger
	config *util.AppConfig
	wGroup sync.WaitGroup

	modules map[uint8]*baseModule

	sgnl_exit chan os.Signal
	sgnl_reload chan os.Signal

	ctx_close context.CancelFunc

	module_err chan struct{}
}
func (self *Application) Initialize(cnf *util.AppConfig) *Application {
	// TODO: Added some cnf variable checks
	self.config = cnf

	self.logger = logrus.New()
	self.logger.Out = os.Stdout
	self.logger.Level = logrus.DebugLevel
	self.logger.Formatter = &logrus.TextFormatter{
		ForceColors: self.config.LogColorized,
		FullTimestamp: self.config.LogTimestamps,
		TimestampFormat: self.config.LogFormat,
	}

	self.modules = make(map[uint8]*baseModule)
	self.module_err = make(chan struct{})

	return self
}
func (self *Application) ConfigureAndLaunch() error {
	var e error
	var ctx context.Context

	ctx, self.ctx_close = context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, util.CTX_MAIN_LOGGER, self.logger)
	ctx = context.WithValue(ctx, util.CTX_MAIN_WGROUP, self.wGroup)
	ctx = context.WithValue(ctx, util.CTX_MAIN_CONFIG, self.config)

	for { // WARNING! Save order! See module const ids!
		// Ex: if e = self.PreloadService(new(p2p.P2PService).Configure(ctx)); e != nil { break }
		if e = self.preloadModule(new(model.ModelModule).Configure(ctx)); e != nil { break }
		if e = self.preloadModule(new(system.SystemModule).Configure(ctx)); e != nil { break }
		break
	}
	if e != nil { self.logger.WithField("module error", e).Errorln(util.Err_App_ModuleError) }

	self.sgnl_exit = make(chan os.Signal)
	self.sgnl_reload = make(chan os.Signal)
	signal.Notify(self.sgnl_exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	signal.Notify(self.sgnl_reload, syscall.SIGHUP)
	self.logger.Debugln("Kernel signal catcher has been initialized!")

	// launch applicatio:n
	self.launch()
	return nil
}
func (self *Application) GetLogger() *logrus.Logger {
	return self.logger
}
func (self *Application) launch() {
	var modules_count uint8 = uint8(len(self.modules))
	for i := uint8(0); i < modules_count; i++ {
		go self.bootstrapModule(i)
		self.logger.Infoln("Module " + util.AppModules[i] + " has been bootstraped!")
	}

DSTR:
	for {
		select {
			// TODO: REFACTOR!!!
		case <-self.module_err:
			self.ctx_close()
			break DSTR
		case <-self.sgnl_exit:
			self.logger.Warnln("Catched signal from kernel! Started modules unload!")
			self.ctx_close()
			break DSTR
		case <-self.sgnl_reload:
			// reload configuration;
			// stop services;
			// start services;
		}
	}

	// unload modules:
	self.wGroup.Wait()
	for i := uint8(0); i < modules_count; i++ {
		self.modules[i].Unload()
		self.logger.Infoln("Module " + util.AppModules[i] + " has been successfully unloaded!")
	}
}
func (self *Application) preloadModule(mPtr util.AppModule, mError error) error {
	if mError != nil { return mError }

	self.modules[uint8(len(self.modules))] = &baseModule{
		AppModule: mPtr,
		error_ch: make(chan error),
	}

	return nil
}
func (self *Application) bootstrapModule(id uint8) {
	if _,ok := (<-self.module_err); !ok { return }
	self.wGroup.Add(1)

	go func(self *Application, id uint8) {
		self.wGroup.Add(1)

		if e := self.modules[id].Load(); e != nil { self.modules[id].error_ch <- e }
		close(self.modules[id].error_ch)

		self.wGroup.Done()
	}(self, id)

	// catch error or close() event:
	e := <-self.modules[id].error_ch
	if e != nil {
		self.logger.WithField("module error", e).Errorln("Module " + util.AppModules[id] + " has been unexpectedly exited!")
		close(self.module_err)
	}

	self.wGroup.Done()
}
