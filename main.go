package main

import (
	"github.com/tsmcalister/go-chip-8/emulator"
)

func main() {
	emulator.InitMemory()
	emulator.LoadProgram("programs/life.ch8")
	for {
		emulator.EmulateStep()
	}
}
