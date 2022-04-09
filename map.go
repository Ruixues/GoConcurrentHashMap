package ConcurrentHashMap

import (
	"sync/atomic"

	"github.com/mitchellh/hashstructure/v2"
)

const maxSegment = 65536

type pair[K comparable, V any] struct {
	Key   K
	Value V
}

type Map[K comparable, V any] struct {
	segments     []*segment[K, V]
	segmentShift uint64
	segmentMask  uint64
	size         int64
}

func NewMap[K comparable, V any](concurrentLevel uint64, initialCapacity uint64) *Map[K, V] {
	var segmentSize uint64 = 1
	if concurrentLevel > maxSegment {
		segmentSize = maxSegment
	}
	var segmentShift uint64 = 64
	for segmentSize < concurrentLevel {
		segmentSize = segmentSize << 1
		segmentShift--
	}
	segmentMask := segmentSize - 1

	var size uint64 = initialCapacity / segmentSize
	if size*segmentSize < initialCapacity {
		size++
	}
	var realSize uint64 = 1
	for realSize < size {
		realSize = realSize << 1
	}
	size = realSize

	segments := make([]*segment[K, V], segmentSize)
	for _, segment := range segments {
		segment.values = make([]*hashEntry[K, V], size)
		segment.len = size
		segment.lenMask = size - 1
	}
	return &Map[K, V]{
		segments:     segments,
		segmentShift: segmentShift,
		segmentMask:  segmentMask,
	}
}
func (m *Map[K, V]) Put(k K, v V) error {
	hash, err := hashstructure.Hash(k, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}
	segmentIndex := (hash >> m.segmentShift) & m.segmentMask
	s := m.segments[segmentIndex]
	delta, err := s.put(hash, k, v)
	if err != nil {
		return err
	}
	atomic.AddInt64(&m.size, delta)
	return nil
}
func (m *Map[K, V]) Remove(key K) error {
	hash, err := hashstructure.Hash(key, hashstructure.FormatV2, nil)
	if err != nil {
		return err
	}
	segmentIndex := (hash >> m.segmentShift) & m.segmentMask
	s := m.segments[segmentIndex]
	delta, err := s.remove(hash, key)
	if err != nil {
		return err
	}
	atomic.AddInt64(&m.size, delta)
	return nil
}
func (m *Map[K, V]) Size() int64 {
	return atomic.LoadInt64(&m.size)
}
