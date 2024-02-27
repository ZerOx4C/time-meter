package winapi

import (
	"syscall"

	"github.com/cwchiu/go-winapi"
)

const (
	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002
)

var (
	// Library
	libgdi32  uintptr
	libuser32 uintptr

	// Functions
	createCompatibleBitmap     uintptr
	setLayeredWindowAttributes uintptr
)

func init() {
	// Library
	libgdi32 = winapi.MustLoadLibrary("gdi32.dll")
	libuser32 = winapi.MustLoadLibrary("user32.dll")

	// Functions
	createCompatibleBitmap = winapi.MustGetProcAddress(libgdi32, "CreateCompatibleBitmap")
	setLayeredWindowAttributes = winapi.MustGetProcAddress(libuser32, "SetLayeredWindowAttributes")
}

func CreateCompatibleBitmap(hdc winapi.HDC, cx, cy int32) winapi.HBITMAP {
	ret, _, _ := syscall.SyscallN(createCompatibleBitmap,
		uintptr(hdc),
		uintptr(cx),
		uintptr(cy),
	)

	return winapi.HBITMAP(ret)
}

func SetLayeredWindowAttributes(hwnd winapi.HWND, crKey winapi.COLORREF, bAlpha byte, dwFlags winapi.DWORD) bool {
	ret, _, _ := syscall.SyscallN(setLayeredWindowAttributes,
		uintptr(hwnd),
		uintptr(crKey),
		uintptr(bAlpha),
		uintptr(dwFlags),
	)

	return ret != 0
}
