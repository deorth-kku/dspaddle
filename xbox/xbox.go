package xbox

import (
	"fmt"
	"sync/atomic"

	"github.com/deorth-kku/go-common"
	"github.com/stamp/hid"
)

const (
	VendorID  = 0x3537
	ProductID = 0x1093

	PaddleByte = 14
	LeftBit    = 2
	RightBit   = 3
)

func IsBitSet(data []byte, byteIdx int, bitIdx uint8) bool {
	if byteIdx >= len(data) || bitIdx > 7 {
		return false
	}
	return (data[byteIdx] & (1 << bitIdx)) != 0
}

type Callback = func()

const (
	EventLeftPaddlePress = iota
	EventLeftPaddleRelease
	EventRightPaddlePress
	EventRightPaddleRelease

	bufSize = 64
)

type Device struct {
	*hid.Device
	*hid.DeviceInfo
	errch     <-chan error
	exit      atomic.Bool
	callbacks [4]Callback

	lastdata [bufSize]byte
}

func Find() []*Device {
	infos := hid.Enumerate(VendorID, ProductID)
	if len(infos) == 0 {
		return nil
	}
	devs := make([]*Device, 0, len(infos))
	for _, info := range infos {
		if info.UsagePage == 1 {
			continue
		}
		devs = append(devs, &Device{
			DeviceInfo: &info,
		})
	}
	return devs
}

func (d *Device) Connect() error {
	if d.Device != nil {
		return common.ErrorString("already connected")
	}
	f, err := d.DeviceInfo.Open()
	if err != nil {
		return fmt.Errorf("failed to open hid device: %w", err)
	}
	d.Device = f
	d.DeviceInfo = &f.DeviceInfo
	return nil
}

type keyoffset = struct {
	B int
	b uint8
}

var keys = [...]keyoffset{
	{
		B: PaddleByte,
		b: LeftBit,
	},
	{
		B: PaddleByte,
		b: RightBit,
	},
}

func (d *Device) Listen() error {
	if d.Device == nil {
		return common.ErrorString("not connected")
	}
	if d.errch != nil {
		return common.ErrorString("already listening")
	}
	errch := make(chan error, 1)
	d.errch = errch

	for !d.exit.Load() {
		var buf [bufSize]byte
		n, err := d.Read(buf[:])
		if err != nil {
			errch <- err
			return err
		}
		if n != bufSize {
			err := fmt.Errorf("incomplate read,only %d bytes was read", n)
			errch <- err
			return err
		}
		for i, key := range keys {
			last := IsBitSet(d.lastdata[:], key.B, key.b)
			this := IsBitSet(buf[:], key.B, key.b)
			switch {
			case !last && this: // press
				cb := d.callbacks[i*2]
				if cb != nil {
					cb()
				}
			case last && !this: // release
				cb := d.callbacks[i*2+1]
				if cb != nil {
					cb()
				}
			}
		}
		d.lastdata = buf
	}
	d.exit.Store(false)
	close(errch)
	return nil
}

func (d *Device) On(event int, cb Callback) {
	if event > len(d.callbacks) {
		return
	}
	d.callbacks[event] = cb
}

func (d *Device) Disconnect() error {
	d.exit.Store(true)
	return <-d.errch
}
