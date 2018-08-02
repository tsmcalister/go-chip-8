package emulator

import "io/ioutil"

/*
	RAM Implementation
*/

// The original Chip-8 interpreter used up the first 512 bytes of memory. Therefore user memory starts at 0x200 in most programs.
// 0x050-0x0A0 is used for fonts
var memory [4096]byte

//ReadMemory reads the memory at address [addr] and [addr+1] and combines them into a Chip8Instruction
func ReadMemory(addr uint32) Chip8Instruction {
	byte1 := memory[addr]
	byte2 := memory[addr+1]
	return Chip8Instruction(int32(byte1)<<8 + int32(byte2))
}

func ReadMemoryByte(addr uint32) byte {
	return memory[addr]
}

func WriteInstructionMemory(addr uint32, value Chip8Instruction) {
	byte1, byte2 := InstructionToBytes(value)
	memory[addr] = byte1
	memory[addr+1] = byte2
}

func WriteByteMemory(addr uint32, data byte) {
	memory[addr] = data
}

func BytesToInstruction(byte1, byte2 byte) Chip8Instruction {
	return Chip8Instruction(uint32(byte1)<<8 + uint32(byte2))
}

func InstructionToBytes(instr Chip8Instruction) (byte1, byte2 byte) {
	value := int32(instr)
	byte1 = byte(value >> 8)
	byte2 = byte(value)
	return
}

func InitMemory() {
	clearMemory()
	loadFontData()
}

func clearMemory() {
	for i := range memory {
		memory[i] = 0x0
	}
}

func ReadSprite(addr uint32, n byte) []byte {
	sprite := make([]byte, n, n)
	for i := 0; uint32(i) < uint32(n); i++ {
		sprite[i] = ReadMemoryByte(addr + uint32(i))
	}
	return sprite
}

func loadFontData() {
	fontData := [80]byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80} // F

	for i := range fontData {
		WriteByteMemory(uint32(i), fontData[i])
	}

}

func GetFontSpriteAddress(character byte) uint32 {
	return uint32(character * 5)
}

/*
	Registers Implementation
*/

var registers = make(map[byte]byte)

// Only VX values from 0-F used (1 nibble)
func StoreRegisterValue(VX byte, value byte) {
	registers[VX] = value
}

func ReadRegisterValue(VX byte) byte {
	return registers[VX]
}

func RegistersToMemory(addr uint32) {
	for i, val := range registers {
		WriteByteMemory(addr+uint32(i), val)
	}
}

func MemoryToRegisters(addr uint32) {
	for i := range registers {
		StoreRegisterValue(i, ReadMemoryByte(addr+uint32(i)))
	}
}

func DrawFlagIsSet() bool {
	return false
}

func LoadProgramHex(program []uint16) {
	startAddress := uint32(0x200)
	for i, programInstr := range program {
		WriteInstructionMemory(startAddress+uint32(i)*2, Chip8Instruction(programInstr))
	}
}

func LoadProgram(file string) {
	program, _ := ioutil.ReadFile(file)
	startAddress := uint32(0x200)
	for i, programByte := range program {
		WriteByteMemory(startAddress+uint32(i), programByte)
	}
}

/*
	Stack Implementation
*/
