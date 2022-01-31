# GHIP-8
 CHIP-8 Emulator Written in Go/Emulador de CHIP-8 escrito em Go
###### Tetris (from https://multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
![Screenshot_20220131_213532](https://user-images.githubusercontent.com/88759306/151891311-a9e3d20b-4270-499b-800f-46c8ccf6283c.png)

###### CHIP-8 test rom (from https://github.com/corax89/chip8-test-rom)
![Screenshot_20220131_213655](https://user-images.githubusercontent.com/88759306/151891320-114e45dd-f8a1-4571-bbd4-8ff8ad11c163.png)

###### Pumpkin "Dreess" Up (from https://johnearnest.github.io/chip8Archive/play.html?p=pumpkindressup) 
![Screenshot_20220131_213842](https://user-images.githubusercontent.com/88759306/151891464-64978899-324b-4819-9524-0987e4c443ab.png)

### Launching roms
Run the emulator from the commandline followed by the path to the rom. Ex:
```
./ghip8 roms/tetris.ch8
```
Optional arguments:
- -p X: sets the colour palette X (Default: Random)
- -f X: sets the desired framerate to X (Default: 60)
- -c X: sets the amount of CPU cycles ran each frame to X (Default: 20)
- -h : prints the console's memory at startup, starting from address 0x200
- -hC : prints the console's memory at startup, starting from address 0x0
- -d : prints the current opcode every cycle
