package emulator

import "fmt"

// screen has 2048 pixels (black and white)
var screen [64 * 32]bool

func updateScreen() {

}

//ClearScreen sets all pixels on screen to off
func ClearScreen() {
	for i := range screen {
		screen[i] = false
	}
}

func WritePixel(x byte, y byte, value bool) bool {
	if x > 32 {
		x -= 32
	}
	if y > 64 {
		y -= 64
	}
	index := x + 64*y
	pixelSet := screen[index]
	screen[index] = screen[index] != value
	if pixelSet && !screen[index] {
		return true
	}
	return false

}

func WriteSpriteByte(x byte, y byte, spriteByte byte) bool {
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
	return flip
}

func PutSprite(x byte, y byte, sprite []byte) bool {
	flip := false
	for i, spriteByte := range sprite {
		if WriteSpriteByte(x, y+byte(i), spriteByte) {
			flip = true
		}
	}
	return flip
}

func PrintScreen() {
	for i, pixel := range screen {
		if i%64 == 0 {
			fmt.Println()
		}
		if pixel {
			fmt.Print("█")
		} else {
			fmt.Print(" ")
		}
	}
}
