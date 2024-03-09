package winapi

import (
	"syscall"
	"unsafe"

	"github.com/cwchiu/go-winapi"
)

const (
	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002
	MB_TOPMOST   = 0x00040000
)

var (
	// Library
	libgdi32  uintptr
	libuser32 uintptr

	// Functions
	createCompatibleBitmap     uintptr
	enumDisplayMonitors        uintptr
	setLayeredWindowAttributes uintptr
)

func init() {
	// Library
	libgdi32 = winapi.MustLoadLibrary("gdi32.dll")
	libuser32 = winapi.MustLoadLibrary("user32.dll")

	// Functions
	createCompatibleBitmap = winapi.MustGetProcAddress(libgdi32, "CreateCompatibleBitmap")
	enumDisplayMonitors = winapi.MustGetProcAddress(libuser32, "EnumDisplayMonitors")
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

func EnumDisplayMonitors(hdc winapi.HDC, lprcClip *winapi.RECT, lpfnEnum uintptr, dwData uintptr) bool {
	ret, _, _ := syscall.SyscallN(enumDisplayMonitors,
		uintptr(hdc),
		uintptr(unsafe.Pointer(lprcClip)),
		uintptr(lpfnEnum),
		uintptr(dwData),
	)

	return ret != 0
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
