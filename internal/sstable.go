package kv

import (
	"encoding/json"
	"os"
)

type SSTableRecord struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func createSSTable(records []SSTableRecord) error {
	file, err := os.Create("data/sstable-1.db")

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
