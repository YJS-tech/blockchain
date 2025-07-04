package maker

import (
	"cxchain-2023131080/blockchain"
	"cxchain-2023131080/common"
	"cxchain-2023131080/statedb"
	"cxchain-2023131080/statemachine"
	"cxchain-2023131080/txpool"
	"cxchain-2023131080/utils/xtime"
	"time"
)

type ChainConfig struct {
	Duration   time.Duration
	Coinbase   common.Address
	Difficulty uint64
}

type BlockMaker struct {
	txpool txpool.Txpool
	state  statedb.StateDB
	exec   statemachine.IMachine

	config ChainConfig
	chain  blockchain.Blockchain

	nextHeader *blockchain.Header
	nextBody   *blockchain.Body

	interupt chan bool
}

func NewBlockMaker(txpool txpool.Txpool, state statedb.StateDB, exec statemachine.StateMachine) *BlockMaker {
	return &BlockMaker{
		txpool: txpool,
		state:  state,
		exec:   exec,
	}
}

func (maker BlockMaker) NewBlock() {
	maker.nextBody = blockchain.NewBlock()
	maker.nextHeader = blockchain.NewHeader(maker.chain.CurrentHeader)
	maker.nextHeader.Coinbase = maker.config.Coinbase
}

func (maker BlockMaker) Pack() {
	end := time.After(maker.config.Duration)
	for {
		select {
		case <-maker.interupt:
			break
		case <-end:
			break
		default:
			maker.pack()
		}
	}
}

func (maker BlockMaker) pack() {
	tx := maker.txpool.Pop()
	receiption := maker.exec.Execute1(maker.state, *tx)
	maker.nextBody.Transactions = append(maker.nextBody.Transactions, *tx)
	maker.nextBody.Receiptions = append(maker.nextBody.Receiptions, *receiption)
}

func (maker BlockMaker) Interupt() {
	maker.interupt <- true
}

func (maker BlockMaker) Finalize() (*blockchain.Header, *blockchain.Body) {
	maker.nextHeader.Timestamp = xtime.Now()
	maker.nextHeader.Nonce = 0
	// TODO
	// for n := 0; ; n++ {
	// 	maker.nextHeader.Nonce = 0
	// 	if maker.nextHeader.Hash() {
	// 		break
	// 	}
	// }

	return maker.nextHeader, maker.nextBody
}
