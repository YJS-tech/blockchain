package statedb

import "cxchain-2023131080/common"

type StateDB interface {
	GetAccount(addr common.Address) *common.Account
	UpdateAccount(acc *common.Account)
	DeleteAccount(addr common.Address)
	Commit() []byte    // 返回根哈希 Bytes
	CommitHex() string // 返回状态根 Hex 字符串
	SetRoot(rootHash []byte) error
}
