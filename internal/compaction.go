package kv

import (
	"os"
	"path/filepath"
	"sort"
)

func getAllSSTables() []string {
	tables, _ := filepath.Glob("data/sstable-*.db")

	return tables
}


func (s *MemTable) MergeSSTables() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//get all sstables
	tables := getAllSSTables()
	if len(tables) == 0 {
		return nil
	}

	var records []SSTableRecord

	for _, table := range tables {
		tablerecords, err := readSSTable(table)
		if err != nil {
			return err
		}

		records = append(records, tablerecords...)
	}

	//sort the records by key & sequence, larger the seq no, newer the record
	sort.Slice(records, func(i, j int) bool {
		if records[i].Key == records[j].Key {
			return records[i].Sequence > records[j].Sequence
		}

		return records[i].Key < records[j].Key
	})

	//then keep the newest version of each key
	seen := make(map[string]bool)

	var compacted []SSTableRecord
	for _, record := range records {
		if seen[record.Key] {
			continue
		}

		seen[record.Key] = true
		// as we are doing full compaction, tombstones can be dropped
		if record.Tombstone {
			continue
		}

		compacted = append(compacted, record)
	}

	// delete old sstables first, the merged output is always name sstable-1.db
	for _, table := range tables {
		if err := os.Remove(table); err != nil {
			return err
		}
	}

	if len(compacted) == 0 {
		s.sstableCounter = 0
		return nil
	}

	//create new sstable
	if err := createSSTable(compacted, 1); err != nil {
		return err
	}

	s.sstableCounter = 1

	return nil
}
