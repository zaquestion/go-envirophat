// leds is direct port of https://github.com/pimoroni/enviro-phat/blob/master/library/envirophat/leds.py
package leds

import (
	"github.com/stianeikeland/go-rpio"
	"log"
)

var (
	state rpio.State
	pin   rpio.Pin
)

func init() {
	if err := rpio.Open(); err != nil {
		log.Fatal(err)
	}

	pin = rpio.Pin(4)
	pin.Output()
	state = pin.Read()
}

func On() {
	pin.High()
	state = rpio.High
}

func Off() {
	pin.Low()
	state = rpio.Low
}

func Toggle() {
	pin.Toggle()
}

func IsOff() bool {
	return state == rpio.Low
}

func IsOn() bool {
	return state == rpio.High
}
