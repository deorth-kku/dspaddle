package main

import "syscall"

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procMapVirtualKey       = user32.NewProc("MapVirtualKeyW")
	procSendInput           = user32.NewProc("SendInput")
	procPostMessage         = user32.NewProc("PostMessageW")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
)
