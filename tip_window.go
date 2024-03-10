package main

import (
	"errors"
	"syscall"
	"time-meter/setting"
	"time-meter/util"
	"unsafe"

	"github.com/cwchiu/go-winapi"
)

type TipWindow struct {
	hInstance winapi.HINSTANCE
	hWnd      winapi.HWND
	settings  *setting.Settings
	boundLeft int32
	onPaint   util.EventHandler
}

func (tw *TipWindow) Initialize() error {
	hInstance := winapi.GetModuleHandle(nil)
	windowClass := tw.createWindowClass(hInstance)

	if winapi.RegisterClassEx(&windowClass) == 0 {
		return errors.New("RegisterClassEx failed")
	}

	hWnd := tw.createWindow(hInstance, windowClass)
	if hWnd == 0 {
		return errors.New("createWindow failed")
	}

	tw.hInstance = hInstance
	tw.hWnd = hWnd

	return nil
}

func (tw *TipWindow) Finalize() error {
	return nil
}

func (tw *TipWindow) Show() {
	winapi.ShowWindow(tw.hWnd, winapi.SW_SHOW)
}

func (tw *TipWindow) Hide() {
	winapi.ShowWindow(tw.hWnd, winapi.SW_HIDE)
}

func (tw *TipWindow) Update() {
	var pos winapi.POINT
	winapi.GetCursorPos(&pos)

	winapi.SetWindowPos(
		tw.hWnd, winapi.HWND_TOPMOST,
		tw.boundLeft, pos.Y, 0, 0,
		winapi.SWP_NOACTIVATE|winapi.SWP_NOSIZE)

	winapi.InvalidateRect(tw.hWnd, nil, true)
}

func (tw *TipWindow) createWindowClass(hInstance winapi.HINSTANCE) winapi.WNDCLASSEX {
	var ret winapi.WNDCLASSEX

	ret.CbSize = uint32(unsafe.Sizeof(ret))
	ret.Style = 0
	ret.LpfnWndProc = syscall.NewCallback(func(hWnd winapi.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
		return tw.wndProc(hWnd, msg, wParam, lParam)
	})
	ret.CbClsExtra = 0
	ret.CbWndExtra = 0
	ret.HInstance = hInstance
	ret.HIcon = winapi.LoadIcon(hInstance, winapi.MAKEINTRESOURCE(132))
	ret.HCursor = winapi.LoadCursor(0, winapi.MAKEINTRESOURCE(winapi.IDC_ARROW))
	ret.HbrBackground = 0
	ret.LpszMenuName = nil
	ret.LpszClassName, _ = syscall.UTF16PtrFromString("meter-tip")

	return ret
}

func (tw *TipWindow) createWindow(hInstance winapi.HINSTANCE, windowClass winapi.WNDCLASSEX) winapi.HWND {
	windowTitlePtr, _ := syscall.UTF16PtrFromString("meter")
	return winapi.CreateWindowEx(
		winapi.WS_EX_NOACTIVATE,
		windowClass.LpszClassName,
		windowTitlePtr,
		winapi.WS_POPUP,
		0, 0, 1, 1,
		0, 0, hInstance, nil)
}

func (tw *TipWindow) wndProc(hWnd winapi.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case winapi.WM_PAINT:
		tw.onPaint.Invoke()

	case winapi.WM_DESTROY:
		winapi.PostQuitMessage(0)

	default:
		return winapi.DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
