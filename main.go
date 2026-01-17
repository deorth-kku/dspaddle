package main

import (
	"context"
	"dspaddle/xbox"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/deorth-kku/go-common"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			slog.Info("recived signal, exiting normally")
			return
		default:
			xboxs := xbox.Find()
			if len(xboxs) == 0 {
				const interval = 1 * time.Second
				slog.Warn("No connected xbox controllers found", "sleep", interval)
				time.Sleep(interval)
				continue
			}

			var wg sync.WaitGroup
			for i, xb := range xboxs {
				if len(cfg.Slots) > 0 && !slices.Contains(cfg.Slots, i) {
					slog.Info("skipping xbox controller", "num", i)
					continue
				}
				slog.Info("use xbox conntroller", "num", i, "name", xb.Product)

				xb.On(xbox.EventLeftPaddlePress, act.touchleft)
				xb.On(xbox.EventLeftPaddleRelease, act.releaseleft)
				xb.On(xbox.EventRightPaddlePress, act.touchright)
				xb.On(xbox.EventRightPaddleRelease, act.releaseright)

				for k, v := range buttons.Range {
					slog.Debug("register action", "event", k)
					event := xbox.ParseEvent(k)
					if event < 0 {
						slog.Warn("invalid event, skipping", "event", event)
						continue
					}
					xb.On(event, func() {
						v(nil)
					})
				}

				wg.Go(func() {
					slog.Info("xbox controller listen thread exit", "error", xb.Listen(ctx))
				})
			}
			wg.Wait()
		}
	}
}
