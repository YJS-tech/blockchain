package mpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"cxchain-2023131080/kvstore"
	"cxchain-2023131080/mpt/utils"
)

type Trie struct {
	db   kvstore.KVStore
	Root Node
}

func NewTrie(db kvstore.KVStore) *Trie {
	return &Trie{db: db}
}

func (t *Trie) Update(key, value []byte) {
	nibbles := utils.HexToNibbles(key)
	t.Root = insert(t, t.Root, nibbles, value)
}

func (t *Trie) Get(key []byte) ([]byte, bool) {
	nibbles := utils.HexToNibbles(key)
	val, ok := get(t.Root, nibbles)
	return val, ok
}

func (t *Trie) Delete(key []byte) {
	nibbles := utils.HexToNibbles(key)
	t.Root, _ = delete(t, t.Root, nibbles)
}

func (t *Trie) saveNode(n Node) {
	if n == nil {
		return
	}
	hash := n.Hash()
	data := n.Serialize()
	t.db.Put(hash, data)
}

func NewTrieWithRootHash(db kvstore.KVStore, rootHash []byte) (*Trie, error) {
	data, err := db.Get(rootHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get root node: %w", err)
	}
	node := deserializeNode(db, data)
	if node == nil {
		return nil, errors.New("failed to deserialize root node")
	}
	return &Trie{db: db, Root: node}, nil
}

func get(n Node, key []byte) ([]byte, bool) {
	switch node := n.(type) {
	case nil:
		return nil, false
	case *LeafNode:
		if bytes.Equal(key, node.Path) && len(node.Value) > 0 {
			return node.Value, true
		}
	case *ExtensionNode:
		if bytes.HasPrefix(key, node.Path) {
			return get(node.Next, key[len(node.Path):])
		}
	case *BranchNode:
		if len(key) == 0 {
			if node.Value != nil && len(node.Value) > 0 {
				return node.Value, true
			}
			return nil, false
		}
		return get(node.Children[key[0]], key[1:])
	}
	return nil, false
}

func insert(t *Trie, n Node, key []byte, value []byte) Node {
	if n == nil {
		leaf := &LeafNode{Type: "leaf", Path: key, Value: value}
		t.saveNode(leaf)
		return leaf
	}

	switch node := n.(type) {
	case *LeafNode:
		if bytes.Equal(key, node.Path) {
			node.Value = value
			t.saveNode(node)
			return node
		}

		commonLen := commonPrefixLen(node.Path, key)
		bn := &BranchNode{Type: "branch"}

		if commonLen > 0 {
			if commonLen < len(node.Path) {
				bn.Children[node.Path[commonLen]] = insert(t, nil, safeSlice(node.Path, commonLen+1), node.Value)
			} else {
				bn.Value = node.Value
			}
			if commonLen < len(key) {
				bn.Children[key[commonLen]] = insert(t, nil, safeSlice(key, commonLen+1), value)
			} else {
				bn.Value = value
			}
			ext := &ExtensionNode{Type: "ext", Path: key[:commonLen], Next: bn}
			t.saveNode(bn)
			t.saveNode(ext)
			return ext
		}

		bn.Children[node.Path[0]] = insert(t, nil, safeSlice(node.Path, 1), node.Value)
		bn.Children[key[0]] = insert(t, nil, safeSlice(key, 1), value)
		t.saveNode(bn)
		return bn

	case *ExtensionNode:
		commonLen := commonPrefixLen(node.Path, key)
		if commonLen == len(node.Path) {
			node.Next = insert(t, node.Next, safeSlice(key, commonLen), value)
			t.saveNode(node)
			return node
		}

		newBranch := &BranchNode{Type: "branch"}
		newBranch.Children[node.Path[commonLen]] = insert(t, node.Next, safeSlice(node.Path, commonLen+1), node.Next.(*LeafNode).Value)
		if commonLen < len(key) {
			newBranch.Children[key[commonLen]] = insert(t, nil, safeSlice(key, commonLen+1), value)
		} else {
			newBranch.Value = value
		}

		if commonLen == 0 {
			t.saveNode(newBranch)
			return newBranch
		}

		ext := &ExtensionNode{Type: "ext", Path: key[:commonLen], Next: newBranch}
		t.saveNode(newBranch)
		t.saveNode(ext)
		return ext

	case *BranchNode:
		if len(key) == 0 {
			node.Value = value
		} else {
			idx := key[0]
			node.Children[idx] = insert(t, node.Children[idx], safeSlice(key, 1), value)
		}
		t.saveNode(node)
		return node
	}
	return n
}

func delete(t *Trie, n Node, key []byte) (Node, bool) {
	// TODO delete 的持久化版本
	return n, false
}

func commonPrefixLen(a, b []byte) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

func safeSlice(b []byte, start int) []byte {
	if start >= len(b) {
		return []byte{}
	}
	return b[start:]
}

func deserializeNode(db kvstore.KVStore, data []byte) Node {
	var m map[string]interface{}
	_ = json.Unmarshal(data, &m)
	switch m["type"] {
	case "leaf":
		var l LeafNode
		_ = json.Unmarshal(data, &l)
		return &l
	case "branch":
		var b BranchNode
		_ = json.Unmarshal(data, &b)
		return &b
	case "ext":
		var e ExtensionNode
		_ = json.Unmarshal(data, &e)
		return &e
	}
	return nil
}

func (t *Trie) RootHash() []byte {
	if t.Root == nil {
		return nil
	}
	return t.Root.Hash()
}
