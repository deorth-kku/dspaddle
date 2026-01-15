package xbox

import (
	"fmt"
	"sync/atomic"

	"github.com/deorth-kku/go-common"
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
	hid.Device
	*hid.DeviceInfo
	errch     <-chan error
	exit      atomic.Bool
	callbacks [keysCount * 2]Callback

	lastdata [bufSize]byte
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

func (d *Device) Connect() error {
	if d.Device != nil {
		return common.ErrorString("already connected")
	}
	f, err := d.DeviceInfo.Open()
	if err != nil {
		return fmt.Errorf("failed to open hid device: %w", err)
	}
	d.Device = f
	return nil
}

const timeout = 3 // ms

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
		n, err := d.ReadTimeout(buf[:], timeout)
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

func (d *Device) On(event Event, cb Callback) {
	if int(event) > len(d.callbacks) {
		return
	}
	d.callbacks[event] = cb
}

func (d *Device) Disconnect() error {
	if d.Device == nil {
		return common.ErrorString("not connected")
	}
	err := d.Device.Close()
	if err != nil {
		return err
	}
	if d.errch == nil {
		return nil
	}
	d.exit.Store(true)
	return <-d.errch
}
