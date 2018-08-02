package emulator

import (
	"bytes"
	"fmt"
	"sync"
)

// screen has 2048 pixels (black and white)
var screen [64 * 32]bool

func ReadScreen(x byte, y byte) bool {
	return screen[x+y*64]
}

//ClearScreen sets all pixels on screen to off
func ClearScreen() {
	for i := range screen {
		screen[i] = false
	}
}

func WritePixel(x byte, y byte, value bool) bool {
	if x > 64 {
		x -= 64
	}
	if y > 32 {
		y -= 32
	}
	index := uint32(x) + 64*uint32(y)
	pixelSet := screen[index]
	screen[index] = screen[index] != value
	if pixelSet && !screen[index] {
		return true
	}
	return false

}

func WriteSpriteByte(x byte, y byte, spriteByte byte, flipped chan bool) {
	flip := false
	// binary 0b10000000
	slidingMask := byte(0x80)
	i := 0
	for slidingMask > 0 {
		pixelValue := slidingMask&spriteByte > 0
		if WritePixel(x+byte(i), y, pixelValue) {
			flip = true
		}
		slidingMask = slidingMask >> 1
		i++
	}
	flipped <- flip
}

func PutSprite(x byte, y byte, sprite []byte) bool {
	flip := false
	flipped := make(chan bool, 100)
	var wg sync.WaitGroup
	for i, spriteByte := range sprite {
		wg.Add(1)
		go WriteSpriteByte(x, y+byte(i), spriteByte, flipped)
	}
	for range sprite {
		msg := <-flipped
		if msg {
			flip = true
		}
	}
	return flip
}

func PrintScreen() {
	var buffer bytes.Buffer
	for i, pixel := range screen {
		if i%64 == 0 {
			buffer.WriteString("\n")
		}
		if pixel {
			buffer.WriteString("█")
		} else {
			buffer.WriteString(" ")
		}
	}
	fmt.Printf(buffer.String())
}
