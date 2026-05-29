# subwaycli

Endless terminal ASCII renderer with an embedded default video of a subway surfers gameplay. Resizable according to your terminal windoo.

![Demo](./subway_cli_demo.gif)

## Requirements

- Go 1.22+
- `ffmpeg` in `PATH`

## Quick Start

```bash
go run ./cmd/subwaycli -fps 24
```

## Build

```bash
make build
```

Run:

```bash
./bin/subwaycli -fps 24
```

## Global Install (macOS)

```bash
sudo install -m 755 ./bin/subwaycli /opt/homebrew/bin/subwaycli
```

Then use from anywhere:

```bash
subwaycli -fps 24
```

## Use Your Own Video

Pass any local video file with `-video`:

```bash
subwaycli -video /absolute/path/to/video.mp4 -fps 24
```

Example:

```bash
subwaycli -video ~/Downloads/my-clip.mov -fps 30 -width 140 -height 45
```
