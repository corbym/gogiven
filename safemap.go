package gogiven

import "sync"

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

//Delete a key from the map
func (rm *safeMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
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

// Len reports the length of the map, same as len(myMap) for the primitive go map.
func (rm *safeMap) Len() int {
	rm.RLock()
	defer rm.RUnlock()
	return len(rm.internal)
}

// AsMapOfSome copies the safeMap into a normal map[string]*Some type
func (rm *safeMap) AsMapOfSome() map[string]*Some {
	rm.RLock()
	defer rm.RUnlock()
	newMap := make(map[string]*Some, len(rm.internal))
	for k, v := range rm.internal {
		newMap[k] = v.(*Some)
	}
	return newMap
}
