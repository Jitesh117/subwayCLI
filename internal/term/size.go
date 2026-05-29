package term

import (
	"fmt"
	"os"
	"strconv"
)

func Size() (int, int, error) {
	cols, rows, err := platformSize()
	if err == nil && cols > 0 && rows > 0 {
		return cols, rows, nil
	}

	if envCols, envRows, ok := envSize(); ok {
		return envCols, envRows, nil
	}

	if err == nil {
		err = fmt.Errorf("terminal size unavailable")
	}
	return 0, 0, err
}

func envSize() (int, int, bool) {
	c, cerr := strconv.Atoi(os.Getenv("COLUMNS"))
	r, rerr := strconv.Atoi(os.Getenv("LINES"))
	if cerr != nil || rerr != nil || c <= 0 || r <= 0 {
		return 0, 0, false
	}
	return c, r, true
}
