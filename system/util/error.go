package util

import "errors"

var (
	Err_Glob_InvalidSelf = errors.New("Detected undefined self struct in called method!")
	Err_Glob_InvalidContext = errors.New("Detected undefined inout context in called method!")

	Err_Main_ErroredExit = errors.New("Could not close application correctly!")
)
