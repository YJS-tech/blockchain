package utils

// HexToNibbles 将byte数组转化成nibbles数组
func HexToNibbles(key []byte) []byte {
	var nibbles []byte
	for _, b := range key {
		nibbles = append(nibbles, b>>4, b&0x0F)
	}
	return nibbles
}
