package lfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

type LockedFile struct {
	file *os.File
	lock *flock.Flock
}

func Open(p string) (*LockedFile, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	l := flock.New(p)
	if err := l.Lock(); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return &LockedFile{file: f, lock: l}, nil
}

func (lf *LockedFile) Read() ([]byte, error) {
	if _, err := lf.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}
	info, err := lf.file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat: %w", err)
	}
	buf := make([]byte, info.Size())
	if _, err := lf.file.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	return buf, nil
}

func (lf *LockedFile) Write(buf []byte) error {
	if _, err := lf.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}
	if err := lf.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate: %w", err)
	}
	if _, err := lf.file.Write(buf); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	if err := lf.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	return nil
}

func (lf *LockedFile) Close() error {
	var errs []error
	if err := lf.lock.Unlock(); err != nil {
		errs = append(errs, fmt.Errorf("failed to unlock: %w", err))
	}
	if err := lf.file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close: %w", err))
	}
	return errors.Join(errs...)
}
