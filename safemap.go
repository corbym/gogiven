package gogiven

import (
	"github.com/corbym/gogiven/base"
	"sync"
)

// safeMap is used internally to hold a threadsafe copy of the global test state.
type safeMap struct {
	sync.RWMutex
	internal map[string]interface{}
}

func newSafeMap() *safeMap {
	return &safeMap{
		internal: make(map[string]interface{}),
	}
}

//Load a key from the map
func (rm *safeMap) Load(key string) (value interface{}, ok bool) {
	rm.RLock()
	defer rm.RUnlock()
	result, ok := rm.internal[key]
	return result, ok
}

//Store a value against a key from the map
func (rm *safeMap) Store(key string, value interface{}) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}

//Keys returns an array of keys that the map contains
func (rm *safeMap) Keys() []string {
	rm.RLock()
	defer rm.RUnlock()
	keys := make([]string, 0, len(rm.internal))
	for k := range rm.internal {
		keys = append(keys, k)
	}
	return keys
}

// AsMapOfSome copies the safeMap into a normal map[string]*Some type
func (rm *safeMap) AsMapOfSome() *base.SomeMap {
	rm.RLock()
	defer rm.RUnlock()
	base.CopyLock.Lock()
	defer base.CopyLock.Unlock()
	var newMap = &base.SomeMap{}
	for k, v := range rm.internal {
		(*newMap)[k] = v.(*base.Some)
	}
	return newMap
}
