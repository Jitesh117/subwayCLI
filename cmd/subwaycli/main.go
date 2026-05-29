package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"subwaycli/internal/ascii"
	"subwaycli/internal/assets"
	"subwaycli/internal/player"
)

func main() {
	videoPath := flag.String("video", "", "path to source video (empty = embedded default)")
	width := flag.Int("width", 0, "output width in characters (0 = auto)")
	height := flag.Int("height", 0, "output height in characters (0 = auto)")
	fps := flag.Int("fps", 24, "playback frames per second")
	chars := flag.String("chars", " .,:;irsXA253hMHGS#9B&@", "dark-to-bright ASCII ramp")
	color := flag.Bool("color", true, "enable truecolor ANSI output")
	flag.Parse()

	if *width < 0 || *height < 0 || *fps <= 0 {
		fmt.Fprintln(os.Stderr, "width/height must be >= 0 and fps must be > 0")
		os.Exit(1)
	}

	renderer, err := ascii.NewRenderer(*chars, *color)
	if err != nil {
		fmt.Fprintf(os.Stderr, "renderer setup failed: %v\n", err)
		os.Exit(1)
	}

	resolvedVideoPath := *videoPath
	if resolvedVideoPath == "" {
		resolvedVideoPath, err = assets.DefaultVideoPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to prepare embedded video: %v\n", err)
			os.Exit(1)
		}
	}

	runner := player.New(player.Config{
		VideoPath: resolvedVideoPath,
		Width:     *width,
		Height:    *height,
		FPS:       *fps,
		Renderer:  renderer,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := runner.RunForever(ctx); err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintf(os.Stderr, "render loop failed: %v\n", err)
		os.Exit(1)
	}
}
