package main

import (
	"errors"
	"syscall"
	"time-meter/setting"
	winapi2 "time-meter/winapi"
	"unsafe"

	"github.com/cwchiu/go-winapi"
)

type MeterWindow struct {
	hInstance          winapi.HINSTANCE
	hWnd               winapi.HWND
	settings           *setting.Settings
	bound              RECT
	lastCursorPos      POINT
	lastMenuId         MenuId
	onPaint            EventHandler
	onMouseMove        EventHandler
	onMouseEnter       EventHandler
	onMouseLeave       EventHandler
	onMouseRightClick  EventHandler
	onPopupMenuCommand EventHandler
}

const (
	EID_UPDATE_CHART = 1 + iota
	EID_WATCH_MOUSE
)

func (mw *MeterWindow) Initialize() error {
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

func (mw *MeterWindow) Finalize() error {
	return nil
}

func (mw *MeterWindow) Show() {
	winapi.ShowWindow(mw.hWnd, winapi.SW_SHOW)
	mw.updateLayout()

	winapi.SetTimer(mw.hWnd, EID_UPDATE_CHART, 1000/2, 0)
	winapi.SetTimer(mw.hWnd, EID_WATCH_MOUSE, 1000/30, 0)
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
	ret.HbrBackground = 0
	ret.LpszMenuName = nil
	ret.LpszClassName, _ = syscall.UTF16PtrFromString("meter")

	return ret
}

func (mw *MeterWindow) createWindow(hInstance winapi.HINSTANCE, windowClass winapi.WNDCLASSEX) winapi.HWND {
	windowTitlePtr, _ := syscall.UTF16PtrFromString("meter")
	return winapi.CreateWindowEx(
		winapi.WS_EX_NOACTIVATE|winapi.WS_EX_LAYERED,
		windowClass.LpszClassName,
		windowTitlePtr,
		winapi.WS_POPUP,
		0, 0, 0, 0,
		0, 0, hInstance, nil)
}

func (mw *MeterWindow) updateLayout() {
	workareaList := mw.getMonitorRectList()
	workarea := workareaList[0]

	if index := mw.settings.TargetDisplayIndex; 0 <= index && index < len(workareaList) {
		workarea = workareaList[index]
	}

	mw.bound.Left = workarea.Left
	mw.bound.Top = workarea.Top
	mw.bound.Right = workarea.Left + int32(mw.settings.MeterWidth)
	mw.bound.Bottom = workarea.Bottom

	winapi.SetWindowPos(
		mw.hWnd,
		winapi.HWND_TOPMOST,
		mw.bound.Left,
		mw.bound.Top,
		mw.bound.Width(),
		mw.bound.Height(),
		winapi.SWP_NOACTIVATE)
}

func (mw *MeterWindow) getMonitorRectList() []RECT {
	ret := []RECT{}
	rectChan := make(chan RECT)

	go winapi2.EnumDisplayMonitors(0, nil, syscall.NewCallback(func(hMonitor winapi.HMONITOR) uintptr {
		var info winapi.MONITORINFO
		info.CbSize = uint32(unsafe.Sizeof(info))
		winapi.GetMonitorInfo(hMonitor, &info)

		rectChan <- (RECT)(info.RcWork)

		return 1
	}), 0)

	count := int(winapi.GetSystemMetrics(winapi.SM_CMONITORS))

	for len(ret) < count {
		ret = append(ret, <-rectChan)
	}

	close(rectChan)

	return ret
}

func (mw *MeterWindow) watchMouse() {
	var currentCursorPos POINT
	winapi.GetCursorPos(currentCursorPos.Unwrap())

	isHit := mw.bound.Contains(currentCursorPos)
	wasHit := mw.bound.Contains(mw.lastCursorPos)
	mw.lastCursorPos = currentCursorPos

	if isHit == wasHit {
		return
	}

	if isHit {
		winapi2.SetLayeredWindowAttributes(mw.hWnd, 0, 255, winapi2.LWA_ALPHA)
		mw.onMouseEnter.Invoke()

	} else {
		winapi2.SetLayeredWindowAttributes(mw.hWnd, 0, mw.settings.MeterOpacity, winapi2.LWA_ALPHA)
		mw.onMouseLeave.Invoke()
	}
}

func (mw *MeterWindow) wndProc(hWnd winapi.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	switch msg {
	case winapi.WM_PAINT:
		mw.onPaint.Invoke()

	case winapi.WM_MOUSEMOVE:
		mw.onMouseMove.Invoke()

	case winapi.WM_RBUTTONUP:
		mw.onMouseRightClick.Invoke()

	case winapi.WM_COMMAND:
		fromMenu := winapi.HIWORD(uint32(wParam)) == 0
		menuId := MenuId(winapi.LOWORD(uint32(wParam)))

		if fromMenu {
			mw.lastMenuId = menuId
			mw.onPopupMenuCommand.Invoke()
		}

	case winapi.WM_DISPLAYCHANGE:
		mw.updateLayout()

	case winapi.WM_TIMER:
		switch wParam {
		case EID_UPDATE_CHART:
			winapi.InvalidateRect(hWnd, nil, true)

		case EID_WATCH_MOUSE:
			mw.watchMouse()
		}

	case winapi.WM_DESTROY:
		winapi.PostQuitMessage(0)

	default:
		return winapi.DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
