package core

import (
	"YH-FireWall/internal/core/config"
	"YH-FireWall/internal/core/group"
	"YH-FireWall/internal/core/rule"
	"log"
	"os"
)

const (
	ConfigPath = ".config.json"
)

var (
	cfg *config.Config
	//  组 规则
	groupName2group = make(map[string]*group.Group)
	ruleName2rule   = make(map[string]*rule.Rule)
	ruleName2group  = make(map[*rule.Rule]*group.Group)
)

func Start() (err error) {
	// 检测当前用户是否为 root 用户
	if os.Geteuid() != 0 {
		log.Fatal("当前用户非 root 用户")
	}
	// 读取配置文件
	cfg, err = config.Load(ConfigPath)
	if err != nil {
		return err
	}
	// 加载组配置
	for _, gc := range cfg.Groups {
		if g, err := group.Load(&gc); err == nil {
			groupName2group[g.Name()] = g
		}
	}
	//

	return nil
}

func Close() error {
	if cfg != nil {
		if err := cfg.Store(ConfigPath); err != nil {
			log.Printf("保存配置文件失败: %v", err)
		}
	}
	return nil
}
