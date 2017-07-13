package main

import  "flag"
import "os"
import "os/signal"
import "syscall"

import "golucky/model"
import "golucky/system"
import "golucky/controller"
import "golucky/system/util"

import "github.com/sirupsen/logrus"
import "golang.org/x/net/context"

var app *util.Application

func init() {
	// main initialization:
	app = new(util.Application)

	// log initialization:
	app.Logout = logrus.New()
	app.Logout.Out = os.Stdout
	app.Logout.Level = logrus.DebugLevel
	app.Logout.Formatter = &logrus.TextFormatter{
		ForceColors: true,
		FullTimestamp: true,
		TimestampFormat: "Mon, 02 Jan 2006 15:04:05 -0700",
	}

	// config initialization:
	var inpConfig *util.AppConfig = new(util.AppConfig)
	flag.StringVar(&app.Config.P2PInfoHash, "p2p_hash", "", "P2P info hash for autodescovery.")
	// ...
	// ...
	// etc.

	var e error
	if app.Config, e = new(util.AppConfig).Configure(inpConfig); e != nil {
		app.Logout.WithField("error", e).Println("Could not parse configuration!")
		os.Exit(1)
	}

	// info messages:
	app.Logout.Debugln("Application logger has been initialized!")
	app.Logout.Infoln("Application has been initialized! Starting subsystems...")
}
func main() {
	var e error = nil // if War will happen =)
	var ctx context.Context

	// create context for next modules initialization
	ctx = context.WithValue(context.Background(), util.CTX_MAIN_LOGGER, app.Logout)

	// application main module initialization:
	for {	// "module error catcher":
		if app.PTR_controller, e = new(controller.ControllerModule).Configure(ctx); e != nil { break }
		if app.PTR_system, e = new(system.SystemModule).Configure(ctx); e != nil { break }
		if app.PTR_model, e = new(model.ModelModule).Configure(ctx); e != nil { break }
		break
	}
	if e != nil { app.Logout.WithField("error", e).Errorln(util.Err_Main_ModuleError) }

	var sgn = make(chan os.Signal)
	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	app.Logout.Debugln("Kernel signal catcher has been initialized!")

	// some subsystem startings ...
	// TODO

	// application kernel signal catcher:
	<-sgn
	app.Logout.Warnln("Catched QUIT signal from kernel! Stopping all systems...")

	// close and destroy all subsystems:
	e = nil // if War will happen
	for {	// "module error catcher":
		if e = app.PTR_model.Destroy(); e != nil { break }
		if e = app.PTR_system.Destroy(); e != nil { break }
		if e = app.PTR_controller.Destroy(); e != nil { break }
		break
	}
	if e != nil {
		app.Logout.WithField("error", e.Error()).Error(util.Err_Main_ErroredExit)
		app.Logout.Warnln("Application has been destroyed with errors!")
		os.Exit(1)
	}
	app.Logout.Infoln("Application has been successfully destroyed!")
	os.Exit(0)
}
