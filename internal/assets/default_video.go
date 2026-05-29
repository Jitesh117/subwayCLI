package assets

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const embeddedVideoName = "subway-surfers.mp4"

var (
	//go:embed subway-surfers.mp4
	embeddedVideo []byte

	embeddedPathOnce sync.Once
	embeddedPath     string
	embeddedPathErr  error
)

// DefaultVideoPath materializes the embedded demo video to a temp file
// and returns its path.
func DefaultVideoPath() (string, error) {
	embeddedPathOnce.Do(func() {
		tmpDir, err := os.MkdirTemp("", "subwaycli-video-*")
		if err != nil {
			embeddedPathErr = fmt.Errorf("create temp dir for embedded video: %w", err)
			return
		}

		path := filepath.Join(tmpDir, embeddedVideoName)
		if err := os.WriteFile(path, embeddedVideo, 0o644); err != nil {
			embeddedPathErr = fmt.Errorf("write embedded video: %w", err)
			return
		}

		embeddedPath = path
	})

	if embeddedPathErr != nil {
		return "", embeddedPathErr
	}
	return embeddedPath, nil
}
