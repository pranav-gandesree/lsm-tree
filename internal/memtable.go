package kv

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
)

const MemTableLimit = 5

type MemTable struct {
	mu   sync.RWMutex
	data map[string]string
}

//it flushes the active mem table while holding its read lock
// func (m *MemTable) FlushToSSTable() error {
// 	m.mu.RLock()
// 	defer m.mu.RUnlock()

// 	records := make([]SSTableRecord, 0, len(m.data))

// 	for key, value := range m.data {
// 		records = append(records, SSTableRecord{
// 			Key:   key,
// 			Value: value,
// 		})
// 	}

// 	sort.Slice(records, func(i, j int) bool {
// 		return records[i].Key < records[j].Key
// 	})

// 	return createSSTable(records)
// }

func FlushToSSTable(data map[string]string) error {
	records := make([]SSTableRecord, 0, len(data))

	for key, value := range data {
		records = append(records, SSTableRecord{
			Key:   key,
			Value: value,
		})
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Key < records[j].Key
	})

	return createSSTable(records)
}

func CreateMemTable() (*MemTable, error) {
	records, err := ReplayWal()

	if err != nil {
		return nil, fmt.Errorf("ReplayWAL failed: %w", err)
	}

	memtable := &MemTable{
		data: make(map[string]string),
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

			memtable.data[record.Key] = *record.Value

		case "DELETE":
			delete(memtable.data, record.Key)

		default:
			return nil, fmt.Errorf(
				"unknown WAL operation: %q",
				record.Operation,
			)
		}
	}

	return memtable, nil
}

func (s *MemTable) PutData(key string, value string) error {
	s.mu.Lock() //only 1 goroutine can write at a time

	record := WALRecord{
		Operation: "PUT",
		Key:       key,
		Value:     &value,
	}

	data, err := json.Marshal(record)

	if err != nil {
		s.mu.Unlock()
		return err
	}

	if err := AppendData(string(data)); err != nil {
		s.mu.Unlock()
		return err
	}

	s.data[key] = value

	//check if memtable is full,
	if len(s.data) < MemTableLimit {
		s.mu.Unlock()
		return nil
	}

	//if memtable is full, make the current table immutable
	immutableData := s.data

	// then create a new active table
	s.data = make(map[string]string)

	s.mu.Unlock()

	// flush immutable memtable
	return FlushToSSTable(immutableData)

}

func (s *MemTable) GetData(key string) (string, bool) {
	s.mu.RLock() //multiple go routines can read at a time
	defer s.mu.RUnlock()
	value, ok := s.data[key]

	return value, ok
}

func (s *MemTable) DeleteData(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := WALRecord{
		Operation: "DELETE",
		Key:       key,
		Value:     nil,
	}
	// wal.AppendData(fmt.Sprintf("%v", record))

	data, err := json.Marshal(record)
	if err != nil {
		return err
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

func (s *MemTable) PrintMap() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("Current map:")

	for key, value := range s.data {
		fmt.Printf("key=%v, value=%+v\n", key, value)
	}
}
