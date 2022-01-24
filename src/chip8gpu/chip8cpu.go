package chip8cpu

import (
	"math/rand"
)

type CHIP8 struct {
	opcode     uint16
	i          uint16
	pc         uint16
	sp         uint16
	stack      [16]uint16
	delayTimer byte
	soundTimer byte
	memory     [4095]byte
	v          [16]byte
	graphics   [64 * 32]byte
	gamePad    [16]byte
	drawFlag   bool
}

type opcodeFunc func(c *CHIP8)

var chip8_fontset = [80]byte{
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
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

var opcodeMap = map[uint16]opcodeFunc{
	0x0000: func(c *CHIP8) {
		switch c.opcode & 0x000F {
		//00E0: Clears the screen.
		case 0x0000:
			for i := 0; i < 64*32; i++ {
				c.graphics[i] = 0x0
			}
			c.drawFlag = true
		//00EE: Returns from a subroutine.
		case 0x000E:
			c.pc = c.stack[c.sp]
			c.sp--
		}
	},
	//1NNN: Jumps to address NNN.
	0x1000: func(c *CHIP8) {
		c.pc = c.opcode & 0x0FFF
	},
	//2NNN: Calls subroutine at NNN.
	0x2000: func(c *CHIP8) {
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.opcode & 0x0FFF
	},
	//3XNN: Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block);
	0x3000: func(c *CHIP8) {
		if c.v[(c.opcode&0x0F00)>>8] == byte(0x00FF&c.opcode) {
			c.pc += 2
		}
	},
	//4XNN: Skips the next instruction if VX does not equal NN. (Usually the next instruction is a jump to skip a code block);
	0x4000: func(c *CHIP8) {
		if c.v[c.opcode&0x0F00] != byte(0x00FF&c.opcode) {
			c.pc += 2
		}
	},
	//5XY0: Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block);
	0x5000: func(c *CHIP8) {
		if c.v[(c.opcode&0x0F00)>>8] == c.v[(c.opcode&0x00F0)>>4] {
			c.pc += 2
		}
	},
	//6XNN: Sets VX to NN.
	0x6000: func(c *CHIP8) {
		c.v[(c.opcode&0x0F00)>>8] = byte(c.opcode & 0x00FF)
	},
	//7XNN: Adds NN to VX. (Carry flag is not changed);
	0x7000: func(c *CHIP8) {
		c.v[(c.opcode&0x0F00)>>8] += byte(c.opcode & 0x00FF)
	},
	//8XY_ :
	0x8000: func(c *CHIP8) {
		switch c.opcode & 0x000F {
		//0: Sets VX to the value of VY;
		case 0x0000:
			c.v[(c.opcode&0x0F00)>>8] = c.v[(c.opcode&0x00F0)>>4]
		//1: Sets VX to VX or VY. (Bitwise OR operation);
		case 0x0001:
			c.v[(c.opcode&0x0F00)>>8] |= c.v[(c.opcode&0x00F0)>>4]
		//2: Sets VX to VX and VY. (Bitwise AND operation);
		case 0x0002:
			c.v[(c.opcode&0x0F00)>>8] &= c.v[(c.opcode&0x00F0)>>4]
		//3: Sets VX to VX xor VY..
		case 0x0003:
			c.v[(c.opcode&0x0F00)>>8] ^= c.v[(c.opcode&0x00F0)>>4]
		//4: Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there is not.
		case 0x0004:
			c.v[(c.opcode&0x0F00)>>8] += c.v[(c.opcode&0x00F0)>>4]
		//5: VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there is not.
		case 0x0005:
			c.v[(c.opcode&0x0F00)>>8] -= c.v[(c.opcode&0x00F0)>>4]
		//6: Stores the least significant bit of VX in VF and then shifts VX to the right by 1
		case 0x0006:
			c.v[(c.opcode&0x0F00)>>8] = c.v[(c.opcode&0x0F00)>>8] >> 1
		//7: Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there is not.
		case 0x0007:
			c.v[(c.opcode&0x0F00)>>8] = c.v[(c.opcode&0x00F0)>>4] - c.v[(c.opcode&0x0F00)>>8]
		//E: Stores the most significant bit of VX in VF and then shifts VX to the left by 1
		case 0x000E:
			c.v[(c.opcode&0x0F00)>>8] = c.v[(c.opcode&0x0F00)>>8] << 1
		}
	},
	//9XY0: Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block);
	0x9000: func(c *CHIP8) {
		if c.v[(c.opcode&0x0F00)>>8] != c.v[(c.opcode&0x00F0)>>4] {
			c.pc += 2
		}
	},
	//ANNN: Sets I to the address NNN.
	0xA000: func(c *CHIP8) {
		c.i = c.opcode & 0x0FFF
	},
	//BNNN: Jumps to the address NNN plus V0.
	0xB000: func(c *CHIP8) {
		c.pc = uint16(c.v[0]) + ((c.opcode) & 0x0FFF)
	},
	//CXNN: Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN.
	0xC000: func(c *CHIP8) {
		c.v[(c.opcode&0x0F00)>>8] = byte(rand.Int()%255) & byte(c.opcode&0x00FF)
	},
	//DXYN: Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
	//Each row of 8 pixels is read as bit-coded starting from memory location I;
	//I value does not change after the execution of this instruction.
	//As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
	0xD000: func(c *CHIP8) {
		x := uint16(c.v[c.opcode&0x0F00>>8])
		y := uint16(c.v[c.opcode&0x00F0>>4])
		height := uint16(c.opcode & 0x000F)
		pixel := uint16(0)
		c.v[0xF] = 0
		for yline := 0; yline < int(height); yline++ {
			pixel = uint16(c.memory[c.i+uint16(yline)])
			for xline := 0; xline < 8; xline++ {
				if pixel&(0x80>>xline) != 0 {
					if c.graphics[(x+uint16(xline)+(y+uint16(yline)*64))] == 1 {
						c.v[0xF] = 1
					}
					c.graphics[(x + uint16(xline) + (y + uint16(yline)*64))] ^= 1
				}
			}
		}
		c.drawFlag = true
	},
	//EX__:
	0xE000: func(c *CHIP8) {
		switch c.opcode & 0x00FF {
		//9E: Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block);
		case 0x009E:
			if c.gamePad[c.v[(c.opcode&0x0F00)>>8]] != 0 {
				c.pc += 2
			}
		//A1: Skips the next instruction if the key stored in VX is not pressed. (Usually the next instruction is a jump to skip a code block);
		case 0x00A1:
			if c.gamePad[c.v[(c.opcode&0x0F00)>>8]] == 0 {
				c.pc += 2
			}
		}
	},
	//FX__:
	0xF000: func(c *CHIP8) {
		switch c.opcode & 0x00FF {
		case 0x0007:
			c.v[(c.opcode&0x0F00)>>8] = c.delayTimer
		case 0x000A:
			keyPressed := false
			for !keyPressed {
				for i := 0; i < 16; i++ {
					if c.gamePad[i] != 0 {
						c.v[(c.opcode&0x0F00)>>8] = byte(i)
						keyPressed = true
					}
				}
			}
			c.v[(c.opcode&0x0F00)>>8] = c.delayTimer
		case 0x0015:
			c.delayTimer = c.v[(c.opcode&0x0F00)>>8]
		case 0x0018:
			c.soundTimer = c.v[(c.opcode&0x0F00)>>8]
		case 0x001E:
			c.i += uint16(c.v[(c.opcode&0x0F00)>>8])
		case 0x0029:
			c.i = uint16(c.v[(c.opcode&0x0F00)>>8] * 0x5)
		case 0x0033:
			c.memory[c.i] = c.v[(c.opcode&0x0F00)>>8] / 100
			c.memory[c.i+1] = (c.v[(c.opcode&0x0F00)>>8] / 10) % 10
			c.memory[c.i+2] = (c.v[(c.opcode&0x0F00)>>8] % 100) % 10
		case 0x0055:
			for i := 0; uint16(i) <= ((c.opcode & 0x0F00) >> 8); i++ {
				c.memory[c.i+uint16(i)] = c.v[i]
			}
		case 0x0065:
			for i := 0; uint16(i) <= ((c.opcode & 0x0F00) >> 8); i++ {
				c.v[i] = c.memory[c.i+uint16(i)]
			}
		}
	},
}

func (c *CHIP8) Initialize() {
	c.pc = 0x200
	c.opcode = 0
	c.i = 0
	c.sp = 0
	for i := 0; i < 64*32; i++ {
		c.graphics[i] = 0x0
	}
	for i := 0; i < 16; i++ {
		c.v[i] = 0x0
		c.stack[i] = 0x0
	}
	for i := 0; i < 80; i++ {
		c.memory[i] = chip8_fontset[i]
	}
}

func (c *CHIP8) RunCycle() {
	// Loads the current instruction into opcode
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.memory[c.pc+1]])
	opcodeMap[c.opcode&0xF000](c)
	c.pc += 2
}

func (c CHIP8) GetCanvas() [64 * 32]byte {
	return c.graphics
}
