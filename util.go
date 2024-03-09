package main

import (
	"syscall"
	winapi2 "time-meter/winapi"

	"github.com/cwchiu/go-winapi"
)

func showErrorMessageBox(hWnd winapi.HWND, caption, message string) {
	captionPtr, _ := syscall.UTF16PtrFromString(caption)
	messagePtr, _ := syscall.UTF16PtrFromString(message)
	winapi.MessageBox(hWnd, messagePtr, captionPtr, winapi.MB_ICONERROR|winapi2.MB_TOPMOST)
}
