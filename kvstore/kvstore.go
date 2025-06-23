package kvstore

import "io"

// KVStore 定义键值存储的基本接口
type KVStore interface {
	// Get 检索与给定键关联的值
	Get(key []byte) ([]byte, error)
	// Put 存储键值对
	Put(key, value []byte) error
	// Delete 删除与给定键关联的值
	Delete(key []byte) error
	// Has 检查键是否存在
	Has(key []byte) (bool, error)
	io.Closer
}
