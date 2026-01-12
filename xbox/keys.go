package xbox

import "maps"

//go:generate stringer -type=Event -linecomment
type Event int

const (
	EventCrossPress         Event = iota // cross.press
	EventCrossRelease                    // cross.release
	EventCirclePress                     // circle.press
	EventCircleRelease                   // circle.release
	EventSquarePress                     // square.press
	EventSquareRelease                   // square.release
	EventTrianglePress                   // triangle.press
	EventTriangleRelease                 // triangle.release
	EventL1Press                         // l1.press
	EventL1Release                       // l1.release
	EventL2Press                         // l2.press
	EventL2Release                       // l2.release
	EventL3Press                         // l3.press
	EventL3Release                       // l3.release
	EventR1Press                         // r1.press
	EventR1Release                       // r1.release
	EventR2Press                         // r2.press
	EventR2Release                       // r2.release
	EventR3Press                         // r3.press
	EventR3Release                       // r3.release
	EventDpadUpPress                     // dpad_up.press
	EventDpadUpRelease                   // dpad_up.release
	EventDpadDownPress                   // dpad_down.press
	EventDpadDownRelease                 // dpad_down.release
	EventDpadLeftPress                   // dpad_left.press
	EventDpadLeftRelease                 // dpad_left.release
	EventDpadRightPress                  // dpad_right.press
	EventDpadRightRelease                // dpad_right.release
	EventSharePress                      // share.press
	EventShareRelease                    // share.release
	EventOptionsPress                    // options.press
	EventOptionsRelease                  // options.release
	EventTouchpadPress                   // touchpad.press
	EventTouchpadRelease                 // touchpad.release
	EventPSPress                         // ps.press
	EventPSRelease                       // ps.release
	EventLeftPaddlePress                 // left_paddle.press
	EventLeftPaddleRelease               // left_paddle.release
	EventRightPaddlePress                // right_paddle.press
	EventRightPaddleRelease              // right_paddle.release
)

const (
	FaceByte    = 3
	CrossBit    = 4
	CircleBit   = 5
	SquareBit   = 6
	TriangleBit = 7

	L1Byte = 3
	L1Bit  = 0
	L2Byte = 4
	L2Bit  = 1
	L3Byte = 2
	L3Bit  = 6

	R1Byte = 3
	R1Bit  = 1
	R2Byte = 5
	R2Bit  = 2
	R3Byte = 2
	R3Bit  = 7

	DpadByte     = 2
	DpadUpBit    = 0
	DpadDownBit  = 1
	DpadLeftBit  = 2
	DpadRightBit = 3

	ShareByte    = 14
	ShareBit     = 1
	OptionsByte  = 2
	OptionBit    = 4
	TouchPadByte = 2
	TouchPadBit  = 5
	PSByte       = 3
	PSbit        = 2

	PaddleByte     = 14
	LeftPaddleBit  = 2
	RightPaddleBit = 3
)

const keysCount = len(keys)

type keyoffset = struct {
	B int
	b uint8
}

var keys = [...]keyoffset{
	{
		B: FaceByte,
		b: CrossBit,
	},
	{
		B: FaceByte,
		b: CircleBit,
	},
	{
		B: FaceByte,
		b: SquareBit,
	},
	{
		B: FaceByte,
		b: TriangleBit,
	},

	{
		B: L1Byte,
		b: L1Bit,
	},
	{
		B: L2Byte,
		b: L2Bit,
	},
	{
		B: L3Byte,
		b: L3Bit,
	},

	{
		B: R1Byte,
		b: R1Bit,
	},
	{
		B: R2Byte,
		b: R2Bit,
	},
	{
		B: R3Byte,
		b: R3Bit,
	},

	{
		B: DpadByte,
		b: DpadUpBit,
	},
	{
		B: DpadByte,
		b: DpadDownBit,
	},
	{
		B: DpadByte,
		b: DpadLeftBit,
	},
	{
		B: DpadByte,
		b: DpadRightBit,
	},

	{
		B: ShareByte,
		b: ShareBit,
	},
	{
		B: OptionsByte,
		b: OptionBit,
	},
	{
		B: TouchPadByte,
		b: TouchPadBit,
	},
	{
		B: PSByte,
		b: PSbit,
	},

	{
		B: PaddleByte,
		b: LeftPaddleBit,
	},
	{
		B: PaddleByte,
		b: RightPaddleBit,
	},
}

var keysMap = maps.Collect(func(yield func(string, Event) bool) {
	for i := range Event(keysCount * 2) {
		if !yield(i.String(), i) {
			return
		}
	}
})

func ParseEvent[T ~string](s T) Event {
	v, ok := keysMap[string(s)]
	if !ok {
		return -1
	}
	return v
}
