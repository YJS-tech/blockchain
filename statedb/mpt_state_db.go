package statedb

import (
	"cxchain-2023131080/common"
	"cxchain-2023131080/kvstore"
	"cxchain-2023131080/mpt"
	"cxchain-2023131080/utils/rlp"
	"encoding/hex"
)

// MPTStateDB 使用Merkle Patricia Tree实现StateDB接口
type MPTStateDB struct {
	trie *mpt.Trie
	db   kvstore.KVStore
}

// accountSerializable 用于JSON序列化的辅助结构体
type accountSerializable struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

func NewMPTStateDB(db kvstore.KVStore) *MPTStateDB {
	return &MPTStateDB{
		trie: mpt.NewTrie(db),
		db:   db,
	}
}

// GetAccount 获取账户信息
func (s *MPTStateDB) GetAccount(addr common.Address) *common.Account {
	key := addr.Bytes()
	data, ok := s.trie.Get(key)
	if !ok {
		return nil
	}

	var acc common.Account
	if err := rlp.DecodeBytes(data, &acc); err != nil {
		return nil
	}

	return &acc
}

// UpdateAccount 更新账户信息
func (s *MPTStateDB) UpdateAccount(acc *common.Account) {
	key := acc.Address.Bytes()
	data, err := rlp.EncodeToBytes(acc)
	if err != nil {
		panic(err) // 或者返回错误
	}
	s.trie.Update(key, data)
}

// DeleteAccount 删除账户
func (s *MPTStateDB) DeleteAccount(addr common.Address) {
	key := addr.Bytes()
	s.trie.Delete(key)
}

// Commit 提交状态并返回根哈希
func (s *MPTStateDB) Commit() []byte {
	if s.trie.Root == nil {
		return nil
	}
	return s.trie.Root.Hash()
}

// CommitHex 提交状态并返回根哈希的十六进制字符串
func (s *MPTStateDB) CommitHex() string {
	hash := s.Commit()
	if hash == nil {
		return ""
	}
	return hex.EncodeToString(hash)
}

// SetRoot 根据 rootHash 重建 trie（需要 MPT 支持）
func (s *MPTStateDB) SetRoot(rootHash []byte) error {
	newTrie, err := mpt.NewTrieWithRootHash(s.db, rootHash)
	if err != nil {
		return err
	}
	s.trie = newTrie
	return nil
}
