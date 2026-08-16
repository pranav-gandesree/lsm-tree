package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type WALRecord[K any, V any] struct {
	Operation string `json:"operation"`
	Key       K      `json:"key"`
	Value     *V     `json:"value,omitempty"` // pointer makes value optional
}

func AppendData(data string) error {

	// f, err := os.Create("data/wal.log") //destroys the existing logs too
	f, err := os.OpenFile(
		"data/wal.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err // donot use log.fatal, it ends and doesnt run f.close
	}

	defer f.Close()

	_, err = f.WriteString(data + "\n")

	fmt.Println("appended data to log")

	return err
}

func ReplayWal[K comparable, V any]() ([]WALRecord[K, V], error) {
	f, err := os.Open("data/wal.log")

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []WALRecord[K, V]{}, nil
		}
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var records []WALRecord[K, V]
	for scanner.Scan() {
		line := scanner.Text()
		var record WALRecord[K, V]
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}
