package main

import "flag"

import "golucky/system/util"
import "golucky/system/application"

import "github.com/sirupsen/logrus"

var app *application.Application

func main() {
	// config initialization:
	var inpConfig *util.AppConfig = new(util.AppConfig)
	flag.StringVar(&inpConfig.P2PInfoHash, "p2p_hash", "", "P2P info hash for autodescovery.")
	// ...
	// etc.
	flag.Parse()

	outConfig, e := new(util.AppConfig).Configure(inpConfig); if e != nil || !flag.Parsed() {
		logrus.WithField("error", e).Println("Could not parse configuration!")
		return
	}

	// app initialization && configuration:
	app = new(application.Application).Initialize(outConfig)
	app.GetLogger().Infoln("Application has been initialized and configured!")

	if e := app.ConfigureAndLaunch(); e != nil {
		app.GetLogger().WithField("error", e).Errorln("Could not configure application!")
		return
	}
	app.GetLogger().Infoln("Application has been successfully unloaded!")
}
