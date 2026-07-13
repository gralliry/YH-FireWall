//go:build linux

package flock

import "golang.org/x/sys/unix"

func Lock(fd uintptr) error {
	return unix.Flock(int(fd), unix.LOCK_EX|unix.LOCK_NB)
}

func Unlock(fd uintptr) error {
	return unix.Flock(int(fd), unix.LOCK_UN)
}
