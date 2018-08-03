package keyboard

import (
	"fmt"
	"time"

	"github.com/MarinX/keylogger"
)

var keyPressed = false
var keyBoardMap byte

func StartKeyboardDeamon() {
	ch := make(chan byte)
	go sniffKeyBoard(ch)
	for {
		select {
		case msg := <-ch:
			keyBoardMap = msg
			keyPressed = true
		default:
			keyPressed = false
			time.Sleep(8 * time.Millisecond)
		}
	}
}

func ReadKeyMap() (bool, byte) {
	return keyPressed, keyBoardMap
}

//StartKeyboardDeamon sends binary encoded values through channel
func sniffKeyBoard(ch chan byte) {
	devs, err := keylogger.NewDevices()
	if err != nil {
		fmt.Println(err)
		return
	}
	//our keyboard..on your system, it will be diffrent
	rd := keylogger.NewKeyLogger(devs[2])
	in, err := rd.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range in {
		//listen only key stroke event
		if i.Type == keylogger.EV_KEY {
			switch i.KeyString() {
			case "1":
				ch <- 0x1
			case "2":
				ch <- 0x2
			case "3":
				ch <- 0x3
			case "Q":
				ch <- 0x4
			case "W":
				ch <- 0x5
			case "E":
				ch <- 0x6
			case "A":
				ch <- 0x7
			case "S":
				ch <- 0x8
			case "D":
				ch <- 0x9
			case "Z":
				ch <- 0xA
			case "X":
				ch <- 0x0
			case "C":
				ch <- 0xB
			case "4":
				ch <- 0xC
			case "R":
				ch <- 0xD
			case "F":
				ch <- 0xE
			case "V":
				ch <- 0xF
			}
		}
	}
}
