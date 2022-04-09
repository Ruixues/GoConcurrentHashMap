package ConcurrentHashMap

import "sync"

type segment[K comparable, V any] struct {
	values  []*hashEntry[K, V]
	len     uint64
	lock    sync.RWMutex
	lenMask uint64
}

func (s *segment[K, V]) put(hash uint64, key K, value V) (int64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.values[(hash&s.lenMask)].Put(key, value)
}
func (s *segment[K, V]) remove(hash uint64, key K) (int64, error) {
	return s.values[(hash & s.lenMask)].Remove(key)
}
