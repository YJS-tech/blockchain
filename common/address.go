package common

import (
	"cxchain-2023131080/crypto/sha3"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Address 表示20字节的账户地址类型
type Address [20]byte

// Bytes 将Address转换为字节切片
func (a Address) Bytes() []byte {
	return a[:]
}

// Hex 返回地址的十六进制字符串表示，带0x前缀
func (a Address) Hex() string {
	return "0x" + hex.EncodeToString(a[:])
}

// String 实现fmt.Stringer接口
func (a Address) String() string {
	return a.Hex()
}

// SetBytes 从字节切片设置地址值
func (a *Address) SetBytes(b []byte) error {
	if len(b) != 20 {
		return fmt.Errorf("address must be 20 bytes long, got %d", len(b))
	}
	copy(a[:], b)
	return nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (a *Address) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}
	// 移除0x前缀（如果存在）
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	return a.SetBytes(bytes)
}

// MarshalJSON 实现json.Marshaler接口
func (a Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Hex())
}

// AddressFromHex 从十六进制字符串创建Address
func AddressFromHex(s string) (Address, error) {
	var addr Address
	if len(s) > 2 && s[:2] == "0x" {
		s = s[2:]
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return addr, err
	}
	if err := addr.SetBytes(b); err != nil {
		return addr, err
	}
	return addr, nil
}

// PubKeyToAddress 计算地址（Keccak256(pubkey的x,y拼接，去掉首字节0x04），取后20字节）
func PubKeyToAddress(pub []byte) Address {
	hash := sha3.Keccak256(pub)
	var addr Address
	copy(addr[:], hash[len(hash)-20:])
	return addr
}
