package test

import (
	"testing"

	"github.com/tsmcalister/go-chip-8/emulator"
)

func TestInstructionByteConversion(t *testing.T) {
	val := emulator.Chip8Instruction(0x1234)
	byte1, byte2 := emulator.InstructionToBytes(val)
	if byte1 != byte(0x12) || byte2 != byte(0x34) {
		t.Errorf("bytes should be 0x12, 0x34 but are %v and %v", byte1, byte2)
	}
}

func TestByteInstructionConversion(t *testing.T) {
	byte1 := byte(0x12)
	byte2 := byte(0x34)
	instr := emulator.BytesToInstruction(byte1, byte2)

	if uint32(instr) != uint32(0x1234) {
		t.Errorf("instruction should be 0x1234 but is %v", uint32(instr))
	}
}
