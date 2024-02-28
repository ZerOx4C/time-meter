package main

import (
	"syscall"

	"github.com/cwchiu/go-winapi"
)

type PopupMenu struct {
	hMenu winapi.HMENU
}

func (pm *PopupMenu) initialize() error {
	pm.hMenu = winapi.CreatePopupMenu()
	return nil
}

func (pm *PopupMenu) finalize() error {
	winapi.DestroyMenu(pm.hMenu)
	pm.hMenu = 0
	return nil
}

func (pm *PopupMenu) appendStringItem(menuId MenuId, title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	winapi.AppendMenu(pm.hMenu, winapi.MF_STRING, winapi.UINT_PTR(menuId), titlePtr)
}

func (pm *PopupMenu) popup(hWnd winapi.HWND) {
	var pos POINT
	winapi.GetCursorPos(pos.unwrap())
	winapi.TrackPopupMenu(pm.hMenu, 0, pos.X, pos.Y, 0, hWnd, nil)
}
