package kv

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
)

const MemTableLimit = 2

var ErrNotFound = errors.New("key not found")
var ErrDeleted = errors.New("key deleted")

type MemTableValue struct {
	Value     string
	Tombstone bool
}
type MemTable struct {
	mu             sync.RWMutex
	data           map[string]MemTableValue
	sstableCounter int
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

func FlushToSSTable(data map[string]MemTableValue, ssTableCounter int) error {
	records := make([]SSTableRecord, 0, len(data))

	for key, value := range data {
		records = append(records, SSTableRecord{
			Key:       key,
			Value:     value.Value,
			Tombstone: value.Tombstone,
		})
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Key < records[j].Key
	})

	return createSSTable(records, ssTableCounter)
}

func CreateMemTable() (*MemTable, error) {
	records, err := ReplayWal()

	if err != nil {
		return nil, fmt.Errorf("ReplayWAL failed: %w", err)
	}

	memtable := &MemTable{
		data:           make(map[string]MemTableValue),
		sstableCounter: 0,
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

			memtable.data[record.Key] = MemTableValue{
				Value:     *record.Value,
				Tombstone: false,
			}

		case "DELETE":
			memtable.data[record.Key] = MemTableValue{
				Tombstone: true,
			}

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

	s.data[key] = MemTableValue{
		Value:     value,
		Tombstone: false,
	}

	//check if memtable is full,
	if len(s.data) < MemTableLimit {
		s.mu.Unlock()
		return nil
	}

	//if memtable is full, make the current table immutable
	immutableData := s.data

	// then create a new active table
	s.data = make(map[string]MemTableValue)
	//increment the counter if memtable is full, so it can create new sstable file
	s.sstableCounter++
	sstableId := s.sstableCounter

	s.mu.Unlock()

	// flush immutable memtable
	return FlushToSSTable(immutableData, sstableId)

}

func (s *MemTable) GetData(key string) (string, error) {
	s.mu.RLock() //multiple go routines can read at a time
	defer s.mu.RUnlock()

	entry, ok := s.data[key]

	if !ok {
		return "", ErrNotFound
	}

	if entry.Tombstone {
		return "", ErrDeleted
	}

	return entry.Value, nil
}

func (s *MemTable) DeleteData(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := WALRecord{
		Operation: "DELETE",
		Key:       key,
		Value:     nil,
	}

	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	if err := AppendData(string(data)); err != nil {
		return err
	}

	s.data[key] = MemTableValue{
		Value:     "",
		Tombstone: true,
	}

	return nil
}

// LatestSSTableID returns the id of the most recently flushed SSTable.
// Callers use it as the starting point for Get's newest-to-oldest search.
func (s *MemTable) LatestSSTableID() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sstableCounter
}

func (s *MemTable) PrintMap() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println("Current map:")

	for key, value := range s.data {
		fmt.Printf("key=%v, value=%+v\n", key, value)
	}
}
