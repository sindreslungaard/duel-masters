package internal

import "sync"

type ConcurrentDictionary[T any] struct {
	sync.RWMutex
	src         []*T
	keyValStore map[string]*T
}

func NewConcurrentDictionary[T any]() ConcurrentDictionary[T] {
	return ConcurrentDictionary[T]{
		src:         []*T{},
		keyValStore: make(map[string]*T),
	}
}

func (d *ConcurrentDictionary[T]) Iter() []*T {

	d.RLock()
	defer d.RUnlock()

	tmp := make([]*T, len(d.src))

	copy(tmp, d.src)

	return tmp

}

func (d *ConcurrentDictionary[T]) Add(key string, value *T) {
	d.Lock()
	defer d.Unlock()

	_, ok := d.keyValStore[key]

	// already added
	if ok {
		return
	}

	d.src = append(d.src, value)
	d.keyValStore[key] = value
}

func (d *ConcurrentDictionary[T]) Remove(key string) {
	d.Lock()
	defer d.Unlock()

	value, ok := d.keyValStore[key]

	if !ok {
		return
	}

	toRemove := -1
	for i, v := range d.src {
		if v == value {
			toRemove = i
			break
		}
	}

	if toRemove < 0 {
		return
	}

	d.src = append(d.src[:toRemove], d.src[toRemove+1:]...)
	delete(d.keyValStore, key)
}

func (d *ConcurrentDictionary[T]) Find(key string) (*T, bool) {
	d.RLock()
	defer d.RUnlock()

	value, ok := d.keyValStore[key]

	return value, ok
}
