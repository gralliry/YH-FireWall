package rtable

import (
	"YH-FireWall/core/rule"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/gopacket/layers"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"path"
	"sort"
	"sync"
	"syscall"
)

var (
	ruleList        []*rule.Rule
	ruleMap         map[string]*rule.Rule
	ruleIsListDirty bool
	mutex           sync.RWMutex
	file            *os.File

	defaultAccept = true
)

type Config struct {
	Path          string `json:"path"`
	DefaultAccept bool   `json:"default_accept"`
}

func Load(config Config) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	//
	defaultAccept = config.DefaultAccept
	// 确保目录存在
	if err = os.MkdirAll(path.Dir(config.Path), 0755); err != nil {
		return fmt.Errorf("failed to create directory for rule table: %w", err)
	}
	// 打开文件
	file, err = os.OpenFile(config.Path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", config.Path, err)
	}
	// 尝试独占锁（非阻塞）
	if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = file.Close()
		return fmt.Errorf("failed to lock file %s: %w", config.Path, err)
	}
	// 设置默认规则
	rules := make([]rule.Config, 0)
	// 用 Decoder 直接解析
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&rules); err != nil && err != io.EOF {
		return fmt.Errorf("failed to decode rule file: %w", err)
	}
	// --- 处理规则 ---
	ruleList = make([]*rule.Rule, 0)
	ruleMap = make(map[string]*rule.Rule)
	// 加载配置
	var rr *rule.Rule
	for _, rc := range rules {
		// 匹配
		if _, exists := ruleMap[rc.Id]; exists {
			log.Errorf("rule %s exists", rc.Id)
			continue
		}
		if rr, err = rule.Parse(&rc); err != nil {
			log.Errorf("failed to parse rule %s: %v", rc.Id, err)
			continue
		}
		// 如果都没有，就添加
		ruleMap[rc.Id] = rr
		ruleList = append(ruleList, rr)
	}
	// 标记为非脏数据 // 排序为match函数负责
	ruleIsListDirty = true
	return nil
}

func Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	// 重置文件指针
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}
	// 清空文件
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}
	// 存储
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // "" 前缀，"  " 缩进
	if err := encoder.Encode(getAll()); err != nil {
		return fmt.Errorf("failed to encode ruleList: %w", err)
	}
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}
	var errs []error
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_UN); err != nil {
		errs = append(errs, err)
	}
	if err := file.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// Append 添加或更新规则
func Append(ro *rule.Option) (string, error) {
	// uuid不可能重复
	rc := ro.Default()
	//
	mutex.Lock()
	defer mutex.Unlock()
	//
	rr, err := rule.Parse(rc)
	if err != nil {
		return "", err
	}
	// 标记
	ruleIsListDirty = true
	// 如果都没有，就添加
	ruleMap[rc.Id] = rr
	ruleList = append(ruleList, rr)
	//
	return rr.Id(), nil
}

// Update 更新规则
func Update(id string, ro *rule.Option) error {
	mutex.RLock()
	defer mutex.RUnlock()
	rr, exists := ruleMap[id]
	if !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	if err := rr.Update(*ro); err != nil {
		return fmt.Errorf("update rule %s error: %v", id, err)
	}
	return nil
}

// Delete 删除规则
func Delete(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	// 获取
	if _, exists := ruleMap[id]; !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	// 重新构造
	index := -1
	for i, r := range ruleList {
		if r.Id() == id {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("rule %s not exists", id)
	}
	// 标记
	ruleIsListDirty = true
	// 删除
	delete(ruleMap, id)
	// 移动
	ruleList[index] = ruleList[len(ruleList)-1]
	ruleList = ruleList[:len(ruleList)-1]
	return nil
}

// Match 匹配：按优先级从高到低
func Match(srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16, inDev *uint32, outDev *uint32, protocol layers.IPProtocol) bool {
	mutex.Lock()
	defer mutex.Unlock()
	// 脏数据
	if ruleIsListDirty {
		sort.SliceStable(ruleList, func(i, j int) bool {
			return ruleList[i].Priority() < ruleList[j].Priority()
		})
	}
	// 匹配
	for _, r := range ruleList {
		if r.Match(srcIP, srcPort, dstIP, dstPort, inDev, outDev, protocol) {
			return r.Accept()
		}
	}
	return defaultAccept
}

func Get(rid string) *rule.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	if rr, exists := ruleMap[rid]; !exists {
		return nil
	} else {
		return rr.Unparse()
	}
}

func getAll() []rule.Config {
	rules := make([]rule.Config, len(ruleList))
	for i, r := range ruleList {
		rules[i] = *r.Unparse()
	}
	return rules
}

func GetAll() []rule.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return getAll()
}

func Enable(id string, enable bool) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	rr, exists := ruleMap[id]
	if !exists {
		return false
	}
	rr.SetEnable(enable)
	return true
}
