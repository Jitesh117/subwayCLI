package ascii

import (
	"errors"
	"strings"
)

var errInvalidFrameSize = errors.New("invalid frame byte length")

type Renderer struct {
	ramp  []byte
	color bool
}

func NewRenderer(chars string, color bool) (*Renderer, error) {
	if len(chars) < 2 {
		return nil, errors.New("chars ramp must contain at least 2 characters")
	}

	return &Renderer{ramp: []byte(chars), color: color}, nil
}

func (r *Renderer) Render(frame []byte, width, height int) (string, error) {
	expected := width * height * 3
	if len(frame) != expected {
		return "", errInvalidFrameSize
	}

	var b strings.Builder
	if r.color {
		b.Grow(width * height * 20)
	} else {
		b.Grow(width*height + height)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := (y*width + x) * 3
			rv := frame[i]
			gv := frame[i+1]
			bv := frame[i+2]

			lum := (2126*int(rv) + 7152*int(gv) + 722*int(bv)) / 10000
			idx := lum * (len(r.ramp) - 1) / 255
			ch := r.ramp[idx]

			if r.color {
				b.WriteString("\x1b[38;2;")
				b.WriteString(itoa(int(rv)))
				b.WriteByte(';')
				b.WriteString(itoa(int(gv)))
				b.WriteByte(';')
				b.WriteString(itoa(int(bv)))
				b.WriteByte('m')
				b.WriteByte(ch)
			} else {
				b.WriteByte(ch)
			}
		}
		if r.color {
			b.WriteString("\x1b[0m")
		}
		b.WriteByte('\n')
	}

	return b.String(), nil
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	var buf [3]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + (v % 10))
		v /= 10
	}
	return string(buf[i:])
}
