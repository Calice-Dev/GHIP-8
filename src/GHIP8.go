package main

import (
	cpu "chip8cpu"
	"fmt"
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

var colorPalettes = [...]colorPalette{
	{0x22, 0x23, 0x23, 0xf0, 0xf6, 0xf0}, // https://lospec.com/palette-list/1bit-monitor-glow
	{0x38, 0x2b, 0x26, 0xb8, 0xc2, 0xb9}, // https://lospec.com/palette-list/paperback-2
	{0x1e, 0x1c, 0x32, 0xc6, 0xba, 0xac}, // https://lospec.com/palette-list/noire-truth
	{0x3e, 0x23, 0x2c, 0xed, 0xf6, 0xd6}, // https://lospec.com/palette-list/pixel-ink
	{0x2e, 0x30, 0x37, 0xeb, 0xe5, 0xce}, // https://lospec.com/palette-list/obra-dinn-ibm-8503
	{0x21, 0x2c, 0x28, 0x72, 0xa4, 0x88}, // https://lospec.com/palette-list/knockia3310
	{0x22, 0x2a, 0x3d, 0xed, 0xf2, 0xe2}, // https://lospec.com/palette-list/note-2c
	{0x0a, 0x2e, 0x44, 0xfc, 0xff, 0xcc}, // https://lospec.com/palette-list/gato-roboto-starboard
}

func main() {
	var keyMap = map[sdl.Scancode]byte{
		sdl.SCANCODE_1:     0x1,
		sdl.SCANCODE_2:     0x2,
		sdl.SCANCODE_3:     0x3,
		sdl.SCANCODE_4:     0xC,
		sdl.SCANCODE_Q:     0x4,
		sdl.SCANCODE_W:     0x5,
		sdl.SCANCODE_UP:    0x5,
		sdl.SCANCODE_E:     0x6,
		sdl.SCANCODE_R:     0xD,
		sdl.SCANCODE_A:     0x7,
		sdl.SCANCODE_LEFT:  0x7,
		sdl.SCANCODE_S:     0x8,
		sdl.SCANCODE_DOWN:  0x8,
		sdl.SCANCODE_D:     0x9,
		sdl.SCANCODE_RIGHT: 0x9,
		sdl.SCANCODE_F:     0xE,
		sdl.SCANCODE_Z:     0xA,
		sdl.SCANCODE_X:     0x0,
		sdl.SCANCODE_C:     0xB,
		sdl.SCANCODE_V:     0xF,
	}
	fmt.Println("Initializing GHIP-8: CHIP-8 Emulator")
	fmt.Println("...")
	var chip8 cpu.CHIP8
	chip8.Initialize()
	fmt.Println("Initialization succesful")
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("A FilePath to a rom is required")
		return
	}
	var romName string
	var paletteIndex int
	frequency := 1.0 / 60.0
	cyclesPerFrame := 20
	rand.Seed(time.Now().UnixNano())
	paletteIndex = rand.Int() % len(colorPalettes)
	pOpcode := false
	for i, v := range args {
		if v == "-h" {
			fmt.Println("Debug: Dumping memory: ")
			chip8.MemoryHexDump(512)
		} else if v == "-hC" {
			fmt.Println("Debug: Dumping memory: ")
			chip8.MemoryHexDump(0)
		} else if v == "-p" {
			var err error
			paletteIndex, err = strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
		} else if v == "-f" {
			framerate, err := strconv.Atoi(args[i+1])
			frequency = 1.0 / float64(framerate)
			if err != nil {
				panic(err)
			}
		} else if v == "-d" {
			fmt.Println("Debug: Printing executed opcodes")
			pOpcode = true
		} else if v == "-c" {
			var err error
			cyclesPerFrame, err = strconv.Atoi(args[i+1])
			if err != nil {
				panic(err)
			}
		} else if v == "-sX" || v == "-shiftX" {
			chip8.ShiftY = false
		} else if v == "-sY" || v == "-shiftY" {
			chip8.ShiftY = true
		} else {
			romName = v
		}
	}
	fmt.Println("Loading ROM: ", romName)
	fmt.Println("...")
	chip8.ReadRom(romName)
	fmt.Println("Loading succesful")
	fmt.Println("Framerate: ", int(1/frequency), "FPS")
	fmt.Println("Speed: ", cyclesPerFrame, " Cycles per Frame")
	fmt.Println("Palette: ", paletteIndex)
	window, renderer, err := sdl.CreateWindowAndRenderer(640, 320, sdl.WINDOW_BORDERLESS|sdl.WINDOW_RESIZABLE)
	renderer.SetLogicalSize(64, 32)
	if err != nil {
		panic(err)
	}
	mainLoop(&chip8, window, renderer, keyMap, &paletteIndex, frequency, pOpcode, cyclesPerFrame)
}

func mainLoop(chip8 *cpu.CHIP8, window *sdl.Window, renderer *sdl.Renderer, keyMap map[sdl.Scancode]byte, paletteIndex *int, frequency float64, p bool, c int) {
	paused := false
	runOneStep := false
	for {
		frameStartTime := time.Now()
		for j := 0; j < c; j++ {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch eventType := event.(type) {
				case *sdl.QuitEvent:
					fmt.Println("Closing Emulator")
					return
				case *sdl.KeyboardEvent:
					scanCode := eventType.Keysym.Scancode
					if scanCode == sdl.SCANCODE_P && eventType.State == sdl.PRESSED {
						*paletteIndex = (*paletteIndex + 1) % len(colorPalettes)
					}
					if scanCode == sdl.SCANCODE_ESCAPE && eventType.State == sdl.PRESSED {
						paused = !paused
					}
					if scanCode == sdl.SCANCODE_N && eventType.State == sdl.PRESSED {
						runOneStep = true
					}
					if eventType.State == sdl.PRESSED {
						_, ok := keyMap[scanCode]
						if ok {
							chip8.GamePad[keyMap[scanCode]] = 1
						}
					}
					if eventType.State == sdl.RELEASED {
						_, ok := keyMap[scanCode]
						if ok {
							chip8.GamePad[keyMap[scanCode]] = 0
						}
					}
				}
			}
			if !paused || runOneStep {
				cycle(chip8, renderer, *paletteIndex, p)
			}
		}
		if !paused || runOneStep {
			chip8.ReduceTimers()
			frameEndTime := time.Now()
			elapsed := frameEndTime.Sub(frameStartTime)
			s := fmt.Sprintf("%fs", frequency)
			//fmt.Println(s)
			//frameTime := time.Duration(elapsed.Milliseconds())
			desiredFrameTime, err := time.ParseDuration(s)
			//fmt.Println(desiredFrameTime)
			if err != nil {
				panic(err)
			}
			sleepTime := (desiredFrameTime - elapsed)
			//fmt.Println(sleepTime)
			if sleepTime > 0 {
				time.Sleep(sleepTime)
			}
			if p {
				fmt.Println()
			}
		}
		runOneStep = false
	}
}

func cycle(chip8 *cpu.CHIP8, renderer *sdl.Renderer, paletteIndex int, p bool) {
	chip8.RunCycle(p)
	if chip8.DrawFlag() {
		UpdateGraphics(renderer, chip8.FrameBuffer, paletteIndex)
		chip8.ResetDrawFlag()
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
