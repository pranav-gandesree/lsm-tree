package kv

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestMergeSSTables(t *testing.T) {
	chdirToTemp(t)
	if err := os.MkdirAll("data", 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	// id 1: a, b
	must(t, FlushToSSTable(map[string]MemTableValue{
		"a": {Value: "apple"},
		"b": {Value: "banana"},
	}, 1))

	// id 2: tombstones b, adds c (v1)
	must(t, FlushToSSTable(map[string]MemTableValue{
		"b": {Tombstone: true},
		"c": {Value: "cherry-v1"},
	}, 2))

	// id 3: c (v2, should win over v1), adds d
	must(t, FlushToSSTable(map[string]MemTableValue{
		"c": {Value: "cherry-v2"},
		"d": {Value: "date"},
	}, 3))

	mt := &MemTable{sstableCounter: 3}

	if err := mt.MergeSSTables(); err != nil {
		t.Fatalf("MergeSSTables: %v", err)
	}

	files, err := filepath.Glob("data/sstable-*.db")
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}
	if len(files) != 1 || files[0] != "data/sstable-1.db" {
		t.Fatalf("expected only data/sstable-1.db to remain, got %v", files)
	}

	if got := mt.LatestSSTableID(); got != 1 {
		t.Fatalf("LatestSSTableID() = %d, want 1", got)
	}

	assertValue(t, "a", "apple")
	assertNotFound(t, "b") // tombstoned, dropped by full compaction
	assertValue(t, "c", "cherry-v2") // proves Sequence tiebreak, not just key sort
	assertValue(t, "d", "date")
}

func TestMergeSSTablesAllTombstoned(t *testing.T) {
	chdirToTemp(t)
	if err := os.MkdirAll("data", 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	must(t, FlushToSSTable(map[string]MemTableValue{
		"x": {Tombstone: true},
	}, 1))

	mt := &MemTable{sstableCounter: 1}

	if err := mt.MergeSSTables(); err != nil {
		t.Fatalf("MergeSSTables: %v", err)
	}

	files, err := filepath.Glob("data/sstable-*.db")
	if err != nil {
		t.Fatalf("Glob: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected no sstable files to remain, got %v", files)
	}

	if got := mt.LatestSSTableID(); got != 0 {
		t.Fatalf("LatestSSTableID() = %d, want 0", got)
	}
}

// chdirToTemp changes the working directory to a fresh temp dir for the
// duration of the test, restoring the original on cleanup. Every path in
// this package is a hardcoded "data/..." literal, so working-directory
// isolation is the only way to sandbox a test.
func chdirToTemp(t *testing.T) {
	t.Helper()

	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}

	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatalf("Chdir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(original); err != nil {
			t.Fatalf("Chdir restore: %v", err)
		}
	})
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}
}

func assertValue(t *testing.T, key, want string) {
	t.Helper()
	got, err := GetFromSSTables(key, 1)
	if err != nil {
		t.Fatalf("GetFromSSTables(%q) error = %v, want nil", key, err)
	}
	if got != want {
		t.Fatalf("GetFromSSTables(%q) = %q, want %q", key, got, want)
	}
}

func assertNotFound(t *testing.T, key string) {
	t.Helper()
	_, err := GetFromSSTables(key, 1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetFromSSTables(%q) error = %v, want ErrNotFound", key, err)
	}
}
