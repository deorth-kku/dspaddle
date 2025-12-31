package main

import (
	"log/slog"
	"testing"
	"time"

	"github.com/deorth-kku/go-common"
	"github.com/kpeu3i/gods4"
)

func TestMain(t *testing.T) {
	main()
}

type testAction struct{}

var (
	altdown = common.MustOk(ParseScanCodeKey("LALT"))
	altup   = altdown.ToUp()
	tabdown = common.MustOk(ParseScanCodeKey("TAB"))
	tabup   = tabdown.ToUp()
)

func (testAction) touchleft() {
	slog.Info("touchleft")
	SendScanCode(altdown)

}

func (testAction) releaseleft() {
	slog.Info("releaseleft")
	SendScanCode(altup)

}

func (testAction) touchright() {
	slog.Info("touchright")
	SendScanCode(tabdown)
}

func (testAction) releaseright() {
	slog.Info("releaseright")
	SendScanCode(tabup)
}

func TestAltTab(t *testing.T) {
	common.SetLog("", "DEBUG", "TEXT")
	controller := gods4.Find()[0]
	err := controller.Connect()
	if err != nil {
		t.Error(err)
	}
	controller.On(gods4.EventTouchpadSwipe, NewActionState(testAction{}).Callback)
	time.AfterFunc(29*time.Second, func() {
		controller.Disconnect()
	})
	err = controller.Listen()
	if err != nil {
		t.Error(err)
	}
}

func TestComboKey(t *testing.T) {
	down := []INPUT{altdown.ToInput(), tabdown.ToInput()}
	up := []INPUT{altup.ToInput(), altup.ToInput()}
	SendInput(down)
	time.Sleep(time.Microsecond)
	SendInput(up)

	time.Sleep(time.Second)

	SendInput(down)
	time.Sleep(time.Microsecond)
	SendInput(up)
}
