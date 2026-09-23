package main

import (
	"errors"
	"fmt"
	"log"

	kv "github.com/pranav-gandesree/lsm-tree/internal"
)

func main() {
	memtable, err := kv.CreateMemTable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("== active memtable after WAL replay ==")
	memtable.PrintMap()

	// MemTableLimit is 2, so this sequence of puts triggers two flushes
	// to SSTable: (a,b) after the 2nd put, (c,d) after the 4th put.
	// "e" stays in the active memtable.
	puts := []struct{ key, value string }{
		{"a", "apple"},
		{"b", "banana"},
		{"c", "cherry"},
		{"d", "date"},
		{"e", "elderberry"},
	}

	for _, p := range puts {
		if err := memtable.PutData(p.key, p.value); err != nil {
			log.Fatal(err)
		}
	}

	// "b" was already flushed to an SSTable, so this writes a tombstone
	// into the active memtable that must shadow the flushed value.
	if err := memtable.DeleteData("b"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n== active memtable after puts + delete ==")
	memtable.PrintMap()

	fmt.Println("\n== reads across memtable + SSTables ==")
	keys := []string{"a", "b", "c", "d", "e", "z"}
	for _, key := range keys {
		value, err := memtable.Get(key, memtable.LatestSSTableID())

		switch {
		case err == nil:
			fmt.Printf("key=%s value=%s\n", key, value)
		case errors.Is(err, kv.ErrNotFound):
			fmt.Printf("key=%s not found\n", key)
		default:
			log.Fatalf("get %s: %v", key, err)
		}
	}
}
