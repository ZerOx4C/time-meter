package main

import (
	"syscall"

	"github.com/cwchiu/go-winapi"
)

func showErrorMessageBox(hWnd winapi.HWND, caption, message string) {
	captionPtr, _ := syscall.UTF16PtrFromString(caption)
	messagePtr, _ := syscall.UTF16PtrFromString(message)
	winapi.MessageBox(hWnd, messagePtr, captionPtr, winapi.MB_ICONERROR|MB_TOPMOST)
}
