package initializer

/**
 * 延迟执行初始化方法
 */

import (
	"sync"
)

var initialFunc = make([]func(), 0, 32)
var lock = &sync.RWMutex{}

func Registry(f func()) {
	lock.Lock()
	defer lock.Unlock()

	initialFunc = append(initialFunc, f)
}

func InitAll() {
	lock.RLock()
	defer lock.RUnlock()

	for _, f := range initialFunc {
		f()
	}
}
