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
	SendScanCode(altdown)

}

func (testAction) releaseleft() {
	SendScanCode(altup)

}

func (testAction) touchright() {
	SendScanCode(tabdown)
}

func (testAction) releaseright() {
	SendScanCode(tabup)
}

func TestAltTab(t *testing.T) {
	controller := gods4.Find()[0]
	err := controller.Connect()
	if err != nil {
		t.Error(err)
	}
	controller.On(gods4.EventTouchpadSwipe, NewActionState(testAction{}).Callback)
	time.AfterFunc(10*time.Second, func() {
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

func TestSendVK(t *testing.T) {
	hwnd, _, _ := procGetForegroundWindow.Call()

	key, _ := ParseVK("T")

	for range 10 {
		SendVK(hwnd, key)
		time.Sleep(time.Millisecond)
		SendVK(hwnd, key.ToUp())
		time.Sleep(time.Second)
	}
}

func TestSendChar(t *testing.T) {
	hwnd, _, _ := procGetForegroundWindow.Call()

	key, _ := ParseVK("T")

	for range 10 {
		SendChar(hwnd, key)
		time.Sleep(time.Second)
	}
}

func TestButtions(t *testing.T) {
	controller := gods4.Find()[0]
	err := controller.Connect()
	if err != nil {
		t.Error(err)
	}
	keys := []gods4.Event{
		"cross.press",
		"cross.release",
		"circle.press",
		"circle.release",
		"square.press",
		"square.release",
		"triangle.press",
		"triangle.release",
		"l1.press",
		"l1.release",
		"l2.press",
		"l2.release",
		"l3.press",
		"l3.release",
		"r1.press",
		"r1.release",
		"r2.press",
		"r2.release",
		"r3.press",
		"r3.release",
		"dpad_up.press",
		"dpad_up.release",
		"dpad_down.press",
		"dpad_down.release",
		"dpad_left.press",
		"dpad_left.release",
		"dpad_right.press",
		"dpad_right.release",
		"share.press",
		"share.release",
		"options.press",
		"options.release",
		"touchpad.swipe",
		"touchpad.press",
		"touchpad.release",
		"ps.press",
		"ps.release",
		"left_stick.move",
		"right_stick.move",
		"accelerometer.update",
		"gyroscope.update",
		"battery.update",
	}
	for _, e := range keys {
		controller.On(e, func(_ any) error {
			slog.Warn(string(e))
			return nil
		})
	}

	time.AfterFunc(10*time.Second, func() {
		controller.Disconnect()
	})
	err = controller.Listen()
	if err != nil {
		t.Error(err)
	}
}
