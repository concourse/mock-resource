package resource

import (
	"fmt"
	"io"
	"os"
)

func IsPrivileged() (bool, error) {
	uids, err := os.Open("/proc/self/uid_map")
	if err != nil {
		return false, err
	}

	for {
		var innerStart int
		var outerStart int
		var length int
		n, err := fmt.Fscanf(uids, "%d %d %d", &innerStart, &outerStart, &length)
		if err != nil {
			if err == io.EOF {
				break
			}

			return false, err
		}

		if n < 3 {
			return false, fmt.Errorf("too few fields scanned: %d (want 3)", n)
		}

		if innerStart == 0 && outerStart == 0 {
			return true, nil
		}
	}

	return false, nil
}
