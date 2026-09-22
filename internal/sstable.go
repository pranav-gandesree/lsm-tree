package kv

import (
	"encoding/json"
	"fmt"
	"os"
)

type SSTableRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

func ReadSSTable(filename string) ([]SSTableRecord, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []SSTableRecord

	decoder := json.NewDecoder(file)

	for decoder.More() {
		var record SSTableRecord

		if err := decoder.Decode(&record); err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil

}

func GetFromSSTable(filename string, key string) (string, bool, error) {
	records, err := ReadSSTable(filename)
	if err != nil {
		return "", false, err
	}

	for _, record := range records {
		if record.Key == key {
			return record.Value, true, nil
		}
	}

	return "", false, nil
}

func GetFromSSTables(key string, latestSSTableId int) (string, bool, error) {
	for id := latestSSTableId; id >= 1; id-- {
		filename := fmt.Sprintf("data/sstable-%d.db", id)

		value, found, err := GetFromSSTable(filename, key)
		if err != nil {
			return "", false, err
		}

		if found {
			return value, true, nil
		}
	}

	return "", false, nil
}
