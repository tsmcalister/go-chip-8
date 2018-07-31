package emulator

type Chip8Instruction uint16

// Instructions from http://mattmik.com/files/chip8/mastering/chip8.html
func (instr Chip8Instruction) getOpCode() string {
	switch instr {
	case 0:
		return "test"
	}
	return ""
}

func (instr Chip8Instruction) interpretInstruction() {

}

//GetNibbles Translates a 16 bit instruction (abcd efgh ijkl mnop) into 4 nibbles each wastefully represented by a byte (0000abcd) (0000efgh) (0000ijkl) (0000mnop)
func GetNibbles(instr Chip8Instruction) (nibble [4]byte) {
	nibble[0] = byte(instr >> 12)
	nibble[1] = byte(instr>>8) & 0xf
	nibble[2] = byte(instr>>4) & 0xf
	nibble[3] = byte(instr) & 0xf
	return
}
