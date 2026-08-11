package kv

import (
	"encoding/json"
	"log"
	"lsm/internal/wal"
	"sync"
)

type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}
type WALRecord[K any, V any] struct {
	Operation string `json:"operation"`
	Key       K      `json:"key"`
	Value     V      `json:"value"`
}

func CreateStore[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{
		data: make(map[K]V),
	}
}

func (s *Store[K, V]) PutData(key K, value V) error {
	s.mu.Lock() //only 1 goroutine can write at a time
	defer s.mu.Unlock()

	record := WALRecord[K, V]{
		Operation: "PUT",
		Key:       key,
		Value:     value,
	}
	// wal.AppendData(fmt.Sprintf("%v", record))

	data, err := json.Marshal(record)

	if err != nil {
		return err
	}

	if err := wal.AppendData(string(data)); err != nil {
		return err
	}

	wal.AppendData(string(data))

	s.data[key] = value

	return nil
}

func (s *Store[K, V]) GetData(key K) (V, bool) {
	s.mu.RLock() //multiple go routines can read at a time
	defer s.mu.RUnlock()
	value, ok := s.data[key]

	return value, ok
}

func (s *Store[K, V]) DeleteData(key K) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := WALRecord[K, V]{
		Operation: "DELETE",
		Key:       key,
	}
	// wal.AppendData(fmt.Sprintf("%v", record))

	data, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}

	// 	err = wal.AppendData(string(data))
	// if err != nil {
	// 	return err
	// }

	if err := wal.AppendData(string(data)); err != nil {
		return err
	}

	wal.AppendData(string(data))

	delete(s.data, key)

	return nil
}
