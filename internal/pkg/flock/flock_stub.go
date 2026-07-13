//go:build !linux

package flock

func Lock(fd uintptr) error {
	return nil
}

func Unlock(fd uintptr) error {
	return nil
}
