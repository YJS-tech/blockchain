// mpt/types.go
package mpt

import (
	"cxchain-2023131080/crypto/sha3"
	"encoding/json"
	"fmt"
)

type Node interface {
	Hash() []byte
	String() string
	Serialize() []byte
}

// BranchNode 分支节点

type BranchNode struct {
	Type     string   `json:"type"`
	Children [16]Node `json:"-"`
	RawChild [][]byte `json:"children"`
	Value    []byte   `json:"value"`
}

func (b *BranchNode) String() string {
	return fmt.Sprintf("BranchNode{...}")
}

func (b *BranchNode) Hash() []byte {
	h := []byte{}
	for _, c := range b.Children {
		if c != nil {
			h = append(h, c.Hash()...)
		} else {
			h = append(h, []byte{}...)
		}
	}
	h = append(h, b.Value...)
	return sha3.Keccak256(h)
}

func (b *BranchNode) Serialize() []byte {
	b.Type = "branch"
	raw := make([][]byte, 16)
	for i, child := range b.Children {
		if child != nil {
			raw[i] = child.Hash()
		}
	}
	b.RawChild = raw
	data, _ := json.Marshal(b)
	return data
}

// LeafNode 叶子节点

type LeafNode struct {
	Type  string `json:"type"`
	Path  []byte `json:"path"`
	Value []byte `json:"value"`
}

func (l *LeafNode) String() string {
	return fmt.Sprintf("LeafNode{Path: %x, Value: %x}", l.Path, l.Value)
}

func (l *LeafNode) Hash() []byte {
	return sha3.Keccak256([]byte(fmt.Sprintf("leaf:%x:%x", l.Path, l.Value)))
}

func (l *LeafNode) Serialize() []byte {
	l.Type = "leaf"
	data, _ := json.Marshal(l)
	return data
}

// ExtensionNode 扩展节点

type ExtensionNode struct {
	Type     string `json:"type"`
	Path     []byte `json:"path"`
	NextHash []byte `json:"next"`
	Next     Node   `json:"-"`
}

func (e *ExtensionNode) String() string {
	return fmt.Sprintf("ExtensionNode{Path: %x}", e.Path)
}

func (e *ExtensionNode) Hash() []byte {
	return sha3.Keccak256([]byte(fmt.Sprintf("ext:%x", e.Path)))
}

func (e *ExtensionNode) Serialize() []byte {
	e.Type = "ext"
	if e.Next != nil {
		e.NextHash = e.Next.Hash()
	}
	data, _ := json.Marshal(e)
	return data
}
