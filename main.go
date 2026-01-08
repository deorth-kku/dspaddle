package main

import (
	"log/slog"
	"os"
	"os/signal"
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

	act := NewActionState(NewScanCodeAction(cfg.Keys))
	buttons := NewButtonActions(cfg.Keys)
	controllers := gods4.Find()
	if len(controllers) == 0 {
		slog.Warn("No connected DS4 controllers found")
		return
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		for i, c := range controllers {
			slog.Info("disconnecting controller", "num", i, "name", c.Name(), "type", c.ConnectionType(), "err", c.Disconnect())
		}
	}()

	var wg sync.WaitGroup
	for i, controller := range controllers {
		err := controller.Connect()
		if err != nil {
			slog.Warn("failed to connect to controller", "num", i, "err", err)
			continue
		}
		if controller.ConnectionType() != gods4.ConnectionTypeUSB {
			slog.Info("disconnecting non-USB controller", "num", i, "name", controller.Name(), "type", controller.ConnectionType(), "err", controller.Disconnect())
			continue
		}

		slog.Info("connected conntroller", "num", i, "name", controller.Name())

		controller.On(gods4.EventTouchpadSwipe, NewActionState(act).Callback)
		for k, v := range buttons.Range {
			controller.On(k, v)
		}
		wg.Go(func() {
			slog.Info("controller listen thread exit", "error", controller.Listen())
		})
	}
	wg.Wait()
}
