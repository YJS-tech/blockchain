package leveldb_test

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"testing"
)

func TestLevelDBExample(t *testing.T) {
	//create database temporary
	dbPath := "testdb"
	defer os.RemoveAll(dbPath)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		t.Errorf("database open failed:%v", err)
	}

	defer db.Close()

	//put data
	err = db.Put([]byte("key1"), []byte("value1"), nil)
	if err != nil {
		t.Errorf("database put failed:%v", err)
	}

	//get data
	val, err := db.Get([]byte("key1"), nil)
	if err != nil || string(val) != "value1" {
		t.Errorf("get failed or value not the same: got %s", val)
	}

	//has data
	exists, _ := db.Has([]byte("key1"), nil)
	if !exists {
		t.Error("should exists key 'key1'")
	}

	//delete data
	_ = db.Delete([]byte("key1"), nil)
	exists, _ = db.Has([]byte("key1"), nil)
	if exists {
		t.Error("should NOT exists key 'key1'")
	}

	fmt.Printf("levelDB test passed!")
}
