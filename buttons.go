package main

import (
	"github.com/deorth-kku/go-common"
	"github.com/kpeu3i/gods4"
)

type ButtonAction interface {
	Press(data any) error
	Release(data any) error
}

type scanCodeButton struct {
	press, release []INPUT
}

func (s *scanCodeButton) Press(_ any) error {
	SendInput(s.press)
	return nil
}

func (s *scanCodeButton) Release(_ any) error {
	SendInput(s.release)
	return nil
}

func NewButtonAction(key StringList) ButtonAction {
	act := scanCodeButton{
		press: parseINPUTs(key),
	}
	act.release = GroupConv(INPUT.ToUP, act.press)
	return &act
}

type gods4CB = func(any) error

type ButtonActions map[string]ButtonAction

func (ba ButtonActions) Range(y common.Yield2[gods4.Event, gods4CB]) {
	for k, v := range ba {
		if !y(gods4.Event(k+".press"), v.Press) {
			return
		}
		if !y(gods4.Event(k+".release"), v.Release) {
			return
		}
	}
}

func NewButtonActions(keys KeysConfig) ButtonActions {
	if keys.Buttons == nil {
		return nil
	}
	m := make(ButtonActions, len(keys.Buttons))
	for k, v := range keys.Buttons {
		m[k] = NewButtonAction(v)
	}
	return m
}
