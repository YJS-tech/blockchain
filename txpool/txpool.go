package txpool

import (
	"cxchain-2023131080/common"
	"cxchain-2023131080/crypto/sha3"
)

type Txpool interface {
	NewTx(tx *common.Transaction)
	Pop() *common.Transaction
	Nonce(addr common.Address) uint64
	SetStateRoot(stateRoot sha3.Hash)
	NotifyTxEvent(txs []*common.Transaction)
}
