package player

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"subwaycli/internal/ascii"
	"subwaycli/internal/term"
	"subwaycli/internal/video"
)

var errTerminalResized = errors.New("terminal resized")

type Config struct {
	VideoPath string
	Width     int
	Height    int
	FPS       int
	Renderer  *ascii.Renderer
}

type Runner struct {
	cfg Config
}

func New(cfg Config) *Runner {
	return &Runner{cfg: cfg}
}

func (r *Runner) RunForever(ctx context.Context) error {
	if r.cfg.Renderer == nil {
		return errors.New("renderer is required")
	}

	hideCursor()
	clearScreen()
	defer showCursor()
	defer resetColor()

	frameDelay := time.Second / time.Duration(r.cfg.FPS)

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		width, height, err := r.effectiveDimensions()
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(250 * time.Millisecond):
			}
			continue
		}

		dec := video.NewDecoder(r.cfg.VideoPath, width, height, r.cfg.FPS)
		if err := dec.Start(); err != nil {
			if errors.Is(err, video.ErrFFmpegMissing) {
				return fmt.Errorf("%w (install ffmpeg and retry)", err)
			}
			// If ffmpeg fails transiently, keep retrying.
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(500 * time.Millisecond):
			}
			continue
		}

		err = r.playDecoder(ctx, dec, frameDelay, width, height)
		dec.Close()

		if err == nil || errors.Is(err, io.EOF) {
			// Finished video: loop from frame zero forever.
			continue
		}
		if errors.Is(err, errTerminalResized) {
			continue
		}
		if errors.Is(err, context.Canceled) {
			return err
		}

		// Any other read/render issue retries by starting ffmpeg again.
		time.Sleep(250 * time.Millisecond)
	}
}

func (r *Runner) playDecoder(ctx context.Context, dec *video.Decoder, frameDelay time.Duration, width, height int) error {
	next := time.Now()
	lastResizeCheck := time.Now()

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		if r.isResponsive() && time.Since(lastResizeCheck) >= 250*time.Millisecond {
			newWidth, newHeight, err := r.effectiveDimensions()
			if err == nil && (newWidth != width || newHeight != height) {
				return errTerminalResized
			}
			lastResizeCheck = time.Now()
		}

		frame, err := dec.ReadFrame()
		if err != nil {
			return err
		}

		asciiFrame, err := r.cfg.Renderer.Render(frame, width, height)
		if err != nil {
			return fmt.Errorf("render frame: %w", err)
		}

		homeCursor()
		if _, err := os.Stdout.WriteString(asciiFrame); err != nil {
			return fmt.Errorf("write frame: %w", err)
		}

		next = next.Add(frameDelay)
		sleep := time.Until(next)
		if sleep > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(sleep):
			}
		} else {
			next = time.Now()
		}
	}
}

func (r *Runner) isResponsive() bool {
	return r.cfg.Width == 0 || r.cfg.Height == 0
}

func (r *Runner) effectiveDimensions() (int, int, error) {
	width := r.cfg.Width
	height := r.cfg.Height

	if width > 0 && height > 0 {
		return width, height, nil
	}

	cols, rows, err := term.Size()
	if err != nil {
		return 0, 0, fmt.Errorf("detect terminal size: %w", err)
	}

	if width <= 0 {
		width = cols
	}
	if height <= 0 {
		height = rows - 1
	}

	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	return width, height, nil
}

func clearScreen() {
	_, _ = os.Stdout.WriteString("\x1b[2J")
}

func homeCursor() {
	_, _ = os.Stdout.WriteString("\x1b[H")
}

func hideCursor() {
	_, _ = os.Stdout.WriteString("\x1b[?25l")
}

func showCursor() {
	_, _ = os.Stdout.WriteString("\x1b[?25h")
}

func resetColor() {
	_, _ = os.Stdout.WriteString("\x1b[0m")
}
