package statemachine

import (
	"cxchain-2023131080/common"
	"cxchain-2023131080/statedb"
	"fmt"
)

type IMachine interface {
	Execute(state statedb.StateDB, tx common.Transaction)
	Execute1(state statedb.StateDB, tx common.Transaction) *common.Receiption
}

type StateMachine struct{}

// Execute1 返回交易执行回执
func (m StateMachine) Execute1(state statedb.StateDB, tx common.Transaction) *common.Receiption {
	from := tx.From()
	value := tx.TxData.Value
	gas := tx.TxData.Gas
	gasPrice := tx.TxData.GasPrice

	if gas < 21000 {
		fmt.Println("Gas too low")
		return &common.Receiption{
			TxHash:  tx.Hash(),
			Status:  0,
			GasUsed: 0,
		}
	}

	gasUsed := uint64(21000)
	cost := value + gasUsed*gasPrice
	fromAcc := state.GetAccount(from)
	if fromAcc == nil {
		fmt.Println("From account not found")
	}
	if fromAcc != nil {
		fmt.Println("From amount =", fromAcc.Amount, "required cost =", cost)
	}
	if fromAcc == nil || fromAcc.Amount < cost {
		fmt.Println("Insufficient balance")
		return &common.Receiption{TxHash: tx.Hash(), Status: 0, GasUsed: 0}
	}

	// 执行交易
	m.Execute(state, tx)

	return &common.Receiption{
		TxHash:  tx.Hash(),
		Status:  1,
		GasUsed: gasUsed,
	}
}

// Execute 处理交易逻辑
func (m StateMachine) Execute(state statedb.StateDB, tx common.Transaction) {
	from := tx.From()
	to := tx.TxData.To
	value := tx.TxData.Value
	gas := tx.TxData.Gas
	gasPrice := tx.TxData.GasPrice

	if gas < 21000 {
		return
	}
	gasUsed := uint64(21000)
	cost := value + gasUsed*gasPrice

	fromAcc := state.GetAccount(from)
	if fromAcc == nil || fromAcc.Amount < cost {
		return
	}
	fromAcc.Amount -= cost
	fromAcc.Nonce += 1
	state.UpdateAccount(fromAcc)

	toAcc := state.GetAccount(to)
	if toAcc == nil {
		toAcc = &common.Account{Address: to}
	}
	toAcc.Amount += value
	state.UpdateAccount(toAcc)
}
