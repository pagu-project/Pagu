package utils

func AtomicToCoin(amount int64) int64 {
	return amount / 1e9
}

func CoinToAtomic(amount int64) int64 {
	return amount * 1e9
}
