package kv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type SSTableRecord struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Tombstone bool   `json:"tombstone"`
}

func createSSTable(records []SSTableRecord, ssTableCounter int) error {
	filename := fmt.Sprintf("data/sstable-%d.db", ssTableCounter)

	file, err := os.Create(filename)

	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)

	for _, record := range records {
		if err := encoder.Encode(record); err != nil {
			return err
		}
	}
	return nil
}

func (s *MemTable) Get(key string, latestSSTableID int) (string, error) {
	// check active memtable first
	value, err := s.GetData(key)

	if err == nil {
		return value, nil
	}

	// key was deleted in active memtable
	if errors.Is(err, ErrDeleted) {
		return "", ErrNotFound
	}

	// if not found, continue with sstables
	if !errors.Is(err, ErrNotFound) {
		return "", err
	}

	// search ssables newest to olldest
	return GetFromSSTables(key, latestSSTableID)
}

func GetFromSSTable(filename string, key string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	for {
		var record SSTableRecord

		err := decoder.Decode(&record)

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", err
		}

		if record.Key == key {
			if record.Tombstone {
				return "", ErrDeleted
			}

			return record.Value, nil
		}

		if record.Key > key {
			break
		}
	}

	return "", ErrNotFound
}

func GetFromSSTables(key string, latestSSTableID int) (string, error) {
	for id := latestSSTableID; id >= 1; id-- {
		filename := fmt.Sprintf("data/sstable-%d.db", id)

		value, err := GetFromSSTable(filename, key)

		if err == nil {
			return value, nil
		}

		if errors.Is(err, ErrDeleted) {
			return "", ErrNotFound
		}

		if errors.Is(err, ErrNotFound) {
			continue
		}

		return "", err
	}

	return "", ErrNotFound
}
