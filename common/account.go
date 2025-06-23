package common

type Account struct {
	Address Address // 账户地址
	Amount  uint64  // 账户余额，使用 *big.Int 方便大数处理
	Nonce   uint64  // 账户交易计数器
}
