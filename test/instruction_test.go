package test

import (
	"testing"

	"github.com/tsmcalister/go-chip-8/emulator"
)

func TestNibble(t *testing.T) {
	nibbles := emulator.GetNibbles(emulator.Chip8Instruction(0x1234))
	for i, nibble := range nibbles {
		if nibble != byte(i+1) {
			t.Error("Nibble translations not working properly")
		}
	}
}
