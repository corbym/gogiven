package gogiven

import "sync"
// SafeMap is used internally to hold a threadsafe copy of the global test state.
type SafeMap struct {
	sync.RWMutex
	internal map[string]interface{}
}

func newSafeMap() *SafeMap {
	return &SafeMap{
		internal: make(map[string]interface{}),
	}
}
//Load a key from the map
func (rm *SafeMap) Load(key string) (value interface{}, ok bool) {
	rm.RLock()
	defer rm.RUnlock()
	result, ok := rm.internal[key]
	return result, ok
}
//Delete a key from the map
func (rm *SafeMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}
//Store a value against a key from the map
func (rm *SafeMap) Store(key string, value interface{}) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
//Keysreturns an array of keys that the map contains
func (rm *SafeMap) Keys() []string {
	rm.RLock()
	defer rm.RUnlock()
	keys := make([]string, 0, len(rm.internal))
	for k := range rm.internal {
		keys = append(keys, k)
	}
	return keys
}
// Len reports the lenght of the map, same as len(myMap) for the primitive go map.
func (rm *SafeMap) Len() int {
	rm.RLock()
	defer rm.RUnlock()
	return len(rm.internal)
}
