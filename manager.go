package filter

import (
	"sync"
	"time"
)

const (
	// DefaultCheckInterval 敏感词检查频率（默认5秒检查一次）
	DefaultCheckInterval = time.Second * 5
)

// NewDirtyManager 使用敏感词存储接口创建敏感词管理的实例
func NewDirtyManager(store DirtyStore, intrusion map[rune]bool, checkInterval ...time.Duration) *DirtyManager {
	interval := DefaultCheckInterval
	if len(checkInterval) > 0 {
		interval = checkInterval[0]
	}
	manage := &DirtyManager{
		store:    store,
		version:  store.Version(),
		filter:   NewNodeChanFilter(store.Read(), intrusion),
		interval: interval,
		intrusion:intrusion,
	}
	go func() {
		manage.checkVersion()
	}()
	return manage
}

// DirtyManager 提供敏感词的管理
type DirtyManager struct {
	store     DirtyStore
	filter    DirtyFilter
	filterMux sync.RWMutex
	version   uint64
	interval  time.Duration
	intrusion  map[rune]bool
}

func (dm *DirtyManager) checkVersion() {
	time.AfterFunc(dm.interval, func() {
		storeVersion := dm.store.Version()
		if dm.version < storeVersion {
			dm.filterMux.Lock()
			dm.filter = NewNodeChanFilter(dm.store.Read(), dm.intrusion)
			dm.filterMux.Unlock()
			dm.version = storeVersion
		}
		dm.checkVersion()
	})
}

// Store 获取敏感词存储接口
func (dm *DirtyManager) Store() DirtyStore {
	return dm.store
}

// Filter 获取敏感词过滤接口
func (dm *DirtyManager) Filter() DirtyFilter {
	dm.filterMux.RLock()
	filter := dm.filter
	dm.filterMux.RUnlock()
	return filter
}
