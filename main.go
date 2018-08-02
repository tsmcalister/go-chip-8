package main

import (
	"os"
	"runtime"

	"github.com/tsmcalister/go-chip-8/emulator"
)

func main() {
	runtime.GOMAXPROCS(4)
	filePath := "programs/"
	if len(os.Args) > 1 {
		filePath += os.Args[1]
	} else {
		filePath += "pic.ch8"
	}

	emulator.InitMemory()
	emulator.LoadProgram(filePath)
	for {
		emulator.EmulateStep()
	}
	/*

		emulator.ResetEmulator()
		program := []uint16{
			0xD015, // draw sprite starting at index
			0x6005, // set vx = 0 to 5
			0xA006, // set index to 6
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA00A, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA00F, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA014, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA019, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA01E, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA023, // set index to 12
			0xD015, // draw sprite
			0x7106, // add 6 to y
			0x6000, // set vx = to 0
			0xA028, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA02D, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA032, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA037, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA03C, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA041, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA046, // set index to 12
			0xD015, // draw sprite
			0x7005, // add 5 to vx
			0xA04B, // set index to 12
			0xD015, // draw sprite

		}
		emulator.LoadProgramHex(program)
		i := 0
		for i < len(program) {
			emulator.EmulateStep()
			i++
		}
	*/

}
