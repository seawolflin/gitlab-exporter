package initializer

/**
 * 延迟执行初始化方法
 */

import (
	"log"
	"sort"
	"sync"
)

const (
	MAX_ORDER     = 0
	MIN_ORDER     = 65535
	DEFUALT_ORDER = 1000
)

type funcList = []func()

var initialFuncMap = make(map[int]funcList, 32)
var lock = &sync.RWMutex{}

// Registry 注册初始化方法，默认优先级为1000
func Registry(f func()) {
	RegistryByOrder(f, DEFUALT_ORDER)
}

// RegistryByOrder 注册初始化方法，指定优先级
func RegistryByOrder(f func(), order int) {
	if order > MIN_ORDER || order < MAX_ORDER {
		log.Panicln("非法的优先级值!")
	}

	lock.Lock()
	defer lock.Unlock()

	if fns, ok := initialFuncMap[order]; ok {
		initialFuncMap[order] = append(fns, f)
	} else {
		initialFuncMap[order] = funcList{f}
	}
}

func InitAll() {
	lock.RLock()
	defer lock.RUnlock()

	sortedKeys := make([]int, 0, len(initialFuncMap))

	for key := range initialFuncMap {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Ints(sortedKeys)

	for i := 0; i < len(sortedKeys); i++ {
		fns := initialFuncMap[sortedKeys[i]]

		for _, fn := range fns {
			fn()
		}
	}
}
