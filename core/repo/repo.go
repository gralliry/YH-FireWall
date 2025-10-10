package repo

import (
	"YH-FireWall/core/config"
	"encoding/json"
	"os"
	"path"
	"syscall"
)

var (
	DefaultConfigPath = "/etc/yfw/config.json"
	file              *os.File
)

func Start() error {
	// 确保目录存在
	if err := os.MkdirAll(path.Dir(DefaultConfigPath), 0755); err != nil {
		return err
	}
	// 打开文件
	var err error
	file, err = os.OpenFile(DefaultConfigPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	// 加锁
	if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return err
	}
	return nil
}

func Load(cfg *config.Config) (err error) {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.Size() == 0 {
		return nil
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}

	return nil
}

func Store(cfg *config.Config) error {
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(cfg); err != nil {
		return err
	}

	return file.Sync()
}

func Close() error {
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}
