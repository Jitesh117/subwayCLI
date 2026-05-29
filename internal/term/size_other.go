//go:build !darwin && !linux

package term

import "fmt"

func platformSize() (int, int, error) {
	return 0, 0, fmt.Errorf("platform terminal-size detection not implemented")
}
