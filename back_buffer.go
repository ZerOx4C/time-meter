package main

import (
	winapi2 "time-meter/winapi"

	"github.com/cwchiu/go-winapi"
)

type BackBuffer struct {
	frontDc       winapi.HDC
	clientRect    RECT
	backDc        winapi.HDC
	backBitmap    winapi.HBITMAP
	oldBackBitmap winapi.HGDIOBJ
	began         bool
}

func (bb *BackBuffer) begin(hWnd winapi.HWND, hdc winapi.HDC) winapi.HDC {
	if bb.began {
		panic("invalid operation.")
	}

	bb.frontDc = hdc
	winapi.GetClientRect(hWnd, bb.clientRect.unwrap())

	bb.backDc = winapi.CreateCompatibleDC(hdc)
	bb.backBitmap = winapi2.CreateCompatibleBitmap(hdc, bb.clientRect.width(), bb.clientRect.height())
	bb.oldBackBitmap = winapi.SelectObject(bb.backDc, winapi.HGDIOBJ(bb.backBitmap))
	bb.began = true

	return bb.backDc
}

func (bb *BackBuffer) end() {
	if !bb.began {
		panic("invalid operation.")
	}

	winapi.BitBlt(
		bb.frontDc, 0, 0, bb.clientRect.width(), bb.clientRect.height(),
		bb.backDc, 0, 0,
		winapi.SRCCOPY)

	winapi.SelectObject(bb.backDc, bb.oldBackBitmap)
	winapi.DeleteDC(bb.backDc)
	winapi.DeleteObject(winapi.HGDIOBJ(bb.backBitmap))

	bb.frontDc = 0
	bb.backDc = 0
	bb.backBitmap = 0
	bb.oldBackBitmap = 0
	bb.began = false
}
