package main

var keyMap = map[string]uint16{
	// Letters
	"A": 'A', "B": 'B', "C": 'C', "D": 'D',
	"E": 'E', "F": 'F', "G": 'G', "H": 'H',
	"I": 'I', "J": 'J', "K": 'K', "L": 'L',
	"M": 'M', "N": 'N', "O": 'O', "P": 'P',
	"Q": 'Q', "R": 'R', "S": 'S', "T": 'T',
	"U": 'U', "V": 'V', "W": 'W', "X": 'X',
	"Y": 'Y', "Z": 'Z',

	// Number row
	"0": '0', "1": '1', "2": '2', "3": '3', "4": '4',
	"5": '5', "6": '6', "7": '7', "8": '8', "9": '9',

	// Function keys
	"F1": 0x70, "F2": 0x71, "F3": 0x72, "F4": 0x73,
	"F5": 0x74, "F6": 0x75, "F7": 0x76, "F8": 0x77,
	"F9": 0x78, "F10": 0x79, "F11": 0x7A, "F12": 0x7B,

	// Control keys
	"ESC":       0x1B,
	"TAB":       0x09,
	"CAPSLOCK":  0x14,
	"LSHIFT":    0xA0,
	"RSHIFT":    0xA1,
	"LCTRL":     0xA2,
	"RCTRL":     0xA3,
	"LALT":      0xA4,
	"RALT":      0xA5,
	"SPACE":     0x20,
	"ENTER":     0x0D,
	"BACKSPACE": 0x08,

	// Navigation
	"INSERT":   0x2D,
	"DELETE":   0x2E,
	"HOME":     0x24,
	"END":      0x23,
	"PAGEUP":   0x21,
	"PAGEDOWN": 0x22,

	// Arrow keys
	"LEFT":  0x25,
	"UP":    0x26,
	"RIGHT": 0x27,
	"DOWN":  0x28,

	// System keys
	"PRINTSCREEN": 0x2C,
	"SCROLLLOCK":  0x91,
	"PAUSE":       0x13,

	// OEM / punctuation
	";":  0xBA,
	"=":  0xBB,
	",":  0xBC,
	"-":  0xBD,
	".":  0xBE,
	"/":  0xBF,
	"`":  0xC0,
	"[":  0xDB,
	"\\": 0xDC,
	"]":  0xDD,
	"'":  0xDE,

	// Numpad
	"NUM 0": 0x60,
	"NUM 1": 0x61,
	"NUM 2": 0x62,
	"NUM 3": 0x63,
	"NUM 4": 0x64,
	"NUM 5": 0x65,
	"NUM 6": 0x66,
	"NUM 7": 0x67,
	"NUM 8": 0x68,
	"NUM 9": 0x69,

	"NUM *": 0x6A,
	"NUM +": 0x6B,
	"NUM -": 0x6D,
	"NUM .": 0x6E,
	"NUM /": 0x6F,

	"NUM ENTER": 0x0D,
}
