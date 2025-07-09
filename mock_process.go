package main

import (
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
	getWindowTextW      = user32.NewProc("GetWindowTextW")
)

func GetMockProcesses() []string {
	// Get the handle of the active window
	hwnd, _, _ := getForegroundWindow.Call()

	// Prepare buffer for window title
	buf := make([]uint16, 200)
	getWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))

	// Convert UTF-16 to string
	runes := utf16.Decode(buf)
	windowTitle := string(runes)

	var result []string

	lower := strings.ToLower(windowTitle)

	if strings.Contains(lower, "chrome") {
		result = append(result, "Chrome")
	}
	if strings.Contains(lower, "youtube") {
		result = append(result, "YouTube")
	}
	if strings.Contains(lower, "instagram") {
		result = append(result, "Instagram")
	}

	return result
}
