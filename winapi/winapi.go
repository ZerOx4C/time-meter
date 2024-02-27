package winapi

import (
	"syscall"

	"github.com/cwchiu/go-winapi"
)

var (
	// Library
	libgdi32 uintptr

	// Functions
	createCompatibleBitmap uintptr
)

func init() {
	// Library
	libgdi32 = winapi.MustLoadLibrary("gdi32.dll")

	// Functions
	createCompatibleBitmap = winapi.MustGetProcAddress(libgdi32, "CreateCompatibleBitmap")
}

func CreateCompatibleBitmap(hdc winapi.HDC, cx, cy int32) winapi.HBITMAP {
	ret, _, _ := syscall.SyscallN(createCompatibleBitmap,
		uintptr(hdc),
		uintptr(cx),
		uintptr(cy),
	)

	return winapi.HBITMAP(ret)
}
