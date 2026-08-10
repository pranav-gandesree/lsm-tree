package kv

import "sync"

type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

func CreateStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) PutData(key, value string) {
	s.mu.Lock() //only 1 goroutine can write at a time
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) GetData(key string) (string, bool) {
	s.mu.RLock() //multiple go routines can read at a time
	defer s.mu.RUnlock()
	value, ok := s.data[key]

	return value, ok
}

func (s *Store) DeleteData(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
