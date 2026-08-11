package kv

import "sync"

type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func CreateStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		data: make(map[K]V),
	}
}

func (s *Store[K, V]) PutData(key K, value V) {
	s.mu.Lock() //only 1 goroutine can write at a time
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store[K, V]) GetData(key K) (V, bool) {
	s.mu.RLock() //multiple go routines can read at a time
	defer s.mu.RUnlock()
	value, ok := s.data[key]

	return value, ok
}

func (s *Store[K, V]) DeleteData(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
