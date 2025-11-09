//go:build linux || netbsd || openbsd || solaris
// +build linux netbsd openbsd solaris

package fdlimit

import "syscall"

func Get() (int, error) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return 0, err
	}
	return int(limit.Cur), nil
}
