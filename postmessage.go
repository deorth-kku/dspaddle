package main

import "strings"

const (
	WM_KEYDOWN = 0x0100
	WM_KEYUP   = 0x0101
	WM_CHAR    = 0x0102
)

type VirtualKey ScanCodeKey

func (vk VirtualKey) ToUp() VirtualKey {
	vk.IsKeyUp = true
	return vk
}

func ParseVK(s string) (VirtualKey, bool) {
	s = strings.ToUpper(s)
	vk, ok := keyMap[s]
	if !ok {
		return VirtualKey{}, false
	}
	if s == "NUM ENTER" {
		return VirtualKey{Code: vk, IsExtKey: true}, true
	}
	return VirtualKey{Code: vk, IsExtKey: IsExtendedKey(vk)}, true
}

func SendChar(hwnd uintptr, vks ...VirtualKey) {
	for _, vk := range vks {
		procPostMessage.Call(hwnd, WM_CHAR, uintptr(vk.Code), 0)
	}

}

func SendVK(hwnd uintptr, vks ...VirtualKey) {
	for _, vk := range vks {
		if vk.IsKeyUp {
			procPostMessage.Call(hwnd, WM_KEYUP, uintptr(vk.Code), 0)
		} else {
			procPostMessage.Call(hwnd, WM_KEYDOWN, uintptr(vk.Code), 0)
		}
	}
}
