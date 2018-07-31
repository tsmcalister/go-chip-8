package test

import (
	"testing"

	"github.com/tsmcalister/go-chip-8/emulator"
)

// write to all possible registers and retrieve value
func TestRegisters(t *testing.T) {
	for i := 0; i < 16; i++ {
		emulator.StoreRegisterValue(byte(i), byte(i))
		if val := emulator.ReadRegisterValue(byte(i)); val != byte(i) {
			t.Errorf("Register V%v should be equal to %v but is %v.", i, i, val)
		}
	}
}
