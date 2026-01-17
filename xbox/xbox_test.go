package xbox

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/deorth-kku/go-common"
)

func TestXbox(t *testing.T) {
	devs := Find()
	if len(devs) == 0 {
		t.Error("failed to find devices")
		return
	}
	dev := devs[0]
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := dev.Listen(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestXboxKey(t *testing.T) {
	devs := Find()
	if len(devs) == 0 {
		t.Error("failed to find devices")
		return
	}
	dev := devs[0]
	f, err := dev.Open()
	if err != nil {
		t.Error(err)
		return
	}
	defer f.Close()

	var data, newdata [64]byte
	common.Must(f.Read(data[:]))

	newdata = data
	for newdata == data {
		common.Must(f.Read(newdata[:]))
	}
	for i, v := range newdata {
		if v == data[i] {
			continue
		}
		newv := data[i]
		for j := range 8 {
			mask := byte(1 << j)
			if v&mask != newv&mask {
				fmt.Printf("you've press key on byte:%d, bit:%d \n", i, j)
			}
		}

	}
}
