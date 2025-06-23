package common

import (
	"crypto/ecdsa"
	"cxchain-2023131080/crypto"
	"cxchain-2023131080/crypto/secp256k1"
	"cxchain-2023131080/crypto/sha3"
	"cxchain-2023131080/utils/hash"
	"cxchain-2023131080/utils/rlp"
	"math/big"
)

type Receiption struct {
	TxHash  hash.Hash
	Status  int
	GasUsed uint64
	// Logs
}
type Transaction struct {
	TxData    TxData
	Signature signature
}

type TxData struct {
	To       Address
	Nonce    uint64
	Value    uint64
	Gas      uint64
	GasPrice uint64
	Input    []byte
}

type signature struct {
	R, S *big.Int
	V    uint8
}

// SignTransaction 用标准库签名（返回R,S,V）
func SignTransaction(tx *Transaction, priv *ecdsa.PrivateKey) error {
	data, err := rlp.EncodeToBytes(tx.TxData)
	if err != nil {
		return err
	}
	msg := sha3.Keccak256(data)

	privBytes := crypto.FromECDSA(priv) // 转为 []byte
	sig, err := secp256k1.Sign(msg, privBytes)
	if err != nil {
		return err
	}

	tx.Signature.R = new(big.Int).SetBytes(sig[:32])
	tx.Signature.S = new(big.Int).SetBytes(sig[32:64])
	tx.Signature.V = sig[64] // 通常是 0 或 1

	return nil
}

func (tx Transaction) From() Address {
	// 检查签名字段是否为空
	if tx.Signature.R == nil || tx.Signature.S == nil {
		return Address{}
	}

	// 编码交易数据用于签名
	toSign, err := rlp.EncodeToBytes(tx.TxData)
	if err != nil {
		return Address{}
	}
	msg := sha3.Keccak256(toSign)

	// 构造 65 字节签名（不要修改 tx.Signature.V）
	sig := make([]byte, 65)
	copy(sig[0:32], tx.Signature.R.FillBytes(make([]byte, 32)))
	copy(sig[32:64], tx.Signature.S.FillBytes(make([]byte, 32)))
	sig[64] = tx.Signature.V

	// 恢复公钥
	pubKey, err := secp256k1.RecoverPubkey(msg, sig)
	if err != nil {
		return Address{}
	}

	return PubKeyToAddress(pubKey)
}

func (tx Transaction) Hash() hash.Hash {
	data, _ := rlp.EncodeToBytes(tx.TxData)
	return hash.Hash(sha3.Keccak256(data))
}
