package sha3

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
)

// Hash 表示32字节的哈希值类型
type Hash [32]byte

// Keccak256 计算数据的Keccak-256哈希
func Keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	return h.Sum(nil)
}

// Bytes 将Hash转换为字节切片
func (h Hash) Bytes() []byte {
	return h[:]
}

// Hex 返回哈希的十六进制字符串表示
func (h Hash) Hex() string {
	return hex.EncodeToString(h[:])
}

// String 实现fmt.Stringer接口
func (h Hash) String() string {
	return h.Hex()
}

// SetBytes 从字节切片设置哈希值
func (h *Hash) SetBytes(b []byte) error {
	if len(b) != 32 {
		return fmt.Errorf("hash must be 32 bytes long, got %d", len(b))
	}
	copy(h[:], b)
	return nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (h *Hash) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	return h.SetBytes(bytes)
}

// MarshalJSON 实现json.Marshaler接口
func (h Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.Hex())
}
