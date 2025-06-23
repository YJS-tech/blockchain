package maker

import (
	"cxchain-2023131080/blockchain"
	"cxchain-2023131080/common"
	"cxchain-2023131080/statedb"
	"cxchain-2023131080/statemachine"
	"cxchain-2023131080/txpool"
	"math/big"
	"time"
)

type BlockProducerConfig struct {
	Duration   time.Duration
	Difficulty big.Int
	MaxTx      int64
	Coinbase   common.Address
}

type BlockProducer struct {
	txpool txpool.Txpool
	statdb statedb.StateDB
	config BlockProducerConfig

	chain blockchain.Blockchain
	m     statemachine.IMachine

	header *blockchain.Header
	block  *blockchain.Body

	interupt chan bool
}

func (producer BlockProducer) NewBlock() {
	producer.header = blockchain.NewHeader(producer.chain.CurrentHeader)
	// new Body
	// producer.statdb =
}

func (producer BlockProducer) pack() {
	t := time.After(producer.config.Duration)
	for {
		select {
		case <-producer.interupt:
			break
		case <-t:
			break
		// TODO 数量
		default:
			tx := producer.txpool.Pop()
			producer.m.Execute(producer.statdb, *tx)

		}
	}
}

func (producer BlockProducer) Interupt() {
	producer.interupt <- true
}

func Seal() {

}
