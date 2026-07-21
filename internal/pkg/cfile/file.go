package cfile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofrs/flock"
)

type CacheFile struct {
	mutex sync.RWMutex
	waitg sync.WaitGroup

	path  string
	lock  *flock.Flock
	cache []byte
	dirty chan struct{}

	close  sync.Once
	closed bool
}

func Open(path string) (lf *CacheFile, err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	// 获取文件锁
	lock := flock.New(path + ".lock")
	locked, err := lock.TryLock()
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		return nil, fmt.Errorf("file is locked by another process")
	}
	defer func() {
		if err != nil {
			lock.Unlock()
		}
	}()
	// 打开文件
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	// 读取文件内容
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	lf = &CacheFile{
		path:   path,
		lock:   lock,
		cache:  buf,
		dirty:  make(chan struct{}, 1),
		closed: false,
	}
	lf.waitg.Go(lf.handle)
	return lf, nil
}

func (f *CacheFile) Read() []byte {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	// 复制流
	return bytes.Clone(f.cache)
}

func (f *CacheFile) Write(buf []byte) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	//
	if f.closed {
		return fmt.Errorf("file is closed")
	}
	// 复制缓存
	f.cache = bytes.Clone(buf)
	// 标记为脏
	select {
	case f.dirty <- struct{}{}:
	default:
	}
	return nil
}

func (f *CacheFile) handle() {
	tmpfile := f.path + ".tmp"
	for range f.dirty {
		f.mutex.RLock()
		buf := bytes.Clone(f.cache)
		f.mutex.RUnlock()

		if err := os.WriteFile(tmpfile, buf, 0644); err != nil {
			continue
		}
		if err := os.Rename(tmpfile, f.path); err != nil {
			continue
		}
	}
	os.Remove(tmpfile)
}

func (f *CacheFile) Close() {
	f.close.Do(func() {
		f.mutex.Lock()
		f.closed = true
		close(f.dirty)
		f.mutex.Unlock()
		f.waitg.Wait()

		f.lock.Unlock()
	})
}
