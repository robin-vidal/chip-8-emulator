# chip-8-emulator

A CHIP-8 emulator written in Go, using [Ebiten](https://ebitengine.org/) for rendering.

![Screenshot from IBM ROM](./assets/ibm-rom-screenshot.png)

## Features

- Fetch / decode / execute pipeline
- Built-in 4×5 font sprites (digits 0–F)
- Sprite rendering with clipping
- ROM loading from any `io.Reader`

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

## Opcodes implemented

| Opcode | Mnemonic | Description          |
|--------|----------|----------------------|
| `00E0` | CLS      | Clear display        |
| `1NNN` | JP NNN   | Jump to address      |
| `6XNN` | LD Vx    | Set register         |
| `7XNN` | ADD Vx   | Add to register      |
| `ANNN` | LD I     | Set index register   |
| `DXYN` | DRW      | Draw sprite          |
