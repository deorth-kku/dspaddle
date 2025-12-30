package main

import (
	"sync/atomic"

	"github.com/kpeu3i/gods4"
)

type ActionIf interface {
	touchleft()
	releaseleft()
	touchright()
	releaseright()
}

type ActionState struct {
	leftPressed  atomic.Bool
	rightPressed atomic.Bool
	cbs          ActionIf
}

func (a *ActionState) touchleft() {
	if !a.leftPressed.Load() {
		a.leftPressed.Store(true)
		a.cbs.touchleft()
	}
}

func (a *ActionState) releaseleft() {
	if a.leftPressed.Load() {
		a.leftPressed.Store(false)
		a.cbs.releaseleft()
	}
}

func (a *ActionState) touchright() {
	if !a.rightPressed.Load() {
		a.rightPressed.Store(true)
		a.cbs.touchright()
	}
}

func (a *ActionState) releaseright() {
	if a.rightPressed.Load() {
		a.rightPressed.Store(false)
		a.cbs.releaseright()
	}
}

var (
	touchLeft    = gods4.Touch{IsActive: true, X: 199, Y: 80}
	releaseLeft  = gods4.Touch{IsActive: false, X: 199, Y: 80}
	touchRight   = gods4.Touch{IsActive: true, X: 175, Y: 80}
	releaseRight = gods4.Touch{IsActive: false, X: 175, Y: 80}
)

func (a *ActionState) Callback(data any) error {
	touch := data.(gods4.Touchpad)
	if touch.Swipe[0].X == touch.Swipe[1].X && touch.Swipe[0].Y == touch.Swipe[1].Y {
		touch.Swipe = touch.Swipe[:1]
	}
	for _, swipe := range touch.Swipe {
		switch swipe {
		case touchLeft:
			a.touchleft()
		case touchRight:
			a.touchright()
		case releaseLeft:
			a.releaseleft()
		case releaseRight:
			a.releaseright()
		}
	}
	return nil
}

func NewActionState(do ActionIf) *ActionState {
	return &ActionState{cbs: do}
}

type scanCodeAction struct {
	leftdown, leftup, rightdown, rightup []INPUT
}

func (s scanCodeAction) touchleft() {
	SendInput(s.leftdown)
}

func (s scanCodeAction) touchright() {
	SendInput(s.rightdown)
}

func (s scanCodeAction) releaseleft() {
	SendInput(s.leftup)
}

func (s scanCodeAction) releaseright() {
	SendInput(s.rightup)
}

func parseINPUTs(list StringList) []INPUT {
	if list == nil {
		return nil
	}
	inputs := make([]INPUT, 0, len(list))
	for _, v := range list {
		key, ok := ParseScanCodeKey(v)
		if !ok {
			continue
		}
		inputs = append(inputs, key.ToInput())
	}
	return inputs
}

func NewScanCodeAction(keys KeysConfig) ActionIf {
	act := scanCodeAction{
		leftdown:  parseINPUTs(keys.Left),
		rightdown: parseINPUTs(keys.Right),
	}
	act.leftup = GroupConv(INPUT.ToUP, act.leftdown)
	act.rightup = GroupConv(INPUT.ToUP, act.rightdown)
	return act
}
