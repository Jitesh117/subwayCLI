package video

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

var ErrFFmpegMissing = errors.New("ffmpeg not found in PATH")

type Decoder struct {
	videoPath string
	width     int
	height    int
	fps       int

	cmd      *exec.Cmd
	stdout   io.ReadCloser
	frameBuf []byte
}

func NewDecoder(videoPath string, width, height, fps int) *Decoder {
	return &Decoder{
		videoPath: videoPath,
		width:     width,
		height:    height,
		fps:       fps,
		frameBuf:  make([]byte, width*height*3),
	}
}

func (d *Decoder) Start() error {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
		"-nostdin",
		"-i", d.videoPath,
		"-vf", fmt.Sprintf("fps=%d,scale=%d:%d:flags=lanczos", d.fps, d.width, d.height),
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"pipe:1",
	}

	d.cmd = exec.Command("ffmpeg", args...)

	stderr := &strings.Builder{}
	d.cmd.Stderr = stderr

	stdout, err := d.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create stdout pipe: %w", err)
	}
	d.stdout = stdout

	if err := d.cmd.Start(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return ErrFFmpegMissing
		}
		return fmt.Errorf("start ffmpeg: %w", err)
	}

	return nil
}

func (d *Decoder) ReadFrame() ([]byte, error) {
	if d.stdout == nil {
		return nil, errors.New("decoder not started")
	}

	_, err := io.ReadFull(d.stdout, d.frameBuf)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, io.EOF
		}
		return nil, fmt.Errorf("read frame: %w", err)
	}

	return d.frameBuf, nil
}

func (d *Decoder) Close() {
	if d.stdout != nil {
		_ = d.stdout.Close()
		d.stdout = nil
	}
	if d.cmd != nil && d.cmd.Process != nil {
		_ = d.cmd.Process.Kill()
		_, _ = d.cmd.Process.Wait()
	}
}
