package main

import (
	"fmt"
)

// the set of fonts to be loaded into memory
var fontSet = [80]byte{
	0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
	0xf0, 0x10, 0xf0, 0x10, 0x10, // 3
	0x90, 0x90, 0xf0, 0x10, 0x10, // 4
	0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
	0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
	0xf0, 0x10, 0x20, 0x40, 0x40, // 7
	0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
	0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
	0xf0, 0x90, 0xf0, 0x90, 0x90, // A
	0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
	0xf0, 0x80, 0x80, 0x80, 0xf0, // C
	0xe0, 0x90, 0x90, 0x90, 0xe0, // D
	0xF0, 0x80, 0xf0, 0x80, 0xf0, // E
	0xf0, 0x80, 0xf0, 0x80, 0x80, // F
}

// This file contains the definition of the 'hardware' of the Chip8
type chip8 struct {
	// An array to hold all 4096 bytes of ram
	ram [4096]byte
	// 2D array of pixels that can be either on or off
	gfx [64][32]bool
	// chip8 has 16 inputs that are either on or off
	input [16]bool
	// The currently executing opcode which is 16Bits
	opcode uint16
	// the array of 16 general purpose registers V0 - VE with a 16th Carry register
	v [16]byte
	// the index register that can have a max value of 0xFFF (12 bits)
	i uint16
	// the program counter register that can have a max value of 0xFFF (12 bits)
	pc uint16
	// two timer registers that count down at 60hz if they are set above 0
	delayTimer, soundTimer byte
	// 16 level stack of 16bit values
	stack [16]uint16
	// points to the last location on the stack
	stackPointer uint16
	// determines if the graphics should be updated
	draw bool
}

func newChip8(rom []byte) *chip8 {
	fmt.Println("Initializing Chip8")
	c := &chip8{
		pc:           0x200,
		opcode:       0,
		i:            0,
		stackPointer: 0,
	}

	fmt.Println("Loading font...")
	// load the font into memory
	for i := 0; i < 80; i++ {
		c.ram[i] = fontSet[i]
	}

	fmt.Println("Loading ROM")
	// load the rom into memory from byte 512 onwards
	for i := 0; i < len(rom); i++ {
		c.ram[i+512] = rom[i]
	}

	fmt.Println("Chip8 Initialized. Printing memory dump:")
	fmt.Println(c.ram)

	return c
}

// This emulates one cycle of the cpu
func (c *chip8) cycle() {
	// As opcodes are two bytes long we fetch two bytes from memory and merge them by shifting the first byte left 8 bits and or-ing it with the next byte
	c.opcode = uint16(c.ram[c.pc])<<8 | uint16(c.ram[c.pc+1])

	// fetch ex
	switch c.opcode & 0xF000 {
	case 0xA000: // ANNN sets I to address NNN
		c.i = c.opcode & 0x0FFF
		c.pc += 2
	case 0x2000:
		c.stack[c.stackPointer] = c.pc
		c.stackPointer++
		c.pc = c.opcode & 0x0FFF
	case 0x0000: // For more opcodes
		switch c.opcode & 0x000F {
		case 0x000: // Clear the screen
			// Execute
		case 0x0004: // 0x8XY4 sets adds vX to vY
			vx := (c.opcode & 0x00f0) >> 4
			vy := (c.opcode & 0x0f00) >> 8
			if c.v[vx] > (0xFF - c.v[vy]) {
				c.v[0xf] = 1
			} else {
				c.v[0xf] = 0
			}
			c.v[vy] += c.v[vx]
			c.pc += 2
		case 0x00E:
			// Execute
		}
	default:
		fmt.Printf("No opcode of: %x\n", c.opcode)
	}

}

func (c *chip8) keys() {
	fmt.Print("Not Implemented")
}