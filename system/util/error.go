package util

import "errors"

var (
	Err_Glob_InvalidSelf = errors.New("Detected undefined self struct in called method!")
	Err_Glob_InvalidContext = errors.New("Detected undefined inout context in called method!")

	Err_Main_ErroredExit = errors.New("Could not close application correctly!")
	Err_Main_ModuleError = errors.New("Could not start application module normaly!")

	Err_Service_ConfigureError = errors.New("Could not configure service!")

	Err_Controller_InvalidRouter = errors.New("Self Router has been already defined!")
	Err_Controller_InvalidMiddle = errors.New("Self Middlewares have been already defined!")
	Err_Controller_ImportedMiddles = errors.New("Self Middles has been already defined!")
	Err_Controller_NotImpMiddles = errors.New("Self Middles are not defined! Did you import it?")
)
