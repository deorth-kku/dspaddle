package xbox

import (
	"log/slog"
	"testing"
	"time"
)

func TestXbox(t *testing.T) {
	devs := Find()
	if len(devs) == 0 {
		t.Error("failed to find devices")
		return
	}
	dev := devs[0]
	err := dev.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	dev.On(EventLeftPaddlePress, func() {
		slog.Warn("left pressed")
	})
	dev.On(EventLeftPaddleRelease, func() {
		slog.Warn("left released")
	})
	dev.On(EventRightPaddlePress, func() {
		slog.Warn("right pressed")
	})
	dev.On(EventRightPaddleRelease, func() {
		slog.Warn("right release")
	})

	time.AfterFunc(10*time.Second, func() {
		err := dev.Disconnect()
		if err != nil {
			t.Error(err)
		}
	})
	err = dev.Listen()
	if err != nil {
		t.Error(err)
	}
}
