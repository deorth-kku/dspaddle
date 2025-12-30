package main

import (
	"strings"
	"unsafe"
)

const (
	KEYEVENTF_EXTENDEDKEY = 1 << iota
	KEYEVENTF_KEYUP
	KEYEVENTF_UNICODE
	KEYEVENTF_SCANCODE

	MAPVK_VK_TO_VSC = 0 // 从 VK 转换到 ScanCode
	MAPVK_VSC_TO_VK = 1 // 从 ScanCode 转换到 VK

	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2
)

func VkToScanCode(vk uint16) uint16 {
	ret, _, _ := procMapVirtualKey.Call(uintptr(vk), MAPVK_VK_TO_VSC)
	return uint16(ret)
}

func IsExtendedKey(vk uint16) bool {
	switch vk {
	case 0x21, // VK_PRIOR (Page Up)
		0x22, // VK_NEXT (Page Down)
		0x23, // VK_END
		0x24, // VK_HOME
		0x25, // VK_LEFT
		0x26, // VK_UP
		0x27, // VK_RIGHT
		0x28, // VK_DOWN
		0x2D, // VK_INSERT
		0x2E, // VK_DELETE
		0xA3, // VK_RCONTROL (右Ctrl)
		0xA5, // VK_RMENU (右Alt)
		0x90: // VK_NUMLOCK
		return true
	}
	return false
}

type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// KEYBDINPUT 结构体
type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// HARDWAREINPUT 结构体
type HARDWAREINPUT struct {
	UMsg    uint32
	WParamL uint16
	WParamH uint16
}

// INPUT 结构体 - Go中模拟 C Union
// 32位和64位下大小不同，这里使用Field对其进行内存布局模拟
type INPUT struct {
	Type uint32
	// 下面是 Union 的最大部分。
	// 在 64位 Windows 上，INPUT 结构体大小通常是 40 字节
	// Type(4) + Padding(4) + Union(32)
	// 我们直接用一个足够的字节数组或最大的结构体来填充
	Ki KEYBDINPUT
	// 注意：如果是用于鼠标，这里需要复杂的 unsafe.Pointer 转换，
	// 但仅用于键盘时，直接放 Ki 即可，因为它们共享内存起始位置。
	padding [8]byte // 额外的填充，确保结构体大小足够
}

const inputLen = unsafe.Sizeof(INPUT{})

func SendInput(inputs []INPUT) uint32 {
	if len(inputs) == 0 {
		return 0
	}
	ret, _, _ := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(inputLen),
	)
	return uint32(ret)
}

type ScanCodeKey struct {
	ScanCode uint16
	IsExtKey bool
	IsKeyUp  bool
}

func (s ScanCodeKey) ToUp() ScanCodeKey {
	s.IsKeyUp = true
	return s
}

func (i INPUT) ToUP() INPUT {
	i.Ki.DwFlags |= KEYEVENTF_KEYUP
	return i
}

func ParseScanCodeKey(s string) (ScanCodeKey, bool) {
	s = strings.ToUpper(s)
	vk, ok := keyMap[s]
	if !ok {
		return ScanCodeKey{}, false
	}
	if s == "NUM ENTER" {
		return ScanCodeKey{ScanCode: VkToScanCode(vk), IsExtKey: true}, true
	}
	return ScanCodeKey{ScanCode: VkToScanCode(vk), IsExtKey: IsExtendedKey(vk)}, true
}

func GroupConv[I, O any](parse func(i I) O, ins []I) (os []O) {
	if ins == nil {
		return
	}
	os = make([]O, len(ins))
	for i, in := range ins {
		os[i] = parse(in)
	}
	return
}

func GroupGet[I, O any](get func(i I) (O, bool), ins []I) (os []O, ok bool) {
	if ins == nil {
		return nil, true
	}
	os = make([]O, len(ins))
	for i, in := range ins {
		os[i], ok = get(in)
		if !ok {
			return nil, ok
		}
	}
	return os, true
}

func (key ScanCodeKey) ToInput() INPUT {
	var dwFlags uint32 = KEYEVENTF_SCANCODE
	if key.IsKeyUp {
		dwFlags |= KEYEVENTF_KEYUP
	}
	if key.IsExtKey {
		dwFlags |= KEYEVENTF_EXTENDEDKEY
	}

	return INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			WVk:     0,            // 使用扫描码时，Virtual Key 设为 0
			WScan:   key.ScanCode, // 这里填入硬件扫描码，如 0x11 代表 'W'
			DwFlags: dwFlags,
		},
	}
}

func SendScanCode(keys ...ScanCodeKey) {
	SendInput(GroupConv(ScanCodeKey.ToInput, keys))
}
