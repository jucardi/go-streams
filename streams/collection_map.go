package streams

import "sync"

var (
	// To ensure implementations on build
	_ IList[*KeyValuePair[string, string]]               = (*mapCollection[string, string])(nil)
	_ IMap[string, string]                               = (*mapCollection[string, string])(nil)
	_ IAbstractCollection[*KeyValuePair[string, string]] = (*mapCollection[string, string])(nil)
)

type mapCollection[K, V comparable] struct {
	*CollectionBase[*KeyValuePair[K, V]]
	m    map[K]V
	keys []K
	mx   sync.RWMutex
}

func (col *mapCollection[K, V]) Get(k K) (val V, ok bool) {
	col.mx.RLock()
	defer col.mx.RUnlock()
	val, ok = col.m[k]
	return
}

func (col *mapCollection[K, V]) Set(key K, value V) bool {
	col.mx.Lock()
	defer col.mx.Unlock()
	l := len(col.m)
	col.m[key] = value
	if len(col.m) > l {
		col.keys = append(col.keys, key)
	}
	return true
}

func (col *mapCollection[K, V]) ContainsKey(k K) bool {
	_, ok := col.Get(k)
	return ok
}

func (col *mapCollection[K, V]) Keys() []K {
	col.mx.RLock()
	defer col.mx.RUnlock()
	if len(col.m) == len(col.keys) {
		return col.keys
	}

	var ret []K
	for k := range col.m {
		ret = append(ret, k)
	}
	col.keys = ret
	return ret
}

func (col *mapCollection[K, V]) Delete(k K) bool {
	col.mx.Lock()
	defer col.mx.Unlock()

	l := len(col.m)
	delete(col.m, k)
	if l == len(col.m) {
		return false
	}
	for i, key := range col.keys {
		if k != key {
			continue
		}
		col.keys = append(col.keys[:i], col.keys[i:]...)
		return true
	}
	return false
}

func (col *mapCollection[K, V]) ToMap() map[K]V {
	return col.m
}

func (col *mapCollection[K, V]) Index(index int) (val *KeyValuePair[K, V], exists bool) {
	keys := col.Keys()
	if index < 0 || len(keys) <= index {
		return
	}
	key := keys[index]
	if v, exists := col.Get(key); exists {
		return &KeyValuePair[K, V]{Key: key, Value: v}, true
	}
	return
}

func (col *mapCollection[K, V]) Add(items ...*KeyValuePair[K, V]) (ret bool) {
	col.mx.Lock()
	defer col.mx.Unlock()
	for _, item := range items {
		col.m[item.Key] = item.Value
	}
	return len(items) > 0
}

func (col *mapCollection[K, V]) RemoveAt(index int, _ ...bool) bool {
	keys := col.Keys()
	if index < 0 || len(keys) >= index {
		return false
	}
	return col.Delete(keys[index])
}

func (col *mapCollection[K, V]) Len() int {
	col.mx.RLock()
	defer col.mx.RUnlock()
	return len(col.m)
}

func (col *mapCollection[K, V]) Clear() {
	col.mx.Lock()
	defer col.mx.Unlock()
	col.m = make(map[K]V)
	col.keys = nil
}

func (col *mapCollection[K, V]) ToArray() (ret []*KeyValuePair[K, V]) {
	col.ForEach(func(item *KeyValuePair[K, V]) {
		ret = append(ret, item)
	})
	return
}
