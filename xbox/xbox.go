package xbox

import (
	"context"
	"fmt"

	"github.com/karalabe/hid"
)

const (
	VendorID  = 0x3537
	ProductID = 0x1093
)

func IsBitSet(data []byte, byteIdx int, bitIdx uint8) bool {
	if byteIdx >= len(data) || bitIdx > 7 {
		return false
	}
	return (data[byteIdx] & (1 << bitIdx)) != 0
}

type Callback = func()

const bufSize = 64

type Device struct {
	*hid.DeviceInfo
	callbacks [keysCount * 2]Callback
	lastdata  [bufSize]byte
}

func Find() []*Device {
	infos, err := hid.Enumerate(VendorID, ProductID)
	if err != nil || len(infos) == 0 {
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

const timeout = 1000 // ms

func (d *Device) Listen(ctx context.Context) error {
	dev, err := d.Open()
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var buf [bufSize]byte
			n, err := dev.ReadTimeout(buf[:], timeout)
			if err != nil {
				return err
			}
			if n != bufSize {
				err := fmt.Errorf("incomplate read,only %d bytes was read", n)
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
	}
}

func (d *Device) On(event Event, cb Callback) {
	if int(event) > len(d.callbacks) {
		return
	}
	d.callbacks[event] = cb
}
