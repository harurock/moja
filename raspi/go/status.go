package main

import "fmt"

const (
	maxHistory = 15
)

type AddCallback func(PrefId)
type DelCallback func(PrefId)
type SpecialCallback func(PrefId)

type LedStatus struct {
	history        [maxHistory]PrefId
	status         [PrefOkinawa + 1]int // valid range 1 - 47 (0:unused)
	idx            int
	addHandler     AddCallback
	delHandler     DelCallback
	specialHandler SpecialCallback
}

func NewLed(add AddCallback, del DelCallback, special SpecialCallback) *LedStatus {
	return &LedStatus{addHandler: add, delHandler: del, specialHandler: special}
}

func (led *LedStatus) Add(newPid PrefId) {
	targetPid := led.history[led.idx]
	if targetPid != PrefInvalid {
		led.status[targetPid]--
		if led.status[targetPid] == 0 {
			led.delHandler(targetPid)
		} else if led.status[targetPid] == 1 {
			led.addHandler(targetPid)
		}
	}
	led.history[led.idx] = newPid
	led.status[newPid]++
	if led.status[newPid] == 1 {
		led.addHandler(newPid)
	} else if led.status[newPid] == 2 {
		led.specialHandler(newPid)
	}
	led.idx++
	if led.idx == maxHistory {
		led.idx = 0
	}
}

func (led *LedStatus) String() string {
	return fmt.Sprintf("history:%v status:%v idx:%d \n", led.history, led.status, led.idx)
}
