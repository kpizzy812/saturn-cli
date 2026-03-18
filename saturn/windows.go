//go:build windows

package main

import (
	"syscall"
	"unsafe"
)

func init() {
	// Enable Virtual Terminal Processing so ANSI color codes render in cmd.exe / PowerShell.
	// gookit/color skips this when it detects True-Color support (Windows 10 build >= 14931),
	// but cmd.exe still needs the flag set explicitly.
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")

	const enableVirtualTerminalProcessing uint32 = 0x0004

	for _, handle := range []syscall.Handle{syscall.Stdout, syscall.Stderr} {
		var mode uint32
		r, _, _ := getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
		if r != 0 {
			setConsoleMode.Call(uintptr(handle), uintptr(mode|enableVirtualTerminalProcessing))
		}
	}
}
