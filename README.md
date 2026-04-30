# chip-8-emulator

CHIP-8 emulator in Go, [Ebiten](https://ebitengine.org/) for rendering.

![Screenshot from IBM ROM](./assets/ibm-rom-screenshot.png)

## Requirements

- Go 1.25+
- X11 libraries (Linux), provided via `nix-shell` on NixOS

## Build

```sh
nix-shell
go build -o chip8 .
```

## Usage

```sh
./chip8 -rom path/to/rom.ch8
```

Window title = ROM filename. Resizable.

## Keypad

Original CHIP-8 hex keypad mapped to keyboard:

```
CHIP-8   Keyboard
1 2 3 C  1 2 3 4
4 5 6 D  Q W E R
7 8 9 E  A S D F
A 0 B F  Z X C V
```

## Quirks

`VM` struct exposes flags for compatibility tuning:

| Flag | Default | Behavior |
|------|---------|----------|
| `ShiftInPlace` | `false` | `false`: `VX = VY` then shift (original CHIP-8). `true`: shift `VX` direct (CHIP-48). |
| `JumpOffsetVX` | `false` | `false`: jump `NNN + V0` (original CHIP-8). `true`: jump `XNN + VX` (CHIP-48). |

## Opcodes

| Opcode | Description |
|--------|-------------|
| `00E0` | Clear display |
| `00EE` | Return from subroutine |
| `1NNN` | Jump |
| `2NNN` | Call subroutine |
| `3XNN` | Skip if `VX == NN` |
| `4XNN` | Skip if `VX != NN` |
| `5XY0` | Skip if `VX == VY` |
| `6XNN` | Set `VX = NN` |
| `7XNN` | Add `NN` to `VX` |
| `8XY0` | Set `VX = VY` |
| `8XY1` | `VX \|= VY` |
| `8XY2` | `VX &= VY` |
| `8XY3` | `VX ^= VY` |
| `8XY4` | `VX += VY`, `VF` = carry |
| `8XY5` | `VX -= VY`, `VF` = no borrow |
| `8XY6` | Shift right (quirk: `ShiftInPlace`) |
| `8XY7` | `VX = VY - VX`, `VF` = no borrow |
| `8XYE` | Shift left (quirk: `ShiftInPlace`) |
| `9XY0` | Skip if `VX != VY` |
| `ANNN` | Set `I = NNN` |
| `BNNN` | Jump `NNN + V0` (quirk: `JumpOffsetVX`) |
| `CXNN` | `VX = rand & NN` |
| `DXYN` | Draw sprite, `VF` = collision |
| `EX9E` | Skip if key `VX` pressed |
| `EXA1` | Skip if key `VX` not pressed |
| `FX07` | `VX = delay timer` |
| `FX0A` | Wait for key, store in `VX` |
| `FX15` | Set delay timer = `VX` |
| `FX18` | Set sound timer = `VX` |
| `FX1E` | `I += VX` |
| `FX29` | `I` = font sprite for `VX` |
| `FX33` | BCD of `VX` → `[I..I+2]` |
| `FX55` | Store `V0..VX` to memory at `I` |
| `FX65` | Load `V0..VX` from memory at `I` |
