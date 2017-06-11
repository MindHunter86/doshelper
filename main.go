package main

import "os"
import "os/signal"
import "syscall"

import "golucky/system"
import "golucky/system/util"

import "github.com/sirupsen/logrus"
import "golang.org/x/net/context"

//
//func init() {
//	log.SetOutput(os.Stdout)
//	log.SetLevel(log.DebugLevel)
////	log.SetFormatter(&log.JSONFormatter{})
//	log.SetFormatter(&log.TextFormatter{
//		ForceColors: true,
//		FullTimestamp: true,
//		TimestampFormat: "Mon, 02 Jan 2006 15:04:05 -0700",
//	})
//}
//func main() {
//	log.Println("Hello world?")
//	log.WithField("error", "Fuck you!").Error("WHAT?!")
//	log.Info("Yes")
//	log.Debug("Useful debugging information.")
//log.Info("Something noteworthy happened!")
//log.Warn("You should probably take a look at this.")
//log.Error("Something failed but I'm not quitting.")
//log.WithFields(log.Fields{"One": 1, "Two": 2}).Fatal("SHIT!")
//}

var app *application
type application struct {
	logout *logrus.Logger
	ptr_module_system *system.SystemModule
}

func init() {
	// main initialization:
	app = new(application)

	// log initialization:
	app.logout = logrus.New()
	app.logout.Out = os.Stdout
	app.logout.Level = logrus.DebugLevel
	app.logout.Formatter = &logrus.TextFormatter{
		ForceColors: true,
		FullTimestamp: true,
		TimestampFormat: "Mon, 02 Jan 2006 15:04:05 -0700",
	}

	// info messages:
	app.logout.Debugln("Application logger has been initialized!")
	app.logout.Infoln("Application has been initialized! Starting subsystems...")
}
func main() {
	var e error
	var ctx context.Context

	// create context for next modules initialization
	ctx = context.WithValue(context.Background(), util.CTX_MAIN_LOGGER, app.logout)

	// application main module initialization:
	app.ptr_module_system, e = new(system.SystemModule).Configure(ctx)

	var sgn = make(chan os.Signal)
	signal.Notify(sgn, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	app.logout.Debugln("Kernel signal catcher has been initialized!")

	// some subsystem startings ...
	// TODO

	// application kernel signal catcher:
	<-sgn
	app.logout.Warnln("Catched QUIT signal from kernel! Stopping all systems...")

	// close and destroy all subsystems:
	e = app.ptr_module_system.Destroy(); if e != nil {
		app.logout.WithField("error", e.Error()).Error(util.Err_Main_ErroredExit)
		app.logout.Warnln("Application has been destroyed with errors!")
	} else {
		app.logout.Infoln("Application has been successfully destroyed!")
		os.Exit(0)
	}

	os.Exit(1)
}
