package main

import (
	"runtime"
	"syscall"
	"unsafe"
)

var (
	modcomdlg32          = syscall.NewLazyDLL("comdlg32.dll")
	procGetOpenFileNameW = modcomdlg32.NewProc("GetOpenFileNameW")
)

type openFileNameW struct {
	LStructSize       uint32
	HwndOwner         uintptr
	HInstance         uintptr
	LpstrFilter       *uint16
	LpstrCustomFilter *uint16
	NMaxCustFilter    uint32
	NFilterIndex      uint32
	LpstrFile         *uint16
	NMaxFile          uint32
	LpstrFileTitle    *uint16
	NMaxFileTitle     uint32
	LpstrInitialDir   *uint16
	LpstrTitle        *uint16
	Flags             uint32
	NFileOffset       uint16
	NFileExtension    uint16
	LpstrDefExt       *uint16
	LCustData         uintptr
	LpfnHook          uintptr
	LpTemplateName    *uint16
}

const (
	ofnPathMustExist = 0x00000800
	ofnFileMustExist = 0x00001000
	ofnNoChangeDir   = 0x00000008
)

// OpenFileDialog opens a Windows file picker for audio files.
// Returns the selected file path or "" if cancelled.
func OpenFileDialog() string {
	fileBuf := make([]uint16, 260)

	// Build filter with embedded null separators for Win32 API.
	filterStr := "Audio Files (*.wav;*.mp3)\x00*.wav;*.mp3\x00All Files (*.*)\x00*.*\x00\x00"
	filter := make([]uint16, len(filterStr))
	for i := 0; i < len(filterStr); i++ {
		filter[i] = uint16(filterStr[i])
	}

	title, _ := syscall.UTF16PtrFromString("Select Alarm Sound")

	ofn := openFileNameW{
		LpstrFilter: &filter[0],
		LpstrFile:   &fileBuf[0],
		NMaxFile:    260,
		LpstrTitle:  title,
		Flags:       ofnPathMustExist | ofnFileMustExist | ofnNoChangeDir,
	}
	ofn.LStructSize = uint32(unsafe.Sizeof(ofn))

	ret, _, _ := procGetOpenFileNameW.Call(uintptr(unsafe.Pointer(&ofn)))
	runtime.KeepAlive(fileBuf)
	runtime.KeepAlive(filter)
	runtime.KeepAlive(title)

	if ret == 0 {
		return ""
	}
	return syscall.UTF16ToString(fileBuf)
}
