package leveldb_test

import (
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"strconv"
	"testing"
)

func BenchmarkLevelDB_Put(b *testing.B) {
	dbPath := "benchdb"
	defer os.Remove(dbPath)

	db, _ := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	for i := 0; i < b.N; i++ {
		key := []byte("key" + strconv.Itoa(i))
		value := []byte("value" + strconv.Itoa(i))
		db.Put(key, value, nil)
	}
}

func BenchmarkLevelDB_Get(b *testing.B) {
	dbPath := "benchdb"
	defer os.Remove(dbPath)

	db, _ := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	_ = db.Put([]byte("key1"), []byte("value1"), nil)

	for i := 0; i < b.N; i++ {
		db.Get([]byte("key1"), nil)
	}
}

func BenchmarkLevelDB_Delete(b *testing.B) {
	dbPath := "benchdb"
	defer os.Remove(dbPath)

	db, _ := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	key := []byte("key")
	value := []byte("value")
	for i := 0; i < b.N; i++ {
		_ = db.Put(key, value, nil)
		_ = db.Delete(key, nil)
	}
}

func BenchmarkLevelDB_Has(b *testing.B) {
	dbPath := "benchdb"
	defer os.Remove(dbPath)

	db, _ := leveldb.OpenFile(dbPath, nil)
	defer db.Close()

	key := []byte("key")
	_ = db.Put(key, []byte("value"), nil)

	for i := 0; i < b.N; i++ {
		_, _ = db.Has(key, nil)
	}
}
