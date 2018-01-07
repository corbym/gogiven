package gogiven

import "sync"

type SafeMap struct {
	sync.RWMutex
	internal map[string]interface{}
}

func newSafeMap() *SafeMap {
	return &SafeMap{
		internal: make(map[string]interface{}),
	}
}

func (rm *SafeMap) Load(key string) (value interface{}, ok bool) {
	rm.RLock()
	defer rm.RUnlock()
	result, ok := rm.internal[key]
	return result, ok
}

func (rm *SafeMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *SafeMap) Store(key string, value interface {}) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
func (rm *SafeMap) Keys() []string {
	rm.RLock()
	defer rm.RUnlock()
	keys := make([]string, 0, len(rm.internal))
	for k := range rm.internal {
		keys = append(keys, k)
	}
	return keys
}
func (rm *SafeMap) Len() int{
	rm.RLock()
	defer rm.RUnlock()
	return len(rm.internal)
}