package main

import (
	cpu "chip8cpu"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type colorPalette struct {
	darkR  byte
	darkG  byte
	darkB  byte
	lightR byte
	lightG byte
	lightB byte
}

var colorPalettes [8]colorPalette = [...]colorPalette{
	colorPalette{0x22, 0x23, 0x23, 0xf0, 0xf6, 0xf0}, // https://lospec.com/palette-list/1bit-monitor-glow
	colorPalette{0x38, 0x2b, 0x26, 0xb8, 0xc2, 0xb9}, // https://lospec.com/palette-list/paperback-2
	colorPalette{0x1e, 0x1c, 0x32, 0xc6, 0xba, 0xac}, // https://lospec.com/palette-list/noire-truth
	colorPalette{0x3e, 0x23, 0x2c, 0xed, 0xf6, 0xd6}, // https://lospec.com/palette-list/pixel-ink
	colorPalette{0x2e, 0x30, 0x37, 0xeb, 0xe5, 0xce}, // https://lospec.com/palette-list/obra-dinn-ibm-8503
	colorPalette{0x21, 0x2c, 0x28, 0x72, 0xa4, 0x88}, // https://lospec.com/palette-list/knockia3310
	colorPalette{0x22, 0x2a, 0x3d, 0xed, 0xf2, 0xe2}, // https://lospec.com/palette-list/note-2c
	colorPalette{0x0a, 0x2e, 0x44, 0xfc, 0xff, 0xcc}, // https://lospec.com/palette-list/gato-roboto-starboard

}

func main() {
	var keyMap = map[sdl.Scancode]byte{
		sdl.SCANCODE_1: 0x1,
		sdl.SCANCODE_2: 0x2,
		sdl.SCANCODE_3: 0x3,
		sdl.SCANCODE_4: 0xC,
		sdl.SCANCODE_Q: 0x4,
		sdl.SCANCODE_W: 0x5,
		sdl.SCANCODE_E: 0x6,
		sdl.SCANCODE_R: 0xD,
		sdl.SCANCODE_A: 0x7,
		sdl.SCANCODE_S: 0x8,
		sdl.SCANCODE_D: 0x9,
		sdl.SCANCODE_F: 0xE,
		sdl.SCANCODE_Z: 0xA,
		sdl.SCANCODE_X: 0x0,
		sdl.SCANCODE_C: 0xB,
		sdl.SCANCODE_V: 0xF,
	}
	args := os.Args[1:]
	romName := args[0]
	var chip8 cpu.CHIP8
	chip8.Initialize(romName)
	window, renderer, err := sdl.CreateWindowAndRenderer(640, 320, sdl.WINDOW_BORDERLESS|sdl.WINDOW_RESIZABLE)
	renderer.SetLogicalSize(64, 32)
	if err != nil {
		panic(err)
	}

	var paletteIndex int
	if len(args) == 2 {
		index, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		index = index % len(colorPalettes)
		paletteIndex = index
	} else {
		rand.Seed(time.Now().UnixNano())
		paletteIndex = rand.Int() % len(colorPalettes)
	}
	mainLoop(&chip8, window, renderer, keyMap, &paletteIndex)
}

func mainLoop(chip8 *cpu.CHIP8, window *sdl.Window, renderer *sdl.Renderer, keyMap map[sdl.Scancode]byte, paletteIndex *int) {
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch eventType := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				scanCode := eventType.Keysym.Scancode
				if scanCode == sdl.SCANCODE_P && eventType.State == sdl.PRESSED {
					*paletteIndex = (*paletteIndex + 1) % len(colorPalettes)
				}
				if eventType.State == sdl.PRESSED {
					chip8.GamePad[keyMap[scanCode]] = 1
				}
				if eventType.State == sdl.RELEASED {
					chip8.GamePad[keyMap[scanCode]] = 0
				}
			}
		}
		cycle(chip8, renderer, *paletteIndex)
	}
}

func cycle(chip8 *cpu.CHIP8, renderer *sdl.Renderer, paletteIndex int) {
	frameStartTime := time.Now()
	chip8.RunCycle()
	if chip8.DrawFlag() {
		UpdateGraphics(renderer, chip8.FrameBuffer, paletteIndex)
		chip8.ResetDrawFlag()
	}
	frameEndTime := time.Now()
	elapsed := frameEndTime.Sub(frameStartTime)
	sleepTime, _ := time.ParseDuration("1.66666ms")
	sleepTime -= elapsed
	if sleepTime > 0 {
		time.Sleep(sleepTime - elapsed)
	}
}

func UpdateGraphics(renderer *sdl.Renderer, graphics [64 * 32]byte, paletteIndex int) {
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if graphics[(y*64)+x] == 0 {
				renderer.SetDrawColor(colorPalettes[paletteIndex].darkR, colorPalettes[paletteIndex].darkG, colorPalettes[paletteIndex].darkB, 255)
			} else {
				renderer.SetDrawColor(colorPalettes[paletteIndex].lightR, colorPalettes[paletteIndex].lightG, colorPalettes[paletteIndex].lightB, 255)
			}
			renderer.DrawPoint(int32(x), int32(y))

		}
	}
	renderer.Present()
}
