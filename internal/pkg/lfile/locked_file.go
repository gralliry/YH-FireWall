package lfile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofrs/flock"
)

type LockedFile struct {
	file *os.File
	lock *flock.Flock

	mutex sync.Mutex
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	locked, err := l.TryLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		_ = f.Close()
		return nil, fmt.Errorf("timed out waiting for file lock (another instance may be running)")
	}
	return &LockedFile{file: f, lock: l}, nil
}

func (f *LockedFile) Read() ([]byte, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if _, err := f.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek: %w", err)
	}
	info, err := f.file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to stat: %w", err)
	}
	buf := make([]byte, info.Size())
	if _, err := f.file.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}
	return buf, nil
}

func (f *LockedFile) Write(buf []byte) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	if _, err := f.file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek: %w", err)
	}
	if err := f.file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate: %w", err)
	}
	if _, err := f.file.Write(buf); err != nil {
		return fmt.Errorf("failed to write: %w", err)
	}
	if err := f.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync: %w", err)
	}
	return nil
}

func (f *LockedFile) Close() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	var errs []error
	if err := f.lock.Unlock(); err != nil {
		errs = append(errs, fmt.Errorf("failed to unlock: %w", err))
	}
	if err := f.file.Close(); err != nil {
		errs = append(errs, fmt.Errorf("failed to close: %w", err))
	}
	return errors.Join(errs...)
}
