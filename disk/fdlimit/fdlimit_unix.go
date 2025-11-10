//go:build linux || netbsd || openbsd || solaris

package fdlimit

import (
	"math"
	"syscall"
)

func Get() (int, error) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return 0, err
	}
	if limit.Cur > math.MaxInt {
		return math.MaxInt, nil
	}
	return int(limit.Cur), nil
}
