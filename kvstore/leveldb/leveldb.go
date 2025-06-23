package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// LevelDB 实现KVStore接口
type LevelDB struct {
	db *leveldb.DB
}

// Open 创建 LevelDB 实例
func Open(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{db: db}, nil
}

// Get 实现KVStore接口的Get方法
func (l *LevelDB) Get(key []byte) ([]byte, error) {
	return l.db.Get(key, nil)
}

// Put 实现KVStore接口的Put方法
func (l *LevelDB) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

// Delete 实现KVStore接口的Delete方法
func (l *LevelDB) Delete(key []byte) error {
	return l.db.Delete(key, nil) // 保持参数签名与接口一致
}

// Has 实现KVStore接口的Has方法
func (l *LevelDB) Has(key []byte) (bool, error) {
	return l.db.Has(key, nil)
}

// Close 实现KVStore接口的Close方法
func (l *LevelDB) Close() error {
	return l.db.Close()
}
