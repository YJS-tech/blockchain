package txpool

import (
	"cxchain-2023131080/common"
	"cxchain-2023131080/crypto/sha3"
	"cxchain-2023131080/statedb"
	"sort"
)

type SortedTxs interface {
	GasPrice() uint64
	Push(tx *common.Transaction)
	Pop() *common.Transaction
	Nonce() uint64
	Len() int
}

type pandingTx []SortedTxs

type DefaultSortedTxs []*common.Transaction

func (txs DefaultSortedTxs) Len() int {
	return len(txs)
}

func (sorted DefaultSortedTxs) GasPrice() uint64 {
	return sorted[0].GasPrice
}

func (sorted *DefaultSortedTxs) Push(tx *common.Transaction) {
	*sorted = append(*sorted, tx)
	sort.Slice(*sorted, func(i, j int) bool {
		return (*sorted)[i].Nonce < (*sorted)[j].Nonce
	})
}

func (sorted *DefaultSortedTxs) Pop() *common.Transaction {
	if len(*sorted) == 0 {
		return nil
	}
	tx := (*sorted)[0]
	*sorted = (*sorted)[1:]
	return tx
}

func (sorted DefaultSortedTxs) Nonce() uint64 {
	return sorted[len(sorted)-1].Nonce + 1
}

func (p pandingTx) Len() int { return len(p) }

func (p pandingTx) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p pandingTx) Less(i, j int) bool {
	return p[i].GasPrice() < p[j].GasPrice()
}

type DefualtPool struct {
	StateDB statedb.StateDB
	all     map[sha3.Hash]bool
	txs     pandingTx
	panding map[common.Address][]SortedTxs
	queue   map[common.Address][]*common.Transaction
}

func NewDefualtPool(stateDB statedb.StateDB) *DefualtPool {
	return &DefualtPool{
		StateDB: stateDB,
		all:     make(map[sha3.Hash]bool),
		txs:     make(pandingTx, 0),
		panding: make(map[common.Address][]SortedTxs),
		queue:   make(map[common.Address][]*common.Transaction),
	}
}

func (pool *DefualtPool) NewTx(tx *common.Transaction) {
	account := pool.StateDB.GetAccount(tx.From())
	if account.Nonce >= tx.Nonce {
		return
	}

	blks := pool.panding[tx.From()]

	nonce := account.Nonce
	if len(blks) > 0 {
		lastBlk := blks[len(blks)-1]
		nonce = lastBlk.Nonce()
	}

	if tx.Nonce > nonce+1 {
		pool.addQueueTx(tx)
	} else if tx.Nonce == nonce+1 {
		//push
		pool.pushPandingTx(blks, tx)
	} else {
		//替换
	}
}

func (pool *DefualtPool) replacePandingTx(blks []SortedTxs, tx *common.Transaction) {
	for _, blk := range blks {
		txs := blk.(*DefaultSortedTxs)
		for j, oldTx := range *txs {
			if oldTx.Nonce == tx.Nonce && oldTx.GasPrice < tx.GasPrice {
				(*txs)[j] = tx
				sort.Slice(*txs, func(m, n int) bool {
					return (*txs)[m].Nonce < (*txs)[n].Nonce
				})
				return
			}
		}
	}
}

func (pool *DefualtPool) pushPandingTx(blks []SortedTxs, tx *common.Transaction) {
	if len(blks) == 0 {
		blk := &DefaultSortedTxs{}
		blk.Push(tx)
		blks = append(blks, blk)
		pool.panding[tx.From()] = blks
		pool.txs = append(pool.txs, blk)
		sort.Sort(pool.txs)
	} else {
		lastBlk := blks[len(blks)-1]
		if lastBlk.GasPrice() <= tx.GasPrice {
			lastBlk.Push(tx)
		} else {
			blk := &DefaultSortedTxs{}
			blk.Push(tx)
			blks = append(blks, blk)
			pool.panding[tx.From()] = blks
			pool.txs = append(pool.txs, blk)
			sort.Sort(pool.txs)
		}
	}
}

func (pool *DefualtPool) addQueueTx(tx *common.Transaction) {
	list := pool.queue[tx.From()]
	list = append(list, tx)
	// 按照Nonce值排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Nonce < list[j].Nonce
	})
	pool.queue[tx.From()] = list
}

func (pool *DefualtPool) Pop() *common.Transaction {
	if len(pool.txs) == 0 {
		return nil
	}
	// 从 gasPrice 最大的块中取出一个交易
	blk := pool.txs[len(pool.txs)-1]
	tx := blk.Pop()

	// 如果该 block 已空，则移除
	if blk.(*DefaultSortedTxs).Len() == 0 {
		pool.txs = pool.txs[:len(pool.txs)-1]
	}

	return tx
}

func (pool *DefualtPool) SetStateRoot(stateRoot sha3.Hash) {
	pool.StateDB.SetRoot(stateRoot[:]) // 注意这里转成 []byte
	// 清空旧交易池（可选：清空queue、panding）
	pool.all = make(map[sha3.Hash]bool)
	pool.txs = nil
	pool.panding = make(map[common.Address][]SortedTxs)
	pool.queue = make(map[common.Address][]*common.Transaction)
}

func (pool *DefualtPool) NotifyTxEvent(txs []*common.Transaction) {
	for _, tx := range txs {
		pool.NewTx(tx)
	}
}
