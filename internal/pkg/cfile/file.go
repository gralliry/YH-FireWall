package cfile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type CacheFile struct {
	mutex sync.RWMutex
	waitg sync.WaitGroup

	path  string
	cache []byte
	dirty chan struct{}

	close  sync.Once
	closed bool
}

func Open(path string) (lf *CacheFile, err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	lf = &CacheFile{
		path:   path,
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
	for range f.dirty {
		// 不要使用Read()去替换这里
		f.mutex.RLock()
		buf := bytes.Clone(f.cache)
		f.mutex.RUnlock()
		// 文件写入不用加锁，f.dirty保证了唯一操作
		tmpfile := f.path + ".tmp"
		if err := os.WriteFile(tmpfile, buf, 0644); err != nil {
			continue
		}
		// Linux 系统 中 rename会覆盖
		if err := os.Rename(tmpfile, f.path); err != nil {
			continue
		}
	}
}

func (f *CacheFile) Close() {
	f.close.Do(func() {
		f.mutex.Lock()
		f.closed = false
		close(f.dirty)
		f.mutex.Unlock()
		f.waitg.Wait()
	})
}
