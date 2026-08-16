package store

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func CreateStore[K comparable, V any]() (*Store[K, V], error) {
	records, err := ReplayWal[K, V]()
	if err != nil {
		return nil, fmt.Errorf("ReplayWAL failed: %w", err)
	}

	store := &Store[K, V]{
		data: make(map[K]V),
	}

	for _, record := range records {
		fmt.Printf(
			"operation=%s key=%v value=%v\n",
			record.Operation,
			record.Key,
			record.Value,
		)

		switch record.Operation {
		case "PUT":
			if record.Value == nil {
				return nil, fmt.Errorf(
					"invalid PUT record: key=%v has nil value",
					record.Key,
				)
			}

			store.data[record.Key] = *record.Value

		case "DELETE":
			delete(store.data, record.Key)

		default:
			return nil, fmt.Errorf(
				"unknown WAL operation: %q",
				record.Operation,
			)
		}
	}

	return store, nil
}

func (s *Store[K, V]) PutData(key K, value V) error {
	s.mu.Lock() //only 1 goroutine can write at a time
	defer s.mu.Unlock()

	record := WALRecord[K, V]{
		Operation: "PUT",
		Key:       key,
		Value:     &value,
	}
	// wal.AppendData(fmt.Sprintf("%v", record))

	data, err := json.Marshal(record)

	if err != nil {
		return err
	}

	if err := AppendData(string(data)); err != nil {
		return err
	}

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
		Value:     nil,
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

	if err := AppendData(string(data)); err != nil {
		return err
	}

	delete(s.data, key)

	return nil
}

func (s *Store[K, V]) PrintMap() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("Current map:")

	for key, value := range s.data {
		fmt.Printf("key=%v, value=%+v\n", key, value)
	}
}
