package util

import "errors"

var (
	Err_Glob_InvalidSelf = errors.New("Detected undefined self struct in called method!")
	Err_Glob_InvalidContext = errors.New("Detected undefined inout context in called method!")

	Err_Main_ErroredExit = errors.New("Could not close application correctly!")
	Err_Main_ModuleError = errors.New("Could not start application module normaly!")

	Err_Controller_Api_InvalidRouter = errors.New("Self Router is not defined! Did you initialize submodule correctly?")
)
