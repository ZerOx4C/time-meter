package ui

import (
	"syscall"
	"time-meter/wrapped"

	"github.com/cwchiu/go-winapi"
)

type PopupMenu struct {
	hMenu winapi.HMENU
}

func (pm *PopupMenu) Initialize() error {
	pm.hMenu = winapi.CreatePopupMenu()
	return nil
}

func (pm *PopupMenu) Finalize() error {
	winapi.DestroyMenu(pm.hMenu)
	pm.hMenu = 0
	return nil
}

func (pm *PopupMenu) AppendStringItem(menuId MenuId, title string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	winapi.AppendMenu(pm.hMenu, winapi.MF_STRING, winapi.UINT_PTR(menuId), titlePtr)
}

func (pm *PopupMenu) Popup(hWnd winapi.HWND) {
	var pos wrapped.POINT
	winapi.GetCursorPos(pos.Unwrap())
	winapi.TrackPopupMenu(pm.hMenu, 0, pos.X, pos.Y, 0, hWnd, nil)
}
