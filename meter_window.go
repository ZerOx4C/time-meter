package main

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/cwchiu/go-winapi"
)

type MeterWindow struct {
	hInstance winapi.HINSTANCE
	hWnd      winapi.HWND
}

func (mw *MeterWindow) initialize() error {
	hInstance := winapi.GetModuleHandle(nil)
	windowClass := mw.createWindowClass(hInstance)

	if winapi.RegisterClassEx(&windowClass) == 0 {
		return errors.New("RegisterClassEx failed")
	}

	hWnd := mw.createWindow(hInstance, windowClass)
	if hWnd == 0 {
		return errors.New("createWindow failed")
	}

	mw.hInstance = hInstance
	mw.hWnd = hWnd

	return nil
}

func (mw *MeterWindow) finalize() error {
	return nil
}

func (mw *MeterWindow) show() {
	winapi.ShowWindow(mw.hWnd, winapi.SW_SHOW)
	winapi.SetTimer(mw.hWnd, 1, 30, 0)
}

func (mw *MeterWindow) createWindowClass(hInstance winapi.HINSTANCE) winapi.WNDCLASSEX {
	var ret winapi.WNDCLASSEX

	ret.CbSize = uint32(unsafe.Sizeof(ret))
	ret.Style = 0
	ret.LpfnWndProc = syscall.NewCallback(func(hWnd winapi.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
		return mw.wndProc(hWnd, msg, wParam, lParam)
	})
	ret.CbClsExtra = 0
	ret.CbWndExtra = 0
	ret.HInstance = hInstance
	ret.HIcon = winapi.LoadIcon(hInstance, winapi.MAKEINTRESOURCE(132))
	ret.HCursor = winapi.LoadCursor(0, winapi.MAKEINTRESOURCE(winapi.IDC_ARROW))
	ret.HbrBackground = winapi.HBRUSH(winapi.GetStockObject(winapi.BLACK_BRUSH))
	ret.LpszMenuName = nil
	ret.LpszClassName, _ = syscall.UTF16PtrFromString("meter")

	return ret
}

func (mw *MeterWindow) createWindow(hInstance winapi.HINSTANCE, windowClass winapi.WNDCLASSEX) winapi.HWND {
	windowTitlePtr, _ := syscall.UTF16PtrFromString("meter")
	return winapi.CreateWindowEx(
		winapi.WS_EX_NOACTIVATE,
		windowClass.LpszClassName,
		windowTitlePtr,
		winapi.WS_POPUP,
		0, 0, 0, 0,
		0, 0, hInstance, nil)
}

func (mw *MeterWindow) updateWindowLayout() {
	var workarea winapi.RECT
	winapi.SystemParametersInfo(winapi.SPI_GETWORKAREA, 0, unsafe.Pointer(&workarea), 0)
	winapi.SetWindowPos(
		mw.hWnd,
		winapi.HWND_TOPMOST,
		workarea.Left,
		workarea.Top,
		50,
		workarea.Bottom-workarea.Top,
		winapi.SWP_NOACTIVATE)
}

func (mw *MeterWindow) wndProc(hWnd winapi.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case winapi.WM_TIMER:
		mw.updateWindowLayout()
	case winapi.WM_DESTROY:
		winapi.PostQuitMessage(0)
	default:
		return winapi.DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
