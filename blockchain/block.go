package blockchain

import (
	"cxchain-2023131080/common"
	"cxchain-2023131080/crypto/sha3"
	"cxchain-2023131080/mpt"
	"cxchain-2023131080/txpool"
	"cxchain-2023131080/utils/hash"
	"cxchain-2023131080/utils/rlp"
)

type Header struct {
	Root       hash.Hash
	ParentHash hash.Hash
	Height     uint64
	Coinbase   common.Address
	Timestamp  uint64

	Nonce uint64
}

type Body struct {
	Transactions []common.Transaction
	Receiptions  []common.Receiption
}

func (header Header) Hash() hash.Hash {
	data, _ := rlp.EncodeToBytes(header)
	return sha3.Keccak256(data)
}

func NewHeader(parent Header) *Header {
	return &Header{
		Root:       parent.Root,
		ParentHash: parent.Hash(),
		Height:     parent.Height + 1,
	}
}

func NewBlock() *Body {
	return &Body{
		Transactions: make([]common.Transaction, 0),
		Receiptions:  make([]common.Receiption, 0),
	}
}

type Blockchain struct {
	CurrentHeader Header
	Statedb       mpt.Trie
	Txpool        txpool.Txpool
}
