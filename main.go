package main

import (
	"crypto/elliptic"
	"cxchain-2023131080/crypto"
	"cxchain-2023131080/statedb"
	"fmt"
	"os"

	"cxchain-2023131080/common"
	"cxchain-2023131080/kvstore/leveldb"
	"cxchain-2023131080/statemachine"
)

func main() {
	dbPath := "./testdb"
	_ = os.RemoveAll(dbPath)

	db, err := leveldb.Open(dbPath)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	state := statedb.NewMPTStateDB(db)

	// 生成密钥对
	privKey, _ := crypto.GenerateKey()
	pubBytes := elliptic.Marshal(privKey.PublicKey.Curve, privKey.PublicKey.X, privKey.PublicKey.Y)
	fromAddr := common.PubKeyToAddress(pubBytes)

	// 这里给收款地址随便用一个固定的地址（确保格式正确，长度20字节）
	toAddr, err := common.AddressFromHex("0xdef2def2def2def2def2def2def2def2def2def2")
	if err != nil {
		panic(err)
	}

	// 初始化账户（给from账户足够余额）
	fromAcc := &common.Account{Address: fromAddr, Amount: 1000000, Nonce: 0}
	toAcc := &common.Account{Address: toAddr, Amount: 0, Nonce: 0}

	state.UpdateAccount(fromAcc)
	state.UpdateAccount(toAcc)

	// 构造交易
	tx := common.Transaction{
		TxData: common.TxData{
			To:       toAddr,
			Nonce:    0,
			Value:    100000,
			Gas:      21000,
			GasPrice: 1,
		},
	}

	// 签名交易
	err = common.SignTransaction(&tx, privKey)
	if err != nil {
		panic(err)
	}

	// 执行交易
	machine := statemachine.StateMachine{}
	receipt := machine.Execute1(state, tx)

	fmt.Println("=== Transaction Receipt ===")
	fmt.Printf("TxHash:  %x\n", receipt.TxHash)
	fmt.Printf("Status:  %d\n", receipt.Status)
	fmt.Printf("GasUsed: %d\n", receipt.GasUsed)

	fmt.Println("=== Final Account States ===")
	fmt.Printf("From: %+v\n", state.GetAccount(fromAddr))
	fmt.Printf("To:   %+v\n", state.GetAccount(toAddr))
	fmt.Printf("State Root: %x\n", state.Commit())
}
