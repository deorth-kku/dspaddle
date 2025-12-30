package main

import (
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
