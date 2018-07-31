package emulator

import (
	"fmt"
	"math/rand"
	"time"
)

var stack [16]uint32
var stackPointer byte

// Are actually 24 bits
var index uint32
var programCounter uint32 = 0x200 // execution starts at this address

var delayTimer byte = 0x0
var soundTimer byte = 0x0

func EmulateStep() {
	print("\033[H\033[2J")
	instr := ReadMemory(programCounter)
	fmt.Printf("%x\n", instr)
	ExecuteInstruction(instr)
	UpdateTimers()
	PrintScreen()
	time.Sleep(100 * time.Millisecond)
}

func UpdateTimers() {
	if delayTimer > 0 {
		delayTimer--
	}
	if soundTimer > 0 {
		soundTimer--
	}
}

func ExecuteInstruction(instr Chip8Instruction) {
	nibbles := GetNibbles(instr)
	switch nibbles[0] {
	case 0x0:
		switch nibbles[2] {
		case 0xE:
			switch nibbles[3] {
			// 0x00E0
			// Clear The Screen
			case 0x0:
				ClearScreen()
				programCounter += 2
			// 0x00EE
			// Return from Subroutine
			case 0xE:
				programCounter = stack[stackPointer]
				stackPointer--
			}
		default:
			// 0x0NNN
			// Execute machine language subroutine at address NNN
			// Ignored by modern interpreters according to http://devernay.free.fr/hacks/chip8/C8TECH10.HTM#0nnn
		}
	case 0x1:
		// 0x1NNN
		// Jump to address NNN
		addr := uint32(nibbles[1])<<8 + uint32(nibbles[2])<<4 + uint32(nibbles[3])
		programCounter = addr
	case 0x2:
		// 0x2NNN
		// Execute subroutine starting at address NNN
		stackPointer++
		stack[stackPointer] = programCounter
		addr := uint32(nibbles[1])<<8 + uint32(nibbles[2])<<4 + uint32(nibbles[3])
		programCounter = addr
	case 0x3:
		// 0x3XNN
		// Skip the following instruction if the value of register VX equals NN
		VX := ReadRegisterValue(nibbles[1])
		if VX == nibbles[2]<<4+nibbles[3] {
			programCounter += 4
		} else {
			programCounter += 2
		}
	case 0x4:
		// 0x4XNN
		// Skip the following instruction if the value of register VX is not equal to NN
		VX := ReadRegisterValue(nibbles[1])
		if VX != nibbles[2]<<4+nibbles[3] {
			programCounter += 4
		} else {
			programCounter += 2
		}
	case 0x5:
		// 0x5XY0
		// Skip the following instruction if the value of register VX is equal to the value of register VY
		VX := ReadRegisterValue(nibbles[1])
		VY := ReadRegisterValue(nibbles[2])
		if VX == VY {
			programCounter += 4
		} else {
			programCounter += 2
		}
	case 0x6:
		// 0x6XNN
		// Store number NN in register VX
		StoreRegisterValue(nibbles[1], nibbles[2]<<4+nibbles[3])
		programCounter += 2
	case 0x7:
		// 0x7XNN
		// Add the value NN to register VX
		VX := ReadRegisterValue(nibbles[1])
		VX += nibbles[2]<<4 + nibbles[3]
		StoreRegisterValue(nibbles[1], VX)
		programCounter += 2
	case 0x8:
		switch nibbles[3] {
		case 0x0:
			// 0x8XY0
			// Store the value of register VY in register VX
			VY := ReadRegisterValue(nibbles[2])
			StoreRegisterValue(nibbles[1], VY)
			programCounter += 2
		case 0x1:
			// 0x8XY1
			// Set VX to VX OR VY
			VX := ReadRegisterValue(nibbles[1])
			VY := ReadRegisterValue(nibbles[2])
			VX = VX | VY
			StoreRegisterValue(nibbles[1], VX)
			programCounter += 2
		case 0x2:
			// 8XY2
			// Set VX to VX AND VY
			VX := ReadRegisterValue(nibbles[1])
			VY := ReadRegisterValue(nibbles[2])
			VX = VX & VY
			StoreRegisterValue(nibbles[1], VX)
			programCounter += 2
		case 0x3:
			// 8XY3
			// Set VX to VX XOR VY
			VX := ReadRegisterValue(nibbles[1])
			VY := ReadRegisterValue(nibbles[2])
			VX = VX ^ VY
			StoreRegisterValue(nibbles[1], VX)
			programCounter += 2
		case 0x4:
			// 8XY4
			// Add the value of register VY to register VX
			// Set VF to 01 if a carry occurs
			// Set VF to 00 if a carry does not occur
			VX := ReadRegisterValue(nibbles[1])
			VY := ReadRegisterValue(nibbles[2])
			StoreRegisterValue(nibbles[1], VX+VY)
			if uint32(VX)+uint32(VY) > 0xff {
				StoreRegisterValue(0xf, 0x1)
			} else {
				StoreRegisterValue(0xf, 0x0)
			}
			programCounter += 2
		case 0x5:
			// 8XY5
			// Subtract the value of register VY from register VX
			// Set VF to 00 if a borrow occurs
			// Set VF to 01 if a borrow does not occur
			VX := ReadRegisterValue(nibbles[1])
			VY := ReadRegisterValue(nibbles[2])
			StoreRegisterValue(nibbles[1], VX-VY)
			if VY > VX {
				StoreRegisterValue(0xf, 0x1)
			} else {
				StoreRegisterValue(0xf, 0x0)
			}
			programCounter += 2
		case 0x6:
			// 8XY6
			// Store the value of register VY shifted right one bit in register VX
			// Set register VF to the least significant bit prior to the shift
			VX := ReadRegisterValue(nibbles[1])
			if VX&0x1 == 1 {
				StoreRegisterValue(0xF, 0x1)
			} else {
				StoreRegisterValue(0xF, 0x0)
			}
			StoreRegisterValue(nibbles[1], VX>>1)
			programCounter += 2

		case 0x7:
			// 8XY7
			// Set register VX to the value of VY minus VX
			// Set VF to 00 if a borrow occurs
			// Set VF to 01 if a borrow does not occur
			VX := ReadRegisterValue(nibbles[1])
			VY := ReadRegisterValue(nibbles[2])
			StoreRegisterValue(nibbles[1], VY-VX)
			if VX > VY {
				StoreRegisterValue(0xf, 0x0)
			} else {
				StoreRegisterValue(0xf, 0x1)
			}
			programCounter += 2
		case 0xE:
			// 8XYE
			// Store the value of register VY shifted left one bit in register VX
			// Set register VF to the most significant bit prior to the shift
			VX := ReadRegisterValue(nibbles[1])
			// binary 0b10000000
			if VX&0x80 == 0x80 {
				StoreRegisterValue(0xf, 0x1)
			} else {
				StoreRegisterValue(0xf, 0x0)
			}
			StoreRegisterValue(nibbles[1], VX<<1)
			programCounter += 2
		}
	case 0x9:
		// 9XY0
		// Skip the following instruction if the value of register VX is not equal to the value of register VY
		VX := ReadRegisterValue(nibbles[1])
		VY := ReadRegisterValue(nibbles[2])
		if VX != VY {
			programCounter += 4
		} else {
			programCounter += 2
		}
	case 0xA:
		// ANNN
		// Store memory address NNN in register I
		index = uint32(nibbles[1])<<8 + uint32(nibbles[2])<<4 + uint32(nibbles[3])
		programCounter += 2
	case 0xB:
		// BNNN
		// Jump to address NNN + V0
		NNN := uint32(nibbles[1])<<8 + uint32(nibbles[2])<<4 + uint32(nibbles[3])
		V0 := ReadRegisterValue(0x0)
		programCounter = NNN + uint32(V0)
	case 0xC:
		// CXNN
		// Set VX to a random number with a mask of NN
		rNum := rand.Uint32()&0xff&uint32(nibbles[2])<<4 + uint32(nibbles[3])
		StoreRegisterValue(nibbles[1], byte(rNum))
		programCounter += 2
	case 0xD:
		// DXYN
		// Draw a sprite at position VX, VY with N bytes of sprite data starting at the address stored in I
		// Set VF to 01 if any set pixels are changed to unset, and 00 otherwise
		VX := ReadRegisterValue(nibbles[1])
		VY := ReadRegisterValue(nibbles[2])
		sprite := ReadSprite(index, nibbles[3])
		flipped := PutSprite(VX, VY, sprite)
		if flipped {
			StoreRegisterValue(0xf, 0x1)
		} else {
			StoreRegisterValue(0xf, 0x0)
		}
		programCounter += 2

	case 0xE:
		// EX9E
		// Skip the following instruction if the key corresponding to the hex value currently stored in register VX is pressed
		/*
			TODO IMPLEMENT
		*/
		programCounter += 2
	case 0xF:
		lastByte := nibbles[2]<<4 + nibbles[3]
		switch lastByte {
		case 0x07:
			// FX07
			// Store the current value of the delay timer in register VX
			StoreRegisterValue(nibbles[1], delayTimer)
			programCounter += 2
		case 0x0A:
			// FX0A
			// Wait for a keypress and store the result in register VX
			// TODO IMPLEMENT
			programCounter += 2
		case 0x15:
			// FX15
			// Set the delay timer to the value of register VX
			VX := ReadRegisterValue(nibbles[1])
			delayTimer = VX
			programCounter += 2
		case 0x18:
			// FX18
			// Set the sound timer to the value of register VX
			VX := ReadRegisterValue(nibbles[1])
			soundTimer = VX
			programCounter += 2
		case 0x1E:
			// FX1E
			// Add the value stored in register VX to register I
			VX := ReadRegisterValue(nibbles[1])
			index = index + uint32(VX)
		case 0x29:
			// FX29
			// Set I to the memory address of the sprite data corresponding to the hexadecimal digit stored in register VX
			VX := ReadRegisterValue(nibbles[1])
			index = GetFontSpriteAddress(VX)
		case 0x33:
			// FX33
			// Store the binary-coded decimal equivalent of the value stored in register VX at addresses I, I+1, and I+2
			VX := ReadRegisterValue(nibbles[1])
			dec_1 := VX % 10
			dec_10 := (VX%100 - dec_1) / 10
			dec_100 := (VX - dec_10 - dec_1) / 100
			WriteByteMemory(index, dec_100)
			WriteByteMemory(index+1, dec_10)
			WriteByteMemory(index+2, dec_1)
		case 0x55:
			// FX55
			// Store the values of registers V0 to VX inclusive in memory starting at address I
			// I is set to I + X + 1 after operation
			RegistersToMemory(index)
			index = index + uint32(nibbles[1]) + 1

		case 0x65:
			// FX65
			// Fill registers V0 to VX inclusive with the values stored in memory starting at address I
			// I is set to I + X + 1 after operation
			MemoryToRegisters(index)
			index = index + uint32(nibbles[1]) + 1
		}
	}
}
