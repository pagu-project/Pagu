package utils_test

import (
	"testing"

	"github.com/kehiy/RoboPac/utils"
	"github.com/stretchr/testify/assert"
)

func TestAtomic(t *testing.T) {
	atomic := int64(10000000000)

	coin := utils.AtomicToCoin(int64(atomic))
	assert.Equal(t, int64(10), coin)

	atomic = utils.CoinToAtomic(int64(100))
	assert.Equal(t, int64(100000000000), atomic)
}
