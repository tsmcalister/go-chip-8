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

// calls a subroutine and returns from it
func Test00EE_2NNN(t *testing.T) {
	emulator.ResetEmulator()
	program := []uint16{
		0x2214, // 0x200 Call subroutine at 0x20A ->
		0xA000, // 0x202
		0xA000, // 0x204
		0xA000, // 0x206
		0xA000, // 0x208
		0xA000, // 0x20A
		0xA000, // 0x20C
		0xA000, // 0x20E
		0xA000, // 0x210
		0xA000, // 0x212
		0x00E0, // 0x214 Clear Screen
		0x00EE, // 0x216 Return from subroutine
	}
	emulator.LoadProgramHex(program)

	emulator.EmulateStep() // Call subroutine
	if emulator.Stack()[0] == 0x0 {
		t.Error("No address added to stack")
	}

	if emulator.StackPointer() != 0x1 {
		t.Error("Stack pointer not incremented")
	}

	if emulator.ProgramCounter() != 0x214 {
		t.Errorf("Program counter not set to 0x20A, it is: %x\n", emulator.ProgramCounter())
	}

	emulator.EmulateStep() // Clear Screen
	emulator.EmulateStep() // Return from Subroutine
	if emulator.StackPointer() != 0 {
		t.Error("Stack pointer wasn't decreased.")
	}
	if emulator.ProgramCounter() != 0x202 {
		t.Error("didn't return to proper address")
	}
	emulator.EmulateStep()
}

// Should clear the screen
func Test00E0(t *testing.T) {
	emulator.ResetEmulator()
	program := []uint16{
		0x00E0,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			if emulator.ReadScreen(byte(i), byte(j)) {
				t.Error("Screen wasn't cleared.")
			}
		}
	}
}

func Test3xkk(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x66)
	program := []uint16{
		0x3066,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ProgramCounter() != 0x204 {
		t.Error("3xkk Values are not compared properly")
	}
}

func Test4xkkk(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x65)
	program := []uint16{
		0x4066,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ProgramCounter() != 0x204 {
		t.Error("4xkk Values are not compared properly")
	}
}

func Test5xy0(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x66)
	emulator.StoreRegisterValue(0x1, 0x66)
	program := []uint16{
		0x5010,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ProgramCounter() != 0x204 {
		t.Error("5xy0: Values are not compared properly")
	}

	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x66)
	emulator.StoreRegisterValue(0x1, 0x65)
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ProgramCounter() != 0x202 {
		t.Error("5xy0: Values are not compared properly")
	}

}

func Test6xkk(t *testing.T) {
	emulator.ResetEmulator()
	program := []uint16{
		0x6066,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x66 {
		t.Error("6xkk: Value not stored in register properly")
	}
}

func Test7xkk(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x1)
	program := []uint16{
		0x7001,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x2 {
		t.Error("7xkk: Value not added to register properly")
	}
}

func Test8xy0(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x1, 0x66)
	program := []uint16{
		0x8010,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()

	if emulator.ReadRegisterValue(0x0) != emulator.ReadRegisterValue(0x1) {
		t.Error("8xy0 values not copied properly between registers")
	}
}

func Test8xy1(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAA) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0x55) // 0b01010101
	program := []uint16{
		0x8011,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xFf {
		t.Error("8xy1: bitwise or between registers not working properly")
	}

}

func Test8xy2(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAA) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0x55) // 0b01010101
	program := []uint16{
		0x8012,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x00 {
		t.Error("8xy2: bitwise and between registers not working properly")
	}
}

func Test8xy3(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xFF) // 0b11111111
	emulator.StoreRegisterValue(0x1, 0x55) // 0b01010101
	program := []uint16{
		0x8013,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xAA {
		t.Error("8xy3: bitwise xor between registers not working properly")
	}
}

func Test8xy4(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAA) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0x55) // 0b01010101
	program := []uint16{
		0x8014,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xFF {
		t.Error("8xy4: adding registers not working properly")
	}
	if emulator.ReadRegisterValue(0xF) != 0x0 {
		t.Error("8xy4: Carry bit error")
	}
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAB) // 0b10101011
	emulator.StoreRegisterValue(0x1, 0x55) // 0b01010101
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x0 {
		t.Error("8xy: Overflow handled incorrectly")
	}
	if emulator.ReadRegisterValue(0xf) != 0x1 {
		t.Error("8xy4: Carry bit error")
	}
}

func Test8xy5(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAA) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0xAA) // 0b01010101
	program := []uint16{
		0x8015,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x0 {
		t.Error("8xy5: Subtraction not working properly")
	}
	if emulator.ReadRegisterValue(0xf) != 0x0 {
		t.Error("8xy5: Subtraction carry bit not working properly")
	}

	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAA) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0xAB) // 0b01010101

	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xFF {
		t.Error("8xy5: Subtraction underflow handled incorrectly")
	}
	if emulator.ReadRegisterValue(0xF) != 0x1 {
		t.Error("8xy5: Subtraction underflow carry bit not handled correctly")
	}
}

func Test8xy6(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x11)
	program := []uint16{
		0x8016,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x8 {
		t.Error("8xy6: SHR not working properly")
	}
	if emulator.ReadRegisterValue(0xF) != 0x1 {
		t.Error("8xy6: SHR carry not working properly")
	}
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x10)
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()

	if emulator.ReadRegisterValue(0x0) != 0x8 {
		t.Error("8xy6: SHR not working properly")
	}
	if emulator.ReadRegisterValue(0xF) != 0x0 {
		t.Error("8xy6: SHR carry not working properly")
	}

}

func Test8xy7(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAA) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0xAA) // 0b01010101

	program := []uint16{
		0x8017,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0x0 {
		t.Error("8xy7: Subtraction not working properly")
	}
	if emulator.ReadRegisterValue(0xf) != 0x0 {
		t.Error("8xy7: Subtraction carry bit not working properly")
	}

	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xAB) // 0b10101010
	emulator.StoreRegisterValue(0x1, 0xAA) // 0b01010101

	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xFF {
		t.Error("8xy7: Subtraction underflow handled incorrectly")
	}
	if emulator.ReadRegisterValue(0xF) != 0x1 {
		t.Error("8xy7: Subtraction underflow carry bit not handled correctly")
	}
}

func Test8xyE(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xFF)
	program := []uint16{
		0x800E,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xFE {
		t.Error("8xyE: multiplying by two not working properly")
	}
	if emulator.ReadRegisterValue(0xf) != 0x1 {
		t.Error("8xyE: multiplying by two carry not working properly")
	}
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x7F)
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) != 0xFE {
		t.Error("8xyE: multiplying by two not working properly")
	}
	if emulator.ReadRegisterValue(0xF) != 0x0 {
		t.Error("8xyE: multiplying by two carry not working properly")
	}

}

func Test9xy0(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x1)
	emulator.StoreRegisterValue(0x1, 0x0)
	program := []uint16{
		0x9010,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ProgramCounter() != 0x204 {
		t.Error("9xy0: register comparison not handled properly")
	}
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0x1)
	emulator.StoreRegisterValue(0x1, 0x1)
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ProgramCounter() != 0x202 {
		t.Error("9xy0: register comparison not handled properly")
	}
}

func TestAnnn(t *testing.T) {
	emulator.ResetEmulator()
	program := []uint16{
		0xAFFF,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()

	if emulator.Index() != 0xFFF {
		t.Error("Annn: copy to index not working")
	}
}

func TestBnnn(t *testing.T) {
	emulator.ResetEmulator()
	emulator.StoreRegisterValue(0x0, 0xFF)
	program := []uint16{
		0xB0FF,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()

	if emulator.ProgramCounter() != 0x0FF+0x0FF {
		t.Error("Bnnn: Jump and add V0 not working properly")
	}
}

func TestCxkk(t *testing.T) {
	emulator.ResetEmulator()
	program := []uint16{
		0xC0FF,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if emulator.ReadRegisterValue(0x0) == 0x0 {
		t.Error("Cxkk: random number generation not working properly")
	}
}

func TestDxyn(t *testing.T) {
	emulator.ResetEmulator()
	emulator.SetIndex(0x0)
	program := []uint16{
		0xD005,
	}
	emulator.LoadProgramHex(program)
	emulator.EmulateStep()
	if !emulator.ReadScreen(1, 2) {
		t.Error("Dxyn: Sprite Rendering brocken")
	}
}
