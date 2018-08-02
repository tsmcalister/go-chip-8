package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tsmcalister/go-chip-8/emulator"
)

func main() {
	args := os.Args

	if len(args) > 1 {
		defer func() {
			recover()
		}()
		disassemble(args[1])
	} else {
		fmt.Println("please supply chip 8 binary")
	}
}

func disassemble(fileName string) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic("Could not open file specified")
	}
	for i := 0; i < len(dat)/2; i++ {
		instr := emulator.Chip8Instruction(uint16(dat[i*2])<<8 | uint16(dat[i*2+1]))
		fmt.Printf("0x%.4x\n", instr)
	}
}
