package main

import (
	"dspaddle/xbox"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/deorth-kku/go-common"
	"github.com/kpeu3i/gods4"
)

func init() {
	if IsDLL {
		go main()
	}
}

func main() {
	cfg, err := GetConfig()
	if err != nil {
		slog.Error("cannot get config", "err", err)
		return
	}

	err = common.SetLog(cfg.Log.File, cfg.Log.Level, "TEXT")
	if err != nil {
		slog.Warn("failed when setting log", "err", err)
	}

	if IsDLL {
		exe, err := os.Executable()
		if err == nil {
			slog.Info("loaded as dll", "exe", exe)
		} else {
			slog.Warn("failed to read exe path", "err", err)
		}

	}

	act := NewScanCodeAction(cfg.Keys)
	buttons := NewButtonActions(cfg.Keys)
	ds4s := gods4.Find()
	xboxs := xbox.Find()
	if len(ds4s)+len(xboxs) == 0 {
		slog.Warn("No connected DS4 controllers found")
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		for i, c := range ds4s {
			slog.Info("disconnecting dualshock 4", "num", i, "name", c.Name(), "type", c.ConnectionType(), "err", c.Disconnect())
		}
		for i, c := range xboxs {
			slog.Info("disconnecting dualshock 4", "num", i, "name", c.DeviceInfo.Product, "err", c.Disconnect())
		}
	}()

	var wg sync.WaitGroup
	for i, ds4 := range ds4s {
		if len(cfg.Slots) > 0 && !slices.Contains(cfg.Slots, i) {
			slog.Info("skipping ds4 controller", "num", i)
			continue
		}
		err := ds4.Connect()
		if err != nil {
			slog.Warn("failed to connect to ds4 controller", "num", i, "err", err)
			continue
		}
		if ds4.ConnectionType() != gods4.ConnectionTypeUSB {
			slog.Info("disconnecting non-USB ds4 controller", "num", i, "name", ds4.Name(), "type", ds4.ConnectionType(), "err", ds4.Disconnect())
			continue
		}

		slog.Info("connected ds4 conntroller", "num", i, "name", ds4.Name())

		ds4.On(gods4.EventTouchpadSwipe, NewActionState(act).Callback)
		for k, v := range buttons.Range {
			slog.Debug("register action", "event", k)
			ds4.On(k, v)
		}
		wg.Go(func() {
			slog.Info("ds4 controller listen thread exit", "error", ds4.Listen())
		})
	}
	for i, xb := range xboxs {
		if len(cfg.Slots) > 0 && !slices.Contains(cfg.Slots, i) {
			slog.Info("skipping xbox controller", "num", i)
			continue
		}
		err := xb.Connect()
		if err != nil {
			slog.Warn("failed to connect to xbox controller", "num", i, "err", err)
			continue
		}

		slog.Info("connected xbox conntroller", "num", i, "name", xb.Product)

		xb.On(xbox.EventLeftPaddlePress, act.touchleft)
		xb.On(xbox.EventLeftPaddleRelease, act.releaseleft)
		xb.On(xbox.EventRightPaddlePress, act.touchright)
		xb.On(xbox.EventRightPaddleRelease, act.releaseright)

		wg.Go(func() {
			slog.Info("xbox controller listen thread exit", "error", xb.Listen())
		})
	}

	wg.Wait()
}
