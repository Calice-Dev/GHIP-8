module ghip8

go 1.17

replace chip8graphics => ./chip8graphics

replace chip8cpu => ./chip8cpu

require chip8cpu v0.0.0-00010101000000-000000000000

require github.com/veandco/go-sdl2 v0.4.12
